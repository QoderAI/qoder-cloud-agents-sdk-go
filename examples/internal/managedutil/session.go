package managedutil

import (
	"context"
	"errors"
	"fmt"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

func (s *Example) Environment(ctx context.Context, r *live.Run) (string, error) {
	r.Step("创建云端执行环境（Environment）")
	env, err := s.Client.Environments.New(ctx, managed.EnvironmentNewParams{Name: live.Name("env"), Config: managed.EnvironmentNewParamsConfigUnion{OfCloud: &managed.CloudConfigParams{}}})
	if err != nil {
		return "", err
	}
	r.Track("environment", env.ID, func(ctx context.Context) error {
		_, err := s.Client.Environments.Archive(ctx, env.ID, managed.EnvironmentArchiveParams{})
		return err
	})
	return env.ID, nil
}
func (s *Example) Agent(ctx context.Context, r *live.Run, params managed.AgentNewParams) (string, error) {
	if s.Model == "" {
		if err := s.SelectModel(ctx, r); err != nil {
			return "", err
		}
	}
	params.Name = live.Name("agent")
	params.Model = managed.ManagedAgentsModelConfigParams{ID: s.Model}
	params.System = managed.String("你是一个 SDK 示例助手。帮助用户了解会话、文件、技能与记忆的用法，必要时调用工具。仅使用可实际读取的数据回答问题。")
	params.Tools = []managed.AgentNewParamsToolUnion{{OfAgentToolset20260401: &managed.ManagedAgentsAgentToolset20260401Params{Type: "agent_toolset_20260401"}}}
	r.Step("创建助手（Agent），配置模型和工具")
	item, err := s.Client.Agents.New(ctx, params)
	if err != nil {
		return "", err
	}
	r.Track("agent", item.ID, func(ctx context.Context) error {
		_, err := s.Client.Agents.Archive(ctx, item.ID, managed.AgentArchiveParams{})
		return err
	})
	return item.ID, nil
}
func (s *Example) NewSession(ctx context.Context, r *live.Run, params managed.SessionNewParams) (string, error) {
	r.Step("创建会话（Session），关联助手和执行环境")
	item, err := s.Client.Sessions.New(ctx, params)
	if err != nil {
		return "", err
	}
	r.Track("session", item.ID, func(ctx context.Context) error { return s.FinishSession(ctx, item.ID) })
	return item.ID, nil
}
func (s *Example) FinishSession(ctx context.Context, id string) error {
	item, err := s.Client.Sessions.Get(ctx, id, managed.SessionGetParams{})
	if err != nil {
		return err
	}
	if item.Status != "idle" && item.Status != "terminated" {
		if _, err = s.Client.Sessions.Events.Send(ctx, id, managed.SessionEventSendParams{Events: []managed.ManagedAgentsEventParamsUnion{{OfUserInterrupt: &managed.ManagedAgentsUserInterruptEventParams{Type: "user.interrupt"}}}}); err != nil {
			return err
		}
		for {
			item, err = s.Client.Sessions.Get(ctx, id, managed.SessionGetParams{})
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
	_, err = s.Client.Sessions.Delete(ctx, id, managed.SessionDeleteParams{})
	return err
}
func (s *Example) Send(ctx context.Context, id, prompt string) (string, error) {
	response, err := s.Client.Sessions.Events.Send(ctx, id, managed.SessionEventSendParams{Events: []managed.ManagedAgentsEventParamsUnion{{OfUserMessage: &managed.ManagedAgentsUserMessageEventParams{Type: "user.message", Content: []managed.ManagedAgentsUserMessageEventParamsContentUnion{{OfText: &managed.ManagedAgentsTextBlockParam{Type: "text", Text: prompt}}}}}}}, option.WithHeader("Idempotency-Key", live.Name("event")))
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
		stream := s.Client.Sessions.Events.StreamEvents(ctx, id, managed.SessionEventStreamParams{}, option.WithHeader("Last-Event-ID", after), option.WithRequestTimeout(s.Config.Timeout))
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
			params := managed.SessionEventListParams{Order: "asc", Limit: managed.Int(100)}
			if result.LastID != "" {
				params.AfterID = managed.String(result.LastID)
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
