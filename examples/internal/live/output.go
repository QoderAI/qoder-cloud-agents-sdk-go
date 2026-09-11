package live

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"
)

// output renders the same events as a readable walkthrough or JSON lines.
type output struct {
	writer      io.Writer
	format, pat string
	mu          sync.Mutex
}

func (o *output) write(scenario, event, detail string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.pat != "" {
		detail = strings.ReplaceAll(detail, o.pat, "[REDACTED]")
	}
	// Replies and tool names are remote data; do not let terminal escape sequences
	// change previously printed output. Preserve ordinary line breaks and tabs.
	detail = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			return '\uFFFD'
		}
		return r
	}, detail)
	if o.format == "json" {
		_ = json.NewEncoder(o.writer).Encode(map[string]string{"time": time.Now().UTC().Format(time.RFC3339), "scenario": scenario, "event": event, "detail": detail})
		return
	}
	switch event {
	case "header", "start", "cleanup_start", "scenario_result", "summary":
		fmt.Fprintf(o.writer, "\n%s\n", detail)
	case "step":
		fmt.Fprintf(o.writer, "\n%s\n", detail)
	case "user", "assistant", "memory", "preview", "history":
		label := "你"
		if event == "assistant" {
			label = "助手"
		} else if event == "memory" {
			label = "通过 SDK 写入的记忆"
		} else if event == "preview" {
			label = "消息预览更新（替换上一版）"
		} else if event == "history" {
			label = "会话历史"
		}
		fmt.Fprintf(o.writer, "    %s：\n", label)
		for _, line := range strings.Split(strings.TrimSpace(detail), "\n") {
			fmt.Fprintf(o.writer, "      %s\n", line)
		}
	case "failed", "cleanup_failed":
		fmt.Fprintf(o.writer, "    ✗ %s\n", detail)
	case "created", "checked", "verified", "cleaned", "idle":
		fmt.Fprintf(o.writer, "    ✓ %s\n", detail)
	case "tool":
		fmt.Fprintf(o.writer, "    → 调用工具：%s\n", detail)
	case "waiting":
		fmt.Fprintf(o.writer, "    … %s\n", detail)
	default:
		fmt.Fprintf(o.writer, "    %s\n", detail)
	}
}

func (r *Run) Log(event, detail string) {
	if r.output == nil {
		r.output = &output{writer: os.Stdout}
	}
	r.output.write(r.Name, event, detail)
}

func (r *Run) Step(description string) {
	r.step++
	r.action = description
	r.Log("step", fmt.Sprintf("[%d] %s", r.step, description))
}

// Wait keeps a long SSE stream or asynchronous job visibly active. The returned
// function stops and joins its heartbeat before the next step or cleanup begins.
func (r *Run) Wait(description string) func() {
	r.Step(description)
	stop, done := make(chan struct{}), make(chan struct{})
	started := time.Now()
	go func() {
		defer close(done)
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				r.Log("waiting", fmt.Sprintf("%s，已等待 %.0f 秒", description, time.Since(started).Seconds()))
			}
		}
	}()
	var once sync.Once
	return func() { once.Do(func() { close(stop); <-done }) }
}

func (r *Run) Message(role, text string) {
	// Model/tool output may contain presigned links. They are unnecessary in a
	// walkthrough and should not be copied into terminal transcripts.
	text = webURL.ReplaceAllString(text, "[URL REDACTED]")
	if r.output != nil && r.output.pat != "" {
		text = strings.ReplaceAll(text, r.output.pat, "[REDACTED]")
	}
	runes := []rune(text)
	if len(runes) > 4000 {
		text = string(runes[:4000]) + "\n（显示内容已截断，校验仍使用完整回复）"
	}
	r.Log(role, text)
}

func (r *Run) ObserveTurn(result *TurnResult, raw string) error {
	if err := result.Observe(raw); err != nil {
		return err
	}
	var event struct{ Type, Name string }
	if err := json.Unmarshal([]byte(raw), &event); err != nil {
		return err
	}
	switch event.Type {
	case "agent.message":
		if result.Text != "" {
			r.Message("assistant", result.Text)
		}
	case "agent.tool_use", "agent.mcp_tool_use":
		if event.Name == "" {
			event.Name = "工具执行"
		}
		r.Log("tool", event.Name)
	case "session.status_idle":
		if result.Complete {
			r.Log("idle", "助手已结束本轮回复")
		}
	}
	return nil
}

var resourceNames = map[string]string{
	"environment": "Environment（执行环境）", "identity": "Identity（身份）", "template": "Template（模板）",
	"agent": "Agent（助手）", "session": "Session（会话）", "file": "File（文件）", "input_file": "输入文件",
	"skill": "Skill（技能）", "memory_store": "Memory Store（记忆库）", "input_memory_store": "输入记忆库",
	"memory_mount": "记忆库绑定", "schedule": "Schedule（任务配置）", "schedule_run": "Schedule Run（运行记录）",
	"deployment": "Deployment（部署配置）", "batch": "Batch（批处理任务）", "dream": "Dream（记忆整理任务）",
}

func resourceName(kind string) string {
	if name, ok := resourceNames[kind]; ok {
		return name
	}
	return kind
}
func cleanupLabel(mode, kind, id string) string {
	action := "删除"
	switch kind {
	case "environment", "template", "agent", "schedule", "deployment":
		action = "归档"
	case "session":
		if mode == "forward" {
			action = "归档"
		}
	case "memory_mount":
		action = "解除"
	case "schedule_run", "batch", "dream":
		action = "结束运行并清理"
	}
	return action + " " + resourceName(kind) + "：" + id
}

func Status(status string) string {
	names := map[string]string{"pending": "等待执行", "validating": "校验输入", "queued": "排队中", "running": "执行中", "processing": "处理中", "completed": "已完成", "failed": "失败", "skipped": "已跳过", "cancelled": "已取消", "canceled": "已取消", "expired": "已过期"}
	if name, ok := names[status]; ok {
		return name + "（" + status + "）"
	}
	return status
}
