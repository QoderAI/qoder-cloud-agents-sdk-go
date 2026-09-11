package forwardutil

import (
	"context"
	"errors"
	"fmt"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func (s *Example) Environment(ctx context.Context, r *live.Run) (string, error) {
	r.Step("创建云端执行环境（Environment）")
	env, err := s.Client.Environments.New(ctx, forward.EnvironmentNewParams{Name: live.Name("env"), Config: map[string]any{"type": "cloud"}})
	if err != nil {
		return "", err
	}
	r.Track("environment", env.ID, func(ctx context.Context) error { _, err := s.Client.Environments.Archive(ctx, env.ID); return err })
	return env.ID, nil
}
func (s *Example) Identity(ctx context.Context, r *live.Run) (string, error) {
	r.Step("创建用于本次示例的身份（Identity）")
	identity, err := s.Client.Identities.New(ctx, forward.IdentityNewParams{ExternalID: live.Name("identity"), Metadata: map[string]any{"suite": "sdk-example"}})
	if err != nil {
		return "", err
	}
	r.Track("identity", identity.ID, func(ctx context.Context) error { _, err := s.Client.Identities.Delete(ctx, identity.ID); return err })
	return identity.ID, nil
}
func (s *Example) Template(ctx context.Context, r *live.Run, params forward.TemplateNewParams) (string, error) {
	if s.Model == "" {
		if err := s.SelectModel(ctx, r); err != nil {
			return "", err
		}
	}
	params.Name = live.Name("template")
	params.Model = forward.ModelConfigUnionParam{OfString: forward.String(s.Model)}
	params.System = forward.String("你是一个 SDK 示例助手。帮助用户了解会话、文件、技能与记忆的用法，必要时调用工具。仅使用可实际读取的数据回答问题。")
	params.Tools = []forward.ToolParam{{Type: "agent_toolset_20260401"}}
	r.Step("创建助手模板（Template），配置模型和工具")
	item, err := s.Client.Templates.New(ctx, params)
	if err != nil {
		return "", err
	}
	r.Track("template", item.ID, func(ctx context.Context) error {
		_, err := s.Client.Templates.Archive(ctx, item.ID, forward.TemplateArchiveParams{})
		return err
	})
	return item.ID, nil
}
func (s *Example) NewSession(ctx context.Context, r *live.Run, params forward.SessionNewParams) (string, error) {
	r.Step("创建会话（Session），关联助手和执行环境")
	item, err := s.Client.Sessions.New(ctx, params)
	if err != nil {
		return "", err
	}
	r.Track("session", item.ID, func(ctx context.Context) error { return s.FinishSession(ctx, item.ID) })
	return item.ID, nil
}
func (s *Example) FinishSession(ctx context.Context, id string) error {
	item, err := s.Client.Sessions.Get(ctx, id)
	if err != nil {
		return err
	}
	if item.Status != "idle" && item.Status != "terminated" {
		if _, err = s.Client.Sessions.Cancel(ctx, id, forward.SessionCancelParams{}); err != nil {
			return err
		}
		for {
			item, err = s.Client.Sessions.Get(ctx, id)
			if err != nil {
				return err
			}
			if item.Status == "idle" || item.Status == "terminated" {
				break
			}
			if err = live.Pause(ctx); err != nil {
				return err
			}
		}
	}
	_, err = s.Client.Sessions.Archive(ctx, id, forward.SessionArchiveParams{})
	return err
}
func (s *Example) Send(ctx context.Context, id, prompt string) (string, error) {
	response, err := s.Client.Sessions.Events.Send(ctx, id, forward.SessionEventSendParams{Events: []forward.SessionEventParam{{Type: "user.message", Content: forward.EventContentUnionParam{OfBlocks: []forward.ContentBlockParam{{Type: "text", Text: forward.String(prompt)}}}}}, IdempotencyKey: forward.String(live.Name("event"))})
	if err != nil {
		return "", err
	}
	if len(response.Data) != 1 || response.Data[0].ID == "" {
		return "", errors.New("send did not return one user event ID")
	}
	return response.Data[0].ID, nil
}
func (s *Example) WaitReply(ctx context.Context, r *live.Run, id, after string, streaming bool) (live.TurnResult, error) {
	method := "轮询会话事件，等待助手回复"
	if streaming {
		method = "通过 SSE 接收助手回复"
	}
	stopWaiting := r.Wait(method)
	defer stopWaiting()
	result := live.TurnResult{LastID: after}
	if streaming {
		stream := s.Client.Sessions.Events.StreamEvents(ctx, id, forward.SessionEventStreamParams{LastEventID: forward.String(after), IncludeToolCalls: forward.Bool(true)}, option.WithRequestTimeout(s.Config.Timeout))
		defer stream.Close()
		for stream.Next() {
			if err := r.ObserveTurn(&result, stream.Current().RawJSON()); err != nil {
				return result, err
			}
			if result.Complete {
				break
			}
		}
		if err := stream.Err(); err != nil {
			return result, err
		}
	} else {
		for !result.Complete {
			params := forward.SessionEventListParams{Order: forward.String("asc"), Limit: forward.Int(100), IncludeToolCalls: forward.Bool(true)}
			if result.LastID != "" {
				params.AfterID = forward.String(result.LastID)
			}
			page := s.Client.Sessions.Events.ListAutoPaging(ctx, id, params)
			count := 0
			for page.Next() {
				count++
				if count > 2000 {
					return result, errors.New("event polling exceeded 2000 events")
				}
				if err := r.ObserveTurn(&result, page.Current().RawJSON()); err != nil {
					return result, err
				}
				if result.Complete {
					break
				}
			}
			if err := page.Err(); err != nil {
				return result, err
			}
			if !result.Complete {
				if err := live.Pause(ctx); err != nil {
					return result, err
				}
			}
		}
	}
	return result, nil
}
func (s *Example) WaitTurn(ctx context.Context, r *live.Run, id, after string, expected []string, tool, streaming bool) error {
	result, err := s.WaitReply(ctx, r, id, after, streaming)
	if err != nil {
		return err
	}
	r.Step("校验本轮回复")
	if err := result.Verify(expected, tool); err != nil {
		return err
	}
	r.Log("checked", fmt.Sprintf("回复包含 %d 个预期值，且本轮已正常结束", len(expected)))
	if tool {
		r.Log("checked", "已观察到实际工具调用")
	}
	return nil
}
func (s *Example) Turn(ctx context.Context, r *live.Run, id, prompt string, expected []string, tool, stream bool) error {
	r.Step("向会话发送消息")
	r.Message("user", prompt)
	after, err := s.Send(ctx, id, prompt)
	if err != nil {
		return err
	}
	return s.WaitTurn(ctx, r, id, after, expected, tool, stream)
}
