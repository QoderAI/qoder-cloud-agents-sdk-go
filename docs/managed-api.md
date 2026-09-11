# Managed API 参考（Go）

本文覆盖当前 `managed` 客户端的 **95 个 HTTP 操作**，以及 **21 个自动分页便利方法**。方法签名、请求字段和返回类型以本仓库源码为准；每个操作附实现及类型链接，嵌套结构、联合分支和完整约束可直接进入字段定义查看。所有 HTTP 路由相对于 `/api/v1/cloud`。

[返回 README](../README.md) · [另一模式 API](forward-api.md)

## 初始化、鉴权与配置

要求 Go 1.23 或以上。安装方式见 [README](../README.md#installation)。以下是可直接编译的模型查询程序；PAT 从应用环境读取，不写入源码。

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
    "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

func main() {
    token := os.Getenv("QODER_ACCESS_TOKEN")
    if token == "" { log.Fatal("请设置 QODER_ACCESS_TOKEN") }
    client := managed.NewClient(
        option.WithAccessToken(token),
        option.WithBaseURL("https://api.qoder.com/api/v1/cloud"),
    )
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    models, err := client.Models.List(ctx, managed.ModelListParams{})
    if err != nil { log.Fatal(err) }
    for _, model := range models.Data {
        if model.IsEnabled { fmt.Println(model.ID) }
    }
}
```

| 配置 | 值与优先级 |
| --- | --- |
| PAT | `option.WithAccessToken(token)`；未指定时读取 `QODER_ACCESS_TOKEN`，使用 Bearer 鉴权 |
| 默认地址 | `https://api.qoder.com/api/v1/cloud/`（国际站） |
| CN 地址 | `https://api.qoder.com.cn/api/v1/cloud/`，通过 `WithBaseURL` 指定 |
| 地址环境变量 | `QODER_BASE_URL`；显式 option 优先于环境变量 |
| `.env` | SDK 不自动加载；可执行 examples 的配置加载器另行读取 `.env.live` |
| 动态凭据 | `option.WithCredential(provider)`，实现 `convention.Credential`；已有 Authorization 优先 |

全部资源从 `client` 访问，子资源保留层级，例如 `client.Sessions.Threads.Events`。在共享客户端前完成配置，避免并发修改 `Options`。

## 参数、响应与请求选项

每个操作均接收调用方的 `ctx context.Context`，末尾的 `opts ...option.RequestOption` 可省略。签名中的类型均属于 `managed` 包（`context`、`option`、`pagination`、`ssestream`、`http` 等显式前缀除外）。实际调用时使用 `managed.类型名` 构造请求。

参数表的“必填”来自 `api:"required"` 标记；路径参数也必须提供，且非空。可选不代表任意组合均有效：互斥字段、联合类型、枚举、权限及跨字段约束以链接的类型定义和服务端校验为准。空参数结构仍需在签名指定的位置传入。

- 可选标量使用 `managed.String("value")`、`managed.Int(20)`、`managed.Bool(false)`；`Bool(false)` 会发送 false，零值 `param.Opt` 则省略字段。
- 支持 null 的字段可用 `param.Null[T]()`；仅对声明允许清空的字段使用。数组和 map 的省略、空值、替换或合并规则见对应字段。
- 联合参数仅设置一个 `Of...` 分支；需要 `type` 判别字段时显式赋值，不能依赖零值自动补齐。
- 响应直接访问公开字段；`RawJSON()`、`JSON.Field.Valid()`、`JSON.Field.Raw()` 和 `JSON.ExtraFields` 用于读取原始值及区分缺失、null 与有效值。

| 常用 option | 用途 |
| --- | --- |
| `WithRequestTimeout(d)` | 每次 HTTP 尝试的超时，包含响应正文读取 |
| `WithMaxRetries(n)` | 最多重试次数；默认 2，即最多 3 次尝试 |
| `WithHeader(k, v)` / `WithHeaderAdd(k, v)` | 设置或追加请求头，例如幂等键、断点续流游标 |
| `WithResponseInto(&response)` | 读取原始 HTTP 响应及 Request ID；不保证 Body 尚未被消费 |
| `WithHTTPClient(client)` / `WithMiddleware(fn)` | 定制 transport、代理或 API 请求中间件 |
| `WithQuery` / `WithJSONSet` | 补充服务端已支持的查询参数或 JSON 字段 |

方法级选项覆盖同名客户端设置，追加选项按自身规则组合。默认不设置统一超时；建议设置 context 总期限，覆盖重试及退避。GET/HEAD 和带幂等键的请求可在请求体可重放时重试连接错误、408、429、5xx；未携带幂等键的写请求只对 429 重试，409 不自动重试。`x-should-retry` 在方法和幂等约束内控制响应重试，`Retry-After-Ms` / `Retry-After` 控制等待时间。一次逻辑写操作重试应复用同一 `Idempotency-Key`。

详细实现见 [请求选项](../convention/option/requestoption.go)、[请求执行](../convention/requestconfig.go)、[参数](../convention/param/param.go) 和 [错误类型](../convention/apierror/error.go)。

## 错误处理

HTTP 非成功状态返回 `*convention.Error`（`managed.Error` 是其别名）；网络、context、参数及解码错误可能是其他类型。SSE 的建连和读取错误必须通过 `stream.Err()` 检查。

以下函数片段需导入 `errors`、`fmt`、`context` 和 `github.com/QoderAI/qoder-cloud-agents-sdk-go/convention`：

```go
func describeError(err error) {
    var apiErr *convention.Error
    switch {
    case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
        fmt.Println("调用已取消或超时")
    case errors.As(err, &apiErr):
        fmt.Printf("status=%d code=%s type=%s request_id=%s\n",
            apiErr.StatusCode, apiErr.Code, apiErr.Type(), apiErr.RequestID)
    case err != nil:
        fmt.Println(err)
    }
}
```

`Message` 是错误说明，`RawJSON()` 保留错误正文。取消本地 context 只终止 HTTP 请求或订阅；云端已开始的工作需要调用相应取消 API，并确认状态。

## 创建会话、发送事件与 SSE

以下函数使用已有 Agent 和 Environment；资源准备及清理见 [会话示例](../examples/managed/session/main.go)。片段需导入 `context`、`fmt`、`managed` 和 `convention/option`。

```go
func startSession(ctx context.Context, client managed.Client, agentID, environmentID string) (*managed.ManagedAgentsSession, error) {
    return client.Sessions.New(ctx, managed.SessionNewParams{
        EnvironmentID: environmentID,
        Agent: managed.SessionNewParamsAgentUnion{OfString: managed.String(agentID)},
    })
}

func sendMessage(ctx context.Context, client managed.Client, sessionID, text, key string) (string, error) {
    result, err := client.Sessions.Events.Send(ctx, sessionID, managed.SessionEventSendParams{
        Events: []managed.ManagedAgentsEventParamsUnion{{
            OfUserMessage: &managed.ManagedAgentsUserMessageEventParams{
                Type: "user.message",
                Content: []managed.ManagedAgentsUserMessageEventParamsContentUnion{{
                    OfText: &managed.ManagedAgentsTextBlockParam{Type: "text", Text: text},
                }},
            },
        }},
    }, option.WithHeader("Idempotency-Key", key))
    if err != nil { return "", err }
    if len(result.Data) != 1 { return "", fmt.Errorf("expected one user event") }
    return result.Data[0].ID, nil
}

func readTurn(ctx context.Context, client managed.Client, sessionID, afterID string) error {
    stream := client.Sessions.Events.StreamEvents(ctx, sessionID, managed.SessionEventStreamParams{},
        option.WithHeader("Last-Event-ID", afterID))
    defer stream.Close()
    for stream.Next() {
        event := stream.Current()
        switch event.Type {
        case "agent.message":
            message := event.AsAgentMessage()
            for _, block := range message.Content {
                if block.Type == "text" { fmt.Print(block.Text) }
            }
        case "session.error", "session.status_terminated":
            return fmt.Errorf("session failed: %s", event.RawJSON())
        case "session.status_idle":
            return nil // 业务代码还应校验 stop_reason 和本轮最终回答。
        }
    }
    return stream.Err()
}
```

`key` 是调用方保存的非空幂等键；发送成功只表示消息接收成功。把返回的用户事件 ID 传给 `readTurn`，`Last-Event-ID` 通过请求选项设置。先按事件 `Type` 判断，再调用 `AsAgentMessage()` 等联合响应访问器。增量事件由 `EventDeltas` 配置，线程事件使用 `Sessions.Threads.Events.StreamEvents`；主动结束订阅时关闭流。

## 分页

```go
func listAll(ctx context.Context, client managed.Client) error {
    pager := client.Agents.ListAutoPaging(ctx, managed.AgentListParams{Limit: managed.Int(20)})
    for pager.Next() { fmt.Println(pager.Current().ID) }
    return pager.Err()
}
```

大部分列表返回 `pagination.PageCursor[T]`，通过响应 `NextPage` 设置下一次请求的 `Page`；Models 使用 `pagination.Page[T]` 的 `AfterID` / `BeforeID`，Sessions 使用 `pagination.BidirectionalPageCursor[T]`。自动分页统一使用 `Next()` / `Current()` / `Err()`；单页对象也可用 `GetNextPage()`，结束返回 `(nil, nil)`。以每个方法返回的分页类型为准，不混用游标。

## 上传与下载

以下片段需导入 `context`、`io`、`strings`、`convention` 和 `managed`。上传使用 multipart；上传完成不等于挂载，资源绑定见 [资源示例](../examples/managed/resources/main.go)。

```go
func uploadText(ctx context.Context, client managed.Client, text string) (*managed.FileMetadata, error) {
    return client.Files.Upload(ctx, managed.FileUploadParams{
        File: convention.UploadFile{
            Reader: strings.NewReader(text), Name: "notes.txt", MediaType: "text/plain",
        },
    })
}

func downloadFile(ctx context.Context, client managed.Client, fileID string, destination io.Writer) error {
    response, err := client.Files.Download(ctx, fileID, managed.FileDownloadParams{})
    if err != nil { return err }
    defer response.Body.Close()
    _, err = io.Copy(destination, response.Body)
    return err
}
```

SDK 先获取临时下载地址，再读取内容；不会把 API PAT、请求头或 middleware 附加到存储请求，自定义 HTTP transport 的行为仍由应用负责。`Skills.Versions.Download` 同样返回需关闭 Body 的 HTTP 响应。文件上传参数接受 `io.Reader`，本地文件需在调用结束后关闭。

## 资源与方法目录

下文列出全部 HTTP 操作以及每个资源实际提供的自动分页方法。参数表覆盖每个请求对象的公开一级字段，嵌入对象字段展开；点入类型可查看嵌套字段、允许值和完整说明。源码注释摘要保留原语言。

- [Agents](#agents)：6 个 HTTP 操作。
- [Sessions](#sessions)：19 个 HTTP 操作。
- [Deployments](#deployments)：8 个 HTTP 操作。
- [DeploymentRuns](#deploymentruns)：2 个 HTTP 操作。
- [Dreams](#dreams)：5 个 HTTP 操作。
- [Environments](#environments)：14 个 HTTP 操作。
- [Skills](#skills)：9 个 HTTP 操作。
- [Vaults](#vaults)：12 个 HTTP 操作。
- [Files](#files)：5 个 HTTP 操作。
- [MemoryStores](#memorystores)：14 个 HTTP 操作。
- [Models](#models)：1 个 HTTP 操作。

## Agents

管理 Agent 的模型、指令、工具与版本。更新可携带 Version 防止覆盖并发修改。

### `Agents.New`

```go
func (r *AgentService) New(ctx context.Context, params AgentNewParams, opts ...option.RequestOption) (res *ManagedAgentsAgent, err error)
```

**POST `/agents`** · [实现](../managed/agent.go#L47) · 创建资源。

返回：`*ManagedAgentsAgent, error`。[ManagedAgentsAgent](../managed/agent.go#L161)。

请求对象 `params`：[AgentNewParams](../managed/agent.go#L4281)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Model](../managed/agent.go#L4285)：`ManagedAgentsModelConfigParams` | json / `model` | 是 | Model identifier. Accepts the model string or a `model_config` object for additional configuration control；类型：[ManagedAgentsModelConfigParams](../managed/agent.go#L2996) |
| [Name](../managed/agent.go#L4287)：`string` | json / `name` | 是 | Human-readable name for the agent. |
| [Description](../managed/agent.go#L4289)：`param.Opt[string]` | json / `description` | 否 | Description of what the agent does. |
| [System](../managed/agent.go#L4291)：`param.Opt[string]` | json / `system` | 否 | System prompt for the agent. |
| [WorkspaceID](../managed/agent.go#L4292)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [MCPServers](../managed/agent.go#L4297)：`[]ManagedAgentsURLMCPServerParams` | json / `mcp_servers` | 否 | MCP servers this agent connects to. Maximum 20. Names must be unique within the array. Every server must be referenced by an `mcp_toolset` in `tools`; unreferenced servers are rejected. See the MCP connector guide.；类型：[ManagedAgentsURLMCPServerParams](../managed/agent.go#L3658) |
| [Metadata](../managed/agent.go#L4300)：`map[string]string` | json / `metadata` | 否 | Arbitrary key-value metadata. Maximum 16 pairs, keys up to 64 chars, values up to 512 chars. |
| [Multiagent](../managed/agent.go#L4303)：`ManagedAgentsMultiagentParams` | json / `multiagent` | 否 | A coordinator topology: the session's primary thread orchestrates work by spawning session threads, each running an agent drawn from the `agents` roster.；类型：[ManagedAgentsMultiagentParams](../managed/session.go#L1038) |
| [Skills](../managed/agent.go#L4305)：`[]ManagedAgentsSkillParamsUnion` | json / `skills` | 否 | Skills available to the agent.；类型：[ManagedAgentsSkillParamsUnion](../managed/agent.go#L3595) |
| [Tools](../managed/agent.go#L4308)：`[]AgentNewParamsToolUnion` | json / `tools` | 否 | Tool configurations available to the agent. Maximum of 128 tools across all toolsets allowed.；类型：[AgentNewParamsToolUnion](../managed/agent.go#L4325) |
| [Betas](../managed/agent.go#L4310)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Agents.Get`

```go
func (r *AgentService) Get(ctx context.Context, agentID string, params AgentGetParams, opts ...option.RequestOption) (res *ManagedAgentsAgent, err error)
```

**GET `/agents/{agent_id}`** · [实现](../managed/agent.go#L59) · 获取资源详情。

路径参数：`agentID string`，必填。

返回：`*ManagedAgentsAgent, error`。[ManagedAgentsAgent](../managed/agent.go#L161)。

请求对象 `params`：[AgentGetParams](../managed/agent.go#L4504)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Version](../managed/agent.go#L4507)：`param.Opt[int64]` | query / `version` | 否 | Agent version. Omit for the most recent version. Must be at least 1 if specified. |
| [WorkspaceID](../managed/agent.go#L4508)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/agent.go#L4510)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Agents.Update`

```go
func (r *AgentService) Update(ctx context.Context, agentID string, params AgentUpdateParams, opts ...option.RequestOption) (res *ManagedAgentsAgent, err error)
```

**POST `/agents/{agent_id}`** · [实现](../managed/agent.go#L75) · 更新资源。

路径参数：`agentID string`，必填。

返回：`*ManagedAgentsAgent, error`。[ManagedAgentsAgent](../managed/agent.go#L161)。

请求对象 `params`：[AgentUpdateParams](../managed/agent.go#L4522)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Description](../managed/agent.go#L4524)：`param.Opt[string]` | json / `description` | 否 | Description. Omit to preserve; send empty string or null to clear. |
| [System](../managed/agent.go#L4526)：`param.Opt[string]` | json / `system` | 否 | System prompt. Omit to preserve; send empty string or null to clear. |
| [Name](../managed/agent.go#L4528)：`param.Opt[string]` | json / `name` | 否 | Human-readable name. Must be non-empty. Omit to preserve. Cannot be cleared. |
| [Version](../managed/agent.go#L4533)：`param.Opt[int64]` | json / `version` | 否 | The agent's current version, used to prevent concurrent overwrites.（完整约束见字段源码） |
| [WorkspaceID](../managed/agent.go#L4534)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [MCPServers](../managed/agent.go#L4540)：`[]ManagedAgentsURLMCPServerParams` | json / `mcp_servers` | 否 | MCP servers.（完整约束见字段源码）；类型：[ManagedAgentsURLMCPServerParams](../managed/agent.go#L3658) |
| [Metadata](../managed/agent.go#L4544)：`map[string]any` | json / `metadata` | 否 | Metadata patch. Set a key to a string to upsert it, or to null to delete it. Omit the field to preserve. The stored bag is limited to 16 keys (up to 64 chars each) with values up to 512 chars. |
| [Skills](../managed/agent.go#L4546)：`[]ManagedAgentsSkillParamsUnion` | json / `skills` | 否 | Skills. Full replacement. Omit to preserve; send empty array or null to clear.；类型：[ManagedAgentsSkillParamsUnion](../managed/agent.go#L3595) |
| [Tools](../managed/agent.go#L4550)：`[]AgentUpdateParamsToolUnion` | json / `tools` | 否 | Tool configurations available to the agent. Full replacement. Omit to preserve; send empty array or null to clear. Maximum of 128 tools across all toolsets allowed.；类型：[AgentUpdateParamsToolUnion](../managed/agent.go#L4575) |
| [Model](../managed/agent.go#L4555)：`ManagedAgentsModelConfigParams` | json / `model` | 否 | Model identifier. Accepts the model string or a `model_config` object for additional configuration control. Omit to preserve. Cannot be cleared.；类型：[ManagedAgentsModelConfigParams](../managed/agent.go#L2996) |
| [Multiagent](../managed/agent.go#L4558)：`ManagedAgentsMultiagentParams` | json / `multiagent` | 否 | A coordinator topology: the session's primary thread orchestrates work by spawning session threads, each running an agent drawn from the `agents` roster.；类型：[ManagedAgentsMultiagentParams](../managed/session.go#L1038) |
| [Betas](../managed/agent.go#L4560)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Agents.List`

```go
func (r *AgentService) List(ctx context.Context, params AgentListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsAgent], err error)
```

**GET `/agents`** · [实现](../managed/agent.go#L91) · 列出资源。

返回：`*pagination.PageCursor[ManagedAgentsAgent], error`。[ManagedAgentsAgent](../managed/agent.go#L161)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[AgentListParams](../managed/agent.go#L4754)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [CreatedAtGte](../managed/agent.go#L4756)：`param.Opt[time.Time]` | query / `created_at[gte]` | 否 | Return agents created at or after this time (inclusive). |
| [CreatedAtLte](../managed/agent.go#L4758)：`param.Opt[time.Time]` | query / `created_at[lte]` | 否 | Return agents created at or before this time (inclusive). |
| [IncludeArchived](../managed/agent.go#L4760)：`param.Opt[bool]` | query / `include_archived` | 否 | Include archived agents in results. Defaults to false. |
| [Limit](../managed/agent.go#L4762)：`param.Opt[int64]` | query / `limit` | 否 | Maximum results per page. Default 20, maximum 100. |
| [Page](../managed/agent.go#L4764)：`param.Opt[string]` | query / `page` | 否 | Opaque pagination cursor from a previous response. |
| [WorkspaceID](../managed/agent.go#L4765)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/agent.go#L4767)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Agents.ListAutoPaging`

```go
func (r *AgentService) ListAutoPaging(ctx context.Context, params AgentListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsAgent]
```

[实现](../managed/agent.go#L112)。自动遍历 `Agents.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Agents.Archive`

```go
func (r *AgentService) Archive(ctx context.Context, agentID string, body AgentArchiveParams, opts ...option.RequestOption) (res *ManagedAgentsAgent, err error)
```

**POST `/agents/{agent_id}/archive`** · [实现](../managed/agent.go#L117) · 归档资源。

路径参数：`agentID string`，必填。

返回：`*ManagedAgentsAgent, error`。[ManagedAgentsAgent](../managed/agent.go#L161)。

请求对象 `body`：[AgentArchiveParams](../managed/agent.go#L4779)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/agent.go#L4780)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/agent.go#L4782)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Agents.Versions.List`

```go
func (r *AgentVersionService) List(ctx context.Context, agentID string, params AgentVersionListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsAgent], err error)
```

**GET `/agents/{agent_id}/versions`** · [实现](../managed/agentversion.go#L39) · 列出资源。

路径参数：`agentID string`，必填。

返回：`*pagination.PageCursor[ManagedAgentsAgent], error`。[ManagedAgentsAgent](../managed/agent.go#L161)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[AgentVersionListParams](../managed/agentversion.go#L68)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Limit](../managed/agentversion.go#L70)：`param.Opt[int64]` | query / `limit` | 否 | Maximum results per page. Default 20, maximum 100. |
| [Page](../managed/agentversion.go#L72)：`param.Opt[string]` | query / `page` | 否 | Opaque pagination cursor. |
| [WorkspaceID](../managed/agentversion.go#L73)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/agentversion.go#L75)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Agents.Versions.ListAutoPaging`

```go
func (r *AgentVersionService) ListAutoPaging(ctx context.Context, agentID string, params AgentVersionListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsAgent]
```

[实现](../managed/agentversion.go#L64)。自动遍历 `Agents.Versions.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

## Sessions

管理会话生命周期、消息事件、资源和子线程。发送消息或创建会话成功表示已接收，执行结果应通过事件确认。

### `Sessions.New`

```go
func (r *SessionService) New(ctx context.Context, params SessionNewParams, opts ...option.RequestOption) (res *ManagedAgentsSession, err error)
```

**POST `/sessions`** · [实现](../managed/session.go#L50) · 创建资源。

返回：`*ManagedAgentsSession, error`。[ManagedAgentsSession](../managed/session.go#L1220)。

请求对象 `params`：[SessionNewParams](../managed/session.go#L2465)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [EnvironmentVariables](../managed/session.go#L2466)：`map[string]string` | json / `environment_variables` | 否 | 见字段定义。 |
| [Agent](../managed/session.go#L2470)：`SessionNewParamsAgentUnion` | json / `agent` | 是 | Agent identifier. Accepts the `agent` ID string, which pins the latest version for the session, or an `agent` object with both id and version specified.；类型：[SessionNewParamsAgentUnion](../managed/session.go#L2505) |
| [EnvironmentID](../managed/session.go#L2472)：`string` | json / `environment_id` | 是 | ID of the `environment` defining the container configuration for this session. |
| [Title](../managed/session.go#L2474)：`param.Opt[string]` | json / `title` | 否 | Human-readable session title. |
| [WorkspaceID](../managed/session.go#L2475)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Budget](../managed/session.go#L2478)：`ManagedAgentsBudgetLimitParam` | json / `budget` | 否 | A hard spend ceiling. The session stops issuing new model requests once the tracked list cost reaches `max_list_cost`.；类型：[ManagedAgentsBudgetLimitParam](../managed/session.go#L585) |
| [InitialEvents](../managed/session.go#L2481)：`[]SessionNewParamsInitialEventUnion` | json / `initial_events` | 否 | Initial events to send to the `session` at creation, processed in order. Supports `user.message` and `user.define_outcome` events. Maximum 50 events.；类型：[SessionNewParamsInitialEventUnion](../managed/session.go#L2603) |
| [Metadata](../managed/session.go#L2484)：`map[string]string` | json / `metadata` | 否 | Arbitrary key-value metadata attached to the session. Maximum 16 pairs, keys up to 64 chars, values up to 512 chars. |
| [Resources](../managed/session.go#L2486)：`[]SessionNewParamsResourceUnion` | json / `resources` | 否 | Resources (e.g. repositories, files) to mount into the session's container.；类型：[SessionNewParamsResourceUnion](../managed/session.go#L2678) |
| [VaultIDs](../managed/session.go#L2488)：`[]string` | json / `vault_ids` | 否 | Vault IDs for stored credentials the agent can use during the session. |
| [Betas](../managed/session.go#L2490)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Sessions.Get`

```go
func (r *SessionService) Get(ctx context.Context, sessionID string, query SessionGetParams, opts ...option.RequestOption) (res *ManagedAgentsSession, err error)
```

**GET `/sessions/{session_id}`** · [实现](../managed/session.go#L62) · 获取资源详情。

路径参数：`sessionID string`，必填。

返回：`*ManagedAgentsSession, error`。[ManagedAgentsSession](../managed/session.go#L1220)。

请求对象 `query`：[SessionGetParams](../managed/session.go#L2790)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/session.go#L2791)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/session.go#L2793)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Sessions.Update`

```go
func (r *SessionService) Update(ctx context.Context, sessionID string, params SessionUpdateParams, opts ...option.RequestOption) (res *ManagedAgentsSession, err error)
```

**POST `/sessions/{session_id}`** · [实现](../managed/session.go#L78) · 更新资源。

路径参数：`sessionID string`，必填。

返回：`*ManagedAgentsSession, error`。[ManagedAgentsSession](../managed/session.go#L1220)。

请求对象 `params`：[SessionUpdateParams](../managed/session.go#L2797)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [EnvironmentVariables](../managed/session.go#L2798)：`map[string]string` | json / `environment_variables` | 否 | 见字段定义。 |
| [Title](../managed/session.go#L2801)：`param.Opt[string]` | json / `title` | 否 | Human-readable session title. |
| [WorkspaceID](../managed/session.go#L2802)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Metadata](../managed/session.go#L2805)：`map[string]any` | json / `metadata` | 否 | Metadata patch. Set a key to a string to upsert it, or to null to delete it. Omit the field to preserve. |
| [Agent](../managed/session.go#L2809)：`ManagedAgentsSessionAgentUpdateParam` | json / `agent` | 否 | Mid-session agent configuration update.（完整约束见字段源码）；类型：[ManagedAgentsSessionAgentUpdateParam](../managed/session.go#L1590) |
| [Budget](../managed/session.go#L2812)：`ManagedAgentsBudgetLimitParam` | json / `budget` | 否 | A hard spend ceiling. The session stops issuing new model requests once the tracked list cost reaches `max_list_cost`.；类型：[ManagedAgentsBudgetLimitParam](../managed/session.go#L585) |
| [VaultIDs](../managed/session.go#L2815)：`[]string` | json / `vault_ids` | 否 | Vault IDs (`vlt_*`) to attach to the session. Not yet supported; requests setting this field are rejected. Reserved for future use. |
| [Betas](../managed/session.go#L2817)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Sessions.List`

```go
func (r *SessionService) List(ctx context.Context, params SessionListParams, opts ...option.RequestOption) (res *pagination.BidirectionalPageCursor[ManagedAgentsSession], err error)
```

**GET `/sessions`** · [实现](../managed/session.go#L94) · 列出资源。

返回：`*pagination.BidirectionalPageCursor[ManagedAgentsSession], error`。[ManagedAgentsSession](../managed/session.go#L1220)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[SessionListParams](../managed/session.go#L2829)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [AgentID](../managed/session.go#L2831)：`param.Opt[string]` | query / `agent_id` | 否 | Filter sessions created with this agent ID. |
| [AgentVersion](../managed/session.go#L2833)：`param.Opt[int64]` | query / `agent_version` | 否 | Filter by agent version. Only applies when `agent_id` is also set. |
| [CreatedAtGt](../managed/session.go#L2835)：`param.Opt[time.Time]` | query / `created_at[gt]` | 否 | Return sessions created after this time (exclusive). |
| [CreatedAtGte](../managed/session.go#L2837)：`param.Opt[time.Time]` | query / `created_at[gte]` | 否 | Return sessions created at or after this time (inclusive). |
| [CreatedAtLt](../managed/session.go#L2839)：`param.Opt[time.Time]` | query / `created_at[lt]` | 否 | Return sessions created before this time (exclusive). |
| [CreatedAtLte](../managed/session.go#L2841)：`param.Opt[time.Time]` | query / `created_at[lte]` | 否 | Return sessions created at or before this time (inclusive). |
| [DeploymentID](../managed/session.go#L2843)：`param.Opt[string]` | query / `deployment_id` | 否 | Filter sessions created by this deployment ID. |
| [IncludeArchived](../managed/session.go#L2845)：`param.Opt[bool]` | query / `include_archived` | 否 | When true, includes archived sessions. Default: false (exclude archived). |
| [Limit](../managed/session.go#L2847)：`param.Opt[int64]` | query / `limit` | 否 | Maximum number of results to return. |
| [MemoryStoreID](../managed/session.go#L2850)：`param.Opt[string]` | query / `memory_store_id` | 否 | Filter sessions whose resources contain a `memory_store` with this memory store ID. |
| [Page](../managed/session.go#L2852)：`param.Opt[string]` | query / `page` | 否 | Opaque pagination cursor from a previous response. |
| [WorkspaceID](../managed/session.go#L2853)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Order](../managed/session.go#L2858)：`SessionListParamsOrder` | query / `order` | 否 | Sort direction for results, ordered by `created_at`. Defaults to `desc` (newest first). Any of "asc", "desc".；类型：[SessionListParamsOrder](../managed/session.go#L2879) |
| [Statuses](../managed/session.go#L2863)：`[]string` | query / `statuses` | 否 | Filter by session status. Repeat the parameter to match any of multiple statuses. Any of "rescheduling", "running", "idle", "terminated". |
| [Betas](../managed/session.go#L2865)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Sessions.ListAutoPaging`

```go
func (r *SessionService) ListAutoPaging(ctx context.Context, params SessionListParams, opts ...option.RequestOption) *pagination.BidirectionalPageCursorAutoPager[ManagedAgentsSession]
```

[实现](../managed/session.go#L115)。自动遍历 `Sessions.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Sessions.Delete`

```go
func (r *SessionService) Delete(ctx context.Context, sessionID string, body SessionDeleteParams, opts ...option.RequestOption) (res *ManagedAgentsDeletedSession, err error)
```

**DELETE `/sessions/{session_id}`** · [实现](../managed/session.go#L120) · 删除资源。

路径参数：`sessionID string`，必填。

返回：`*ManagedAgentsDeletedSession, error`。[ManagedAgentsDeletedSession](../managed/session.go#L676)。

请求对象 `body`：[SessionDeleteParams](../managed/session.go#L2886)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/session.go#L2887)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/session.go#L2889)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Sessions.Archive`

```go
func (r *SessionService) Archive(ctx context.Context, sessionID string, body SessionArchiveParams, opts ...option.RequestOption) (res *ManagedAgentsSession, err error)
```

**POST `/sessions/{session_id}/archive`** · [实现](../managed/session.go#L136) · 归档资源。

路径参数：`sessionID string`，必填。

返回：`*ManagedAgentsSession, error`。[ManagedAgentsSession](../managed/session.go#L1220)。

请求对象 `body`：[SessionArchiveParams](../managed/session.go#L2893)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/session.go#L2894)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/session.go#L2896)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Sessions.Events.List`

```go
func (r *SessionEventService) List(ctx context.Context, sessionID string, params SessionEventListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsSessionEventUnion], err error)
```

**GET `/sessions/{session_id}/events`** · [实现](../managed/sessionevent.go#L44) · 列出资源。

路径参数：`sessionID string`，必填。

返回：`*pagination.PageCursor[ManagedAgentsSessionEventUnion], error`。[ManagedAgentsSessionEventUnion](../managed/sessionevent.go#L3756)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[SessionEventListParams](../managed/sessionevent.go#L7560)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [BeforeID](../managed/sessionevent.go#L7561)：`param.Opt[string]` | query / `before_id` | 否 | 见字段定义。 |
| [AfterID](../managed/sessionevent.go#L7562)：`param.Opt[string]` | query / `after_id` | 否 | 见字段定义。 |
| [CreatedAtGt](../managed/sessionevent.go#L7566)：`param.Opt[time.Time]` | query / `created_at[gt]` | 否 | Return events created after this time (exclusive). Compared against the event's `processed_at` value. |
| [CreatedAtGte](../managed/sessionevent.go#L7569)：`param.Opt[time.Time]` | query / `created_at[gte]` | 否 | Return events created at or after this time (inclusive). Compared against the event's `processed_at` value. |
| [CreatedAtLt](../managed/sessionevent.go#L7572)：`param.Opt[time.Time]` | query / `created_at[lt]` | 否 | Return events created before this time (exclusive). Compared against the event's `processed_at` value. |
| [CreatedAtLte](../managed/sessionevent.go#L7575)：`param.Opt[time.Time]` | query / `created_at[lte]` | 否 | Return events created at or before this time (inclusive). Compared against the event's `processed_at` value. |
| [Limit](../managed/sessionevent.go#L7577)：`param.Opt[int64]` | query / `limit` | 否 | Query parameter for limit |
| [Page](../managed/sessionevent.go#L7579)：`param.Opt[string]` | query / `page` | 否 | Opaque pagination cursor from a previous response's `next_page`. |
| [WorkspaceID](../managed/sessionevent.go#L7580)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Order](../managed/sessionevent.go#L7585)：`SessionEventListParamsOrder` | query / `order` | 否 | Sort direction for results, ordered by the event's `processed_at`. Defaults to `asc` (chronological). Any of "asc", "desc".；类型：[SessionEventListParamsOrder](../managed/sessionevent.go#L7605) |
| [Types](../managed/sessionevent.go#L7588)：`[]string` | query / `types` | 否 | Filter by event type. Values match the `type` field on returned events (for example, `user.message` or `agent.tool_use`). Omit to return all event types. |
| [Betas](../managed/sessionevent.go#L7590)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Sessions.Events.ListAutoPaging`

```go
func (r *SessionEventService) ListAutoPaging(ctx context.Context, sessionID string, params SessionEventListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsSessionEventUnion]
```

[实现](../managed/sessionevent.go#L69)。自动遍历 `Sessions.Events.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Sessions.Events.Send`

```go
func (r *SessionEventService) Send(ctx context.Context, sessionID string, params SessionEventSendParams, opts ...option.RequestOption) (res *ManagedAgentsSendSessionEvents, err error)
```

**POST `/sessions/{session_id}/events`** · [实现](../managed/sessionevent.go#L74) · 发送会话事件。

路径参数：`sessionID string`，必填。

返回：`*ManagedAgentsSendSessionEvents, error`。[ManagedAgentsSendSessionEvents](../managed/sessionevent.go#L3230)。

请求对象 `params`：[SessionEventSendParams](../managed/sessionevent.go#L7612)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Events](../managed/sessionevent.go#L7614)：`[]ManagedAgentsEventParamsUnion` | json / `events` | 是 | Events to send to the `session`.；类型：[ManagedAgentsEventParamsUnion](../managed/sessionevent.go#L1751) |
| [WorkspaceID](../managed/sessionevent.go#L7615)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/sessionevent.go#L7617)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Sessions.Events.StreamEvents`

```go
func (r *SessionEventService) StreamEvents(ctx context.Context, sessionID string, params SessionEventStreamParams, opts ...option.RequestOption) (stream *ssestream.Stream[ManagedAgentsStreamSessionEventsUnion])
```

**GET `/sessions/{session_id}/events/stream`** · [实现](../managed/sessionevent.go#L90) · 订阅 SSE 事件流。

路径参数：`sessionID string`，必填。

返回：`*ssestream.Stream[ManagedAgentsStreamSessionEventsUnion]`。[ManagedAgentsStreamSessionEventsUnion](../managed/sessionevent.go#L5229)、[事件流](../convention/ssestream/ssestream.go)。

请求对象 `params`：[SessionEventStreamParams](../managed/sessionevent.go#L7629)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/sessionevent.go#L7630)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [EventDeltas](../managed/sessionevent.go#L7641)：`[]ManagedAgentsDeltaType` | query / `event_deltas` | 否 | When set, this connection also receives streaming deltas (`event_start`, `event_delta`) while an event is being produced, before the event itself arrives.（完整约束见字段源码）；类型：[ManagedAgentsDeltaType](../managed/session.go#L771) |
| [Betas](../managed/sessionevent.go#L7643)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Sessions.Resources.Get`

```go
func (r *SessionResourceService) Get(ctx context.Context, resourceID string, params SessionResourceGetParams, opts ...option.RequestOption) (res *SessionResourceGetResponseUnion, err error)
```

**GET `/sessions/{session_id}/resources/{resource_id}`** · [实现](../managed/sessionresource.go#L43) · 获取资源详情。

路径参数：`resourceID string`，必填。

返回：`*SessionResourceGetResponseUnion, error`。[SessionResourceGetResponseUnion](../managed/sessionresource.go#L483)。

请求对象 `params`：[SessionResourceGetParams](../managed/sessionresource.go#L681)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [SessionID](../managed/sessionresource.go#L682)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [WorkspaceID](../managed/sessionresource.go#L683)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/sessionresource.go#L685)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Sessions.Resources.Update`

```go
func (r *SessionResourceService) Update(ctx context.Context, resourceID string, params SessionResourceUpdateParams, opts ...option.RequestOption) (res *SessionResourceUpdateResponseUnion, err error)
```

**POST `/sessions/{session_id}/resources/{resource_id}`** · [实现](../managed/sessionresource.go#L63) · 更新资源。

路径参数：`resourceID string`，必填。

返回：`*SessionResourceUpdateResponseUnion, error`。[SessionResourceUpdateResponseUnion](../managed/sessionresource.go#L586)。

请求对象 `params`：[SessionResourceUpdateParams](../managed/sessionresource.go#L689)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Password](../managed/sessionresource.go#L690)：`param.Opt[string]` | json / `password` | 否 | 见字段定义。 |
| [SessionID](../managed/sessionresource.go#L692)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [AuthorizationToken](../managed/sessionresource.go#L695)：`string` | json / `authorization_token` | 否 | New authorization token for the resource. Currently only `github_repository` resources support token rotation. |
| [WorkspaceID](../managed/sessionresource.go#L696)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/sessionresource.go#L698)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Sessions.Resources.List`

```go
func (r *SessionResourceService) List(ctx context.Context, sessionID string, params SessionResourceListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsSessionResourceUnion], err error)
```

**GET `/sessions/{session_id}/resources`** · [实现](../managed/sessionresource.go#L83) · 列出资源。

路径参数：`sessionID string`，必填。

返回：`*pagination.PageCursor[ManagedAgentsSessionResourceUnion], error`。[ManagedAgentsSessionResourceUnion](../managed/sessionresource.go#L380)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[SessionResourceListParams](../managed/sessionresource.go#L710)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [BeforeID](../managed/sessionresource.go#L711)：`param.Opt[string]` | query / `before_id` | 否 | 见字段定义。 |
| [AfterID](../managed/sessionresource.go#L712)：`param.Opt[string]` | query / `after_id` | 否 | 见字段定义。 |
| [Limit](../managed/sessionresource.go#L716)：`param.Opt[int64]` | query / `limit` | 否 | Maximum number of resources to return per page (max 1000). If omitted, returns all resources. |
| [Page](../managed/sessionresource.go#L718)：`param.Opt[string]` | query / `page` | 否 | Opaque cursor from a previous response's `next_page` field. |
| [WorkspaceID](../managed/sessionresource.go#L719)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/sessionresource.go#L721)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Sessions.Resources.ListAutoPaging`

```go
func (r *SessionResourceService) ListAutoPaging(ctx context.Context, sessionID string, params SessionResourceListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsSessionResourceUnion]
```

[实现](../managed/sessionresource.go#L108)。自动遍历 `Sessions.Resources.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Sessions.Resources.Delete`

```go
func (r *SessionResourceService) Delete(ctx context.Context, resourceID string, params SessionResourceDeleteParams, opts ...option.RequestOption) (res *ManagedAgentsDeleteSessionResource, err error)
```

**DELETE `/sessions/{session_id}/resources/{resource_id}`** · [实现](../managed/sessionresource.go#L113) · 删除资源。

路径参数：`resourceID string`，必填。

返回：`*ManagedAgentsDeleteSessionResource, error`。[ManagedAgentsDeleteSessionResource](../managed/sessionresource.go#L149)。

请求对象 `params`：[SessionResourceDeleteParams](../managed/sessionresource.go#L734)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [SessionID](../managed/sessionresource.go#L735)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [WorkspaceID](../managed/sessionresource.go#L736)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/sessionresource.go#L738)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Sessions.Resources.Add`

```go
func (r *SessionResourceService) Add(ctx context.Context, sessionID string, params SessionResourceAddParams, opts ...option.RequestOption) (res *ManagedAgentsFileResource, err error)
```

**POST `/sessions/{session_id}/resources`** · [实现](../managed/sessionresource.go#L133) · 添加会话资源。

路径参数：`sessionID string`，必填。

返回：`*ManagedAgentsFileResource, error`。[ManagedAgentsFileResource](../managed/sessionresource.go#L174)。

请求对象 `params`：[SessionResourceAddParams](../managed/sessionresource.go#L742)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [ManagedAgentsFileResourceParams.FileID](../managed/session.go#L783)：`string` | json / `file_id` | 是 | ID of a previously uploaded file. |
| [ManagedAgentsFileResourceParams.Type](../managed/session.go#L785)：`ManagedAgentsFileResourceParamsType` | json / `type` | 是 | Any of "file".；类型：[ManagedAgentsFileResourceParamsType](../managed/session.go#L799) |
| [ManagedAgentsFileResourceParams.MountPath](../managed/session.go#L787)：`param.Opt[string]` | json / `mount_path` | 否 | Mount path in the container. Defaults to `/mnt/session/uploads/<file_id>`. |
| [WorkspaceID](../managed/sessionresource.go#L745)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/sessionresource.go#L747)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Sessions.Threads.Get`

```go
func (r *SessionThreadService) Get(ctx context.Context, threadID string, params SessionThreadGetParams, opts ...option.RequestOption) (res *ManagedAgentsSessionThread, err error)
```

**GET `/sessions/{session_id}/threads/{thread_id}`** · [实现](../managed/sessionthread.go#L45) · 获取资源详情。

路径参数：`threadID string`，必填。

返回：`*ManagedAgentsSessionThread, error`。[ManagedAgentsSessionThread](../managed/sessionthread.go#L116)。

请求对象 `params`：[SessionThreadGetParams](../managed/sessionthread.go#L1057)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [SessionID](../managed/sessionthread.go#L1058)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [WorkspaceID](../managed/sessionthread.go#L1059)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/sessionthread.go#L1061)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Sessions.Threads.List`

```go
func (r *SessionThreadService) List(ctx context.Context, sessionID string, params SessionThreadListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsSessionThread], err error)
```

**GET `/sessions/{session_id}/threads`** · [实现](../managed/sessionthread.go#L65) · 列出资源。

路径参数：`sessionID string`，必填。

返回：`*pagination.PageCursor[ManagedAgentsSessionThread], error`。[ManagedAgentsSessionThread](../managed/sessionthread.go#L116)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[SessionThreadListParams](../managed/sessionthread.go#L1065)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Limit](../managed/sessionthread.go#L1067)：`param.Opt[int64]` | query / `limit` | 否 | Maximum results per page. Defaults to 1000. |
| [Page](../managed/sessionthread.go#L1069)：`param.Opt[string]` | query / `page` | 否 | Opaque pagination cursor from a previous response's `next_page`. Forward-only. |
| [WorkspaceID](../managed/sessionthread.go#L1070)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/sessionthread.go#L1072)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Sessions.Threads.ListAutoPaging`

```go
func (r *SessionThreadService) ListAutoPaging(ctx context.Context, sessionID string, params SessionThreadListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsSessionThread]
```

[实现](../managed/sessionthread.go#L90)。自动遍历 `Sessions.Threads.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Sessions.Threads.Archive`

```go
func (r *SessionThreadService) Archive(ctx context.Context, threadID string, params SessionThreadArchiveParams, opts ...option.RequestOption) (res *ManagedAgentsSessionThread, err error)
```

**POST `/sessions/{session_id}/threads/{thread_id}/archive`** · [实现](../managed/sessionthread.go#L95) · 归档资源。

路径参数：`threadID string`，必填。

返回：`*ManagedAgentsSessionThread, error`。[ManagedAgentsSessionThread](../managed/sessionthread.go#L116)。

请求对象 `params`：[SessionThreadArchiveParams](../managed/sessionthread.go#L1085)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [SessionID](../managed/sessionthread.go#L1086)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [WorkspaceID](../managed/sessionthread.go#L1087)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/sessionthread.go#L1089)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Sessions.Threads.Events.List`

```go
func (r *SessionThreadEventService) List(ctx context.Context, threadID string, params SessionThreadEventListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsSessionEventUnion], err error)
```

**GET `/sessions/{session_id}/threads/{thread_id}/events`** · [实现](../managed/sessionthreadevent.go#L40) · 列出资源。

路径参数：`threadID string`，必填。

返回：`*pagination.PageCursor[ManagedAgentsSessionEventUnion], error`。[ManagedAgentsSessionEventUnion](../managed/sessionevent.go#L3756)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[SessionThreadEventListParams](../managed/sessionthreadevent.go#L97)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [BeforeID](../managed/sessionthreadevent.go#L98)：`param.Opt[string]` | query / `before_id` | 否 | 见字段定义。 |
| [AfterID](../managed/sessionthreadevent.go#L99)：`param.Opt[string]` | query / `after_id` | 否 | 见字段定义。 |
| [SessionID](../managed/sessionthreadevent.go#L101)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [Limit](../managed/sessionthreadevent.go#L103)：`param.Opt[int64]` | query / `limit` | 否 | Query parameter for limit |
| [Page](../managed/sessionthreadevent.go#L105)：`param.Opt[string]` | query / `page` | 否 | Query parameter for page |
| [WorkspaceID](../managed/sessionthreadevent.go#L106)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/sessionthreadevent.go#L108)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Sessions.Threads.Events.ListAutoPaging`

```go
func (r *SessionThreadEventService) ListAutoPaging(ctx context.Context, threadID string, params SessionThreadEventListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsSessionEventUnion]
```

[实现](../managed/sessionthreadevent.go#L69)。自动遍历 `Sessions.Threads.Events.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Sessions.Threads.Events.StreamEvents`

```go
func (r *SessionThreadEventService) StreamEvents(ctx context.Context, threadID string, params SessionThreadEventStreamParams, opts ...option.RequestOption) (stream *ssestream.Stream[ManagedAgentsStreamSessionThreadEventsUnion])
```

**GET `/sessions/{session_id}/threads/{thread_id}/stream`** · [实现](../managed/sessionthreadevent.go#L74) · 订阅 SSE 事件流。

路径参数：`threadID string`，必填。

返回：`*ssestream.Stream[ManagedAgentsStreamSessionThreadEventsUnion]`。[ManagedAgentsStreamSessionThreadEventsUnion](../managed/sessionthread.go#L409)、[事件流](../convention/ssestream/ssestream.go)。

请求对象 `params`：[SessionThreadEventStreamParams](../managed/sessionthreadevent.go#L121)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [SessionID](../managed/sessionthreadevent.go#L122)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [WorkspaceID](../managed/sessionthreadevent.go#L123)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [EventDeltas](../managed/sessionthreadevent.go#L134)：`[]ManagedAgentsDeltaType` | query / `event_deltas` | 否 | When set, this connection also receives streaming deltas (`event_start`, `event_delta`) while an event is being produced, before the event itself arrives.（完整约束见字段源码）；类型：[ManagedAgentsDeltaType](../managed/session.go#L771) |
| [Betas](../managed/sessionthreadevent.go#L136)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

## Deployments

管理部署、暂停、恢复与手动执行；运行记录通过顶层 DeploymentRuns 查询。

### `Deployments.New`

```go
func (r *DeploymentService) New(ctx context.Context, params DeploymentNewParams, opts ...option.RequestOption) (res *ManagedAgentsDeployment, err error)
```

**POST `/deployments`** · [实现](../managed/deployment.go#L43) · 创建资源。

返回：`*ManagedAgentsDeployment, error`。[ManagedAgentsDeployment](../managed/deployment.go#L205)。

请求对象 `params`：[DeploymentNewParams](../managed/deployment.go#L1832)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [EnvironmentVariables](../managed/deployment.go#L1833)：`param.Opt[string]` | json / `environment_variables` | 否 | 见字段定义。 |
| [Agent](../managed/deployment.go#L1838)：`DeploymentNewParamsAgentUnion` | json / `agent` | 是 | Agent to deploy. Accepts the `agent` ID string, which pins the latest version, or an `agent` object with both id and version specified. The agent must exist and not be archived.；类型：[DeploymentNewParamsAgentUnion](../managed/deployment.go#L1881) |
| [EnvironmentID](../managed/deployment.go#L1841)：`string` | json / `environment_id` | 是 | ID of the `environment` defining the container configuration for sessions created from this deployment. |
| [InitialEvents](../managed/deployment.go#L1844)：`[]ManagedAgentsDeploymentInitialEventParamsUnion` | json / `initial_events` | 是 | Events to send to each session immediately after creation. At least 1, maximum 50.；类型：[ManagedAgentsDeploymentInitialEventParamsUnion](../managed/deployment.go#L434) |
| [Name](../managed/deployment.go#L1846)：`string` | json / `name` | 是 | Human-readable name for the deployment. |
| [Description](../managed/deployment.go#L1848)：`param.Opt[string]` | json / `description` | 否 | Description of what the deployment does. |
| [WorkspaceID](../managed/deployment.go#L1849)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Budget](../managed/deployment.go#L1852)：`ManagedAgentsBudgetLimitParam` | json / `budget` | 否 | A hard spend ceiling. The session stops issuing new model requests once the tracked list cost reaches `max_list_cost`.；类型：[ManagedAgentsBudgetLimitParam](../managed/session.go#L585) |
| [Metadata](../managed/deployment.go#L1855)：`map[string]string` | json / `metadata` | 否 | Arbitrary key-value metadata. Maximum 16 pairs, keys up to 64 chars, values up to 512 chars. |
| [Resources](../managed/deployment.go#L1858)：`[]DeploymentNewParamsResourceUnion` | json / `resources` | 否 | Resources (e.g. repositories, files) to mount into each session's container. Maximum 500.；类型：[DeploymentNewParamsResourceUnion](../managed/deployment.go#L1906) |
| [Schedule](../managed/deployment.go#L1861)：`ManagedAgentsScheduleParams` | json / `schedule` | 否 | 5-field POSIX cron schedule. Literal wall-clock matching in the configured timezone.；类型：[ManagedAgentsScheduleParams](../managed/deployment.go#L1528) |
| [VaultIDs](../managed/deployment.go#L1864)：`[]string` | json / `vault_ids` | 否 | Vault IDs for stored credentials the agent can use during sessions created from this deployment. Maximum 50. |
| [Betas](../managed/deployment.go#L1866)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Deployments.Get`

```go
func (r *DeploymentService) Get(ctx context.Context, deploymentID string, query DeploymentGetParams, opts ...option.RequestOption) (res *ManagedAgentsDeployment, err error)
```

**GET `/deployments/{id}`** · [实现](../managed/deployment.go#L55) · 获取资源详情。

路径参数：`deploymentID string`，必填。

返回：`*ManagedAgentsDeployment, error`。[ManagedAgentsDeployment](../managed/deployment.go#L205)。

请求对象 `query`：[DeploymentGetParams](../managed/deployment.go#L2018)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/deployment.go#L2019)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/deployment.go#L2021)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Deployments.Update`

```go
func (r *DeploymentService) Update(ctx context.Context, deploymentID string, params DeploymentUpdateParams, opts ...option.RequestOption) (res *ManagedAgentsDeployment, err error)
```

**POST `/deployments/{id}`** · [实现](../managed/deployment.go#L71) · 更新资源。

路径参数：`deploymentID string`，必填。

返回：`*ManagedAgentsDeployment, error`。[ManagedAgentsDeployment](../managed/deployment.go#L205)。

请求对象 `params`：[DeploymentUpdateParams](../managed/deployment.go#L2025)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [EnvironmentVariables](../managed/deployment.go#L2026)：`param.Opt[string]` | json / `environment_variables` | 否 | 见字段定义。 |
| [Description](../managed/deployment.go#L2029)：`param.Opt[string]` | json / `description` | 否 | Description. Omit to preserve; send empty string or null to clear. |
| [EnvironmentID](../managed/deployment.go#L2031)：`param.Opt[string]` | json / `environment_id` | 否 | ID of the `environment` where sessions run. Omit to preserve. Cannot be cleared. |
| [Name](../managed/deployment.go#L2033)：`param.Opt[string]` | json / `name` | 否 | Human-readable name. Must be non-empty. Omit to preserve. Cannot be cleared. |
| [WorkspaceID](../managed/deployment.go#L2034)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Metadata](../managed/deployment.go#L2038)：`map[string]any` | json / `metadata` | 否 | Metadata patch. Set a key to a string to upsert it, or to null to delete it. Omit the field to preserve. The stored bag is limited to 16 keys (up to 64 chars each) with values up to 512 chars. |
| [Resources](../managed/deployment.go#L2041)：`[]DeploymentUpdateParamsResourceUnion` | json / `resources` | 否 | Session resources. Full replacement. Omit to preserve; send empty array or null to clear. Maximum 500.；类型：[DeploymentUpdateParamsResourceUnion](../managed/deployment.go#L2099) |
| [VaultIDs](../managed/deployment.go#L2044)：`[]string` | json / `vault_ids` | 否 | Vault IDs. Full replacement. Omit to preserve; send empty array or null to clear. Maximum 50. |
| [Agent](../managed/deployment.go#L2048)：`DeploymentUpdateParamsAgentUnion` | json / `agent` | 否 | Agent to deploy. Accepts the `agent` ID string, which re-pins to the latest version, or an `agent` object with both id and version specified. Omit to preserve. Cannot be cleared.；类型：[DeploymentUpdateParamsAgentUnion](../managed/deployment.go#L2074) |
| [Budget](../managed/deployment.go#L2051)：`ManagedAgentsBudgetLimitParam` | json / `budget` | 否 | A hard spend ceiling. The session stops issuing new model requests once the tracked list cost reaches `max_list_cost`.；类型：[ManagedAgentsBudgetLimitParam](../managed/session.go#L585) |
| [InitialEvents](../managed/deployment.go#L2054)：`[]ManagedAgentsDeploymentInitialEventParamsUnion` | json / `initial_events` | 否 | Initial events. Full replacement. Omit to preserve. Cannot be cleared. At least 1, maximum 50.；类型：[ManagedAgentsDeploymentInitialEventParamsUnion](../managed/deployment.go#L434) |
| [Schedule](../managed/deployment.go#L2057)：`ManagedAgentsScheduleParams` | json / `schedule` | 否 | 5-field POSIX cron schedule. Literal wall-clock matching in the configured timezone.；类型：[ManagedAgentsScheduleParams](../managed/deployment.go#L1528) |
| [Betas](../managed/deployment.go#L2059)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Deployments.List`

```go
func (r *DeploymentService) List(ctx context.Context, params DeploymentListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsDeployment], err error)
```

**GET `/deployments`** · [实现](../managed/deployment.go#L87) · 列出资源。

返回：`*pagination.PageCursor[ManagedAgentsDeployment], error`。[ManagedAgentsDeployment](../managed/deployment.go#L205)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[DeploymentListParams](../managed/deployment.go#L2211)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [BeforeID](../managed/deployment.go#L2212)：`param.Opt[string]` | query / `before_id` | 否 | 见字段定义。 |
| [AfterID](../managed/deployment.go#L2213)：`param.Opt[string]` | query / `after_id` | 否 | 见字段定义。 |
| [AgentID](../managed/deployment.go#L2216)：`param.Opt[string]` | query / `agent_id` | 否 | Filter by agent ID. |
| [CreatedAtGte](../managed/deployment.go#L2218)：`param.Opt[time.Time]` | query / `created_at[gte]` | 否 | Return deployments created at or after this time (inclusive). |
| [CreatedAtLte](../managed/deployment.go#L2220)：`param.Opt[time.Time]` | query / `created_at[lte]` | 否 | Return deployments created at or before this time (inclusive). |
| [IncludeArchived](../managed/deployment.go#L2222)：`param.Opt[bool]` | query / `include_archived` | 否 | When true, includes archived deployments. Default: false (exclude archived). |
| [Limit](../managed/deployment.go#L2224)：`param.Opt[int64]` | query / `limit` | 否 | Maximum results per page. Default 20, maximum 100. |
| [Page](../managed/deployment.go#L2226)：`param.Opt[string]` | query / `page` | 否 | Opaque pagination cursor. |
| [WorkspaceID](../managed/deployment.go#L2227)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Status](../managed/deployment.go#L2232)：`ManagedAgentsDeploymentStatus` | query / `status` | 否 | Filter by status: `active` or `paused`. Omit for both. To include archived deployments, use `include_archived` instead; the two cannot be combined. Any of "active", "paused".；类型：[ManagedAgentsDeploymentStatus](../managed/deployment.go#L801) |
| [Betas](../managed/deployment.go#L2234)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Deployments.ListAutoPaging`

```go
func (r *DeploymentService) ListAutoPaging(ctx context.Context, params DeploymentListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsDeployment]
```

[实现](../managed/deployment.go#L108)。自动遍历 `Deployments.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Deployments.Archive`

```go
func (r *DeploymentService) Archive(ctx context.Context, deploymentID string, body DeploymentArchiveParams, opts ...option.RequestOption) (res *ManagedAgentsDeployment, err error)
```

**POST `/deployments/{id}/archive`** · [实现](../managed/deployment.go#L113) · 归档资源。

路径参数：`deploymentID string`，必填。

返回：`*ManagedAgentsDeployment, error`。[ManagedAgentsDeployment](../managed/deployment.go#L205)。

请求对象 `body`：[DeploymentArchiveParams](../managed/deployment.go#L2247)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/deployment.go#L2248)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/deployment.go#L2250)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Deployments.Pause`

```go
func (r *DeploymentService) Pause(ctx context.Context, deploymentID string, body DeploymentPauseParams, opts ...option.RequestOption) (res *ManagedAgentsDeployment, err error)
```

**POST `/deployments/{id}/pause`** · [实现](../managed/deployment.go#L129) · 暂停后续执行。

路径参数：`deploymentID string`，必填。

返回：`*ManagedAgentsDeployment, error`。[ManagedAgentsDeployment](../managed/deployment.go#L205)。

请求对象 `body`：[DeploymentPauseParams](../managed/deployment.go#L2254)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/deployment.go#L2255)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/deployment.go#L2257)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Deployments.Run`

```go
func (r *DeploymentService) Run(ctx context.Context, deploymentID string, body DeploymentRunParams, opts ...option.RequestOption) (res *ManagedAgentsDeploymentRun, err error)
```

**POST `/deployments/{id}/run`** · [实现](../managed/deployment.go#L145) · 立即触发一次执行。

路径参数：`deploymentID string`，必填。

返回：`*ManagedAgentsDeploymentRun, error`。[ManagedAgentsDeploymentRun](../managed/deploymentrun.go#L113)。

请求对象 `body`：[DeploymentRunParams](../managed/deployment.go#L2261)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/deployment.go#L2262)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/deployment.go#L2264)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Deployments.Unpause`

```go
func (r *DeploymentService) Unpause(ctx context.Context, deploymentID string, body DeploymentUnpauseParams, opts ...option.RequestOption) (res *ManagedAgentsDeployment, err error)
```

**POST `/deployments/{id}/unpause`** · [实现](../managed/deployment.go#L161) · 恢复执行。

路径参数：`deploymentID string`，必填。

返回：`*ManagedAgentsDeployment, error`。[ManagedAgentsDeployment](../managed/deployment.go#L205)。

请求对象 `body`：[DeploymentUnpauseParams](../managed/deployment.go#L2268)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/deployment.go#L2269)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/deployment.go#L2271)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

## DeploymentRuns

查询部署运行实例及关联会话；不位于 Deployments 的子路径中。

### `DeploymentRuns.Get`

```go
func (r *DeploymentRunService) Get(ctx context.Context, deploymentRunID string, query DeploymentRunGetParams, opts ...option.RequestOption) (res *ManagedAgentsDeploymentRun, err error)
```

**GET `/deployment_runs/{run_id}`** · [实现](../managed/deploymentrun.go#L43) · 获取资源详情。

路径参数：`deploymentRunID string`，必填。

返回：`*ManagedAgentsDeploymentRun, error`。[ManagedAgentsDeploymentRun](../managed/deploymentrun.go#L113)。

请求对象 `query`：[DeploymentRunGetParams](../managed/deploymentrun.go#L906)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/deploymentrun.go#L907)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/deploymentrun.go#L909)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `DeploymentRuns.List`

```go
func (r *DeploymentRunService) List(ctx context.Context, params DeploymentRunListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsDeploymentRun], err error)
```

**GET `/deployment_runs`** · [实现](../managed/deploymentrun.go#L59) · 列出资源。

返回：`*pagination.PageCursor[ManagedAgentsDeploymentRun], error`。[ManagedAgentsDeploymentRun](../managed/deploymentrun.go#L113)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[DeploymentRunListParams](../managed/deploymentrun.go#L913)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [BeforeID](../managed/deploymentrun.go#L914)：`param.Opt[string]` | query / `before_id` | 否 | 见字段定义。 |
| [AfterID](../managed/deploymentrun.go#L915)：`param.Opt[string]` | query / `after_id` | 否 | 见字段定义。 |
| [CreatedAtGt](../managed/deploymentrun.go#L918)：`param.Opt[time.Time]` | query / `created_at[gt]` | 否 | Return runs created strictly after this time (exclusive). |
| [CreatedAtGte](../managed/deploymentrun.go#L920)：`param.Opt[time.Time]` | query / `created_at[gte]` | 否 | Return runs created at or after this time (inclusive). |
| [CreatedAtLt](../managed/deploymentrun.go#L922)：`param.Opt[time.Time]` | query / `created_at[lt]` | 否 | Return runs created strictly before this time (exclusive). |
| [CreatedAtLte](../managed/deploymentrun.go#L924)：`param.Opt[time.Time]` | query / `created_at[lte]` | 否 | Return runs created at or before this time (inclusive). |
| [DeploymentID](../managed/deploymentrun.go#L928)：`param.Opt[string]` | query / `deployment_id` | 否 | Filter to a specific deployment. Omit to list across all deployments in the workspace. Filtering by a non-existent `deployment_id` returns 200 with empty data. |
| [HasError](../managed/deploymentrun.go#L931)：`param.Opt[bool]` | query / `has_error` | 否 | Filter: true for runs with non-null `error`, false for runs with non-null `session_id`. Omit for all. |
| [Limit](../managed/deploymentrun.go#L933)：`param.Opt[int64]` | query / `limit` | 否 | Maximum results per page. Default 20, maximum 1000. |
| [Page](../managed/deploymentrun.go#L936)：`param.Opt[string]` | query / `page` | 否 | Opaque pagination cursor. Pass `next_page` from the previous response. Invalid or expired cursors return 400. |
| [WorkspaceID](../managed/deploymentrun.go#L937)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [TriggerType](../managed/deploymentrun.go#L941)：`ManagedAgentsTriggerType` | query / `trigger_type` | 否 | Filter runs by what triggered them. Omit to return all runs. Any of "schedule", "manual".；类型：[ManagedAgentsTriggerType](../managed/deploymentrun.go#L790) |
| [Betas](../managed/deploymentrun.go#L943)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `DeploymentRuns.ListAutoPaging`

```go
func (r *DeploymentRunService) ListAutoPaging(ctx context.Context, params DeploymentRunListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsDeploymentRun]
```

[实现](../managed/deploymentrun.go#L80)。自动遍历 `DeploymentRuns.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

## Dreams

启动、查询、取消或归档记忆整理任务；应等待完成并检查输出 Memory Store 内容。

### `Dreams.New`

```go
func (r *DreamService) New(ctx context.Context, params DreamNewParams, opts ...option.RequestOption) (res *Dream, err error)
```

**POST `/dreams`** · [实现](../managed/dream.go#L43) · 创建资源。

返回：`*Dream, error`。[Dream](../managed/dream.go#L134)。

请求对象 `params`：[DreamNewParams](../managed/dream.go#L852)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Inputs](../managed/dream.go#L853)：`[]DreamInputUnionParam` | json / `inputs` | 是 | 类型：[DreamInputUnionParam](../managed/dream.go#L304) |
| [Model](../managed/dream.go#L855)：`DreamNewParamsModelUnion` | json / `model` | 是 | Model identifier and configuration applied to every pipeline stage.；类型：[DreamNewParamsModelUnion](../managed/dream.go#L878) |
| [Instructions](../managed/dream.go#L856)：`param.Opt[string]` | json / `instructions` | 否 | 见字段定义。 |
| [WorkspaceID](../managed/dream.go#L857)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [OutputBehavior](../managed/dream.go#L861)：`OutputBehaviorUnionParam` | json / `output_behavior` | 否 | The default destination: the job creates a new output memory store as a clone of the memory_store input and writes the consolidated memories into it. The input store is never mutated.；类型：[OutputBehaviorUnionParam](../managed/dream.go#L691) |
| [Betas](../managed/dream.go#L863)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Dreams.Get`

```go
func (r *DreamService) Get(ctx context.Context, dreamID string, query DreamGetParams, opts ...option.RequestOption) (res *Dream, err error)
```

**GET `/dreams/{id}`** · [实现](../managed/dream.go#L55) · 获取资源详情。

路径参数：`dreamID string`，必填。

返回：`*Dream, error`。[Dream](../managed/dream.go#L134)。

请求对象 `query`：[DreamGetParams](../managed/dream.go#L900)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/dream.go#L901)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/dream.go#L903)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Dreams.List`

```go
func (r *DreamService) List(ctx context.Context, params DreamListParams, opts ...option.RequestOption) (res *pagination.PageCursor[Dream], err error)
```

**GET `/dreams`** · [实现](../managed/dream.go#L71) · 列出资源。

返回：`*pagination.PageCursor[Dream], error`。[Dream](../managed/dream.go#L134)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[DreamListParams](../managed/dream.go#L907)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [CreatedAtGt](../managed/dream.go#L910)：`param.Opt[time.Time]` | query / `created_at[gt]` | 否 | Return dreams with `created_at` strictly after this timestamp (exclusive lower bound, RFC 3339). Unset applies no lower bound. |
| [CreatedAtLt](../managed/dream.go#L913)：`param.Opt[time.Time]` | query / `created_at[lt]` | 否 | Return dreams with `created_at` strictly before this timestamp (exclusive upper bound, RFC 3339). Unset applies no upper bound. |
| [IncludeArchived](../managed/dream.go#L915)：`param.Opt[bool]` | query / `include_archived` | 否 | Query parameter for include_archived |
| [Limit](../managed/dream.go#L917)：`param.Opt[int64]` | query / `limit` | 否 | Query parameter for limit |
| [Page](../managed/dream.go#L919)：`param.Opt[string]` | query / `page` | 否 | Query parameter for page |
| [WorkspaceID](../managed/dream.go#L920)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Statuses](../managed/dream.go#L923)：`[]DreamStatus` | query / `statuses` | 否 | Filter by lifecycle status. Repeat the parameter to match any of multiple statuses. Empty applies no status filter.；类型：[DreamStatus](../managed/dream.go#L567) |
| [Betas](../managed/dream.go#L925)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Dreams.ListAutoPaging`

```go
func (r *DreamService) ListAutoPaging(ctx context.Context, params DreamListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[Dream]
```

[实现](../managed/dream.go#L92)。自动遍历 `Dreams.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Dreams.Archive`

```go
func (r *DreamService) Archive(ctx context.Context, dreamID string, body DreamArchiveParams, opts ...option.RequestOption) (res *Dream, err error)
```

**POST `/dreams/{id}/archive`** · [实现](../managed/dream.go#L97) · 归档资源。

路径参数：`dreamID string`，必填。

返回：`*Dream, error`。[Dream](../managed/dream.go#L134)。

请求对象 `body`：[DreamArchiveParams](../managed/dream.go#L937)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/dream.go#L938)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/dream.go#L940)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Dreams.Cancel`

```go
func (r *DreamService) Cancel(ctx context.Context, dreamID string, body DreamCancelParams, opts ...option.RequestOption) (res *Dream, err error)
```

**POST `/dreams/{id}/cancel`** · [实现](../managed/dream.go#L113) · 取消执行。

路径参数：`dreamID string`，必填。

返回：`*Dream, error`。[Dream](../managed/dream.go#L134)。

请求对象 `body`：[DreamCancelParams](../managed/dream.go#L944)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/dream.go#L945)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/dream.go#L947)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

## Environments

管理执行环境配置和生命周期。Managed 的 Work 子资源提供自托管 worker 领取、确认、心跳与停止操作。

### `Environments.New`

```go
func (r *EnvironmentService) New(ctx context.Context, params EnvironmentNewParams, opts ...option.RequestOption) (res *Environment, err error)
```

**POST `/environments`** · [实现](../managed/environment.go#L46) · 创建资源。

返回：`*Environment, error`。[Environment](../managed/environment.go#L344)。

请求对象 `params`：[EnvironmentNewParams](../managed/environment.go#L736)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Name](../managed/environment.go#L738)：`string` | json / `name` | 是 | Human-readable name for the environment |
| [Description](../managed/environment.go#L740)：`param.Opt[string]` | json / `description` | 否 | Optional description of the environment |
| [WorkspaceID](../managed/environment.go#L741)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Config](../managed/environment.go#L743)：`EnvironmentNewParamsConfigUnion` | json / `config` | 否 | Environment configuration；类型：[EnvironmentNewParamsConfigUnion](../managed/environment.go#L769) |
| [Scope](../managed/environment.go#L750)：`EnvironmentNewParamsScope` | json / `scope` | 否 | The visibility scope for this environment.（完整约束见字段源码）；类型：[EnvironmentNewParamsScope](../managed/environment.go#L829) |
| [Metadata](../managed/environment.go#L752)：`map[string]string` | json / `metadata` | 否 | User-provided metadata key-value pairs |
| [Betas](../managed/environment.go#L754)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Environments.Get`

```go
func (r *EnvironmentService) Get(ctx context.Context, environmentID string, query EnvironmentGetParams, opts ...option.RequestOption) (res *Environment, err error)
```

**GET `/environments/{environment_id}`** · [实现](../managed/environment.go#L58) · 获取资源详情。

路径参数：`environmentID string`，必填。

返回：`*Environment, error`。[Environment](../managed/environment.go#L344)。

请求对象 `query`：[EnvironmentGetParams](../managed/environment.go#L836)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/environment.go#L837)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/environment.go#L839)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Environments.Update`

```go
func (r *EnvironmentService) Update(ctx context.Context, environmentID string, params EnvironmentUpdateParams, opts ...option.RequestOption) (res *Environment, err error)
```

**POST `/environments/{environment_id}`** · [实现](../managed/environment.go#L74) · 更新资源。

路径参数：`environmentID string`，必填。

返回：`*Environment, error`。[Environment](../managed/environment.go#L344)。

请求对象 `params`：[EnvironmentUpdateParams](../managed/environment.go#L843)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Description](../managed/environment.go#L846)：`param.Opt[string]` | json / `description` | 否 | Updated description of the environment. Omit to preserve; null clears to null; an empty string is stored as an empty string. |
| [Name](../managed/environment.go#L848)：`param.Opt[string]` | json / `name` | 否 | Updated name for the environment |
| [WorkspaceID](../managed/environment.go#L849)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Config](../managed/environment.go#L851)：`EnvironmentUpdateParamsConfigUnion` | json / `config` | 否 | Updated environment configuration；类型：[EnvironmentUpdateParamsConfigUnion](../managed/environment.go#L877) |
| [Scope](../managed/environment.go#L857)：`EnvironmentUpdateParamsScope` | json / `scope` | 否 | The visibility scope for this environment. 'organization' makes the environment visible to all accounts. 'account' restricts visibility to the owning account only. Any of "organization", "account".；类型：[EnvironmentUpdateParamsScope](../managed/environment.go#L936) |
| [Metadata](../managed/environment.go#L860)：`map[string]any` | json / `metadata` | 否 | User-provided metadata key-value pairs. Set a value to null or empty string to delete the key. |
| [Betas](../managed/environment.go#L862)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Environments.List`

```go
func (r *EnvironmentService) List(ctx context.Context, params EnvironmentListParams, opts ...option.RequestOption) (res *pagination.PageCursor[Environment], err error)
```

**GET `/environments`** · [实现](../managed/environment.go#L90) · 列出资源。

返回：`*pagination.PageCursor[Environment], error`。[Environment](../managed/environment.go#L344)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[EnvironmentListParams](../managed/environment.go#L943)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [CreatedAtGte](../managed/environment.go#L944)：`param.Opt[time.Time]` | query / `created_at[gte]` | 否 | 见字段定义。 |
| [CreatedAtLte](../managed/environment.go#L945)：`param.Opt[time.Time]` | query / `created_at[lte]` | 否 | 见字段定义。 |
| [Page](../managed/environment.go#L949)：`param.Opt[string]` | query / `page` | 否 | Opaque cursor from previous response for pagination. Pass the `next_page` value from the previous response. |
| [IncludeArchived](../managed/environment.go#L951)：`param.Opt[bool]` | query / `include_archived` | 否 | Include archived environments in the response |
| [Limit](../managed/environment.go#L953)：`param.Opt[int64]` | query / `limit` | 否 | Maximum number of environments to return |
| [WorkspaceID](../managed/environment.go#L954)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/environment.go#L956)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Environments.ListAutoPaging`

```go
func (r *EnvironmentService) ListAutoPaging(ctx context.Context, params EnvironmentListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[Environment]
```

[实现](../managed/environment.go#L111)。自动遍历 `Environments.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Environments.Delete`

```go
func (r *EnvironmentService) Delete(ctx context.Context, environmentID string, body EnvironmentDeleteParams, opts ...option.RequestOption) (res *EnvironmentDeleteResponse, err error)
```

**DELETE `/environments/{environment_id}`** · [实现](../managed/environment.go#L116) · 删除资源。

路径参数：`environmentID string`，必填。

返回：`*EnvironmentDeleteResponse, error`。[EnvironmentDeleteResponse](../managed/environment.go#L467)。

请求对象 `body`：[EnvironmentDeleteParams](../managed/environment.go#L969)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/environment.go#L970)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/environment.go#L972)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Environments.Archive`

```go
func (r *EnvironmentService) Archive(ctx context.Context, environmentID string, body EnvironmentArchiveParams, opts ...option.RequestOption) (res *Environment, err error)
```

**POST `/environments/{environment_id}/archive`** · [实现](../managed/environment.go#L133) · 归档资源。

路径参数：`environmentID string`，必填。

返回：`*Environment, error`。[Environment](../managed/environment.go#L344)。

请求对象 `body`：[EnvironmentArchiveParams](../managed/environment.go#L976)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/environment.go#L977)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/environment.go#L979)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Environments.Work.Get`

```go
func (r *EnvironmentWorkService) Get(ctx context.Context, workID string, params EnvironmentWorkGetParams, opts ...option.RequestOption) (res *SelfHostedWork, err error)
```

**GET `/environments/{environment_id}/work/{work_id}`** · [实现](../managed/environmentwork.go#L47) · 获取资源详情。

路径参数：`workID string`，必填。

返回：`*SelfHostedWork, error`。[SelfHostedWork](../managed/environmentwork.go#L286)。

请求对象 `params`：[EnvironmentWorkGetParams](../managed/environmentwork.go#L508)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [EnvironmentID](../managed/environmentwork.go#L509)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [WorkspaceID](../managed/environmentwork.go#L510)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/environmentwork.go#L512)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Environments.Work.Update`

```go
func (r *EnvironmentWorkService) Update(ctx context.Context, workID string, params EnvironmentWorkUpdateParams, opts ...option.RequestOption) (res *SelfHostedWork, err error)
```

**POST `/environments/{environment_id}/work/{work_id}`** · [实现](../managed/environmentwork.go#L72) · 更新资源。

路径参数：`workID string`，必填。

返回：`*SelfHostedWork, error`。[SelfHostedWork](../managed/environmentwork.go#L286)。

请求对象 `params`：[EnvironmentWorkUpdateParams](../managed/environmentwork.go#L516)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [EnvironmentID](../managed/environmentwork.go#L517)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [SelfHostedWorkUpdateRequest.Metadata](../managed/environmentwork.go#L472)：`map[string]any` | json / `metadata` | 是 | Metadata patch. Set a key to a string to upsert it, or to null to delete it. Omit the field to preserve existing metadata. |
| [WorkspaceID](../managed/environmentwork.go#L520)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/environmentwork.go#L522)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Environments.Work.List`

```go
func (r *EnvironmentWorkService) List(ctx context.Context, environmentID string, params EnvironmentWorkListParams, opts ...option.RequestOption) (res *pagination.PageCursor[SelfHostedWork], err error)
```

**GET `/environments/{environment_id}/work`** · [实现](../managed/environmentwork.go#L97) · 列出资源。

路径参数：`environmentID string`，必填。

返回：`*pagination.PageCursor[SelfHostedWork], error`。[SelfHostedWork](../managed/environmentwork.go#L286)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[EnvironmentWorkListParams](../managed/environmentwork.go#L533)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [BeforeID](../managed/environmentwork.go#L534)：`param.Opt[string]` | query / `before_id` | 否 | 见字段定义。 |
| [AfterID](../managed/environmentwork.go#L535)：`param.Opt[string]` | query / `after_id` | 否 | 见字段定义。 |
| [Page](../managed/environmentwork.go#L538)：`param.Opt[string]` | query / `page` | 否 | Opaque cursor from previous response for pagination |
| [Limit](../managed/environmentwork.go#L540)：`param.Opt[int64]` | query / `limit` | 否 | Maximum number of work items to return |
| [Betas](../managed/environmentwork.go#L542)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Environments.Work.ListAutoPaging`

```go
func (r *EnvironmentWorkService) ListAutoPaging(ctx context.Context, environmentID string, params EnvironmentWorkListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[SelfHostedWork]
```

[实现](../managed/environmentwork.go#L127)。自动遍历 `Environments.Work.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Environments.Work.Ack`

```go
func (r *EnvironmentWorkService) Ack(ctx context.Context, workID string, params EnvironmentWorkAckParams, opts ...option.RequestOption) (res *SelfHostedWork, err error)
```

**POST `/environments/{environment_id}/work/{work_id}/ack`** · [实现](../managed/environmentwork.go#L138) · 确认工作领取。

路径参数：`workID string`，必填。

返回：`*SelfHostedWork, error`。[SelfHostedWork](../managed/environmentwork.go#L286)。

请求对象 `params`：[EnvironmentWorkAckParams](../managed/environmentwork.go#L555)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [EnvironmentID](../managed/environmentwork.go#L556)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [Betas](../managed/environmentwork.go#L558)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Environments.Work.Heartbeat`

```go
func (r *EnvironmentWorkService) Heartbeat(ctx context.Context, workID string, params EnvironmentWorkHeartbeatParams, opts ...option.RequestOption) (res *SelfHostedWorkHeartbeatResponse, err error)
```

**POST `/environments/{environment_id}/work/{work_id}/heartbeat`** · [实现](../managed/environmentwork.go#L163) · 更新工作心跳。

路径参数：`workID string`，必填。

返回：`*SelfHostedWorkHeartbeatResponse, error`。[SelfHostedWorkHeartbeatResponse](../managed/environmentwork.go#L354)。

请求对象 `params`：[EnvironmentWorkHeartbeatParams](../managed/environmentwork.go#L562)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [EnvironmentID](../managed/environmentwork.go#L563)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [DesiredTTLSeconds](../managed/environmentwork.go#L565)：`param.Opt[int64]` | query / `desired_ttl_seconds` | 否 | Desired TTL in seconds |
| [ExpectedLastHeartbeat](../managed/environmentwork.go#L570)：`param.Opt[string]` | query / `expected_last_heartbeat` | 否 | Expected last_heartbeat for conditional update (optimistic concurrency).（完整约束见字段源码） |
| [Betas](../managed/environmentwork.go#L572)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Environments.Work.Poll`

```go
func (r *EnvironmentWorkService) Poll(ctx context.Context, environmentID string, params EnvironmentWorkPollParams, opts ...option.RequestOption) (res *SelfHostedWork, err error)
```

**GET `/environments/{environment_id}/work/poll`** · [实现](../managed/environmentwork.go#L188) · 长轮询领取待执行工作。

路径参数：`environmentID string`，必填。

返回：`*SelfHostedWork, error`。[SelfHostedWork](../managed/environmentwork.go#L286)。

请求对象 `params`：[EnvironmentWorkPollParams](../managed/environmentwork.go#L585)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [BlockMs](../managed/environmentwork.go#L589)：`param.Opt[int64]` | query / `block_ms` | 否 | How long to wait for work to arrive before returning. Must be 1-999 in milliseconds. Defaults to non-blocking (returns immediately if no work is available). |
| [ReclaimOlderThanMs](../managed/environmentwork.go#L592)：`param.Opt[int64]` | query / `reclaim_older_than_ms` | 否 | Reclaim unacknowledged work items older than this many milliseconds. If omitted, uses the default (5000ms). |
| [QoderWorkerID](../managed/environmentwork.go#L595)：`param.Opt[string]` | header / `Worker-ID` | 否 | Unique identifier for the specific worker polling, used to track aggregated environment-level work metrics in Console |
| [Betas](../managed/environmentwork.go#L597)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Environments.Work.Stats`

```go
func (r *EnvironmentWorkService) Stats(ctx context.Context, environmentID string, query EnvironmentWorkStatsParams, opts ...option.RequestOption) (res *SelfHostedWorkQueueStats, err error)
```

**GET `/environments/{environment_id}/work/stats`** · [实现](../managed/environmentwork.go#L207) · 查询统计。

路径参数：`environmentID string`，必填。

返回：`*SelfHostedWorkQueueStats, error`。[SelfHostedWorkQueueStats](../managed/environmentwork.go#L420)。

请求对象 `query`：[EnvironmentWorkStatsParams](../managed/environmentwork.go#L610)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/environmentwork.go#L611)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/environmentwork.go#L613)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Environments.Work.Stop`

```go
func (r *EnvironmentWorkService) Stop(ctx context.Context, workID string, params EnvironmentWorkStopParams, opts ...option.RequestOption) (res *SelfHostedWork, err error)
```

**POST `/environments/{environment_id}/work/{work_id}/stop`** · [实现](../managed/environmentwork.go#L228) · 停止工作。

路径参数：`workID string`，必填。

返回：`*SelfHostedWork, error`。[SelfHostedWork](../managed/environmentwork.go#L286)。

请求对象 `params`：[EnvironmentWorkStopParams](../managed/environmentwork.go#L617)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [EnvironmentID](../managed/environmentwork.go#L618)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [SelfHostedWorkStopRequest.Force](../managed/environmentwork.go#L454)：`param.Opt[bool]` | json / `force` | 否 | If true, immediately stop work without graceful shutdown |
| [WorkspaceID](../managed/environmentwork.go#L621)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/environmentwork.go#L623)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

## Skills

上传 Skill 文件树并管理版本；文件名使用相对路径，例如 my-skill/SKILL.md。

### `Skills.New`

```go
func (r *SkillService) New(ctx context.Context, params SkillNewParams, opts ...option.RequestOption) (res *Skill, err error)
```

**POST `/skills`** · [实现](../managed/skill.go#L49) · 创建资源。

返回：`*Skill, error`。[Skill](../managed/skill.go#L139)。

请求对象 `params`：[SkillNewParams](../managed/skill.go#L235)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Metadata](../managed/skill.go#L236)：`map[string]string` | multipart / `metadata` | 否 | 见字段定义。 |
| [Files](../managed/skill.go#L242)：`[]io.Reader` | multipart / `files` | 是 | Files to upload for the skill. All files must be in the same top-level directory and must include a SKILL.md file at the root of that directory. |
| [DisplayName](../managed/skill.go#L246)：`param.Opt[string]` | multipart / `display_title` | 否 | Human-readable, single-line label for the Skill. Maximum 255 characters. Always set: derived from the SKILL.md frontmatter `name` when omitted at creation. Not unique. |
| [WorkspaceID](../managed/skill.go#L247)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/skill.go#L249)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Skills.Get`

```go
func (r *SkillService) Get(ctx context.Context, skillID string, query SkillGetParams, opts ...option.RequestOption) (res *Skill, err error)
```

**GET `/skills/{skill_id}`** · [实现](../managed/skill.go#L60) · 获取资源详情。

路径参数：`skillID string`，必填。

返回：`*Skill, error`。[Skill](../managed/skill.go#L139)。

请求对象 `query`：[SkillGetParams](../managed/skill.go#L271)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/skill.go#L272)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/skill.go#L274)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Skills.List`

```go
func (r *SkillService) List(ctx context.Context, params SkillListParams, opts ...option.RequestOption) (res *pagination.PageCursor[Skill], err error)
```

**GET `/skills`** · [实现](../managed/skill.go#L75) · 列出资源。

返回：`*pagination.PageCursor[Skill], error`。[Skill](../managed/skill.go#L139)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[SkillListParams](../managed/skill.go#L278)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [DisplayName](../managed/skill.go#L279)：`param.Opt[string]` | query / `display_title` | 否 | 见字段定义。 |
| [Name](../managed/skill.go#L281)：`param.Opt[string]` | query / `name` | 否 | 见字段定义。 |
| [BeforeID](../managed/skill.go#L283)：`param.Opt[string]` | query / `before_id` | 否 | 见字段定义。 |
| [AfterID](../managed/skill.go#L284)：`param.Opt[string]` | query / `after_id` | 否 | 见字段定义。 |
| [Page](../managed/skill.go#L290)：`param.Opt[string]` | query / `page` | 否 | Pagination token for fetching a specific page of results. Pass the value from a previous response's `next_page` field to get the next page of results. |
| [Source](../managed/skill.go#L297)：`param.Opt[string]` | query / `source` | 否 | Filter skills by source. If provided, only skills from the specified source will be returned: - `"custom"`: only return user-created skills - `"qoder"`: only return Qoder-created skills |
| [Limit](../managed/skill.go#L301)：`param.Opt[int64]` | query / `limit` | 否 | Number of results to return per page. Ranges from `1` to `1000`. Defaults to `20`. |
| [WorkspaceID](../managed/skill.go#L302)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/skill.go#L304)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Skills.ListAutoPaging`

```go
func (r *SkillService) ListAutoPaging(ctx context.Context, params SkillListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[Skill]
```

[实现](../managed/skill.go#L96)。自动遍历 `Skills.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Skills.Delete`

```go
func (r *SkillService) Delete(ctx context.Context, skillID string, body SkillDeleteParams, opts ...option.RequestOption) (res *DeletedSkill, err error)
```

**DELETE `/skills/{skill_id}`** · [实现](../managed/skill.go#L101) · 删除资源。

路径参数：`skillID string`，必填。

返回：`*DeletedSkill, error`。[DeletedSkill](../managed/skill.go#L115)。

请求对象 `body`：[SkillDeleteParams](../managed/skill.go#L316)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/skill.go#L317)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/skill.go#L319)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Skills.Versions.New`

```go
func (r *SkillVersionService) New(ctx context.Context, skillID string, params SkillVersionNewParams, opts ...option.RequestOption) (res *SkillVersion, err error)
```

**POST `/skills/{skill_id}/versions`** · [实现](../managed/skillversion.go#L47) · 创建资源。

路径参数：`skillID string`，必填。

返回：`*SkillVersion, error`。[SkillVersion](../managed/skillversion.go#L172)。

请求对象 `params`：[SkillVersionNewParams](../managed/skillversion.go#L220)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Files](../managed/skillversion.go#L225)：`[]io.Reader` | multipart / `files` | 是 | Files to upload for the skill. All files must be in the same top-level directory and must include a SKILL.md file at the root of that directory. |
| [WorkspaceID](../managed/skillversion.go#L226)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/skillversion.go#L228)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Skills.Versions.Get`

```go
func (r *SkillVersionService) Get(ctx context.Context, version string, params SkillVersionGetParams, opts ...option.RequestOption) (res *SkillVersion, err error)
```

**GET `/skills/{skill_id}/versions/{version}`** · [实现](../managed/skillversion.go#L62) · 获取资源详情。

路径参数：`version string`，必填。

返回：`*SkillVersion, error`。[SkillVersion](../managed/skillversion.go#L172)。

请求对象 `params`：[SkillVersionGetParams](../managed/skillversion.go#L250)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [SkillID](../managed/skillversion.go#L254)：`string` | 嵌入 / `` | 是 | Unique identifier for the skill. The format and length of IDs may change over time. |
| [WorkspaceID](../managed/skillversion.go#L255)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/skillversion.go#L257)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Skills.Versions.List`

```go
func (r *SkillVersionService) List(ctx context.Context, skillID string, params SkillVersionListParams, opts ...option.RequestOption) (res *pagination.PageCursor[SkillVersion], err error)
```

**GET `/skills/{skill_id}/versions`** · [实现](../managed/skillversion.go#L81) · 列出资源。

路径参数：`skillID string`，必填。

返回：`*pagination.PageCursor[SkillVersion], error`。[SkillVersion](../managed/skillversion.go#L172)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[SkillVersionListParams](../managed/skillversion.go#L261)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Limit](../managed/skillversion.go#L265)：`param.Opt[int64]` | query / `limit` | 否 | Number of results to return per page. Ranges from `1` to `1000`. Defaults to `20`. |
| [Page](../managed/skillversion.go#L267)：`param.Opt[string]` | query / `page` | 否 | Optionally set to the `next_page` token from the previous response. |
| [WorkspaceID](../managed/skillversion.go#L268)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/skillversion.go#L270)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Skills.Versions.ListAutoPaging`

```go
func (r *SkillVersionService) ListAutoPaging(ctx context.Context, skillID string, params SkillVersionListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[SkillVersion]
```

[实现](../managed/skillversion.go#L106)。自动遍历 `Skills.Versions.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Skills.Versions.Delete`

```go
func (r *SkillVersionService) Delete(ctx context.Context, version string, params SkillVersionDeleteParams, opts ...option.RequestOption) (res *DeletedSkillVersion, err error)
```

**DELETE `/skills/{skill_id}/versions/{version}`** · [实现](../managed/skillversion.go#L111) · 删除资源。

路径参数：`version string`，必填。

返回：`*DeletedSkillVersion, error`。[DeletedSkillVersion](../managed/skillversion.go#L149)。

请求对象 `params`：[SkillVersionDeleteParams](../managed/skillversion.go#L283)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [SkillID](../managed/skillversion.go#L287)：`string` | 嵌入 / `` | 是 | Unique identifier for the skill. The format and length of IDs may change over time. |
| [WorkspaceID](../managed/skillversion.go#L288)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/skillversion.go#L290)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Skills.Versions.Download`

```go
func (r *SkillVersionService) Download(ctx context.Context, version string, params SkillVersionDownloadParams, opts ...option.RequestOption) (res *http.Response, err error)
```

**GET `/skills/{skill_id}/versions/{version}/content`** · [实现](../managed/skillversion.go#L130) · 下载文件内容。

路径参数：`version string`，必填。

返回：`*http.Response, error`。HTTP 响应；读取后关闭 Body。

请求对象 `params`：[SkillVersionDownloadParams](../managed/skillversion.go#L294)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [SkillID](../managed/skillversion.go#L298)：`string` | 嵌入 / `` | 是 | Unique identifier for the skill. The format and length of IDs may change over time. |
| [WorkspaceID](../managed/skillversion.go#L299)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/skillversion.go#L301)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

## Vaults

管理凭据容器及其凭据，支持归档、更新凭据和 MCP OAuth 验证。

### `Vaults.New`

```go
func (r *VaultService) New(ctx context.Context, params VaultNewParams, opts ...option.RequestOption) (res *ManagedAgentsVault, err error)
```

**POST `/vaults`** · [实现](../managed/vault.go#L44) · 创建资源。

返回：`*ManagedAgentsVault, error`。[ManagedAgentsVault](../managed/vault.go#L157)。

请求对象 `params`：[VaultNewParams](../managed/vault.go#L198)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [DisplayName](../managed/vault.go#L200)：`string` | json / `display_name` | 是 | Human-readable name for the vault. 1-255 characters. |
| [WorkspaceID](../managed/vault.go#L201)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Metadata](../managed/vault.go#L204)：`map[string]string` | json / `metadata` | 否 | Arbitrary key-value metadata to attach to the vault. Maximum 16 pairs, keys up to 64 chars, values up to 512 chars. |
| [Betas](../managed/vault.go#L206)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Vaults.Get`

```go
func (r *VaultService) Get(ctx context.Context, vaultID string, query VaultGetParams, opts ...option.RequestOption) (res *ManagedAgentsVault, err error)
```

**GET `/vaults/{vault_id}`** · [实现](../managed/vault.go#L56) · 获取资源详情。

路径参数：`vaultID string`，必填。

返回：`*ManagedAgentsVault, error`。[ManagedAgentsVault](../managed/vault.go#L157)。

请求对象 `query`：[VaultGetParams](../managed/vault.go#L218)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/vault.go#L219)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/vault.go#L221)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Vaults.List`

```go
func (r *VaultService) List(ctx context.Context, params VaultListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsVault], err error)
```

**GET `/vaults`** · [实现](../managed/vault.go#L72) · 列出资源。

返回：`*pagination.PageCursor[ManagedAgentsVault], error`。[ManagedAgentsVault](../managed/vault.go#L157)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[VaultListParams](../managed/vault.go#L225)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Name](../managed/vault.go#L226)：`param.Opt[string]` | query / `name` | 否 | 见字段定义。 |
| [BeforeID](../managed/vault.go#L228)：`param.Opt[string]` | query / `before_id` | 否 | 见字段定义。 |
| [AfterID](../managed/vault.go#L229)：`param.Opt[string]` | query / `after_id` | 否 | 见字段定义。 |
| [IncludeArchived](../managed/vault.go#L232)：`param.Opt[bool]` | query / `include_archived` | 否 | Whether to include archived vaults in the results. |
| [Limit](../managed/vault.go#L234)：`param.Opt[int64]` | query / `limit` | 否 | Maximum number of vaults to return per page. Defaults to 20, maximum 100. |
| [Page](../managed/vault.go#L236)：`param.Opt[string]` | query / `page` | 否 | Opaque pagination token from a previous `list_vaults` response. |
| [WorkspaceID](../managed/vault.go#L237)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/vault.go#L239)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Vaults.ListAutoPaging`

```go
func (r *VaultService) ListAutoPaging(ctx context.Context, params VaultListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsVault]
```

[实现](../managed/vault.go#L93)。自动遍历 `Vaults.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Vaults.Delete`

```go
func (r *VaultService) Delete(ctx context.Context, vaultID string, body VaultDeleteParams, opts ...option.RequestOption) (res *ManagedAgentsDeletedVault, err error)
```

**DELETE `/vaults/{vault_id}`** · [实现](../managed/vault.go#L98) · 删除资源。

路径参数：`vaultID string`，必填。

返回：`*ManagedAgentsDeletedVault, error`。[ManagedAgentsDeletedVault](../managed/vault.go#L130)。

请求对象 `body`：[VaultDeleteParams](../managed/vault.go#L251)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/vault.go#L252)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/vault.go#L254)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Vaults.Archive`

```go
func (r *VaultService) Archive(ctx context.Context, vaultID string, body VaultArchiveParams, opts ...option.RequestOption) (res *ManagedAgentsVault, err error)
```

**POST `/vaults/{vault_id}/archive`** · [实现](../managed/vault.go#L114) · 归档资源。

路径参数：`vaultID string`，必填。

返回：`*ManagedAgentsVault, error`。[ManagedAgentsVault](../managed/vault.go#L157)。

请求对象 `body`：[VaultArchiveParams](../managed/vault.go#L258)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/vault.go#L259)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/vault.go#L261)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Vaults.Credentials.New`

```go
func (r *VaultCredentialService) New(ctx context.Context, vaultID string, params VaultCredentialNewParams, opts ...option.RequestOption) (res *ManagedAgentsCredential, err error)
```

**POST `/vaults/{vault_id}/credentials`** · [实现](../managed/vaultcredential.go#L43) · 创建资源。

路径参数：`vaultID string`，必填。

返回：`*ManagedAgentsCredential, error`。[ManagedAgentsCredential](../managed/vaultcredential.go#L190)。

请求对象 `params`：[VaultCredentialNewParams](../managed/vaultcredential.go#L1521)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Auth](../managed/vaultcredential.go#L1523)：`VaultCredentialNewParamsAuthUnion` | json / `auth` | 是 | Authentication details for creating a credential.；类型：[VaultCredentialNewParamsAuthUnion](../managed/vaultcredential.go#L1546) |
| [DisplayName](../managed/vaultcredential.go#L1525)：`param.Opt[string]` | json / `display_name` | 否 | Human-readable name for the credential. Up to 255 characters. |
| [WorkspaceID](../managed/vaultcredential.go#L1526)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Metadata](../managed/vaultcredential.go#L1529)：`map[string]string` | json / `metadata` | 否 | Arbitrary key-value metadata to attach to the credential. Maximum 16 pairs, keys up to 64 chars, values up to 512 chars. |
| [Betas](../managed/vaultcredential.go#L1531)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Vaults.Credentials.Get`

```go
func (r *VaultCredentialService) Get(ctx context.Context, credentialID string, params VaultCredentialGetParams, opts ...option.RequestOption) (res *ManagedAgentsCredential, err error)
```

**GET `/vaults/{vault_id}/credentials/{credential_id}`** · [实现](../managed/vaultcredential.go#L59) · 获取资源详情。

路径参数：`credentialID string`，必填。

返回：`*ManagedAgentsCredential, error`。[ManagedAgentsCredential](../managed/vaultcredential.go#L190)。

请求对象 `params`：[VaultCredentialGetParams](../managed/vaultcredential.go#L1666)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [VaultID](../managed/vaultcredential.go#L1667)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [WorkspaceID](../managed/vaultcredential.go#L1668)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/vaultcredential.go#L1670)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Vaults.Credentials.Update`

```go
func (r *VaultCredentialService) Update(ctx context.Context, credentialID string, params VaultCredentialUpdateParams, opts ...option.RequestOption) (res *ManagedAgentsCredential, err error)
```

**POST `/vaults/{vault_id}/credentials/{credential_id}`** · [实现](../managed/vaultcredential.go#L79) · 更新资源。

路径参数：`credentialID string`，必填。

返回：`*ManagedAgentsCredential, error`。[ManagedAgentsCredential](../managed/vaultcredential.go#L190)。

请求对象 `params`：[VaultCredentialUpdateParams](../managed/vaultcredential.go#L1674)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [VaultID](../managed/vaultcredential.go#L1675)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [DisplayName](../managed/vaultcredential.go#L1678)：`param.Opt[string]` | json / `display_name` | 否 | Deprecated: Qoder does not support updating credential display names. Leave this field omitted; setting it returns an error before sending a request. |
| [WorkspaceID](../managed/vaultcredential.go#L1679)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Metadata](../managed/vaultcredential.go#L1682)：`map[string]any` | json / `metadata` | 否 | Metadata patch. Set a key to a string to upsert it, or to null to delete it. Omitted keys are preserved. |
| [Auth](../managed/vaultcredential.go#L1684)：`VaultCredentialUpdateParamsAuthUnion` | json / `auth` | 否 | Updated authentication details for a credential.；类型：[VaultCredentialUpdateParamsAuthUnion](../managed/vaultcredential.go#L1704) |
| [Betas](../managed/vaultcredential.go#L1686)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Vaults.Credentials.List`

```go
func (r *VaultCredentialService) List(ctx context.Context, vaultID string, params VaultCredentialListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsCredential], err error)
```

**GET `/vaults/{vault_id}/credentials`** · [实现](../managed/vaultcredential.go#L99) · 列出资源。

路径参数：`vaultID string`，必填。

返回：`*pagination.PageCursor[ManagedAgentsCredential], error`。[ManagedAgentsCredential](../managed/vaultcredential.go#L190)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[VaultCredentialListParams](../managed/vaultcredential.go#L1806)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Name](../managed/vaultcredential.go#L1807)：`param.Opt[string]` | query / `name` | 否 | 见字段定义。 |
| [BeforeID](../managed/vaultcredential.go#L1809)：`param.Opt[string]` | query / `before_id` | 否 | 见字段定义。 |
| [AfterID](../managed/vaultcredential.go#L1810)：`param.Opt[string]` | query / `after_id` | 否 | 见字段定义。 |
| [IncludeArchived](../managed/vaultcredential.go#L1813)：`param.Opt[bool]` | query / `include_archived` | 否 | Whether to include archived credentials in the results. |
| [Limit](../managed/vaultcredential.go#L1815)：`param.Opt[int64]` | query / `limit` | 否 | Maximum number of credentials to return per page. Defaults to 20, maximum 100. |
| [Page](../managed/vaultcredential.go#L1817)：`param.Opt[string]` | query / `page` | 否 | Opaque pagination token from a previous `list_credentials` response. |
| [WorkspaceID](../managed/vaultcredential.go#L1818)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/vaultcredential.go#L1820)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Vaults.Credentials.ListAutoPaging`

```go
func (r *VaultCredentialService) ListAutoPaging(ctx context.Context, vaultID string, params VaultCredentialListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsCredential]
```

[实现](../managed/vaultcredential.go#L124)。自动遍历 `Vaults.Credentials.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Vaults.Credentials.Delete`

```go
func (r *VaultCredentialService) Delete(ctx context.Context, credentialID string, params VaultCredentialDeleteParams, opts ...option.RequestOption) (res *ManagedAgentsDeletedCredential, err error)
```

**DELETE `/vaults/{vault_id}/credentials/{credential_id}`** · [实现](../managed/vaultcredential.go#L129) · 删除资源。

路径参数：`credentialID string`，必填。

返回：`*ManagedAgentsDeletedCredential, error`。[ManagedAgentsDeletedCredential](../managed/vaultcredential.go#L446)。

请求对象 `params`：[VaultCredentialDeleteParams](../managed/vaultcredential.go#L1833)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [VaultID](../managed/vaultcredential.go#L1834)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [WorkspaceID](../managed/vaultcredential.go#L1835)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/vaultcredential.go#L1837)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Vaults.Credentials.Archive`

```go
func (r *VaultCredentialService) Archive(ctx context.Context, credentialID string, params VaultCredentialArchiveParams, opts ...option.RequestOption) (res *ManagedAgentsCredential, err error)
```

**POST `/vaults/{vault_id}/credentials/{credential_id}/archive`** · [实现](../managed/vaultcredential.go#L149) · 归档资源。

路径参数：`credentialID string`，必填。

返回：`*ManagedAgentsCredential, error`。[ManagedAgentsCredential](../managed/vaultcredential.go#L190)。

请求对象 `params`：[VaultCredentialArchiveParams](../managed/vaultcredential.go#L1841)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [VaultID](../managed/vaultcredential.go#L1842)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [WorkspaceID](../managed/vaultcredential.go#L1843)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/vaultcredential.go#L1845)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Vaults.Credentials.MCPOAuthValidate`

```go
func (r *VaultCredentialService) MCPOAuthValidate(ctx context.Context, credentialID string, params VaultCredentialMCPOAuthValidateParams, opts ...option.RequestOption) (res *ManagedAgentsCredentialValidation, err error)
```

**POST `/vaults/{vault_id}/credentials/{credential_id}/mcp_oauth_validate`** · [实现](../managed/vaultcredential.go#L169) · 验证 MCP OAuth 凭据。

路径参数：`credentialID string`，必填。

返回：`*ManagedAgentsCredentialValidation, error`。[ManagedAgentsCredentialValidation](../managed/vaultcredential.go#L390)。

请求对象 `params`：[VaultCredentialMCPOAuthValidateParams](../managed/vaultcredential.go#L1849)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [VaultID](../managed/vaultcredential.go#L1850)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [WorkspaceID](../managed/vaultcredential.go#L1851)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/vaultcredential.go#L1853)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

## Files

上传文件、查询元数据、列举、下载内容和删除文件。上传后需另外挂载到会话才能供 Agent 使用。

### `Files.List`

```go
func (r *FileService) List(ctx context.Context, params FileListParams, opts ...option.RequestOption) (res *pagination.PageCursor[FileMetadata], err error)
```

**GET `/files`** · [实现](../managed/file.go#L47) · 列出资源。

返回：`*pagination.PageCursor[FileMetadata], error`。[FileMetadata](../managed/file.go#L162)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[FileListParams](../managed/file.go#L236)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Name](../managed/file.go#L237)：`param.Opt[string]` | query / `name` | 否 | 见字段定义。 |
| [BeforeID](../managed/file.go#L239)：`param.Opt[string]` | query / `before_id` | 否 | 见字段定义。 |
| [AfterID](../managed/file.go#L240)：`param.Opt[string]` | query / `after_id` | 否 | 见字段定义。 |
| [Page](../managed/file.go#L244)：`param.Opt[string]` | query / `page` | 否 | Opaque page cursor returned in a prior list response's `next_page`. Prefixed `page_`. |
| [Limit](../managed/file.go#L248)：`param.Opt[int64]` | query / `limit` | 否 | Number of items to return per page. Defaults to `20`. Ranges from `1` to `1000`. |
| [ScopeID](../managed/file.go#L251)：`param.Opt[string]` | query / `scope_id` | 否 | Filter by scope ID. Only returns files associated with the specified scope (e.g., a session ID). |
| [WorkspaceID](../managed/file.go#L252)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [IDs](../managed/file.go#L258)：`[]string` | query / `ids` | 否 | Restrict the result set to Files whose `id` is in this list.（完整约束见字段源码） |
| [Betas](../managed/file.go#L260)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Files.ListAutoPaging`

```go
func (r *FileService) ListAutoPaging(ctx context.Context, params FileListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[FileMetadata]
```

[实现](../managed/file.go#L68)。自动遍历 `Files.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Files.Delete`

```go
func (r *FileService) Delete(ctx context.Context, fileID string, body FileDeleteParams, opts ...option.RequestOption) (res *DeletedFile, err error)
```

**DELETE `/files/{file_id}`** · [实现](../managed/file.go#L73) · 删除资源。

路径参数：`fileID string`，必填。

返回：`*DeletedFile, error`。[DeletedFile](../managed/file.go#L129)。

请求对象 `body`：[FileDeleteParams](../managed/file.go#L272)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/file.go#L273)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/file.go#L275)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Files.Download`

```go
func (r *FileService) Download(ctx context.Context, fileID string, query FileDownloadParams, opts ...option.RequestOption) (res *http.Response, err error)
```

**GET `/files/{file_id}/content`** · [实现](../managed/file.go#L88) · 下载文件内容。

路径参数：`fileID string`，必填。

返回：`*http.Response, error`。HTTP 响应；读取后关闭 Body。

请求对象 `query`：[FileDownloadParams](../managed/file.go#L279)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/file.go#L280)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/file.go#L282)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Files.GetMetadata`

```go
func (r *FileService) GetMetadata(ctx context.Context, fileID string, query FileGetMetadataParams, opts ...option.RequestOption) (res *FileMetadata, err error)
```

**GET `/files/{file_id}`** · [实现](../managed/file.go#L104) · 获取文件元数据。

路径参数：`fileID string`，必填。

返回：`*FileMetadata, error`。[FileMetadata](../managed/file.go#L162)。

请求对象 `query`：[FileGetMetadataParams](../managed/file.go#L286)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/file.go#L287)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/file.go#L289)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Files.Upload`

```go
func (r *FileService) Upload(ctx context.Context, params FileUploadParams, opts ...option.RequestOption) (res *FileMetadata, err error)
```

**POST `/files`** · [实现](../managed/file.go#L119) · 上传文件。

返回：`*FileMetadata, error`。[FileMetadata](../managed/file.go#L162)。

请求对象 `params`：[FileUploadParams](../managed/file.go#L293)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Name](../managed/file.go#L294)：`param.Opt[string]` | multipart / `name` | 否 | 见字段定义。 |
| [Metadata](../managed/file.go#L296)：`map[string]string` | multipart / `metadata` | 否 | 见字段定义。 |
| [File](../managed/file.go#L301)：`io.Reader` | multipart / `file` | 是 | The file to upload. Only the final path component of the part's `filename` is kept; an absent or empty `filename` is replaced with `unnamed` plus the extension for the file's stored `mime_type`, when known. |
| [ExpiresInSeconds](../managed/file.go#L304)：`param.Opt[int64]` | multipart / `expires_in_seconds` | 否 | Seconds from upload until the file expires and its bytes become permanently unavailable. Must be between 3600 (one hour) and 7776000 (ninety days). |
| [WorkspaceID](../managed/file.go#L305)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/file.go#L307)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

## MemoryStores

管理记忆存储、记忆条目及版本。删除条目、删除存储、归档存储和版本脱敏是不同操作。

### `MemoryStores.New`

```go
func (r *MemoryStoreService) New(ctx context.Context, params MemoryStoreNewParams, opts ...option.RequestOption) (res *ManagedAgentsMemoryStore, err error)
```

**POST `/memory_stores`** · [实现](../managed/memorystore.go#L46) · 创建资源。

返回：`*ManagedAgentsMemoryStore, error`。[ManagedAgentsMemoryStore](../managed/memorystore.go#L178)。

请求对象 `params`：[MemoryStoreNewParams](../managed/memorystore.go#L229)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Name](../managed/memorystore.go#L234)：`string` | json / `name` | 是 | Human-readable name for the store.（完整约束见字段源码） |
| [Description](../managed/memorystore.go#L238)：`param.Opt[string]` | json / `description` | 否 | Free-text description of what the store contains, up to 1024 characters. Included in the agent's system prompt when the store is attached, so word it to be useful to the agent. |
| [WorkspaceID](../managed/memorystore.go#L239)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Metadata](../managed/memorystore.go#L243)：`map[string]string` | json / `metadata` | 否 | Arbitrary key-value tags for your own bookkeeping (such as the end user a store belongs to). Up to 16 pairs; keys 1–64 characters; values up to 512 characters. Not visible to the agent. |
| [Betas](../managed/memorystore.go#L245)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `MemoryStores.Get`

```go
func (r *MemoryStoreService) Get(ctx context.Context, memoryStoreID string, query MemoryStoreGetParams, opts ...option.RequestOption) (res *ManagedAgentsMemoryStore, err error)
```

**GET `/memory_stores/{memory_store_id}`** · [实现](../managed/memorystore.go#L58) · 获取资源详情。

路径参数：`memoryStoreID string`，必填。

返回：`*ManagedAgentsMemoryStore, error`。[ManagedAgentsMemoryStore](../managed/memorystore.go#L178)。

请求对象 `query`：[MemoryStoreGetParams](../managed/memorystore.go#L257)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/memorystore.go#L258)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/memorystore.go#L260)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `MemoryStores.Update`

```go
func (r *MemoryStoreService) Update(ctx context.Context, memoryStoreID string, params MemoryStoreUpdateParams, opts ...option.RequestOption) (res *ManagedAgentsMemoryStore, err error)
```

**POST `/memory_stores/{memory_store_id}`** · [实现](../managed/memorystore.go#L74) · 更新资源。

路径参数：`memoryStoreID string`，必填。

返回：`*ManagedAgentsMemoryStore, error`。[ManagedAgentsMemoryStore](../managed/memorystore.go#L178)。

请求对象 `params`：[MemoryStoreUpdateParams](../managed/memorystore.go#L264)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Description](../managed/memorystore.go#L267)：`param.Opt[string]` | json / `description` | 否 | New description for the store, up to 1024 characters. Pass an empty string to clear it. |
| [Name](../managed/memorystore.go#L271)：`param.Opt[string]` | json / `name` | 否 | New human-readable name for the store. 1–255 characters; no control characters. Renaming changes the slug used for the store's `mount_path` in sessions created after the update. |
| [WorkspaceID](../managed/memorystore.go#L272)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Metadata](../managed/memorystore.go#L276)：`map[string]any` | json / `metadata` | 否 | Metadata patch. Set a key to a string to upsert it, or to null to delete it. Omit the field to preserve. The stored bag is limited to 16 keys (up to 64 chars each) with values up to 512 chars. |
| [Betas](../managed/memorystore.go#L278)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `MemoryStores.List`

```go
func (r *MemoryStoreService) List(ctx context.Context, params MemoryStoreListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsMemoryStore], err error)
```

**GET `/memory_stores`** · [实现](../managed/memorystore.go#L90) · 列出资源。

返回：`*pagination.PageCursor[ManagedAgentsMemoryStore], error`。[ManagedAgentsMemoryStore](../managed/memorystore.go#L178)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[MemoryStoreListParams](../managed/memorystore.go#L290)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Name](../managed/memorystore.go#L291)：`param.Opt[string]` | query / `name` | 否 | 见字段定义。 |
| [CreatedAtGte](../managed/memorystore.go#L295)：`param.Opt[time.Time]` | query / `created_at[gte]` | 否 | Return only stores whose `created_at` is at or after this time (inclusive). Sent on the wire as `created_at[gte]`. |
| [CreatedAtLte](../managed/memorystore.go#L298)：`param.Opt[time.Time]` | query / `created_at[lte]` | 否 | Return only stores whose `created_at` is at or before this time (inclusive). Sent on the wire as `created_at[lte]`. |
| [IncludeArchived](../managed/memorystore.go#L301)：`param.Opt[bool]` | query / `include_archived` | 否 | When `true`, archived stores are included in the results. Defaults to `false` (archived stores are excluded). |
| [Limit](../managed/memorystore.go#L304)：`param.Opt[int64]` | query / `limit` | 否 | Maximum number of stores to return per page. Must be between 1 and 100. Defaults to 20 when omitted. |
| [Page](../managed/memorystore.go#L307)：`param.Opt[string]` | query / `page` | 否 | Opaque pagination cursor (a `page_...` value). Pass the `next_page` value from a previous response to fetch the next page; omit for the first page. |
| [WorkspaceID](../managed/memorystore.go#L308)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/memorystore.go#L310)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `MemoryStores.ListAutoPaging`

```go
func (r *MemoryStoreService) ListAutoPaging(ctx context.Context, params MemoryStoreListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsMemoryStore]
```

[实现](../managed/memorystore.go#L111)。自动遍历 `MemoryStores.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `MemoryStores.Delete`

```go
func (r *MemoryStoreService) Delete(ctx context.Context, memoryStoreID string, body MemoryStoreDeleteParams, opts ...option.RequestOption) (res *ManagedAgentsDeletedMemoryStore, err error)
```

**DELETE `/memory_stores/{memory_store_id}`** · [实现](../managed/memorystore.go#L116) · 删除资源。

路径参数：`memoryStoreID string`，必填。

返回：`*ManagedAgentsDeletedMemoryStore, error`。[ManagedAgentsDeletedMemoryStore](../managed/memorystore.go#L148)。

请求对象 `body`：[MemoryStoreDeleteParams](../managed/memorystore.go#L323)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/memorystore.go#L324)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/memorystore.go#L326)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `MemoryStores.Archive`

```go
func (r *MemoryStoreService) Archive(ctx context.Context, memoryStoreID string, body MemoryStoreArchiveParams, opts ...option.RequestOption) (res *ManagedAgentsMemoryStore, err error)
```

**POST `/memory_stores/{memory_store_id}/archive`** · [实现](../managed/memorystore.go#L132) · 归档资源。

路径参数：`memoryStoreID string`，必填。

返回：`*ManagedAgentsMemoryStore, error`。[ManagedAgentsMemoryStore](../managed/memorystore.go#L178)。

请求对象 `body`：[MemoryStoreArchiveParams](../managed/memorystore.go#L330)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [WorkspaceID](../managed/memorystore.go#L331)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/memorystore.go#L333)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `MemoryStores.Memories.New`

```go
func (r *MemoryStoreMemoryService) New(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryNewParams, opts ...option.RequestOption) (res *ManagedAgentsMemory, err error)
```

**POST `/memory_stores/{memory_store_id}/memories`** · [实现](../managed/memorystorememory.go#L43) · 创建资源。

路径参数：`memoryStoreID string`，必填。

返回：`*ManagedAgentsMemory, error`。[ManagedAgentsMemory](../managed/memorystorememory.go#L186)。

请求对象 `params`：[MemoryStoreMemoryNewParams](../managed/memorystorememory.go#L418)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Metadata](../managed/memorystorememory.go#L419)：`map[string]string` | json / `metadata` | 否 | 见字段定义。 |
| [Content](../managed/memorystorememory.go#L423)：`param.Opt[string]` | json / `content` | 是 | UTF-8 text content for the new memory. Maximum 100 kB (102,400 bytes). Required; pass `""` explicitly to create an empty memory. |
| [Path](../managed/memorystorememory.go#L429)：`string` | json / `path` | 是 | Hierarchical path for the new memory, e.g.（完整约束见字段源码） |
| [WorkspaceID](../managed/memorystorememory.go#L430)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [View](../managed/memorystorememory.go#L434)：`ManagedAgentsMemoryView` | query / `view` | 否 | Query parameter for view Any of "basic", "full".；类型：[ManagedAgentsMemoryView](../managed/memorystorememory.go#L378) |
| [Betas](../managed/memorystorememory.go#L436)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `MemoryStores.Memories.Get`

```go
func (r *MemoryStoreMemoryService) Get(ctx context.Context, memoryID string, params MemoryStoreMemoryGetParams, opts ...option.RequestOption) (res *ManagedAgentsMemory, err error)
```

**GET `/memory_stores/{memory_store_id}/memories/{memory_id}`** · [实现](../managed/memorystorememory.go#L59) · 获取资源详情。

路径参数：`memoryID string`，必填。

返回：`*ManagedAgentsMemory, error`。[ManagedAgentsMemory](../managed/memorystorememory.go#L186)。

请求对象 `params`：[MemoryStoreMemoryGetParams](../managed/memorystorememory.go#L457)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [MemoryStoreID](../managed/memorystorememory.go#L458)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [WorkspaceID](../managed/memorystorememory.go#L459)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [View](../managed/memorystorememory.go#L463)：`ManagedAgentsMemoryView` | query / `view` | 否 | Query parameter for view Any of "basic", "full".；类型：[ManagedAgentsMemoryView](../managed/memorystorememory.go#L378) |
| [Betas](../managed/memorystorememory.go#L465)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `MemoryStores.Memories.Update`

```go
func (r *MemoryStoreMemoryService) Update(ctx context.Context, memoryID string, params MemoryStoreMemoryUpdateParams, opts ...option.RequestOption) (res *ManagedAgentsMemory, err error)
```

**POST `/memory_stores/{memory_store_id}/memories/{memory_id}`** · [实现](../managed/memorystorememory.go#L79) · 更新资源。

路径参数：`memoryID string`，必填。

返回：`*ManagedAgentsMemory, error`。[ManagedAgentsMemory](../managed/memorystorememory.go#L186)。

请求对象 `params`：[MemoryStoreMemoryUpdateParams](../managed/memorystorememory.go#L478)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [ContentSha256](../managed/memorystorememory.go#L479)：`param.Opt[string]` | json / `content_sha256` | 否 | 见字段定义。 |
| [Metadata](../managed/memorystorememory.go#L481)：`map[string]any` | json / `metadata` | 否 | 见字段定义。 |
| [MemoryStoreID](../managed/memorystorememory.go#L483)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [Content](../managed/memorystorememory.go#L486)：`param.Opt[string]` | json / `content` | 否 | New UTF-8 text content for the memory. Maximum 100 kB (102,400 bytes). Omit to leave the content unchanged (e.g., for a rename-only update). |
| [Path](../managed/memorystorememory.go#L493)：`param.Opt[string]` | json / `path` | 否 | New path for the memory (a rename).（完整约束见字段源码） |
| [WorkspaceID](../managed/memorystorememory.go#L494)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [View](../managed/memorystorememory.go#L498)：`ManagedAgentsMemoryView` | query / `view` | 否 | Query parameter for view Any of "basic", "full".；类型：[ManagedAgentsMemoryView](../managed/memorystorememory.go#L378) |
| [Precondition](../managed/memorystorememory.go#L505)：`ManagedAgentsPreconditionParam` | json / `precondition` | 否 | Optimistic-concurrency precondition: the update applies only if the memory's stored `content_sha256` equals the supplied value.（完整约束见字段源码）；类型：[ManagedAgentsPreconditionParam](../managed/memorystorememory.go#L393) |
| [Betas](../managed/memorystorememory.go#L507)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `MemoryStores.Memories.List`

```go
func (r *MemoryStoreMemoryService) List(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsMemoryListItemUnion], err error)
```

**GET `/memory_stores/{memory_store_id}/memories`** · [实现](../managed/memorystorememory.go#L99) · 列出资源。

路径参数：`memoryStoreID string`，必填。

返回：`*pagination.PageCursor[ManagedAgentsMemoryListItemUnion], error`。[ManagedAgentsMemoryListItemUnion](../managed/memorystorememory.go#L259)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[MemoryStoreMemoryListParams](../managed/memorystorememory.go#L539)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Depth](../managed/memorystorememory.go#L543)：`param.Opt[int64]` | query / `depth` | 否 | `0` (or omitted) returns all descendants below `path_prefix` (recursive).（完整约束见字段源码） |
| [Limit](../managed/memorystorememory.go#L547)：`param.Opt[int64]` | query / `limit` | 否 | Maximum number of items to return per page. Must be between 1 and 100. Defaults to 20 when omitted. Capped at 20 when `view=full`. Both `memory` and `memory_prefix` items count toward the limit. |
| [Page](../managed/memorystorememory.go#L550)：`param.Opt[string]` | query / `page` | 否 | Opaque pagination cursor (a `page_...` value). Pass the `next_page` value from a previous response to fetch the next page; omit for the first page. |
| [PathPrefix](../managed/memorystorememory.go#L554)：`param.Opt[string]` | query / `path_prefix` | 否 | Optional path prefix filter. Must end with `/` (segment-aligned), e.g., `/notes/`. This value appears in request URLs. Do not include secrets or personally identifiable information. |
| [WorkspaceID](../managed/memorystorememory.go#L555)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [View](../managed/memorystorememory.go#L561)：`ManagedAgentsMemoryView` | query / `view` | 否 | Which projection of each `memory` to return.（完整约束见字段源码）；类型：[ManagedAgentsMemoryView](../managed/memorystorememory.go#L378) |
| [Betas](../managed/memorystorememory.go#L563)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `MemoryStores.Memories.ListAutoPaging`

```go
func (r *MemoryStoreMemoryService) ListAutoPaging(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsMemoryListItemUnion]
```

[实现](../managed/memorystorememory.go#L124)。自动遍历 `MemoryStores.Memories.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `MemoryStores.Memories.Delete`

```go
func (r *MemoryStoreMemoryService) Delete(ctx context.Context, memoryID string, params MemoryStoreMemoryDeleteParams, opts ...option.RequestOption) (res *ManagedAgentsDeletedMemory, err error)
```

**DELETE `/memory_stores/{memory_store_id}/memories/{memory_id}`** · [实现](../managed/memorystorememory.go#L129) · 删除资源。

路径参数：`memoryID string`，必填。

返回：`*ManagedAgentsDeletedMemory, error`。[ManagedAgentsDeletedMemory](../managed/memorystorememory.go#L154)。

请求对象 `params`：[MemoryStoreMemoryDeleteParams](../managed/memorystorememory.go#L576)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [MemoryStoreID](../managed/memorystorememory.go#L577)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [ExpectedContentSha256](../managed/memorystorememory.go#L579)：`param.Opt[string]` | query / `expected_content_sha256` | 否 | Query parameter for expected_content_sha256 |
| [WorkspaceID](../managed/memorystorememory.go#L580)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/memorystorememory.go#L582)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `MemoryStores.MemoryVersions.Get`

```go
func (r *MemoryStoreMemoryVersionService) Get(ctx context.Context, memoryVersionID string, params MemoryStoreMemoryVersionGetParams, opts ...option.RequestOption) (res *ManagedAgentsMemoryVersion, err error)
```

**GET `/memory_stores/{memory_store_id}/memory_versions/{memory_version_id}`** · [实现](../managed/memorystorememoryversion.go#L44) · 获取资源详情。

路径参数：`memoryVersionID string`，必填。

返回：`*ManagedAgentsMemoryVersion, error`。[ManagedAgentsMemoryVersion](../managed/memorystorememoryversion.go#L241)。

请求对象 `params`：[MemoryStoreMemoryVersionGetParams](../managed/memorystorememoryversion.go#L408)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [MemoryStoreID](../managed/memorystorememoryversion.go#L409)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [WorkspaceID](../managed/memorystorememoryversion.go#L410)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [View](../managed/memorystorememoryversion.go#L414)：`ManagedAgentsMemoryView` | query / `view` | 否 | Query parameter for view Any of "basic", "full".；类型：[ManagedAgentsMemoryView](../managed/memorystorememory.go#L378) |
| [Betas](../managed/memorystorememoryversion.go#L416)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `MemoryStores.MemoryVersions.List`

```go
func (r *MemoryStoreMemoryVersionService) List(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryVersionListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsMemoryVersion], err error)
```

**GET `/memory_stores/{memory_store_id}/memory_versions`** · [实现](../managed/memorystorememoryversion.go#L64) · 列出资源。

路径参数：`memoryStoreID string`，必填。

返回：`*pagination.PageCursor[ManagedAgentsMemoryVersion], error`。[ManagedAgentsMemoryVersion](../managed/memorystorememoryversion.go#L241)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[MemoryStoreMemoryVersionListParams](../managed/memorystorememoryversion.go#L429)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [APIKeyID](../managed/memorystorememoryversion.go#L431)：`param.Opt[string]` | query / `api_key_id` | 否 | Query parameter for api_key_id |
| [CreatedAtGte](../managed/memorystorememoryversion.go#L433)：`param.Opt[time.Time]` | query / `created_at[gte]` | 否 | Return versions created at or after this time (inclusive). |
| [CreatedAtLte](../managed/memorystorememoryversion.go#L435)：`param.Opt[time.Time]` | query / `created_at[lte]` | 否 | Return versions created at or before this time (inclusive). |
| [Limit](../managed/memorystorememoryversion.go#L437)：`param.Opt[int64]` | query / `limit` | 否 | Query parameter for limit |
| [MemoryID](../managed/memorystorememoryversion.go#L439)：`param.Opt[string]` | query / `memory_id` | 否 | Query parameter for memory_id |
| [Page](../managed/memorystorememoryversion.go#L441)：`param.Opt[string]` | query / `page` | 否 | Query parameter for page |
| [ServiceAccountID](../managed/memorystorememoryversion.go#L443)：`param.Opt[string]` | query / `service_account_id` | 否 | Query parameter for service_account_id |
| [SessionID](../managed/memorystorememoryversion.go#L445)：`param.Opt[string]` | query / `session_id` | 否 | Query parameter for session_id |
| [WorkspaceID](../managed/memorystorememoryversion.go#L446)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Operation](../managed/memorystorememoryversion.go#L450)：`ManagedAgentsMemoryVersionOperation` | query / `operation` | 否 | Query parameter for operation Any of "created", "modified", "deleted".；类型：[ManagedAgentsMemoryVersionOperation](../managed/memorystorememoryversion.go#L322) |
| [View](../managed/memorystorememoryversion.go#L454)：`ManagedAgentsMemoryView` | query / `view` | 否 | Query parameter for view Any of "basic", "full".；类型：[ManagedAgentsMemoryView](../managed/memorystorememory.go#L378) |
| [Betas](../managed/memorystorememoryversion.go#L456)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `MemoryStores.MemoryVersions.ListAutoPaging`

```go
func (r *MemoryStoreMemoryVersionService) ListAutoPaging(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryVersionListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsMemoryVersion]
```

[实现](../managed/memorystorememoryversion.go#L89)。自动遍历 `MemoryStores.MemoryVersions.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `MemoryStores.MemoryVersions.Redact`

```go
func (r *MemoryStoreMemoryVersionService) Redact(ctx context.Context, memoryVersionID string, params MemoryStoreMemoryVersionRedactParams, opts ...option.RequestOption) (res *ManagedAgentsMemoryVersion, err error)
```

**POST `/memory_stores/{memory_store_id}/memory_versions/{memory_version_id}/redact`** · [实现](../managed/memorystorememoryversion.go#L94) · 对记忆版本脱敏。

路径参数：`memoryVersionID string`，必填。

返回：`*ManagedAgentsMemoryVersion, error`。[ManagedAgentsMemoryVersion](../managed/memorystorememoryversion.go#L241)。

请求对象 `params`：[MemoryStoreMemoryVersionRedactParams](../managed/memorystorememoryversion.go#L469)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [MemoryStoreID](../managed/memorystorememoryversion.go#L470)：`string` | 嵌入 / `` | 是 | 见字段定义。 |
| [WorkspaceID](../managed/memorystorememoryversion.go#L471)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/memorystorememoryversion.go#L473)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

## Models

查询当前账号可用模型和能力；创建模板或 Agent 时应使用列表中启用的模型 ID。

### `Models.List`

```go
func (r *ModelService) List(ctx context.Context, params ModelListParams, opts ...option.RequestOption) (res *pagination.Page[ModelInfo], err error)
```

**GET `/models`** · [实现](../managed/model.go#L46) · 列出资源。

返回：`*pagination.Page[ModelInfo], error`。[ModelInfo](../managed/model.go#L194)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[ModelListParams](../managed/model.go#L297)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [AfterID](../managed/model.go#L300)：`param.Opt[string]` | query / `after_id` | 否 | ID of the object to use as a cursor for pagination. When provided, returns the page of results immediately after this object. |
| [BeforeID](../managed/model.go#L303)：`param.Opt[string]` | query / `before_id` | 否 | ID of the object to use as a cursor for pagination. When provided, returns the page of results immediately before this object. |
| [Limit](../managed/model.go#L307)：`param.Opt[int64]` | query / `limit` | 否 | Number of items to return per page. Defaults to `20`. Ranges from `1` to `1000`. |
| [WorkspaceID](../managed/model.go#L308)：`param.Opt[string]` | header / `qoder-workspace-id` | 否 | 见字段定义。 |
| [Betas](../managed/model.go#L310)：`[]QoderBeta` | header / `x-qoder-beta` | 否 | Optional header to specify the beta version(s) you want to use.；类型：[QoderBeta](../managed/qoder.go#L13) |

### `Models.ListAutoPaging`

```go
func (r *ModelService) ListAutoPaging(ctx context.Context, params ModelListParams, opts ...option.RequestOption) *pagination.PageAutoPager[ModelInfo]
```

[实现](../managed/model.go#L70)。自动遍历 `Models.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

## 生命周期与完整场景

Agent、Environment 与 Session 分别有独立生命周期。Deployment 触发后通过 `DeploymentRuns` 查看运行，Dream 需要确认完成及输出 Memory Store 内容。Work 的 Poll、Ack、Heartbeat、Stop 是自托管 worker 协议。Managed 没有 Sessions.Cancel，使用当前资源表提供的生命周期 API。

可运行的资源准备、对话验证与清理流程见 [Managed 场景目录](../examples/managed)。
