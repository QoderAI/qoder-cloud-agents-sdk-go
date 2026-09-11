# Qoder Cloud Agents Go SDK

使用 Go 调用 Qoder Cloud Agents，支持类型化参数、上下文取消、流式事件和请求配置。

SDK 提供两种客户端，共享鉴权、HTTP、错误处理、参数和分页工具：

| 模式 | Go 包 | 主要资源 |
| --- | --- | --- |
| Forward | `forward` | Identity、Template、Session、Schedule、Batch、Channel |
| Managed | `managed` | Agent、Session、Deployment、Dream |

本页介绍 Go SDK 的用法与配置。

- [Installation](#installation)
- [Requirements](#requirements)
- [Usage](#usage)
- [Request fields](#request-fields)
- [Response objects](#response-objects)
- [Error handling](#error-handling)
- [Retries](#retries)
- [Timeouts](#timeouts)
- [Long-running operations](#long-running-operations)
- [File uploads and downloads](#file-uploads-and-downloads)
- [Pagination](#pagination)
- [Request options](#request-options)
- [HTTP client customization](#http-client-customization)
- [Examples](#examples)
- [Development](#development)
- [License](#license)

## Installation

在应用的 Go module 中安装：

```sh
go get github.com/QoderAI/qoder-cloud-agents-sdk-go
```

应用应保留 `go.mod` 和 `go.sum` 中解析出的版本。

## Requirements

- Go **1.23 或更高版本**。
- 目标环境可用的 PAT，以及所调用资源的访问权限。

## Usage

### 创建客户端

下面是一个完整的 Forward 程序：从环境变量读取 PAT，连接 CN 环境并列出账号启用的模型。

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func main() {
	token := os.Getenv("QODER_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("请设置 QODER_ACCESS_TOKEN")
	}
	client := forward.NewClient(
		option.WithAccessToken(token),
		option.WithBaseURL("https://api.qoder.com.cn/api/v1/forward"),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	models, err := client.Models.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, model := range models.Data {
		if model.IsEnabled {
			fmt.Println(model.ID)
		}
	}
}
```

将有效令牌设置为 `QODER_ACCESS_TOKEN` 后即可运行。Managed 使用 `managed.NewClient` 和 CN 地址 `https://api.qoder.com.cn/api/v1/cloud`；其模型查询为 `client.Models.List(ctx, managed.ModelListParams{})`。

显式传入的 option 优先于环境变量。未指定地址或令牌时，客户端按下表读取：

| 配置 | Forward | Managed |
| --- | --- | --- |
| 访问令牌环境变量 | `QODER_ACCESS_TOKEN` | `QODER_ACCESS_TOKEN` |
| API 地址环境变量 | `QODER_FORWARD_BASE_URL` | `QODER_BASE_URL` |
| 默认 API 地址 | `https://api.qoder.com/api/v1/forward/` | `https://api.qoder.com/api/v1/cloud/` |

默认地址是国际站；访问 CN 时请显式配置 CN 地址。地址必须包含完整 API 根路径，末尾 `/` 可以省略。两种模式使用不同账号时，分别向各自客户端传入对应令牌。

SDK 客户端本身不读取 `.env` 文件。`examples/` 中的程序另有配置加载器，会读取 `.env.live` 的 `QODER_FORWARD_PAT`、`QODER_MANAGED_PAT` 等配置。

也可通过 `option.WithCredential(provider)` 使用实现 `convention.Credential` 的动态令牌提供者。提供者在请求时获取令牌；请求中已有的 `Authorization`（包括来自 `QODER_ACCESS_TOKEN` 的值）优先，因此采用动态凭据时应避免同时配置静态令牌。

以下示例为可放入应用的函数片段，省略 `package` 和 `import`。按片段使用的名称导入标准库；SDK 包路径统一为：

| 包名 | import 路径 |
| --- | --- |
| `forward` | `github.com/QoderAI/qoder-cloud-agents-sdk-go/forward` |
| `managed` | `github.com/QoderAI/qoder-cloud-agents-sdk-go/managed` |
| `convention` | `github.com/QoderAI/qoder-cloud-agents-sdk-go/convention` |
| `option` | `github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option` |
| `param` | `github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param` |

### Forward 会话

Forward 通过 Identity 和 Template 创建 Session。以下函数使用已有的 Identity 和 Template；Template 中需配置可用的执行环境。创建这些资源的完整流程见 [Forward 会话示例](examples/forward/session/main.go)。

```go
func startForwardSession(ctx context.Context, client forward.Client, identityID, templateID string) (*forward.Session, error) {
	return client.Sessions.New(ctx, forward.SessionNewParams{
		IdentityID: identityID,
		TemplateID: templateID,
		Title:      forward.String("项目助手"),
	})
}

func sendForwardMessage(ctx context.Context, client forward.Client, sessionID, text, requestKey string) (string, error) {
	result, err := client.Sessions.Events.Send(ctx, sessionID, forward.SessionEventSendParams{
		Events: []forward.SessionEventParam{{
			Type: "user.message",
			Content: forward.EventContentUnionParam{OfBlocks: []forward.ContentBlockParam{{
				Type: "text", Text: forward.String(text),
			}}},
		}},
		IdempotencyKey: forward.String(requestKey),
	})
	if err != nil {
		return "", err
	}
	if len(result.Data) != 1 {
		return "", fmt.Errorf("expected one user event, got %d", len(result.Data))
	}
	return result.Data[0].ID, nil
}
```

`requestKey` 是调用方为本次逻辑消息生成并保存的非空唯一键；重试同一条消息时复用该键，新消息使用新键。返回值是用户消息事件 ID，可用作接收后续事件的起点。

### Managed 会话

Managed 使用已有的 Agent 和 Environment 创建 Session。创建资源的完整流程见 [Managed 会话示例](examples/managed/session/main.go)。

```go
func startManagedSession(ctx context.Context, client managed.Client, agentID, environmentID string) (*managed.ManagedAgentsSession, error) {
	return client.Sessions.New(ctx, managed.SessionNewParams{
		EnvironmentID: environmentID,
		Agent: managed.SessionNewParamsAgentUnion{
			OfString: managed.String(agentID),
		},
	})
}

func sendManagedMessage(ctx context.Context, client managed.Client, sessionID, text, requestKey string) (string, error) {
	result, err := client.Sessions.Events.Send(ctx, sessionID, managed.SessionEventSendParams{
		Events: []managed.ManagedAgentsEventParamsUnion{{
			OfUserMessage: &managed.ManagedAgentsUserMessageEventParams{
				Type: "user.message",
				Content: []managed.ManagedAgentsUserMessageEventParamsContentUnion{{
					OfText: &managed.ManagedAgentsTextBlockParam{Type: "text", Text: text},
				}},
			},
		}},
	}, option.WithHeader("Idempotency-Key", requestKey))
	if err != nil {
		return "", err
	}
	if len(result.Data) != 1 {
		return "", fmt.Errorf("expected one user event, got %d", len(result.Data))
	}
	return result.Data[0].ID, nil
}
```

同一段连续对话复用 `sessionID`，继续发送新消息。`Events.Send` 成功表示消息已接收；Agent 的回答和执行状态通过事件列表或 SSE 获取。

### System prompts and tools

系统指令、模型和工具在 Forward Template 或 Managed Agent 上配置。`modelID` 应来自对应账号的 `Models.List`，并确认 `IsEnabled`。

```go
func createForwardTemplate(ctx context.Context, client forward.Client, environmentID, modelID string) (*forward.Template, error) {
	return client.Templates.New(ctx, forward.TemplateNewParams{
		Name:          "project-assistant",
		EnvironmentID: environmentID,
		Model:         forward.ModelConfigUnionParam{OfString: forward.String(modelID)},
		System:        forward.String("请使用中文回答，根据可读取的项目资料给出建议。"),
		Tools:         []forward.ToolParam{{Type: "agent_toolset_20260401"}},
	})
}

func createManagedAgent(ctx context.Context, client managed.Client, modelID string) (*managed.ManagedAgentsAgent, error) {
	return client.Agents.New(ctx, managed.AgentNewParams{
		Name:   "project-assistant",
		Model:  managed.ManagedAgentsModelConfigParams{ID: modelID},
		System: managed.String("请使用中文回答，根据可读取的项目资料给出建议。"),
		Tools: []managed.AgentNewParamsToolUnion{{
			OfAgentToolset20260401: &managed.ManagedAgentsAgentToolset20260401Params{
				Type: "agent_toolset_20260401",
			},
		}},
	})
}
```

内置工具由云端运行时执行，调用过程会产生工具事件。上述 Managed 示例中的事件、内容块和工具均显式设置 `Type`；构造请求时不能依赖这些字段自动补默认值。

### Streaming

两种模式都使用 `Sessions.Events.StreamEvents`。下面展示 Forward 接收完整消息事件的方式；`afterID` 传入刚发送的用户消息事件 ID，避免从历史事件中误判当前轮次。

```go
func readForwardTurn(ctx context.Context, client forward.Client, sessionID, afterID string) error {
	stream := client.Sessions.Events.StreamEvents(ctx, sessionID, forward.SessionEventStreamParams{
		LastEventID: forward.String(afterID),
	})
	defer stream.Close()

	for stream.Next() {
		event := stream.Current()
		switch event.Type {
		case "agent.message":
			blocks, err := event.ContentBlocks()
			if err != nil {
				return err
			}
			for _, block := range blocks {
				if block.Type == "text" {
					fmt.Println(block.Text)
				}
			}
		case "session.status_idle":
			fmt.Printf("本轮暂停或结束：%v\n", event.StopReason)
			return nil
		case "session.error", "session.status_terminated":
			return fmt.Errorf("session stopped: type=%s event_id=%s", event.Type, event.ID)
		}
	}
	if err := stream.Err(); err != nil {
		return err
	}
	return fmt.Errorf("stream ended before an idle event; last_event_id=%s", stream.LastEventID())
}
```

Managed 用请求头传入起点，通过响应联合类型访问消息内容：

```go
func readManagedTurn(ctx context.Context, client managed.Client, sessionID, afterID string) error {
	stream := client.Sessions.Events.StreamEvents(ctx, sessionID, managed.SessionEventStreamParams{},
		option.WithHeader("Last-Event-ID", afterID),
	)
	defer stream.Close()

	for stream.Next() {
		event := stream.Current()
		switch event.Type {
		case "agent.message":
			for _, block := range event.AsAgentMessage().Content {
				if block.Type == "text" {
					fmt.Println(block.Text)
				}
			}
		case "session.status_idle":
			fmt.Printf("本轮暂停或结束：%s\n", event.StopReason.Type)
			return nil
		case "session.error", "session.status_terminated":
			return fmt.Errorf("session stopped: type=%s event_id=%s", event.Type, event.ID)
		}
	}
	if err := stream.Err(); err != nil {
		return err
	}
	return fmt.Errorf("stream ended before an idle event; last_event_id=%s", stream.LastEventID())
}
```

调用方需关闭流，并在读取结束时检查 `Err()`。`session.status_idle` 可能表示正常结束，也可能表示等待调用方处理或达到预算等情况；业务成功应结合 `stop_reason` 和最终回答判断。上述函数用于展示事件消费，完整的执行断言见 [场景示例](examples/internal/live/turn.go)。

需要显示生成中的内容时，在 Stream Params 中设置 `EventDeltas`：Forward 使用 `[]string{"agent.message"}`，Managed 使用 `[]managed.ManagedAgentsDeltaType{"agent.message"}`。此时还会收到 `event_start`、`event_delta` 等预览事件；最终完整事件会再次携带完整内容。应用应更新同一条消息的预览，避免把增量和最终内容重复追加。同一个 ID 的多个增量事件不能仅按 ID 去重。

流中断后由调用方重新连接。保存 `stream.LastEventID()`，作为下次的 `LastEventID` 或 `Last-Event-ID`；SDK 不自动重连已建立的 SSE 流。恢复订阅时复用原 Session，不重新发送已被接收的用户消息。

## Request fields

必填字段通常使用普通值；可选标量使用 `param.Opt[T]`。例如 `forward.String`、`managed.String`、`Bool`、`Int` 和 `Float` 用于明确设置可选值。

| Go 值 | 请求中的含义 |
| --- | --- |
| 未设置的 `param.Opt[T]` | 省略该字段 |
| `forward.String("")` | 发送空字符串 |
| `forward.Bool(false)` | 发送 `false` |
| `forward.Int(0)` | 发送数字 `0` |
| `param.Null[string]()` | 发送 JSON `null` |
| nil map / slice（带 `omitzero`） | 省略该字段 |
| 非 nil 的空 map / slice | 发送 `{}` / `[]` |

`param.IsOmitted(value)` 和 `param.IsNull(value)` 可检查参数状态。结构体、map、slice 的显式 null 分别使用 `param.NullStruct[T]()`、`param.NullMap[T]()`、`param.NullSlice[T]()`。能否清空某个字段仍由该 API 的校验规则决定。

下面只演示序列化，不发送请求：

```go
func encodeOptionalFields() ([]byte, error) {
	params := forward.IdentityUpdateParams{
		Enabled: forward.Bool(false),
	}
	// Name 未设置，序列化时省略；Enabled 会明确发送 false。
	return json.Marshal(params)
}
```

### Request unions

联合参数通过 `Of...` 字段选择一种表示，同一个联合值只设置一个分支。例如 Forward 的模型参数可以是模型 ID 字符串，也可以是配置对象：

```go
func configuredModel(modelID string) forward.ModelConfigUnionParam {
	return forward.ModelConfigUnionParam{
		OfConfig: &forward.ModelConfigParam{
			ID:     modelID,
			Effort: forward.String("high"),
		},
	}
}
```

`Effort` 等模型选项需与所选模型支持的能力匹配。只传模型 ID 时使用 `forward.ModelConfigUnionParam{OfString: forward.String(modelID)}`。

### Extra fields and raw JSON

请求对象的 `SetExtraFields` 可补充 SDK 尚未建模、但服务端已支持的字段。同名 extra field 会覆盖结构体字段，应仅使用应用控制的字段名和值。

```go
func removeSessionMetadataKey() forward.SessionUpdateParams {
	params := forward.SessionUpdateParams{}
	params.SetExtraFields(map[string]any{
		"metadata": map[string]any{"obsolete": nil},
	})
	return params
}
```

需要按原始 JSON 发送请求时，可以使用 `param.SetJSON`。它控制整个对象的序列化，结构体已有字段不会再与原始 JSON 合并。服务端仍会校验请求内容。

```go
func paramsFromJSON(raw []byte) (managed.AgentNewParams, error) {
	var params managed.AgentNewParams
	if !json.Valid(raw) {
		return params, fmt.Errorf("invalid JSON")
	}
	param.SetJSON(raw, &params)
	return params, nil
}
```

不要依赖将任意 JSON 反序列化到参数联合类型后，所有 `Of...` 分支都能被还原。需要读取响应数据时使用响应类型；需要保持请求原始表示时使用 `SetJSON`。

## Response objects

响应字段可以直接访问；`JSON` 元信息用于区分未返回、null、无效类型与实际值：

```go
func inspectIdentity(identity forward.Identity) {
	fmt.Println(identity.ID)
	if identity.JSON.Name.Valid() {
		fmt.Println("展示名：", identity.Name)
	}
	switch identity.JSON.Name.Raw() {
	case "":
		fmt.Println("响应未包含 name")
	case "null":
		fmt.Println("name 为 null")
	}
}
```

`RawJSON()` 返回对象对应的原始 JSON。未知字段可通过 `result.JSON.ExtraFields["field_name"].Raw()` 读取。`Valid()` 为 false 时，还应结合 `Raw()` 判断是缺失、null 还是类型不匹配。

### Response unions

Managed 的响应联合类型提供 `AsAny()` 和 `As...()` 方法；应先根据判别字段确认类型。比如 `agent.message` 事件可使用 `event.AsAgentMessage()`，见上面的 Managed SSE 示例。

Forward 的部分开放内容直接保留为 `json.RawMessage`，例如 `SessionEvent.Content`。对于文本内容块数组，可调用 `ContentBlocks()`；处理其他事件时按事件类型解析，或保留 `RawJSON()`。两种模式均保留未知字段，但具体访问方法取决于对应响应类型。

## Error handling

HTTP API 返回错误状态时，SDK 返回 `*convention.Error`；`forward.Error` 和 `managed.Error` 是该类型的别名。通过 `errors.As` 读取状态码、错误码和 Request ID：

```go
func lookupSession(ctx context.Context, client forward.Client, sessionID string) error {
	_, err := client.Sessions.Get(ctx, sessionID)
	if err == nil {
		return nil
	}

	var apiErr *convention.Error
	switch {
	case errors.Is(err, context.Canceled):
		fmt.Println("调用已取消")
	case errors.Is(err, context.DeadlineExceeded):
		fmt.Println("等待请求超时")
	case errors.As(err, &apiErr):
		fmt.Printf("status=%d code=%s type=%s request_id=%s\n",
			apiErr.StatusCode, apiErr.Code, apiErr.Type(), apiErr.RequestID)
	default:
		fmt.Printf("网络、参数或解析错误：%v\n", err)
	}
	return err
}
```

| 字段或方法 | 内容 |
| --- | --- |
| `StatusCode` | HTTP 状态码 |
| `Code` | 服务端错误码，是否返回取决于错误响应 |
| `Message` | 服务端错误说明 |
| `Type()` | 错误类型 |
| `RequestID` | 定位请求的 ID，服务端未提供时可能为空 |
| `Request` / `Response` | 原始 HTTP 请求和响应 |
| `RawJSON()` | 收到的错误响应正文；网关返回非 JSON 时保留原文 |

网络、context、参数构造和解码错误不保证是 `*convention.Error`，需要保留其他错误处理分支。排查问题时优先记录状态码、错误码和 Request ID；原始请求可能包含 PAT，正文也可能包含业务数据，不宜直接输出到通用日志。

SSE 建连失败和读取错误由 `stream.Err()` 返回；业务事件中的 `session.error`、终止状态和异常 `stop_reason` 还需由调用方处理。

## Retries

默认最多重试 **2 次**，即最多进行 3 次 HTTP 尝试；请求体必须可重放，调用方 context 也必须尚未结束。

- `GET` / `HEAD`，以及带 `Idempotency-Key` 的其他请求：通常对连接错误、408、429、5xx 重试。
- 没有幂等键的其他请求：只允许对 429 重试。
- **409 不自动重试**。
- 在上述约束内，服务端 `x-should-retry` 可显式开启或禁止响应重试。

等待时间优先采用有效的 `Retry-After-Ms` / `Retry-After`，否则使用指数退避。重试不会重置调用方 context 的总期限。

```go
func clientWithoutRetries(token string) forward.Client {
	return forward.NewClient(
		option.WithAccessToken(token),
		option.WithBaseURL("https://api.qoder.com.cn/api/v1/forward"),
		option.WithMaxRetries(0),
	)
}

func getWithRetries(ctx context.Context, client forward.Client, sessionID string) (*forward.Session, error) {
	return client.Sessions.Get(ctx, sessionID, option.WithMaxRetries(3))
}
```

幂等键代表一次逻辑操作，应由调用方生成和保存。SDK 不会为所有写操作自动生成幂等键；也不应在一次重试中换用新键。

## Timeouts

SDK 默认不设置统一的请求超时，默认 HTTP 客户端为 `http.DefaultClient`。可以分别限制整个调用和每次 HTTP 尝试：

```go
func getWithTimeout(parent context.Context, client forward.Client, sessionID string) (*forward.Session, error) {
	ctx, cancel := context.WithTimeout(parent, time.Minute)
	defer cancel()
	return client.Sessions.Get(ctx, sessionID,
		option.WithRequestTimeout(15*time.Second),
	)
}
```

此例中，context 的 1 分钟覆盖所有尝试和退避；每次尝试最多 15 秒。先到达的期限生效。单次超时能否重试，仍受前述请求方法、幂等键和请求体规则约束。

SSE 和下载的超时也覆盖响应正文的读取。不要把适用于普通 API 的短超时直接用于需要长期保持的事件流；可在流式方法上覆盖 `WithRequestTimeout`，并为整个订阅设置合适的 context。

取消本地 context 会停止当前 HTTP 调用或事件订阅。若已创建云端执行任务，还需根据业务意图调用对应的取消或中断 API，随后确认服务端状态。

## Long-running operations

Session、Schedule、Batch、Deployment 和 Dream 的创建或触发响应，与执行完成是不同阶段。应通过事件流、Run 状态或任务结果确认最终执行情况。

| 场景 | 等待与结果入口 |
| --- | --- |
| Session | `Sessions.Events` 的流或事件列表，结合 idle 的 `stop_reason` 和最终回答 |
| Forward Schedule | `ScheduleRuns`，以及关联 Session 的结果 |
| Forward Batch | `Batches.Get`、`Batches.Tasks.List`、`Batches.GetOutput` |
| Managed Deployment | `DeploymentRuns`，以及关联 Session 的结果 |
| Managed Dream | Dream 状态，以及输出 Memory Store 的内容 |

例如 Batch 的 `completion_window: "24h"` 是提交给服务端的完成窗口；示例程序的 `-timeout 5m` 是本地等待期限。处于 `queued` 的 Batch 可能尚未开始执行，延长等待时间不会改变服务端调度规则。示例超时后会主动取消其创建的任务，不能用该运行的结果验证后续执行。

## File uploads and downloads

文件上传参数接收 `io.Reader`。使用 `convention.UploadFile` 指定文件名和媒体类型：

```go
func uploadText(ctx context.Context, client forward.Client, text string) (*forward.FileMetadata, error) {
	return client.Files.Upload(ctx, forward.FileUploadParams{
		File: convention.UploadFile{
			Reader:    strings.NewReader(text),
			Name:      "notes.txt",
			MediaType: "text/plain",
		},
		Purpose: forward.String("session_resource"),
	})
}
```

本地文件可以通过 `os.Open` 打开，传入 Reader，并在调用结束后关闭。Managed 的对应方法是 `client.Files.Upload(ctx, managed.FileUploadParams{File: ...})`，同样接受 `convention.UploadFile`。

上传文件与将文件挂载到 Session 是两个步骤；挂载方式见 [Forward 资源示例](examples/forward/resources/main.go) 和 [Managed 资源示例](examples/managed/resources/main.go)。Skill 文件树通过 `Skills.New` / `Skills.Versions.New` 上传，`files` 中的 Reader 使用相对路径文件名，例如 `my-skill/SKILL.md`。

文件下载返回 `*http.Response`，调用方负责读取并关闭 Body：

```go
func downloadFile(ctx context.Context, client forward.Client, fileID string, destination io.Writer) error {
	response, err := client.Files.Download(ctx, fileID)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	_, err = io.Copy(destination, response.Body)
	return err
}
```

Managed 的下载方法为 `client.Files.Download(ctx, fileID, managed.FileDownloadParams{})`。SDK 先从 API 获取临时下载地址，再请求存储内容；API 的 PAT、请求头和 middleware 不会被附加到这次存储请求，但自定义 HTTP transport 的行为仍由调用方控制。

## Pagination

`ListAutoPaging` 自动获取后续分页，使用 `Next()`、`Current()` 和 `Err()` 遍历：

```go
func listTemplates(ctx context.Context, client forward.Client) error {
	pager := client.Templates.ListAutoPaging(ctx, forward.TemplateListParams{
		Limit: forward.Int(20), Status: forward.String("active"),
	})
	for pager.Next() {
		fmt.Println(pager.Current().ID)
	}
	return pager.Err()
}
```

Managed 的列表使用相同的遍历方式，例如 `client.Agents.ListAutoPaging(ctx, managed.AgentListParams{Limit: managed.Int(20)})`。

也可以使用 `List` 获取单页，再通过 `GetNextPage()` 手动翻页：

```go
func listAgentPages(ctx context.Context, client managed.Client) error {
	page, err := client.Agents.List(ctx, managed.AgentListParams{Limit: managed.Int(20)})
	for err == nil && page != nil {
		for _, agent := range page.Data {
			fmt.Println(agent.ID)
		}
		page, err = page.GetNextPage()
	}
	return err
}
```

不同 API 使用 `after_id` / `before_id` 或 `next_page` 等游标。分页工具会维护后续请求所需的参数，并保留筛选条件；请使用对应响应的分页方法，不自行混用游标类型。

## Request options

`option` 包提供函数式配置，可传给客户端构造函数或单次方法调用。方法级选项覆盖客户端或服务级的同名设置；middleware、追加请求头等选项按自身的追加规则组合。共享客户端前应完成配置，避免并发修改 Options。

```go
func inspectResponse(ctx context.Context, client forward.Client, sessionID, traceID string) error {
	var response *http.Response
	_, err := client.Sessions.Get(ctx, sessionID,
		option.WithHeader("X-Trace-ID", traceID),
		option.WithResponseInto(&response),
	)
	if response != nil {
		fmt.Println("Request ID:", response.Header.Get("X-Request-ID"))
	}
	return err
}
```

| 选项 | 用途 |
| --- | --- |
| `WithBaseURL` | 设置完整 API 根地址 |
| `WithAccessToken` / `WithCredential` | 静态令牌或动态凭据 |
| `WithMaxRetries` / `WithRequestTimeout` | 重试次数和单次尝试超时 |
| `WithHeader` / `WithHeaderAdd` / `WithHeaderDel` | 设置、追加或删除请求头 |
| `WithQuery` / `WithQueryAdd` / `WithQueryDel` | 自定义查询参数 |
| `WithJSONSet` / `WithJSONDel` | 修改 JSON 请求体中的字段 |
| `WithResponseInto` | 获取原始 HTTP 响应及响应头 |
| `WithResponseBodyInto` | 替换默认的响应解码目标 |
| `WithRequestBody` | 提供自定义序列化请求体 |

`WithResponseInto` 不保证响应 Body 仍未被读取；需要原始 JSON 时优先使用实体的 `RawJSON()`。完整选项见 [requestoption.go](convention/option/requestoption.go)。

## HTTP client customization

使用 `WithHTTPClient` 配置连接池、代理或自定义 transport。下面保留默认 transport 的设置，并调整空闲连接数：

```go
func clientWithTransport(token string) managed.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = 20
	return managed.NewClient(
		option.WithAccessToken(token),
		option.WithBaseURL("https://api.qoder.com.cn/api/v1/cloud"),
		option.WithHTTPClient(&http.Client{Transport: transport}),
	)
}
```

`WithMiddleware` 可以拦截 API 请求，记录耗时等信息。middleware 在每次 HTTP 尝试中执行；下面只记录方法、路径、状态码和耗时：

```go
func requestTiming(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
	started := time.Now()
	response, err := next(req)
	status := 0
	if response != nil {
		status = response.StatusCode
	}
	log.Printf("method=%s path=%s status=%d duration=%s",
		req.Method, req.URL.Path, status, time.Since(started))
	return response, err
}
```

将 `option.WithMiddleware(requestTiming)` 传给 `NewClient` 即可。API middleware 不用于临时地址的存储下载；需要配置存储请求时使用合适的 HTTP client / transport。

## Examples

`examples/` 按功能提供 18 个独立程序，分别放在 `examples/forward/` 和 `examples/managed/` 下，每个场景目录都有自己的 `main.go`，直接运行目录即可。默认读取 SDK 仓库根目录的 `.env.live`。从 [.env.live.example](.env.live.example) 复制配置模板，并填写两种模式各自的 URL 和 PAT。

```sh
# 在 SDK 仓库根目录执行；查询模型是只读操作。
go run ./examples/forward/models
go run ./examples/managed/models

# 创建资源并执行一次真实对话；示例结束后清理本次资源。
go run ./examples/forward/session
go run ./examples/managed/session

# 写入记忆，在全新会话的回答中验证记忆生效。
go run ./examples/forward/memory
go run ./examples/managed/memory
```

| 场景 | Forward | Managed |
| --- | --- | --- |
| 模型与基础会话 | `models`、`session` | `models`、`session` |
| 多轮对话与文本增量 | `conversation`、`streaming-deltas` | `conversation`、`streaming-deltas` |
| 个性化与业务工具 | `identity-config`：共享模板的用户配置覆盖 | `custom-tools`：Go 函数执行与结果回传 |
| 文件、环境变量、Skill | `resources` | `resources` |
| 记忆 | `memory`：Identity + Template 绑定 | `memory`：Session 资源挂载 |
| 调度与批处理 | `schedule`、`batch` | `deployment`、`dream` |

记忆场景写入正文与 `MEMORY.md` 索引，提问不提供记忆中的答案，检查最终回复是否使用了保存的项目约定。

程序默认展示步骤、消息、助手回复和清理结果；`-output json` 输出结构化日志，`-timeout` 设置每个场景的等待期限。执行场景会产生实际资源和模型用量。

## Development

```sh
make test
make build
go test -race ./...
```

离线测试覆盖请求序列化、API 契约、响应解码、错误、重试、分页、SSE 和清理行为。`*_live_test.go` 需要 `live` build tag 及测试配置；可执行示例不需要该 tag。

- [Forward 使用与迁移说明](forward/README.md)
- [Managed 使用与迁移说明](managed/README.md)
- [公共包范围](managed/PACKAGE_SCOPE.md)

## License

Released under the [MIT License](LICENSE).
