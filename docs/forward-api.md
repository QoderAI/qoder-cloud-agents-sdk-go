# Forward API 参考（Go）

本文覆盖当前 `forward` 客户端的 **110 个 HTTP 操作**，以及 **21 个自动分页便利方法**。方法签名、请求字段和返回类型以本仓库源码为准；每个操作附实现及类型链接，嵌套结构、联合分支和完整约束可直接进入字段定义查看。所有 HTTP 路由相对于 `/api/v1/forward`。

[返回 README](../README.md) · [另一模式 API](managed-api.md)

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
    "github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func main() {
    token := os.Getenv("QODER_ACCESS_TOKEN")
    if token == "" { log.Fatal("请设置 QODER_ACCESS_TOKEN") }
    client := forward.NewClient(
        option.WithAccessToken(token),
        option.WithBaseURL("https://api.qoder.com/api/v1/forward"),
    )
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    models, err := client.Models.List(ctx)
    if err != nil { log.Fatal(err) }
    for _, model := range models.Data {
        if model.IsEnabled { fmt.Println(model.ID) }
    }
}
```

| 配置 | 值与优先级 |
| --- | --- |
| PAT | `option.WithAccessToken(token)`；未指定时读取 `QODER_ACCESS_TOKEN`，使用 Bearer 鉴权 |
| 默认地址 | `https://api.qoder.com/api/v1/forward/`（国际站） |
| CN 地址 | `https://api.qoder.com.cn/api/v1/forward/`，通过 `WithBaseURL` 指定 |
| 地址环境变量 | `QODER_FORWARD_BASE_URL`；显式 option 优先于环境变量 |
| `.env` | SDK 不自动加载；可执行 examples 的配置加载器另行读取 `.env.live` |
| 动态凭据 | `option.WithCredential(provider)`，实现 `convention.Credential`；已有 Authorization 优先 |

全部资源从 `client` 访问，子资源保留层级，例如 `client.Sessions.Threads.Events`。在共享客户端前完成配置，避免并发修改 `Options`。

## 参数、响应与请求选项

每个操作均接收调用方的 `ctx context.Context`，末尾的 `opts ...option.RequestOption` 可省略。签名中的类型均属于 `forward` 包（`context`、`option`、`pagination`、`ssestream`、`http` 等显式前缀除外）。实际调用时使用 `forward.类型名` 构造请求。

参数表的“必填”来自 `api:"required"` 标记；路径参数也必须提供，且非空。可选不代表任意组合均有效：互斥字段、联合类型、枚举、权限及跨字段约束以链接的类型定义和服务端校验为准。空参数结构仍需在签名指定的位置传入。

- 可选标量使用 `forward.String("value")`、`forward.Int(20)`、`forward.Bool(false)`；`Bool(false)` 会发送 false，零值 `param.Opt` 则省略字段。
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

HTTP 非成功状态返回 `*convention.Error`（`forward.Error` 是其别名）；网络、context、参数及解码错误可能是其他类型。SSE 的建连和读取错误必须通过 `stream.Err()` 检查。

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

以下片段使用已有的 Identity 和 Template；Template 应配置可用 Environment。资源准备及清理见 [会话示例](../examples/forward/session/main.go)。函数片段需导入 `context`、`fmt` 和 `forward` 包。

```go
func startSession(ctx context.Context, client forward.Client, identityID, templateID string) (*forward.Session, error) {
    return client.Sessions.New(ctx, forward.SessionNewParams{
        IdentityID: identityID, TemplateID: templateID,
        Title: forward.String("项目助手"),
    })
}

func sendMessage(ctx context.Context, client forward.Client, sessionID, text, key string) (string, error) {
    result, err := client.Sessions.Events.Send(ctx, sessionID, forward.SessionEventSendParams{
        Events: []forward.SessionEventParam{{
            Type: "user.message",
            Content: forward.EventContentUnionParam{OfBlocks: []forward.ContentBlockParam{{
                Type: "text", Text: forward.String(text),
            }}},
        }},
        IdempotencyKey: forward.String(key),
    })
    if err != nil { return "", err }
    if len(result.Data) != 1 { return "", fmt.Errorf("expected one user event") }
    return result.Data[0].ID, nil
}

func readTurn(ctx context.Context, client forward.Client, sessionID, afterID string) error {
    stream := client.Sessions.Events.StreamEvents(ctx, sessionID, forward.SessionEventStreamParams{
        LastEventID: forward.String(afterID),
    })
    defer stream.Close()
    for stream.Next() {
        event := stream.Current()
        switch event.Type {
        case "agent.message":
            blocks, err := event.ContentBlocks()
            if err != nil { return err }
            for _, block := range blocks {
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

`key` 是调用方保存的非空幂等键。`sendMessage` 返回的用户事件 ID 传给 `readTurn`，避免处理旧轮次；续订可把最后已消费的完整事件 ID 放在 `LastEventID`。流不会因某一条消息自动关闭，主动退出并 `Close()`。文本增量、tool/thinking 筛选由 `SessionEventStreamParams` 控制；线程事件使用 `Sessions.Threads.Events.StreamEvents`。

## 分页

```go
func listAll(ctx context.Context, client forward.Client) error {
    pager := client.Templates.ListAutoPaging(ctx, forward.TemplateListParams{
        Limit: forward.Int(20), Status: forward.String("active"),
    })
    for pager.Next() { fmt.Println(pager.Current().ID) }
    return pager.Err()
}
```

`pagination.Page[T]` 使用 `AfterID` / `BeforeID`；`pagination.PageCursor[T]` 使用响应 `NextPage` 传给后续 `Page`。可手动调用 `page.GetNextPage()`，无下一页返回 `(nil, nil)`。使用签名标明的分页类型，不混用游标。`Models.List`、`Identities.ListTemplates` 和 `Identities.MemoryStores.List` 返回独立列表响应，没有 `ListAutoPaging`。

## 上传与下载

以下片段需导入 `context`、`io`、`strings`、`convention` 和 `forward`。上传使用 multipart；上传完成不等于挂载，资源绑定见 [资源示例](../examples/forward/resources/main.go)。

```go
func uploadText(ctx context.Context, client forward.Client, text string) (*forward.FileMetadata, error) {
    return client.Files.Upload(ctx, forward.FileUploadParams{
        File: convention.UploadFile{
            Reader: strings.NewReader(text), Name: "notes.txt", MediaType: "text/plain",
        },
        Purpose: forward.String("session_resource"),
    })
}

func downloadFile(ctx context.Context, client forward.Client, fileID string, destination io.Writer) error {
    response, err := client.Files.Download(ctx, fileID)
    if err != nil { return err }
    defer response.Body.Close()
    _, err = io.Copy(destination, response.Body)
    return err
}
```

SDK 先获取临时下载地址，再读取内容；不会把 API PAT、请求头或 middleware 附加到存储请求，自定义 HTTP transport 的行为仍由应用负责。`Skills.Versions.Download` 同样返回需关闭 Body 的 HTTP 响应。文件上传参数接受 `io.Reader`，本地文件需在调用结束后关闭。

## 资源与方法目录

下文列出全部 HTTP 操作以及每个资源实际提供的自动分页方法。参数表覆盖每个请求对象的公开一级字段，嵌入对象字段展开；点入类型可查看嵌套字段、允许值和完整说明。源码注释摘要保留原语言。

- [Templates](#templates)：6 个 HTTP 操作。
- [Identities](#identities)：18 个 HTTP 操作。
- [Sessions](#sessions)：15 个 HTTP 操作。
- [Schedules](#schedules)：9 个 HTTP 操作。
- [ScheduleRuns](#scheduleruns)：2 个 HTTP 操作。
- [Batches](#batches)：7 个 HTTP 操作。
- [Channels](#channels)：7 个 HTTP 操作。
- [ChannelPairings](#channelpairings)：2 个 HTTP 操作。
- [Environments](#environments)：6 个 HTTP 操作。
- [Files](#files)：5 个 HTTP 操作。
- [Skills](#skills)：10 个 HTTP 操作。
- [Vaults](#vaults)：8 个 HTTP 操作。
- [MemoryStores](#memorystores)：14 个 HTTP 操作。
- [Models](#models)：1 个 HTTP 操作。

## Templates

定义模型、系统指令、工具、资源与默认 Environment；通过 Clone 复制、Archive 归档。

### `Templates.List`

```go
func (r *TemplateService) List(ctx context.Context, params TemplateListParams, opts ...option.RequestOption) (res *pagination.Page[Template], err error)
```

**GET `/templates`** · [实现](../forward/template.go#L30) · 列出资源。

返回：`*pagination.Page[Template], error`。[Template](../forward/template.go#L249)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[TemplateListParams](../forward/template.go#L50)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Status](../forward/template.go#L52)：`param.Opt[string]` | query / `status` | 否 | 按 `active` 或 `archived` 过滤。 |
| [Limit](../forward/template.go#L54)：`param.Opt[int64]` | query / `limit` | 否 | 分页大小，最大 100。 |
| [AfterID](../forward/template.go#L56)：`param.Opt[string]` | query / `after_id` | 否 | 向后翻页游标，不能与 `before_id` 同用。 |
| [BeforeID](../forward/template.go#L58)：`param.Opt[string]` | query / `before_id` | 否 | 向前翻页游标，不能与 `after_id` 同用。 |

### `Templates.ListAutoPaging`

```go
func (r *TemplateService) ListAutoPaging(ctx context.Context, params TemplateListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Template]
```

[实现](../forward/template.go#L46)。自动遍历 `Templates.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Templates.New`

```go
func (r *TemplateService) New(ctx context.Context, params TemplateNewParams, opts ...option.RequestOption) (res *Template, err error)
```

**POST `/templates`** · [实现](../forward/template.go#L67) · 创建资源。

返回：`*Template, error`。[Template](../forward/template.go#L249)。

请求对象 `params`：[TemplateNewParams](../forward/template.go#L80)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Name](../forward/template.go#L82)：`string` | json / `name` | 是 | Template 名称，1-256 个字符，租户内唯一。 |
| [Model](../forward/template.go#L84)：`ModelConfigUnionParam` | json / `model` | 是 | 模型标识。可传 string（如 `"ultimate"`），或传 Agent model 对象以同时配置 `effort` 或 `context_window`。可通过列出模型接口查询可用值。；类型：[ModelConfigUnionParam](../forward/params.go#L59) |
| [EnvironmentID](../forward/template.go#L86)：`string` | json / `environment_id` | 是 | 创建 Session 时默认使用的 Environment ID。 |
| [Description](../forward/template.go#L88)：`param.Opt[string]` | json / `description` | 否 | Template 描述，最多 2048 个字符。 |
| [System](../forward/template.go#L90)：`param.Opt[string]` | json / `system` | 否 | System Prompt，最多 100,000 个字符。 |
| [Tools](../forward/template.go#L92)：`[]ToolParam` | json / `tools` | 否 | 工具配置列表，最多 128 项。；类型：[ToolParam](../forward/configuration_types.go#L146) |
| [MCPServers](../forward/template.go#L94)：`[]MCPServerParam` | json / `mcp_servers` | 否 | MCP Server 配置列表，最多 20 项。；类型：[MCPServerParam](../forward/configuration_types.go#L180) |
| [Skills](../forward/template.go#L96)：`[]SkillBindingParam` | json / `skills` | 否 | Skill 绑定列表，最多 20 项。；类型：[SkillBindingParam](../forward/configuration_types.go#L211) |
| [Multiagent](../forward/template.go#L98)：`MultiagentConfigParam` | json / `multiagent` | 否 | Multi-agent 协作配置。`type` 必须为 `coordinator`；省略或传 `null` 表示不启用。；类型：[MultiagentConfigParam](../forward/configuration_types.go#L270) |
| [Vaults](../forward/template.go#L100)：`map[string]ResourceBindingParam` | json / `vaults` | 否 | 默认 Vault 配置，按 Vault ID 组织。；类型：[ResourceBindingParam](../forward/configuration_types.go#L21) |
| [Files](../forward/template.go#L102)：`map[string]ResourceBindingParam` | json / `files` | 否 | 默认文件资源配置，按 file ID 组织。；类型：[ResourceBindingParam](../forward/configuration_types.go#L21) |
| [GitHubRepositories](../forward/template.go#L104)：`map[string]GitHubRepositoryParam` | json / `github_repositories` | 否 | 默认 GitHub 仓库配置，按调用方指定的 binding key 组织，最多 20 项。；类型：[GitHubRepositoryParam](../forward/configuration_types.go#L50) |
| [EnvironmentVariables](../forward/template.go#L106)：`EnvironmentVariablesUnionParam` | json / `environment_variables` | 否 | 默认 Session 环境变量。；类型：[EnvironmentVariablesUnionParam](../forward/params.go#L74) |
| [Metadata](../forward/template.go#L108)：`map[string]any` | json / `metadata` | 否 | 自定义元数据。 |
| [IdempotencyKey](../forward/template.go#L110)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |
| [QoderBeta](../forward/template.go#L112)：`param.Opt[string]` | header / `X-Qoder-Beta` | 否 | 启用 Browser Use 时必须设置为 `browser-use-2026-07-14`。 |

### `Templates.Get`

```go
func (r *TemplateService) Get(ctx context.Context, templateID string, opts ...option.RequestOption) (res *Template, err error)
```

**GET `/templates/{template_id}`** · [实现](../forward/template.go#L123) · 获取资源详情。

路径参数：`templateID string`，必填。

返回：`*Template, error`。[Template](../forward/template.go#L249)。

无请求参数对象；可直接传入请求选项。

### `Templates.Update`

```go
func (r *TemplateService) Update(ctx context.Context, templateID string, params TemplateUpdateParams, opts ...option.RequestOption) (res *Template, err error)
```

**POST `/templates/{template_id}`** · [实现](../forward/template.go#L135) · 更新资源。

路径参数：`templateID string`，必填。

返回：`*Template, error`。[Template](../forward/template.go#L249)。

请求对象 `params`：[TemplateUpdateParams](../forward/template.go#L151)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Name](../forward/template.go#L153)：`param.Opt[string]` | json / `name` | 否 | 新的 Template 名称。 |
| [Description](../forward/template.go#L155)：`param.Opt[string]` | json / `description` | 否 | 新的 Template 描述。 |
| [Model](../forward/template.go#L157)：`ModelConfigUnionParam` | json / `model` | 否 | 新的模型标识。可传 string，或传 Agent model 对象以同时配置 `effort` 或 `context_window`。可通过列出模型接口查询可用值。；类型：[ModelConfigUnionParam](../forward/params.go#L59) |
| [System](../forward/template.go#L159)：`param.Opt[string]` | json / `system` | 否 | 新的 System Prompt。 |
| [Tools](../forward/template.go#L161)：`[]ToolParam` | json / `tools` | 否 | 整体替换工具配置列表。；类型：[ToolParam](../forward/configuration_types.go#L146) |
| [MCPServers](../forward/template.go#L163)：`[]MCPServerParam` | json / `mcp_servers` | 否 | 整体替换 MCP Server 列表。；类型：[MCPServerParam](../forward/configuration_types.go#L180) |
| [Skills](../forward/template.go#L165)：`[]SkillBindingParam` | json / `skills` | 否 | 整体替换 Skill 绑定列表。；类型：[SkillBindingParam](../forward/configuration_types.go#L211) |
| [Multiagent](../forward/template.go#L167)：`MultiagentConfigParam` | json / `multiagent` | 否 | 整体替换 Multi-agent 协作配置；传 `null` 表示清空，省略则保留当前配置。；类型：[MultiagentConfigParam](../forward/configuration_types.go#L270) |
| [EnvironmentID](../forward/template.go#L169)：`param.Opt[string]` | json / `environment_id` | 否 | 替换默认 Environment ID；`null` 或空字符串表示清空。；允许 null。 |
| [Vaults](../forward/template.go#L171)：`map[string]ResourceBindingParam` | json / `vaults` | 否 | 整体替换默认 Vault 配置；按 Vault ID 组织，`null` 表示清空。；类型：[ResourceBindingParam](../forward/configuration_types.go#L21) |
| [Files](../forward/template.go#L173)：`map[string]ResourceBindingParam` | json / `files` | 否 | 整体替换默认文件资源配置；`null` 表示清空。；类型：[ResourceBindingParam](../forward/configuration_types.go#L21) |
| [GitHubRepositories](../forward/template.go#L175)：`map[string]GitHubRepositoryParam` | json / `github_repositories` | 否 | 整体替换默认 GitHub 仓库配置；按 binding key 组织，`null` 或空 object 表示清空。；类型：[GitHubRepositoryParam](../forward/configuration_types.go#L50) |
| [EnvironmentVariables](../forward/template.go#L177)：`EnvironmentVariablesUnionParam` | json / `environment_variables` | 否 | 整体替换默认环境变量；`null` 表示清空。；类型：[EnvironmentVariablesUnionParam](../forward/params.go#L74) |
| [Metadata](../forward/template.go#L179)：`map[string]any` | json / `metadata` | 否 | 合并更新自定义元数据。 |
| [IdempotencyKey](../forward/template.go#L181)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |
| [QoderBeta](../forward/template.go#L183)：`param.Opt[string]` | header / `X-Qoder-Beta` | 否 | 更新后的 `tools` 中包含 Browser Use 工具集时，必须设置为 `browser-use-2026-07-14`。 |

### `Templates.Archive`

```go
func (r *TemplateService) Archive(ctx context.Context, templateID string, params TemplateArchiveParams, opts ...option.RequestOption) (res *Template, err error)
```

**POST `/templates/{template_id}/archive`** · [实现](../forward/template.go#L196) · 归档资源。

路径参数：`templateID string`，必填。

返回：`*Template, error`。[Template](../forward/template.go#L249)。

请求对象 `params`：[TemplateArchiveParams](../forward/template.go#L209)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [IdempotencyKey](../forward/template.go#L211)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Templates.Clone`

```go
func (r *TemplateService) Clone(ctx context.Context, templateID string, params TemplateCloneParams, opts ...option.RequestOption) (res *Template, err error)
```

**POST `/templates/{template_id}/clone`** · [实现](../forward/template.go#L220) · 克隆资源。

路径参数：`templateID string`，必填。

返回：`*Template, error`。[Template](../forward/template.go#L249)。

请求对象 `params`：[TemplateCloneParams](../forward/template.go#L233)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Name](../forward/template.go#L235)：`param.Opt[string]` | json / `name` | 否 | 新 Template 名称；不传时使用 `<源名称> Copy <随机短 ID>`。 |
| [Description](../forward/template.go#L237)：`param.Opt[string]` | json / `description` | 否 | 新 Template 描述；不传时沿用源描述。 |
| [IdempotencyKey](../forward/template.go#L239)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

## Identities

管理业务身份及其状态；Configs 管理 Identity 对 Template 的配置覆盖，MemoryStores 管理这对身份与模板的记忆挂载。

### `Identities.List`

```go
func (r *IdentityService) List(ctx context.Context, params IdentityListParams, opts ...option.RequestOption) (res *pagination.Page[Identity], err error)
```

**GET `/identities`** · [实现](../forward/identity.go#L32) · 列出资源。

返回：`*pagination.Page[Identity], error`。[Identity](../forward/identity.go#L346)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[IdentityListParams](../forward/identity.go#L52)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [ExternalID](../forward/identity.go#L54)：`param.Opt[string]` | query / `external_id` | 否 | 按集成方终端用户 ID 过滤。 |
| [IdentityIDs](../forward/identity.go#L56)：`[]string` | query / `identity_ids` | 否 | 按多个 Identity ID 过滤；支持逗号分隔或重复 query 参数，去重后最多 100 个。 |
| [Search](../forward/identity.go#L58)：`param.Opt[string]` | query / `search` | 否 | 匹配 Identity ID、名称或外部 ID。 |
| [Enabled](../forward/identity.go#L60)：`param.Opt[bool]` | query / `enabled` | 否 | 按是否启用过滤；非布尔值返回 400。 |
| [Limit](../forward/identity.go#L62)：`param.Opt[int64]` | query / `limit` | 否 | 分页大小，最大 100；超过上限时按最大值处理。 |
| [AfterID](../forward/identity.go#L64)：`param.Opt[string]` | query / `after_id` | 否 | 向后翻页游标，不能与 `before_id` 同用。 |
| [BeforeID](../forward/identity.go#L66)：`param.Opt[string]` | query / `before_id` | 否 | 向前翻页游标，不能与 `after_id` 同用。 |

### `Identities.ListAutoPaging`

```go
func (r *IdentityService) ListAutoPaging(ctx context.Context, params IdentityListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Identity]
```

[实现](../forward/identity.go#L48)。自动遍历 `Identities.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Identities.New`

```go
func (r *IdentityService) New(ctx context.Context, params IdentityNewParams, opts ...option.RequestOption) (res *Identity, err error)
```

**POST `/identities`** · [实现](../forward/identity.go#L75) · 创建资源。

返回：`*Identity, error`。[Identity](../forward/identity.go#L346)。

请求对象 `params`：[IdentityNewParams](../forward/identity.go#L85)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [ExternalID](../forward/identity.go#L87)：`string` | json / `external_id` | 是 | 集成方系统中的终端用户 ID，不能是空串或纯空白。 |
| [Name](../forward/identity.go#L89)：`param.Opt[string]` | json / `name` | 否 | 展示名，传入时不能是空串或纯空白。 |
| [Enabled](../forward/identity.go#L91)：`param.Opt[bool]` | json / `enabled` | 否 | 是否启用该 Identity，默认 `true`。 |
| [Metadata](../forward/identity.go#L93)：`map[string]any` | json / `metadata` | 否 | 业务元数据，建议最多 16 个 key。 |
| [IdempotencyKey](../forward/identity.go#L95)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Identities.EnsureAdmin`

```go
func (r *IdentityService) EnsureAdmin(ctx context.Context, opts ...option.RequestOption) (res *Identity, err error)
```

**POST `/identities/admin/ensure`** · [实现](../forward/identity.go#L106) · 确保管理员身份存在。

返回：`*Identity, error`。[Identity](../forward/identity.go#L346)。

无请求参数对象；可直接传入请求选项。

### `Identities.Stats`

```go
func (r *IdentityService) Stats(ctx context.Context, opts ...option.RequestOption) (res *IdentityStats, err error)
```

**GET `/identities/stats`** · [实现](../forward/identity.go#L115) · 查询统计。

返回：`*IdentityStats, error`。[IdentityStats](../forward/identity.go#L324)。

无请求参数对象；可直接传入请求选项。

### `Identities.Get`

```go
func (r *IdentityService) Get(ctx context.Context, identityID string, opts ...option.RequestOption) (res *Identity, err error)
```

**GET `/identities/{identity_id}`** · [实现](../forward/identity.go#L124) · 获取资源详情。

路径参数：`identityID string`，必填。

返回：`*Identity, error`。[Identity](../forward/identity.go#L346)。

无请求参数对象；可直接传入请求选项。

### `Identities.Update`

```go
func (r *IdentityService) Update(ctx context.Context, identityID string, params IdentityUpdateParams, opts ...option.RequestOption) (res *Identity, err error)
```

**POST `/identities/{identity_id}`** · [实现](../forward/identity.go#L136) · 更新资源。

路径参数：`identityID string`，必填。

返回：`*Identity, error`。[Identity](../forward/identity.go#L346)。

请求对象 `params`：[IdentityUpdateParams](../forward/identity.go#L149)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [ExternalID](../forward/identity.go#L151)：`param.Opt[string]` | json / `external_id` | 否 | 替换原有终端用户 ID。 |
| [Name](../forward/identity.go#L153)：`param.Opt[string]` | json / `name` | 否 | 替换展示名。 |
| [Enabled](../forward/identity.go#L155)：`param.Opt[bool]` | json / `enabled` | 否 | 更新 Identity 是否可用。 |
| [Metadata](../forward/identity.go#L157)：`map[string]any` | json / `metadata` | 否 | 合并更新业务元数据；空字符串 value 删除对应 key。 |
| [IdempotencyKey](../forward/identity.go#L159)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Identities.Delete`

```go
func (r *IdentityService) Delete(ctx context.Context, identityID string, opts ...option.RequestOption) (res *DeletedIdentity, err error)
```

**DELETE `/identities/{identity_id}`** · [实现](../forward/identity.go#L172) · 删除资源。

路径参数：`identityID string`，必填。

返回：`*DeletedIdentity, error`。[DeletedIdentity](../forward/identity.go#L308)。

无请求参数对象；可直接传入请求选项。

### `Identities.ListTemplates`

```go
func (r *IdentityService) ListTemplates(ctx context.Context, identityID string, opts ...option.RequestOption) (res *IdentityListTemplatesResponse, err error)
```

**GET `/identities/{identity_id}/agents`** · [实现](../forward/identity.go#L184) · 列出身份可用模板。

路径参数：`identityID string`，必填。

返回：`*IdentityListTemplatesResponse, error`。[IdentityListTemplatesResponse](../forward/identity.go#L380)。

无请求参数对象；可直接传入请求选项。

### `Identities.Clear`

```go
func (r *IdentityService) Clear(ctx context.Context, identityID string, params IdentityClearParams, opts ...option.RequestOption) (res *IdentityClearResponse, err error)
```

**POST `/identities/{identity_id}/clear`** · [实现](../forward/identity.go#L196) · 清理身份关联资源。

路径参数：`identityID string`，必填。

返回：`*IdentityClearResponse, error`。[IdentityClearResponse](../forward/identity.go#L243)。

请求对象 `params`：[IdentityClearParams](../forward/identity.go#L207)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Reason](../forward/identity.go#L209)：`param.Opt[string]` | json / `reason` | 否 | 清理原因，仅用于记录调用意图。 |

### `Identities.Disable`

```go
func (r *IdentityService) Disable(ctx context.Context, identityID string, opts ...option.RequestOption) (res *Identity, err error)
```

**POST `/identities/{identity_id}/disable`** · [实现](../forward/identity.go#L220) · 停用身份。

路径参数：`identityID string`，必填。

返回：`*Identity, error`。[Identity](../forward/identity.go#L346)。

无请求参数对象；可直接传入请求选项。

### `Identities.Enable`

```go
func (r *IdentityService) Enable(ctx context.Context, identityID string, opts ...option.RequestOption) (res *Identity, err error)
```

**POST `/identities/{identity_id}/enable`** · [实现](../forward/identity.go#L232) · 启用身份。

路径参数：`identityID string`，必填。

返回：`*Identity, error`。[Identity](../forward/identity.go#L346)。

无请求参数对象；可直接传入请求选项。

### `Identities.Configs.List`

```go
func (r *IdentityConfigService) List(ctx context.Context, identityID string, params IdentityConfigListParams, opts ...option.RequestOption) (res *pagination.Page[IdentityConfig], err error)
```

**GET `/identities/{identity_id}/templates`** · [实现](../forward/identityconfig.go#L30) · 列出资源。

路径参数：`identityID string`，必填。

返回：`*pagination.Page[IdentityConfig], error`。[IdentityConfig](../forward/identityconfig.go#L266)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[IdentityConfigListParams](../forward/identityconfig.go#L53)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [TemplateID](../forward/identityconfig.go#L55)：`param.Opt[string]` | query / `template_id` | 否 | 按 Forward Template ID 过滤。 |
| [Status](../forward/identityconfig.go#L57)：`param.Opt[string]` | query / `status` | 否 | 按 `active` 或 `archived` 过滤。 |
| [Limit](../forward/identityconfig.go#L59)：`param.Opt[int64]` | query / `limit` | 否 | 分页大小，最大 100。 |
| [AfterID](../forward/identityconfig.go#L61)：`param.Opt[string]` | query / `after_id` | 否 | 来自上一页响应 `last_id` 的向后游标。 |
| [BeforeID](../forward/identityconfig.go#L63)：`param.Opt[string]` | query / `before_id` | 否 | 来自上一页响应 `first_id` 的向前游标。 |

### `Identities.Configs.ListAutoPaging`

```go
func (r *IdentityConfigService) ListAutoPaging(ctx context.Context, identityID string, params IdentityConfigListParams, opts ...option.RequestOption) *pagination.PageAutoPager[IdentityConfig]
```

[实现](../forward/identityconfig.go#L49)。自动遍历 `Identities.Configs.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Identities.Configs.Get`

```go
func (r *IdentityConfigService) Get(ctx context.Context, identityID string, templateID string, opts ...option.RequestOption) (res *IdentityConfig, err error)
```

**GET `/identities/{identity_id}/templates/{template_id}/config`** · [实现](../forward/identityconfig.go#L72) · 获取资源详情。

路径参数：`identityID string`，必填。

路径参数：`templateID string`，必填。

返回：`*IdentityConfig, error`。[IdentityConfig](../forward/identityconfig.go#L266)。

无请求参数对象；可直接传入请求选项。

### `Identities.Configs.Upsert`

```go
func (r *IdentityConfigService) Upsert(ctx context.Context, identityID string, templateID string, params IdentityConfigUpsertParams, opts ...option.RequestOption) (res *IdentityConfig, err error)
```

**POST `/identities/{identity_id}/templates/{template_id}/config`** · [实现](../forward/identityconfig.go#L87) · 创建或更新配置。

路径参数：`identityID string`，必填。

路径参数：`templateID string`，必填。

返回：`*IdentityConfig, error`。[IdentityConfig](../forward/identityconfig.go#L266)。

请求对象 `params`：[IdentityConfigUpsertParams](../forward/identityconfig.go#L103)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Name](../forward/identityconfig.go#L105)：`param.Opt[string]` | json / `name` | 否 | Config 展示名。 |
| [IdentityConfig](../forward/identityconfig.go#L107)：`IdentityConfigSpecParam` | json / `identity_config` | 是 | 用户级覆盖配置。；类型：[IdentityConfigSpecParam](../forward/configuration_types.go#L556) |
| [Metadata](../forward/identityconfig.go#L109)：`map[string]any` | json / `metadata` | 否 | 业务元数据；传入时整体替换已有 metadata。 |
| [IdempotencyKey](../forward/identityconfig.go#L111)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Identities.Configs.GetEffective`

```go
func (r *IdentityConfigService) GetEffective(ctx context.Context, identityID string, templateID string, opts ...option.RequestOption) (res *EffectiveConfig, err error)
```

**GET `/identities/{identity_id}/templates/{template_id}/effective`** · [实现](../forward/identityconfig.go#L124) · 获取合并后的有效配置。

路径参数：`identityID string`，必填。

路径参数：`templateID string`，必填。

返回：`*EffectiveConfig, error`。[EffectiveConfig](../forward/identityconfig.go#L138)。

无请求参数对象；可直接传入请求选项。

### `Identities.MemoryStores.List`

```go
func (r *IdentityMemoryStoreService) List(ctx context.Context, identityID string, templateID string, opts ...option.RequestOption) (res *IdentityMemoryStoreListResponse, err error)
```

**GET `/identities/{identity_id}/templates/{template_id}/memory_stores`** · [实现](../forward/identitymemorystore.go#L28) · 列出资源。

路径参数：`identityID string`，必填。

路径参数：`templateID string`，必填。

返回：`*IdentityMemoryStoreListResponse, error`。[IdentityMemoryStoreListResponse](../forward/identitymemorystore.go#L138)。

无请求参数对象；可直接传入请求选项。

### `Identities.MemoryStores.Mount`

```go
func (r *IdentityMemoryStoreService) Mount(ctx context.Context, identityID string, templateID string, params IdentityMemoryStoreMountParams, opts ...option.RequestOption) (res *MemoryStoreMount, err error)
```

**POST `/identities/{identity_id}/templates/{template_id}/memory_stores`** · [实现](../forward/identitymemorystore.go#L43) · 挂载记忆存储。

路径参数：`identityID string`，必填。

路径参数：`templateID string`，必填。

返回：`*MemoryStoreMount, error`。[MemoryStoreMount](../forward/identitymemorystore.go#L110)。

请求对象 `params`：[IdentityMemoryStoreMountParams](../forward/identitymemorystore.go#L57)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [MemoryStoreID](../forward/identitymemorystore.go#L59)：`string` | json / `memory_store_id` | 是 | 要挂载的 Memory Store ID（`memstore_...`）。必须是当前调用方可见的 active Store。 |

### `Identities.MemoryStores.Detach`

```go
func (r *IdentityMemoryStoreService) Detach(ctx context.Context, identityID string, templateID string, memoryStoreID string, opts ...option.RequestOption) (res *DeletedMemoryStoreMount, err error)
```

**DELETE `/identities/{identity_id}/templates/{template_id}/memory_stores/{memory_store_id}`** · [实现](../forward/identitymemorystore.go#L72) · 解除记忆存储挂载。

路径参数：`identityID string`，必填。

路径参数：`templateID string`，必填。

路径参数：`memoryStoreID string`，必填。

返回：`*DeletedMemoryStoreMount, error`。[DeletedMemoryStoreMount](../forward/identitymemorystore.go#L89)。

无请求参数对象；可直接传入请求选项。

## Sessions

管理会话生命周期、消息事件、资源和子线程。发送消息或创建会话成功表示已接收，执行结果应通过事件确认。

### `Sessions.List`

```go
func (r *SessionService) List(ctx context.Context, params SessionListParams, opts ...option.RequestOption) (res *pagination.Page[Session], err error)
```

**GET `/sessions`** · [实现](../forward/session.go#L33) · 列出资源。

返回：`*pagination.Page[Session], error`。[Session](../forward/session.go#L243)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[SessionListParams](../forward/session.go#L53)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [IdentityIDs](../forward/session.go#L55)：`[]string` | query / `identity_ids` | 否 | 按一个或多个 Identity ID 过滤，支持逗号分隔。 |
| [TemplateID](../forward/session.go#L57)：`param.Opt[string]` | query / `template_id` | 否 | 按 Forward Template ID 过滤。 |
| [SourceType](../forward/session.go#L59)：`param.Opt[string]` | query / `source_type` | 否 | 按 `api`、`im`、`schedule` 或 `batch` 过滤。 |
| [CreatedAtGt](../forward/session.go#L61)：`param.Opt[time.Time]` | query / `created_at[gt]` | 否 | 创建时间严格大于该 RFC 3339 时间。 |
| [CreatedAtGte](../forward/session.go#L63)：`param.Opt[time.Time]` | query / `created_at[gte]` | 否 | 创建时间大于等于该 RFC 3339 时间。 |
| [CreatedAtLt](../forward/session.go#L65)：`param.Opt[time.Time]` | query / `created_at[lt]` | 否 | 创建时间严格小于该 RFC 3339 时间。 |
| [CreatedAtLte](../forward/session.go#L67)：`param.Opt[time.Time]` | query / `created_at[lte]` | 否 | 创建时间小于等于该 RFC 3339 时间。 |
| [UpdatedAtGt](../forward/session.go#L69)：`param.Opt[time.Time]` | query / `updated_at[gt]` | 否 | 更新时间严格大于该 RFC 3339 时间。 |
| [UpdatedAtGte](../forward/session.go#L71)：`param.Opt[time.Time]` | query / `updated_at[gte]` | 否 | 更新时间大于等于该 RFC 3339 时间。 |
| [UpdatedAtLt](../forward/session.go#L73)：`param.Opt[time.Time]` | query / `updated_at[lt]` | 否 | 更新时间严格小于该 RFC 3339 时间。 |
| [UpdatedAtLte](../forward/session.go#L75)：`param.Opt[time.Time]` | query / `updated_at[lte]` | 否 | 更新时间小于等于该 RFC 3339 时间。 |
| [Limit](../forward/session.go#L77)：`param.Opt[int64]` | query / `limit` | 否 | 分页大小，最大 100。 |
| [AfterID](../forward/session.go#L79)：`param.Opt[string]` | query / `after_id` | 否 | 向后翻页游标，传入上一页响应的 `last_id`。 |
| [BeforeID](../forward/session.go#L81)：`param.Opt[string]` | query / `before_id` | 否 | 向前翻页游标，传入当前页响应的 `first_id`。 |
| [Order](../forward/session.go#L83)：`param.Opt[string]` | query / `order` | 否 | 创建时间排序方向：`desc` 或 `asc`。 |
| [IncludeArchived](../forward/session.go#L85)：`param.Opt[bool]` | query / `include_archived` | 否 | 是否包含已归档 Session。 |

### `Sessions.ListAutoPaging`

```go
func (r *SessionService) ListAutoPaging(ctx context.Context, params SessionListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Session]
```

[实现](../forward/session.go#L49)。自动遍历 `Sessions.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Sessions.New`

```go
func (r *SessionService) New(ctx context.Context, params SessionNewParams, opts ...option.RequestOption) (res *Session, err error)
```

**POST `/sessions`** · [实现](../forward/session.go#L94) · 创建资源。

返回：`*Session, error`。[Session](../forward/session.go#L243)。

请求对象 `params`：[SessionNewParams](../forward/session.go#L104)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [IdentityID](../forward/session.go#L106)：`string` | json / `identity_id` | 是 | Forward Identity ID。 |
| [TemplateID](../forward/session.go#L108)：`string` | json / `template_id` | 是 | Forward Template ID。 |
| [Title](../forward/session.go#L110)：`param.Opt[string]` | json / `title` | 否 | Session 标题。 |
| [Metadata](../forward/session.go#L112)：`map[string]any` | json / `metadata` | 否 | 业务元数据。 |
| [Config](../forward/session.go#L113)：`SessionNewParamsConfigParam` | json / `config` | 否 | 类型：[SessionNewParamsConfigParam](../forward/session.go#L126) |
| [Resources](../forward/session.go#L114)：`[]SessionResourceSpecParam` | json / `resources` | 否 | 类型：[SessionResourceSpecParam](../forward/configuration_types.go#L367) |
| [IdempotencyKey](../forward/session.go#L116)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Sessions.Get`

```go
func (r *SessionService) Get(ctx context.Context, sessionID string, opts ...option.RequestOption) (res *Session, err error)
```

**GET `/sessions/{session_id}`** · [实现](../forward/session.go#L140) · 获取资源详情。

路径参数：`sessionID string`，必填。

返回：`*Session, error`。[Session](../forward/session.go#L243)。

无请求参数对象；可直接传入请求选项。

### `Sessions.Update`

```go
func (r *SessionService) Update(ctx context.Context, sessionID string, params SessionUpdateParams, opts ...option.RequestOption) (res *Session, err error)
```

**POST `/sessions/{session_id}`** · [实现](../forward/session.go#L152) · 更新资源。

路径参数：`sessionID string`，必填。

返回：`*Session, error`。[Session](../forward/session.go#L243)。

请求对象 `params`：[SessionUpdateParams](../forward/session.go#L165)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Title](../forward/session.go#L167)：`param.Opt[string]` | json / `title` | 否 | 新的 Session 标题。 |
| [Metadata](../forward/session.go#L169)：`map[string]any` | json / `metadata` | 否 | metadata merge patch；传入的 key 覆盖已有 key，未出现的 key 保留。 |
| [Config](../forward/session.go#L170)：`SessionUpdateParamsConfigParam` | json / `config` | 否 | 类型：[SessionUpdateParamsConfigParam](../forward/session.go#L182) |
| [IdempotencyKey](../forward/session.go#L172)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Sessions.Archive`

```go
func (r *SessionService) Archive(ctx context.Context, sessionID string, params SessionArchiveParams, opts ...option.RequestOption) (res *Session, err error)
```

**POST `/sessions/{session_id}/archive`** · [实现](../forward/session.go#L196) · 归档资源。

路径参数：`sessionID string`，必填。

返回：`*Session, error`。[Session](../forward/session.go#L243)。

请求对象 `params`：[SessionArchiveParams](../forward/session.go#L209)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [IdempotencyKey](../forward/session.go#L211)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Sessions.Cancel`

```go
func (r *SessionService) Cancel(ctx context.Context, sessionID string, params SessionCancelParams, opts ...option.RequestOption) (res *Session, err error)
```

**POST `/sessions/{session_id}/cancel`** · [实现](../forward/session.go#L220) · 取消执行。

路径参数：`sessionID string`，必填。

返回：`*Session, error`。[Session](../forward/session.go#L243)。

请求对象 `params`：[SessionCancelParams](../forward/session.go#L233)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [IdempotencyKey](../forward/session.go#L235)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Sessions.Events.List`

```go
func (r *SessionEventService) List(ctx context.Context, sessionID string, params SessionEventListParams, opts ...option.RequestOption) (res *pagination.Page[SessionEvent], err error)
```

**GET `/sessions/{session_id}/events`** · [实现](../forward/sessionevent.go#L32) · 列出资源。

路径参数：`sessionID string`，必填。

返回：`*pagination.Page[SessionEvent], error`。[SessionEvent](../forward/sessionevent.go#L139)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[SessionEventListParams](../forward/sessionevent.go#L55)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Limit](../forward/sessionevent.go#L57)：`param.Opt[int64]` | query / `limit` | 否 | 分页大小，最大 100。 |
| [AfterID](../forward/sessionevent.go#L59)：`param.Opt[string]` | query / `after_id` | 否 | 返回该 Event ID 之后的事件。 |
| [BeforeID](../forward/sessionevent.go#L61)：`param.Opt[string]` | query / `before_id` | 否 | 返回该 Event ID 之前的事件。 |
| [Order](../forward/sessionevent.go#L63)：`param.Opt[string]` | query / `order` | 否 | 排序方向：`asc` 或 `desc`。 |
| [Type](../forward/sessionevent.go#L65)：`param.Opt[string]` | query / `type` | 否 | 按 Event 类型过滤，支持逗号分隔。 |
| [Types](../forward/sessionevent.go#L67)：`[]string` | query / `types[]` | 否 | 数组形式的 Event 类型过滤。 |
| [IncludeToolCalls](../forward/sessionevent.go#L69)：`param.Opt[bool]` | query / `include_tool_calls` | 否 | 是否包含工具调用类事件。 |
| [IncludeThinking](../forward/sessionevent.go#L71)：`param.Opt[bool]` | query / `include_thinking` | 否 | 是否包含思考过程事件。 |

### `Sessions.Events.ListAutoPaging`

```go
func (r *SessionEventService) ListAutoPaging(ctx context.Context, sessionID string, params SessionEventListParams, opts ...option.RequestOption) *pagination.PageAutoPager[SessionEvent]
```

[实现](../forward/sessionevent.go#L51)。自动遍历 `Sessions.Events.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Sessions.Events.Send`

```go
func (r *SessionEventService) Send(ctx context.Context, sessionID string, params SessionEventSendParams, opts ...option.RequestOption) (res *SessionEventSendResponse, err error)
```

**POST `/sessions/{session_id}/events`** · [实现](../forward/sessionevent.go#L80) · 发送会话事件。

路径参数：`sessionID string`，必填。

返回：`*SessionEventSendResponse, error`。[SessionEventSendResponse](../forward/sessionevent.go#L231)。

请求对象 `params`：[SessionEventSendParams](../forward/sessionevent.go#L93)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Events](../forward/sessionevent.go#L94)：`[]SessionEventParam` | json / `events` | 是 | 类型：[SessionEventParam](../forward/configuration_types.go#L579) |
| [IdempotencyKey](../forward/sessionevent.go#L96)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Sessions.Events.StreamEvents`

```go
func (r *SessionEventService) StreamEvents(ctx context.Context, sessionID string, params SessionEventStreamParams, opts ...option.RequestOption) *ssestream.Stream[SessionEvent]
```

**GET `/sessions/{session_id}/events/stream`** · [实现](../forward/sessionevent.go#L109) · 订阅 SSE 事件流。

路径参数：`sessionID string`，必填。

返回：`*ssestream.Stream[SessionEvent]`。[SessionEvent](../forward/sessionevent.go#L139)、[事件流](../convention/ssestream/ssestream.go)。

请求对象 `params`：[SessionEventStreamParams](../forward/sessionevent.go#L123)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [EventDeltas](../forward/sessionevent.go#L125)：`[]string` | query / `event_deltas[]` | 否 | 订阅指定公开事件类型的流式增量事件。支持重复传参，取值见 流式增量事件。 |
| [IncludeToolCalls](../forward/sessionevent.go#L127)：`param.Opt[bool]` | query / `include_tool_calls` | 否 | 是否包含工具调用类事件。 |
| [IncludeThinking](../forward/sessionevent.go#L129)：`param.Opt[bool]` | query / `include_thinking` | 否 | 是否包含思考过程事件。 |
| [LastEventID](../forward/sessionevent.go#L131)：`param.Opt[string]` | header / `Last-Event-ID` | 否 | 从该 Event ID 之后恢复订阅。 |

### `Sessions.Resources.Add`

```go
func (r *SessionResourceService) Add(ctx context.Context, sessionID string, params SessionResourceAddParams, opts ...option.RequestOption) (res *SessionResource, err error)
```

**POST `/sessions/{session_id}/resources`** · [实现](../forward/sessionresource.go#L28) · 添加会话资源。

路径参数：`sessionID string`，必填。

返回：`*SessionResource, error`。[SessionResource](../forward/sessionresource.go#L57)。

请求对象 `params`：[SessionResourceAddParams](../forward/sessionresource.go#L39)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Type](../forward/sessionresource.go#L41)：`string` | json / `type` | 是 | 资源类型，必须为 `file`。 |
| [FileID](../forward/sessionresource.go#L43)：`string` | json / `file_id` | 是 | Files API 返回的 File ID，文件必须已上传完成。 |
| [MountPath](../forward/sessionresource.go#L45)：`param.Opt[string]` | json / `mount_path` | 否 | Agent 容器内挂载路径；省略时由 Forward 根据文件名生成，默认挂载到 `/data/workspace/<文件名>`。 |

### `Sessions.Threads.List`

```go
func (r *SessionThreadService) List(ctx context.Context, sessionID string, params SessionThreadListParams, opts ...option.RequestOption) (res *pagination.Page[SessionThread], err error)
```

**GET `/sessions/{session_id}/threads`** · [实现](../forward/sessionthread.go#L31) · 列出资源。

路径参数：`sessionID string`，必填。

返回：`*pagination.Page[SessionThread], error`。[SessionThread](../forward/sessionthread.go#L110)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[SessionThreadListParams](../forward/sessionthread.go#L54)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Limit](../forward/sessionthread.go#L56)：`param.Opt[int64]` | query / `limit` | 否 | 分页大小，范围为 1–100。 |
| [AfterID](../forward/sessionthread.go#L58)：`param.Opt[string]` | query / `after_id` | 否 | 返回该 Thread ID 之后的记录。 |
| [BeforeID](../forward/sessionthread.go#L60)：`param.Opt[string]` | query / `before_id` | 否 | 返回该 Thread ID 之前的记录。 |

### `Sessions.Threads.ListAutoPaging`

```go
func (r *SessionThreadService) ListAutoPaging(ctx context.Context, sessionID string, params SessionThreadListParams, opts ...option.RequestOption) *pagination.PageAutoPager[SessionThread]
```

[实现](../forward/sessionthread.go#L50)。自动遍历 `Sessions.Threads.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Sessions.Threads.Get`

```go
func (r *SessionThreadService) Get(ctx context.Context, sessionID string, threadID string, opts ...option.RequestOption) (res *SessionThread, err error)
```

**GET `/sessions/{session_id}/threads/{thread_id}`** · [实现](../forward/sessionthread.go#L69) · 获取资源详情。

路径参数：`sessionID string`，必填。

路径参数：`threadID string`，必填。

返回：`*SessionThread, error`。[SessionThread](../forward/sessionthread.go#L110)。

无请求参数对象；可直接传入请求选项。

### `Sessions.Threads.Archive`

```go
func (r *SessionThreadService) Archive(ctx context.Context, sessionID string, threadID string, params SessionThreadArchiveParams, opts ...option.RequestOption) (res *SessionThread, err error)
```

**POST `/sessions/{session_id}/threads/{thread_id}/archive`** · [实现](../forward/sessionthread.go#L84) · 归档资源。

路径参数：`sessionID string`，必填。

路径参数：`threadID string`，必填。

返回：`*SessionThread, error`。[SessionThread](../forward/sessionthread.go#L110)。

请求对象 `params`：[SessionThreadArchiveParams](../forward/sessionthread.go#L100)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [IdempotencyKey](../forward/sessionthread.go#L102)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 标识一次逻辑归档尝试；建议为每次新的逻辑尝试生成唯一值。 |

### `Sessions.Threads.Events.List`

```go
func (r *SessionThreadEventService) List(ctx context.Context, sessionID string, threadID string, params SessionThreadEventListParams, opts ...option.RequestOption) (res *pagination.Page[SessionEvent], err error)
```

**GET `/sessions/{session_id}/threads/{thread_id}/events`** · [实现](../forward/sessionthreadevent.go#L28) · 列出资源。

路径参数：`sessionID string`，必填。

路径参数：`threadID string`，必填。

返回：`*pagination.Page[SessionEvent], error`。[SessionEvent](../forward/sessionevent.go#L139)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[SessionThreadEventListParams](../forward/sessionthreadevent.go#L54)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Limit](../forward/sessionthreadevent.go#L56)：`param.Opt[int64]` | query / `limit` | 否 | 分页大小，范围为 1–100。 |
| [AfterID](../forward/sessionthreadevent.go#L58)：`param.Opt[string]` | query / `after_id` | 否 | 返回该 Event ID 之后的记录。 |
| [BeforeID](../forward/sessionthreadevent.go#L60)：`param.Opt[string]` | query / `before_id` | 否 | 返回该 Event ID 之前的记录。 |

### `Sessions.Threads.Events.ListAutoPaging`

```go
func (r *SessionThreadEventService) ListAutoPaging(ctx context.Context, sessionID string, threadID string, params SessionThreadEventListParams, opts ...option.RequestOption) *pagination.PageAutoPager[SessionEvent]
```

[实现](../forward/sessionthreadevent.go#L50)。自动遍历 `Sessions.Threads.Events.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Sessions.Threads.Events.StreamEvents`

```go
func (r *SessionThreadEventService) StreamEvents(ctx context.Context, sessionID string, threadID string, params SessionThreadEventStreamParams, opts ...option.RequestOption) *ssestream.Stream[SessionEvent]
```

**GET `/sessions/{session_id}/threads/{thread_id}/stream`** · [实现](../forward/sessionthreadevent.go#L69) · 订阅 SSE 事件流。

路径参数：`sessionID string`，必填。

路径参数：`threadID string`，必填。

返回：`*ssestream.Stream[SessionEvent]`。[SessionEvent](../forward/sessionevent.go#L139)、[事件流](../convention/ssestream/ssestream.go)。

请求对象 `params`：[SessionThreadEventStreamParams](../forward/sessionthreadevent.go#L86)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [LastEventID](../forward/sessionthreadevent.go#L88)：`param.Opt[string]` | header / `Last-Event-ID` | 否 | 从该 Thread Event 之后继续订阅。 |

## Schedules

管理定时任务及手动触发；执行实例通过顶层 ScheduleRuns 查询。

### `Schedules.List`

```go
func (r *ScheduleService) List(ctx context.Context, params ScheduleListParams, opts ...option.RequestOption) (res *pagination.Page[Schedule], err error)
```

**GET `/schedules`** · [实现](../forward/schedule.go#L30) · 列出资源。

返回：`*pagination.Page[Schedule], error`。[Schedule](../forward/schedule.go#L317)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[ScheduleListParams](../forward/schedule.go#L50)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [IdentityID](../forward/schedule.go#L52)：`param.Opt[string]` | query / `identity_id` | 否 | PAT 或管理员 SAT 可省略，省略时查询当前 owner 全部 Identity；Identity-bound SAT 省略时自动绑定自身，显式传其他 Identity 返回 403。 |
| [TemplateID](../forward/schedule.go#L54)：`param.Opt[string]` | query / `template_id` | 否 | 按 Forward Template ID 过滤。 |
| [Status](../forward/schedule.go#L56)：`param.Opt[string]` | query / `status` | 否 | 按 `active` 或 `paused` 过滤。 |
| [IncludeArchived](../forward/schedule.go#L58)：`param.Opt[bool]` | query / `include_archived` | 否 | 是否包含已归档 Schedule。 |
| [Limit](../forward/schedule.go#L60)：`param.Opt[int64]` | query / `limit` | 否 | 分页大小，最大 100。 |
| [AfterID](../forward/schedule.go#L62)：`param.Opt[string]` | query / `after_id` | 否 | 向后翻页游标。 |
| [BeforeID](../forward/schedule.go#L64)：`param.Opt[string]` | query / `before_id` | 否 | 向前翻页游标。 |
| [SortBy](../forward/schedule.go#L66)：`param.Opt[string]` | query / `sort_by` | 否 | 排序字段：`created_at` 或 `upcoming_runs_at`。 |
| [Order](../forward/schedule.go#L68)：`param.Opt[string]` | query / `order` | 否 | 排序方向：`asc` 或 `desc`。 |

### `Schedules.ListAutoPaging`

```go
func (r *ScheduleService) ListAutoPaging(ctx context.Context, params ScheduleListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Schedule]
```

[实现](../forward/schedule.go#L46)。自动遍历 `Schedules.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Schedules.New`

```go
func (r *ScheduleService) New(ctx context.Context, params ScheduleNewParams, opts ...option.RequestOption) (res *Schedule, err error)
```

**POST `/schedules`** · [实现](../forward/schedule.go#L77) · 创建资源。

返回：`*Schedule, error`。[Schedule](../forward/schedule.go#L317)。

请求对象 `params`：[ScheduleNewParams](../forward/schedule.go#L87)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [IdentityID](../forward/schedule.go#L89)：`string` | json / `identity_id` | 是 | Schedule 所属 Forward Identity ID。 |
| [TemplateID](../forward/schedule.go#L91)：`string` | json / `template_id` | 是 | 要执行的 Forward Template ID。 |
| [Name](../forward/schedule.go#L93)：`string` | json / `name` | 是 | Schedule 名称。 |
| [Description](../forward/schedule.go#L95)：`param.Opt[string]` | json / `description` | 否 | Schedule 描述。 |
| [InitialEvents](../forward/schedule.go#L97)：`[]map[string]any` | json / `initial_events` | 是 | 每次执行注入的初始事件，当前支持 `user.message`。 |
| [Execution](../forward/schedule.go#L99)：`map[string]any` | json / `execution` | 否 | 执行策略；省略时使用服务端默认值。 |
| [TriggerPolicy](../forward/schedule.go#L101)：`map[string]any` | json / `trigger_policy` | 否 | 触发策略；省略或 `null` 时按 `manual` 处理。；允许 null。 |
| [EnvironmentID](../forward/schedule.go#L103)：`string` | json / `environment_id` | 是 | 执行环境。 |
| [Sinks](../forward/schedule.go#L105)：`[]map[string]any` | json / `sinks` | 否 | 执行结果推送目标；为兼容性保留数组形式，当前最多允许一个元素。；允许 null。 |
| [Metadata](../forward/schedule.go#L107)：`map[string]any` | json / `metadata` | 否 | 业务元数据，仅用于标签或透传。 |
| [IdempotencyKey](../forward/schedule.go#L109)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Schedules.ArchiveMany`

```go
func (r *ScheduleService) ArchiveMany(ctx context.Context, params ScheduleArchiveManyParams, opts ...option.RequestOption) (res *ScheduleArchiveManyResponse, err error)
```

**POST `/schedules/archive`** · [实现](../forward/schedule.go#L120) · 批量归档资源。

返回：`*ScheduleArchiveManyResponse, error`。[ScheduleArchiveManyResponse](../forward/schedule.go#L302)。

请求对象 `params`：[ScheduleArchiveManyParams](../forward/schedule.go#L130)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Scope](../forward/schedule.go#L132)：`string` | json / `scope` | 是 | 归档范围，当前仅支持 by_schedule_ids。 |
| [ScheduleIDs](../forward/schedule.go#L134)：`[]string` | json / `schedule_ids` | 是 | 去重后必须包含 1～50 个非空 Schedule ID。 |
| [IdempotencyKey](../forward/schedule.go#L136)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键；相同 owner、路径和请求体可安全重放。 |

### `Schedules.Get`

```go
func (r *ScheduleService) Get(ctx context.Context, scheduleID string, opts ...option.RequestOption) (res *Schedule, err error)
```

**GET `/schedules/{schedule_id}`** · [实现](../forward/schedule.go#L149) · 获取资源详情。

路径参数：`scheduleID string`，必填。

返回：`*Schedule, error`。[Schedule](../forward/schedule.go#L317)。

无请求参数对象；可直接传入请求选项。

### `Schedules.Update`

```go
func (r *ScheduleService) Update(ctx context.Context, scheduleID string, params ScheduleUpdateParams, opts ...option.RequestOption) (res *Schedule, err error)
```

**POST `/schedules/{schedule_id}`** · [实现](../forward/schedule.go#L161) · 更新资源。

路径参数：`scheduleID string`，必填。

返回：`*Schedule, error`。[Schedule](../forward/schedule.go#L317)。

请求对象 `params`：[ScheduleUpdateParams](../forward/schedule.go#L174)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Name](../forward/schedule.go#L176)：`param.Opt[string]` | json / `name` | 否 | 新的 Schedule 名称。 |
| [Description](../forward/schedule.go#L178)：`param.Opt[string]` | json / `description` | 否 | 新的 Schedule 描述。 |
| [TemplateID](../forward/schedule.go#L180)：`param.Opt[string]` | json / `template_id` | 否 | 新的 Forward Template ID。 |
| [InitialEvents](../forward/schedule.go#L182)：`[]map[string]any` | json / `initial_events` | 否 | 替换初始事件列表。 |
| [Execution](../forward/schedule.go#L184)：`map[string]any` | json / `execution` | 否 | 合并更新执行策略。 |
| [TriggerPolicy](../forward/schedule.go#L186)：`map[string]any` | json / `trigger_policy` | 否 | 更新触发策略；`null` 表示改为 manual。；允许 null。 |
| [EnvironmentID](../forward/schedule.go#L188)：`param.Opt[string]` | json / `environment_id` | 否 | 新的执行环境。 |
| [Sinks](../forward/schedule.go#L190)：`[]map[string]any` | json / `sinks` | 否 | 执行结果推送目标；为兼容性保留数组形式，当前最多允许一个元素。；允许 null。 |
| [Metadata](../forward/schedule.go#L192)：`map[string]any` | json / `metadata` | 否 | 合并更新 metadata；value 为 `null` 删除 key。 |
| [IdempotencyKey](../forward/schedule.go#L194)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Schedules.Archive`

```go
func (r *ScheduleService) Archive(ctx context.Context, scheduleID string, params ScheduleArchiveParams, opts ...option.RequestOption) (res *Schedule, err error)
```

**POST `/schedules/{schedule_id}/archive`** · [实现](../forward/schedule.go#L207) · 归档资源。

路径参数：`scheduleID string`，必填。

返回：`*Schedule, error`。[Schedule](../forward/schedule.go#L317)。

请求对象 `params`：[ScheduleArchiveParams](../forward/schedule.go#L220)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [IdempotencyKey](../forward/schedule.go#L222)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Schedules.Pause`

```go
func (r *ScheduleService) Pause(ctx context.Context, scheduleID string, params SchedulePauseParams, opts ...option.RequestOption) (res *Schedule, err error)
```

**POST `/schedules/{schedule_id}/pause`** · [实现](../forward/schedule.go#L231) · 暂停后续执行。

路径参数：`scheduleID string`，必填。

返回：`*Schedule, error`。[Schedule](../forward/schedule.go#L317)。

请求对象 `params`：[SchedulePauseParams](../forward/schedule.go#L244)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [IdempotencyKey](../forward/schedule.go#L246)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Schedules.Run`

```go
func (r *ScheduleService) Run(ctx context.Context, scheduleID string, params ScheduleRunParams, opts ...option.RequestOption) (res *ScheduleRun, err error)
```

**POST `/schedules/{schedule_id}/run`** · [实现](../forward/schedule.go#L255) · 立即触发一次执行。

路径参数：`scheduleID string`，必填。

返回：`*ScheduleRun, error`。[ScheduleRun](../forward/schedulerun.go#L100)。

请求对象 `params`：[ScheduleRunParams](../forward/schedule.go#L268)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [IdempotencyKey](../forward/schedule.go#L270)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Schedules.Unpause`

```go
func (r *ScheduleService) Unpause(ctx context.Context, scheduleID string, params ScheduleUnpauseParams, opts ...option.RequestOption) (res *Schedule, err error)
```

**POST `/schedules/{schedule_id}/unpause`** · [实现](../forward/schedule.go#L279) · 恢复执行。

路径参数：`scheduleID string`，必填。

返回：`*Schedule, error`。[Schedule](../forward/schedule.go#L317)。

请求对象 `params`：[ScheduleUnpauseParams](../forward/schedule.go#L292)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [IdempotencyKey](../forward/schedule.go#L294)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

## ScheduleRuns

查询定时任务执行实例与关联会话。

### `ScheduleRuns.List`

```go
func (r *ScheduleRunService) List(ctx context.Context, params ScheduleRunListParams, opts ...option.RequestOption) (res *pagination.Page[ScheduleRun], err error)
```

**GET `/schedule_runs`** · [实现](../forward/schedulerun.go#L30) · 列出资源。

返回：`*pagination.Page[ScheduleRun], error`。[ScheduleRun](../forward/schedulerun.go#L100)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[ScheduleRunListParams](../forward/schedulerun.go#L50)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [IdentityID](../forward/schedulerun.go#L52)：`string` | query / `identity_id` | 是 | Run 所属 Forward Identity ID。 |
| [ScheduleID](../forward/schedulerun.go#L54)：`param.Opt[string]` | query / `schedule_id` | 否 | 按 Schedule ID 过滤。 |
| [Status](../forward/schedulerun.go#L56)：`param.Opt[string]` | query / `status` | 否 | 按 `pending`、`running`、`completed`、`failed` 或 `skipped` 过滤。 |
| [TriggerType](../forward/schedulerun.go#L58)：`param.Opt[string]` | query / `trigger_type` | 否 | 按 `schedule` 或 `manual` 过滤。 |
| [HasError](../forward/schedulerun.go#L60)：`param.Opt[bool]` | query / `has_error` | 否 | 是否只返回有错误或无错误的 Run。 |
| [Limit](../forward/schedulerun.go#L62)：`param.Opt[int64]` | query / `limit` | 否 | 分页大小，最大 100。 |
| [AfterID](../forward/schedulerun.go#L64)：`param.Opt[string]` | query / `after_id` | 否 | 向后翻页游标。 |
| [BeforeID](../forward/schedulerun.go#L66)：`param.Opt[string]` | query / `before_id` | 否 | 向前翻页游标。 |
| [SortBy](../forward/schedulerun.go#L68)：`param.Opt[string]` | query / `sort_by` | 否 | 排序字段：`created_at` 或 `triggered_at`。 |
| [Order](../forward/schedulerun.go#L70)：`param.Opt[string]` | query / `order` | 否 | 排序方向：`asc` 或 `desc`。 |

### `ScheduleRuns.ListAutoPaging`

```go
func (r *ScheduleRunService) ListAutoPaging(ctx context.Context, params ScheduleRunListParams, opts ...option.RequestOption) *pagination.PageAutoPager[ScheduleRun]
```

[实现](../forward/schedulerun.go#L46)。自动遍历 `ScheduleRuns.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `ScheduleRuns.Get`

```go
func (r *ScheduleRunService) Get(ctx context.Context, runID string, params ScheduleRunGetParams, opts ...option.RequestOption) (res *ScheduleRun, err error)
```

**GET `/schedule_runs/{run_id}`** · [实现](../forward/schedulerun.go#L79) · 获取资源详情。

路径参数：`runID string`，必填。

返回：`*ScheduleRun, error`。[ScheduleRun](../forward/schedulerun.go#L100)。

请求对象 `params`：[ScheduleRunGetParams](../forward/schedulerun.go#L90)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [IdentityID](../forward/schedulerun.go#L92)：`param.Opt[string]` | query / `identity_id` | 否 | 额外归属约束。 |

## Batches

提交 JSONL 批处理并查询任务、输出、错误与取消状态。创建后的 queued 状态不表示执行完成。

### `Batches.List`

```go
func (r *BatchService) List(ctx context.Context, params BatchListParams, opts ...option.RequestOption) (res *pagination.Page[Batch], err error)
```

**GET `/batches`** · [实现](../forward/batch.go#L31) · 列出资源。

返回：`*pagination.Page[Batch], error`。[Batch](../forward/batch.go#L172)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[BatchListParams](../forward/batch.go#L51)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Status](../forward/batch.go#L53)：`param.Opt[string]` | query / `status` | 否 | 按状态过滤。 |
| [Limit](../forward/batch.go#L55)：`param.Opt[int64]` | query / `limit` | 否 | 分页大小，最大 100。 |
| [AfterID](../forward/batch.go#L57)：`param.Opt[string]` | query / `after_id` | 否 | 向后翻页游标。 |
| [BeforeID](../forward/batch.go#L59)：`param.Opt[string]` | query / `before_id` | 否 | 向前翻页游标。 |

### `Batches.ListAutoPaging`

```go
func (r *BatchService) ListAutoPaging(ctx context.Context, params BatchListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Batch]
```

[实现](../forward/batch.go#L47)。自动遍历 `Batches.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Batches.New`

```go
func (r *BatchService) New(ctx context.Context, params BatchNewParams, opts ...option.RequestOption) (res *Batch, err error)
```

**POST `/batches`** · [实现](../forward/batch.go#L68) · 创建资源。

返回：`*Batch, error`。[Batch](../forward/batch.go#L172)。

请求对象 `params`：[BatchNewParams](../forward/batch.go#L78)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [InputFileID](../forward/batch.go#L80)：`string` | json / `input_file_id` | 是 | 通过 Files API 上传的 JSONL 文件 ID。 |
| [CompletionWindow](../forward/batch.go#L82)：`string` | json / `completion_window` | 是 | 完成窗口：`24h`、`48h`、`72h`。超时后 Batch 自动进入 `expired` 状态。 |
| [Metadata](../forward/batch.go#L84)：`map[string]any` | json / `metadata` | 否 | 调用方业务元数据，最多 16 个 key；value 可为任意 JSON 类型；整体序列化后 ≤ 2KB，key ≤ 64 字符，且不得包含 NUL（U+0000）。 |
| [IdempotencyKey](../forward/batch.go#L86)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Batches.Get`

```go
func (r *BatchService) Get(ctx context.Context, batchID string, opts ...option.RequestOption) (res *Batch, err error)
```

**GET `/batches/{batch_id}`** · [实现](../forward/batch.go#L97) · 获取资源详情。

路径参数：`batchID string`，必填。

返回：`*Batch, error`。[Batch](../forward/batch.go#L172)。

无请求参数对象；可直接传入请求选项。

### `Batches.Cancel`

```go
func (r *BatchService) Cancel(ctx context.Context, batchID string, params BatchCancelParams, opts ...option.RequestOption) (res *Batch, err error)
```

**POST `/batches/{batch_id}/cancel`** · [实现](../forward/batch.go#L109) · 取消执行。

路径参数：`batchID string`，必填。

返回：`*Batch, error`。[Batch](../forward/batch.go#L172)。

请求对象 `params`：[BatchCancelParams](../forward/batch.go#L122)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [IdempotencyKey](../forward/batch.go#L124)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Batches.GetError`

```go
func (r *BatchService) GetError(ctx context.Context, batchID string, opts ...option.RequestOption) (res *BatchFile, err error)
```

**GET `/batches/{batch_id}/error`** · [实现](../forward/batch.go#L133) · 获取批处理错误输出。

路径参数：`batchID string`，必填。

返回：`*BatchFile, error`。[BatchFile](../forward/batch.go#L156)。

无请求参数对象；可直接传入请求选项。

### `Batches.GetOutput`

```go
func (r *BatchService) GetOutput(ctx context.Context, batchID string, opts ...option.RequestOption) (res *BatchFile, err error)
```

**GET `/batches/{batch_id}/output`** · [实现](../forward/batch.go#L145) · 获取批处理输出。

路径参数：`batchID string`，必填。

返回：`*BatchFile, error`。[BatchFile](../forward/batch.go#L156)。

无请求参数对象；可直接传入请求选项。

### `Batches.Tasks.List`

```go
func (r *BatchTaskService) List(ctx context.Context, batchID string, params BatchTaskListParams, opts ...option.RequestOption) (res *pagination.Page[BatchTask], err error)
```

**GET `/batches/{batch_id}/tasks`** · [实现](../forward/batchtask.go#L30) · 列出资源。

路径参数：`batchID string`，必填。

返回：`*pagination.Page[BatchTask], error`。[BatchTask](../forward/batchtask.go#L69)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[BatchTaskListParams](../forward/batchtask.go#L53)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Status](../forward/batchtask.go#L55)：`param.Opt[string]` | query / `status` | 否 | 按任务状态过滤：`pending`、`running`、`completed`、`failed`、`cancelled`、`expired`。 |
| [CustomID](../forward/batchtask.go#L57)：`param.Opt[string]` | query / `custom_id` | 否 | 按调用方任务标识精确过滤，仅支持单值；未命中返回空列表。 |
| [Limit](../forward/batchtask.go#L59)：`param.Opt[int64]` | query / `limit` | 否 | 分页大小，最大 100。 |
| [AfterID](../forward/batchtask.go#L61)：`param.Opt[string]` | query / `after_id` | 否 | 向后翻页游标，传上一页响应的 `last_id`；游标必须属于当前 Batch。 |

### `Batches.Tasks.ListAutoPaging`

```go
func (r *BatchTaskService) ListAutoPaging(ctx context.Context, batchID string, params BatchTaskListParams, opts ...option.RequestOption) *pagination.PageAutoPager[BatchTask]
```

[实现](../forward/batchtask.go#L49)。自动遍历 `Batches.Tasks.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

## Channels

管理消息渠道及二维码配对会话。实际渠道连接与用户配对独立管理。

### `Channels.List`

```go
func (r *ChannelService) List(ctx context.Context, params ChannelListParams, opts ...option.RequestOption) (res *pagination.Page[Channel], err error)
```

**GET `/channels`** · [实现](../forward/channel.go#L31) · 列出资源。

返回：`*pagination.Page[Channel], error`。[Channel](../forward/channel.go#L185)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[ChannelListParams](../forward/channel.go#L51)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [ChannelType](../forward/channel.go#L53)：`param.Opt[string]` | query / `channel_type` | 否 | 按 `wechat`、`wecom`、`feishu`、`dingtalk` 或 `teams`（Global）过滤。 |
| [Enabled](../forward/channel.go#L55)：`param.Opt[bool]` | query / `enabled` | 否 | 按人工启停状态过滤。 |
| [BindingStatus](../forward/channel.go#L57)：`param.Opt[string]` | query / `binding_status` | 否 | 按 `unbound`、`bound` 或 `expired` 过滤。 |
| [IdentityID](../forward/channel.go#L59)：`param.Opt[string]` | query / `identity_id` | 否 | 按 Forward Identity ID 过滤。 |
| [TemplateID](../forward/channel.go#L61)：`param.Opt[string]` | query / `template_id` | 否 | 按 Forward Template ID 过滤。 |
| [Limit](../forward/channel.go#L63)：`param.Opt[int64]` | query / `limit` | 否 | 分页大小，最大 100。 |
| [AfterID](../forward/channel.go#L65)：`param.Opt[string]` | query / `after_id` | 否 | 向后翻页游标。 |
| [BeforeID](../forward/channel.go#L67)：`param.Opt[string]` | query / `before_id` | 否 | 向前翻页游标。 |

### `Channels.ListAutoPaging`

```go
func (r *ChannelService) ListAutoPaging(ctx context.Context, params ChannelListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Channel]
```

[实现](../forward/channel.go#L47)。自动遍历 `Channels.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Channels.New`

```go
func (r *ChannelService) New(ctx context.Context, params ChannelNewParams, opts ...option.RequestOption) (res *Channel, err error)
```

**POST `/channels`** · [实现](../forward/channel.go#L76) · 创建资源。

返回：`*Channel, error`。[Channel](../forward/channel.go#L185)。

请求对象 `params`：[ChannelNewParams](../forward/channel.go#L86)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [IdentityID](../forward/channel.go#L88)：`param.Opt[string]` | json / `identity_id` | 否 | `fixed` 模式必填；`pairing` 模式不传。 |
| [IdentityResolution](../forward/channel.go#L89)：`map[string]any` | json / `identity_resolution` | 否 | 见字段定义。 |
| [TemplateID](../forward/channel.go#L91)：`param.Opt[string]` | json / `template_id` | 否 | `fixed` 模式必填；`pairing` 模式不传。 |
| [ChannelType](../forward/channel.go#L93)：`string` | json / `channel_type` | 是 | 渠道类型，当前支持 `wechat`、`wecom`、`feishu`、`dingtalk` 和 `teams`（Global）。 |
| [Name](../forward/channel.go#L95)：`param.Opt[string]` | json / `name` | 否 | Channel 展示名。 |
| [Enabled](../forward/channel.go#L97)：`param.Opt[bool]` | json / `enabled` | 否 | 人工启停开关，默认 `true`。传 `false` 可创建后暂不处理上行消息。 |
| [ChannelConfig](../forward/channel.go#L98)：`map[string]any` | json / `channel_config` | 否 | 见字段定义。 |
| [IdempotencyKey](../forward/channel.go#L100)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Channels.Get`

```go
func (r *ChannelService) Get(ctx context.Context, channelID string, opts ...option.RequestOption) (res *Channel, err error)
```

**GET `/channels/{channel_id}`** · [实现](../forward/channel.go#L111) · 获取资源详情。

路径参数：`channelID string`，必填。

返回：`*Channel, error`。[Channel](../forward/channel.go#L185)。

无请求参数对象；可直接传入请求选项。

### `Channels.Update`

```go
func (r *ChannelService) Update(ctx context.Context, channelID string, params ChannelUpdateParams, opts ...option.RequestOption) (res *Channel, err error)
```

**POST `/channels/{channel_id}`** · [实现](../forward/channel.go#L123) · 更新资源。

路径参数：`channelID string`，必填。

返回：`*Channel, error`。[Channel](../forward/channel.go#L185)。

请求对象 `params`：[ChannelUpdateParams](../forward/channel.go#L136)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Name](../forward/channel.go#L138)：`param.Opt[string]` | json / `name` | 否 | Channel 展示名。 |
| [IdentityID](../forward/channel.go#L140)：`param.Opt[string]` | json / `identity_id` | 否 | `fixed` 模式下新的 Forward Identity ID。 |
| [TemplateID](../forward/channel.go#L142)：`param.Opt[string]` | json / `template_id` | 否 | `fixed` 模式下新的 Forward Template ID。 |
| [Enabled](../forward/channel.go#L144)：`param.Opt[bool]` | json / `enabled` | 否 | 人工启停开关。 |
| [ChannelConfig](../forward/channel.go#L145)：`map[string]any` | json / `channel_config` | 否 | 见字段定义。 |
| [IdempotencyKey](../forward/channel.go#L147)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Channels.Delete`

```go
func (r *ChannelService) Delete(ctx context.Context, channelID string, opts ...option.RequestOption) (res *DeletedChannel, err error)
```

**DELETE `/channels/{channel_id}`** · [实现](../forward/channel.go#L158) · 删除资源。

路径参数：`channelID string`，必填。

返回：`*DeletedChannel, error`。[DeletedChannel](../forward/channel.go#L169)。

无请求参数对象；可直接传入请求选项。

### `Channels.QRSessions.New`

```go
func (r *ChannelQRSessionService) New(ctx context.Context, channelID string, params ChannelQRSessionNewParams, opts ...option.RequestOption) (res *ChannelQRSession, err error)
```

**POST `/channels/{channel_id}/qr_sessions`** · [实现](../forward/channelqrsession.go#L28) · 创建资源。

路径参数：`channelID string`，必填。

返回：`*ChannelQRSession, error`。[ChannelQRSession](../forward/channelqrsession.go#L63)。

请求对象 `params`：[ChannelQRSessionNewParams](../forward/channelqrsession.go#L41)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [IdempotencyKey](../forward/channelqrsession.go#L43)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 有副作用请求可选的幂等键。 |

### `Channels.QRSessions.Get`

```go
func (r *ChannelQRSessionService) Get(ctx context.Context, sessionKey string, opts ...option.RequestOption) (res *ChannelQRSession, err error)
```

**GET `/qr_sessions/{session_key}`** · [实现](../forward/channelqrsession.go#L52) · 获取资源详情。

路径参数：`sessionKey string`，必填。

返回：`*ChannelQRSession, error`。[ChannelQRSession](../forward/channelqrsession.go#L63)。

无请求参数对象；可直接传入请求选项。

## ChannelPairings

建立或删除 Identity 与消息渠道的配对关系。

### `ChannelPairings.New`

```go
func (r *ChannelPairingService) New(ctx context.Context, params ChannelPairingNewParams, opts ...option.RequestOption) (res *ChannelPairing, err error)
```

**POST `/channel_pairings`** · [实现](../forward/channelpairing.go#L27) · 创建资源。

返回：`*ChannelPairing, error`。[ChannelPairing](../forward/channelpairing.go#L87)。

请求对象 `params`：[ChannelPairingNewParams](../forward/channelpairing.go#L37)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Code](../forward/channelpairing.go#L39)：`string` | json / `code` | 是 | Channel 消息中显示的 6 位配对码。 |
| [IdentityID](../forward/channelpairing.go#L41)：`string` | json / `identity_id` | 是 | 要绑定的 Forward Identity ID。 |
| [TemplateID](../forward/channelpairing.go#L43)：`string` | json / `template_id` | 是 | 要绑定的 Forward Template ID。 |
| [IdempotencyKey](../forward/channelpairing.go#L45)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 由客户端生成的唯一幂等键，用于安全重试同一次配对请求。 |

### `ChannelPairings.Delete`

```go
func (r *ChannelPairingService) Delete(ctx context.Context, pairingID string, opts ...option.RequestOption) (res *DeletedChannelPairing, err error)
```

**DELETE `/channel_pairings/{pairing_id}`** · [实现](../forward/channelpairing.go#L58) · 删除资源。

路径参数：`pairingID string`，必填。

返回：`*DeletedChannelPairing, error`。[DeletedChannelPairing](../forward/channelpairing.go#L69)。

无请求参数对象；可直接传入请求选项。

## Environments

管理执行环境配置和生命周期。

### `Environments.List`

```go
func (r *EnvironmentService) List(ctx context.Context, params EnvironmentListParams, opts ...option.RequestOption) (res *pagination.PageCursor[Environment], err error)
```

**GET `/environments`** · [实现](../forward/environment.go#L30) · 列出资源。

返回：`*pagination.PageCursor[Environment], error`。[Environment](../forward/environment.go#L166)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[EnvironmentListParams](../forward/environment.go#L50)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Limit](../forward/environment.go#L52)：`param.Opt[int64]` | query / `limit` | 否 | 分页大小，最大 100。 |
| [Page](../forward/environment.go#L54)：`param.Opt[string]` | query / `page` | 否 | 分页游标（推荐使用），取值来自上一页响应的 `next_page`；与 `after_id`、`before_id` 互斥。 |
| [AfterID](../forward/environment.go#L56)：`param.Opt[string]` | query / `after_id` | 否 | 向后翻页游标；与 `page`、`before_id` 互斥。 |
| [BeforeID](../forward/environment.go#L58)：`param.Opt[string]` | query / `before_id` | 否 | 向前翻页游标；与 `page`、`after_id` 互斥。 |

### `Environments.ListAutoPaging`

```go
func (r *EnvironmentService) ListAutoPaging(ctx context.Context, params EnvironmentListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[Environment]
```

[实现](../forward/environment.go#L46)。自动遍历 `Environments.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Environments.New`

```go
func (r *EnvironmentService) New(ctx context.Context, params EnvironmentNewParams, opts ...option.RequestOption) (res *Environment, err error)
```

**POST `/environments`** · [实现](../forward/environment.go#L67) · 创建资源。

返回：`*Environment, error`。[Environment](../forward/environment.go#L166)。

请求对象 `params`：[EnvironmentNewParams](../forward/environment.go#L77)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Name](../forward/environment.go#L79)：`string` | json / `name` | 是 | Environment 名称；去除首尾空白后不能为空。 |
| [Description](../forward/environment.go#L81)：`param.Opt[string]` | json / `description` | 否 | 描述。 |
| [Config](../forward/environment.go#L83)：`map[string]any` | json / `config` | 否 | Environment 运行时配置对象；省略时默认使用 `{"type":"cloud"}`。显式传入时不能为 `null` 或空对象。字段详见 schemas。 |
| [Metadata](../forward/environment.go#L85)：`map[string]any` | json / `metadata` | 否 | Environment metadata；省略时为 `{}`，显式传入时不能为 `null`。 |
| [IdempotencyKey](../forward/environment.go#L87)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 建议创建请求携带。相同 key 和相同请求可安全重试。 |

### `Environments.Get`

```go
func (r *EnvironmentService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *Environment, err error)
```

**GET `/environments/{id}`** · [实现](../forward/environment.go#L100) · 获取资源详情。

路径参数：`id string`，必填。

返回：`*Environment, error`。[Environment](../forward/environment.go#L166)。

无请求参数对象；可直接传入请求选项。

### `Environments.Update`

```go
func (r *EnvironmentService) Update(ctx context.Context, id string, params EnvironmentUpdateParams, opts ...option.RequestOption) (res *Environment, err error)
```

**POST `/environments/{id}`** · [实现](../forward/environment.go#L112) · 更新资源。

路径参数：`id string`，必填。

返回：`*Environment, error`。[Environment](../forward/environment.go#L166)。

请求对象 `params`：[EnvironmentUpdateParams](../forward/environment.go#L123)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Name](../forward/environment.go#L125)：`param.Opt[string]` | json / `name` | 否 | 新名称。 |
| [Description](../forward/environment.go#L127)：`param.Opt[string]` | json / `description` | 否 | 新描述。 |
| [Config](../forward/environment.go#L129)：`map[string]any` | json / `config` | 否 | 新配置；传入时不能为 `null`，显式 `null` 返回 400。字段详见 schemas。 |
| [Metadata](../forward/environment.go#L131)：`map[string]any` | json / `metadata` | 否 | 要合并的 Environment metadata；传入时不能为 `null`，显式 `null` 返回 400。 |

### `Environments.Archive`

```go
func (r *EnvironmentService) Archive(ctx context.Context, id string, opts ...option.RequestOption) (res *Environment, err error)
```

**POST `/environments/{id}/archive`** · [实现](../forward/environment.go#L144) · 归档资源。

路径参数：`id string`，必填。

返回：`*Environment, error`。[Environment](../forward/environment.go#L166)。

无请求参数对象；可直接传入请求选项。

### `Environments.Delete`

```go
func (r *EnvironmentService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error)
```

**DELETE `/environments/{id}`** · [实现](../forward/environment.go#L155) · 删除资源。

路径参数：`id string`，必填。

返回：`error`。成功时 error 为 nil。

无请求参数对象；可直接传入请求选项。

## Files

上传文件、查询元数据、列举、下载内容和删除文件。上传后需另外挂载到会话才能供 Agent 使用。

### `Files.List`

```go
func (r *FileService) List(ctx context.Context, params FileListParams, opts ...option.RequestOption) (res *pagination.PageCursor[FileMetadata], err error)
```

**GET `/files`** · [实现](../forward/file.go#L31) · 列出资源。

返回：`*pagination.PageCursor[FileMetadata], error`。[FileMetadata](../forward/file.go#L135)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[FileListParams](../forward/file.go#L51)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Limit](../forward/file.go#L53)：`param.Opt[int64]` | query / `limit` | 否 | 分页大小，最大 100。 |
| [Page](../forward/file.go#L55)：`param.Opt[string]` | query / `page` | 否 | 分页游标（推荐使用），取值来自上一页响应的 `next_page`；与 `after_id`、`before_id` 互斥。 |
| [AfterID](../forward/file.go#L57)：`param.Opt[string]` | query / `after_id` | 否 | 向后翻页游标；与 `page`、`before_id` 互斥。 |
| [BeforeID](../forward/file.go#L59)：`param.Opt[string]` | query / `before_id` | 否 | 向前翻页游标；与 `page`、`after_id` 互斥。 |
| [Name](../forward/file.go#L61)：`param.Opt[string]` | query / `name` | 否 | 按文件名搜索。 |
| [ScopeID](../forward/file.go#L63)：`param.Opt[string]` | query / `scope_id` | 否 | 按资源作用域 ID 过滤，常用于 Session 资源文件查询。传入时不要同时使用 `before_id` 或 `after_id`；当前游标参数在该过滤模式下不生效。 |

### `Files.ListAutoPaging`

```go
func (r *FileService) ListAutoPaging(ctx context.Context, params FileListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[FileMetadata]
```

[实现](../forward/file.go#L47)。自动遍历 `Files.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Files.Upload`

```go
func (r *FileService) Upload(ctx context.Context, params FileUploadParams, opts ...option.RequestOption) (res *FileMetadata, err error)
```

**POST `/files`** · [实现](../forward/file.go#L72) · 上传文件。

返回：`*FileMetadata, error`。[FileMetadata](../forward/file.go#L135)。

请求对象 `params`：[FileUploadParams](../forward/file.go#L82)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [File](../forward/file.go#L84)：`io.Reader` | multipart / `file` | 是 | 待上传文件内容。支持类型见支持上传的文件类型。 |
| [Name](../forward/file.go#L86)：`param.Opt[string]` | multipart / `name` | 否 | 文件展示名，未传时使用 multipart 文件名；规范化后长度为 1-255 bytes。 |
| [Purpose](../forward/file.go#L88)：`param.Opt[string]` | multipart / `purpose` | 否 | 文件用途，默认 `user_upload`；作为 Batch 输入文件时必须传 `session_resource`。 |
| [Metadata](../forward/file.go#L90)：`map[string]any` | multipart / `metadata` | 否 | 元数据对象；`created_by` 为保留字段，不可传入（传入返回 400）。 |
| [IdempotencyKey](../forward/file.go#L92)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 可选创建请求幂等键。传入时相同 key 只能用于相同请求；不传时不提供本地幂等重放保护。 |

### `Files.GetMetadata`

```go
func (r *FileService) GetMetadata(ctx context.Context, fileID string, opts ...option.RequestOption) (res *FileMetadata, err error)
```

**GET `/files/{file_id}`** · [实现](../forward/file.go#L101) · 获取文件元数据。

路径参数：`fileID string`，必填。

返回：`*FileMetadata, error`。[FileMetadata](../forward/file.go#L135)。

无请求参数对象；可直接传入请求选项。

### `Files.Delete`

```go
func (r *FileService) Delete(ctx context.Context, fileID string, opts ...option.RequestOption) (err error)
```

**DELETE `/files/{file_id}`** · [实现](../forward/file.go#L113) · 删除资源。

路径参数：`fileID string`，必填。

返回：`error`。成功时 error 为 nil。

无请求参数对象；可直接传入请求选项。

### `Files.Download`

```go
func (r *FileService) Download(ctx context.Context, fileID string, opts ...option.RequestOption) (res *http.Response, err error)
```

**GET `/files/{file_id}/content`** · [实现](../forward/file.go#L125) · 下载文件内容。

路径参数：`fileID string`，必填。

返回：`*http.Response, error`。HTTP 响应；读取后关闭 Body。

无请求参数对象；可直接传入请求选项。

## Skills

上传 Skill 文件树并管理版本；文件名使用相对路径，例如 my-skill/SKILL.md。

### `Skills.List`

```go
func (r *SkillService) List(ctx context.Context, params SkillListParams, opts ...option.RequestOption) (res *pagination.PageCursor[Skill], err error)
```

**GET `/skills`** · [实现](../forward/skill.go#L32) · 列出资源。

返回：`*pagination.PageCursor[Skill], error`。[Skill](../forward/skill.go#L177)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[SkillListParams](../forward/skill.go#L52)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Limit](../forward/skill.go#L54)：`param.Opt[int64]` | query / `limit` | 否 | 分页大小，最大 100。 |
| [Page](../forward/skill.go#L56)：`param.Opt[string]` | query / `page` | 否 | 分页游标（推荐使用），取值来自上一页响应的 `next_page`；与 `after_id`、`before_id` 互斥。 |
| [AfterID](../forward/skill.go#L58)：`param.Opt[string]` | query / `after_id` | 否 | 向后翻页游标；与 `page`、`before_id` 互斥。 |
| [BeforeID](../forward/skill.go#L60)：`param.Opt[string]` | query / `before_id` | 否 | 向前翻页游标；与 `page`、`after_id` 互斥。 |
| [DisplayTitle](../forward/skill.go#L62)：`param.Opt[string]` | query / `display_title` | 否 | 按 Skill 展示名前缀搜索，不区分大小写。 |
| [Source](../forward/skill.go#L64)：`param.Opt[string]` | query / `source` | 否 | 按 Skill 来源过滤，可选 `custom`、`qoder`。传 `source` 时不支持 `before_id`。 |
| [Name](../forward/skill.go#L66)：`param.Opt[string]` | query / `name` | 否 | ⚠️ **已弃用**：`display_title` 的兼容别名，语义完全一致。请使用 `display_title`。 |

### `Skills.ListAutoPaging`

```go
func (r *SkillService) ListAutoPaging(ctx context.Context, params SkillListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[Skill]
```

[实现](../forward/skill.go#L48)。自动遍历 `Skills.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Skills.New`

```go
func (r *SkillService) New(ctx context.Context, params SkillNewParams, opts ...option.RequestOption) (res *Skill, err error)
```

**POST `/skills`** · [实现](../forward/skill.go#L75) · 创建资源。

返回：`*Skill, error`。[Skill](../forward/skill.go#L177)。

请求对象 `params`：[SkillNewParams](../forward/skill.go#L85)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Files](../forward/skill.go#L87)：`[]io.Reader` | multipart / `files` | 否 | 推荐上传字段，可**重复出现**多次。支持两种形态： ① 单个 `.zip` 包； ② 裸文件树——每个 part 独立上传一个文件，`filename` 携带相对路径（如 `code-review/SKILL.md`、`code-review/scripts/run.sh`）。 压缩包本身与解压后总大小均不超过 50 MB。 |
| [Metadata](../forward/skill.go#L89)：`map[string]any` | multipart / `metadata` | 否 | 调用方元数据对象，最多 15 个键；`created_by` 为保留字段，不可传入（传入返回 400）。 |
| [IconID](../forward/skill.go#L91)：`param.Opt[string]` | multipart / `icon_id` | 否 | Forward Resource icon 公开 ID。 |
| [File](../forward/skill.go#L93)：`io.Reader` | multipart / `file` | 否 | ⚠️ **已弃用**：单个 `.zip` 包，宽松包规则。命中时响应头返回 `Deprecation: true`。请迁移到 `files`。 |
| [Name](../forward/skill.go#L95)：`param.Opt[string]` | multipart / `name` | 否 | ⚠️ **已弃用**：最终名称始终从上传包内 `SKILL.md` frontmatter 的 `name` 解析。字段保留仅为兼容，传入将被忽略。 |
| [Description](../forward/skill.go#L97)：`param.Opt[string]` | multipart / `description` | 否 | ⚠️ **已弃用**：最终描述始终从 `SKILL.md` 解析。 |
| [Type](../forward/skill.go#L99)：`param.Opt[string]` | multipart / `type` | 否 | ⚠️ **已弃用**：Skill 创建类型，可选 `custom`、`prebuilt`，默认 `custom`。`prebuilt` 会使响应 `source` 字段返回 `qoder`（其余为 `custom`）。 |
| [IdempotencyKey](../forward/skill.go#L101)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 建议提供。相同 key 且规范化后的 `files` 指纹一致时可安全重试。 |

### `Skills.Get`

```go
func (r *SkillService) Get(ctx context.Context, id string, params SkillGetParams, opts ...option.RequestOption) (res *Skill, err error)
```

**GET `/skills/{id}`** · [实现](../forward/skill.go#L110) · 获取资源详情。

路径参数：`id string`，必填。

返回：`*Skill, error`。[Skill](../forward/skill.go#L177)。

请求对象 `params`：[SkillGetParams](../forward/skill.go#L121)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [IncludeContent](../forward/skill.go#L123)：`param.Opt[bool]` | query / `include_content` | 否 | ⚠️ **已弃用**：为 `true` 时随响应返回 `content` 与 `content_encoding`（base64 zip）。命中时响应头会返回 `Deprecation: true`。请改用 下载 Skill 版本内容。 |

### `Skills.Update`

```go
func (r *SkillService) Update(ctx context.Context, id string, params SkillUpdateParams, opts ...option.RequestOption) (res *Skill, err error)
```

**PUT `/skills/{id}`** · [实现](../forward/skill.go#L132) · 更新资源。

路径参数：`id string`，必填。

返回：`*Skill, error`。[Skill](../forward/skill.go#L177)。

请求对象 `params`：[SkillUpdateParams](../forward/skill.go#L143)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Description](../forward/skill.go#L145)：`param.Opt[string]` | json / `description` | 否 | 新描述。 |
| [Content](../forward/skill.go#L147)：`param.Opt[string]` | json / `content` | 否 | 新内容（zip 包内容）。压缩包本身与解压后总大小均不超过 50 MB，超过返回 400；请求体整体（含 base64 编码与 JSON 信封）上限约 67.7 MB，超过返回 413。 |
| [ContentEncoding](../forward/skill.go#L149)：`param.Opt[string]` | json / `content_encoding` | 否 | `content` 的编码。支持 `base64`、`utf-8`、`utf8`、`plain`、`text`；省略时按 UTF-8 文本处理。传入该字段时必须同时提供非空 `content`。 |
| [Metadata](../forward/skill.go#L151)：`map[string]any` | json / `metadata` | 否 | 元数据对象，会**替换**当前 metadata（非合并）；传入时不能为 `null`，value 必须为 string。`created_by` 为保留字段，不可传入（传入返回 400）。 |
| [IconID](../forward/skill.go#L153)：`param.Opt[string]` | json / `icon_id` | 否 | 更新或清空 Forward icon。；允许 null。 |
| [Name](../forward/skill.go#L155)：`param.Opt[string]` | json / `name` | 否 | ⚠️ **已弃用**：技能名不可修改。传入必须与当前规范名完全一致，否则返回 400；一致时为空操作。 |

### `Skills.Delete`

```go
func (r *SkillService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error)
```

**DELETE `/skills/{id}`** · [实现](../forward/skill.go#L166) · 删除资源。

路径参数：`id string`，必填。

返回：`error`。成功时 error 为 nil。

无请求参数对象；可直接传入请求选项。

### `Skills.Versions.List`

```go
func (r *SkillVersionService) List(ctx context.Context, id string, params SkillVersionListParams, opts ...option.RequestOption) (res *pagination.PageCursor[SkillVersion], err error)
```

**GET `/skills/{id}/versions`** · [实现](../forward/skillversion.go#L31) · 列出资源。

路径参数：`id string`，必填。

返回：`*pagination.PageCursor[SkillVersion], error`。[SkillVersion](../forward/skillversion.go#L152)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[SkillVersionListParams](../forward/skillversion.go#L54)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Limit](../forward/skillversion.go#L56)：`param.Opt[int64]` | query / `limit` | 否 | 分页大小，最大 100，默认 20。 |
| [Page](../forward/skillversion.go#L58)：`param.Opt[string]` | query / `page` | 否 | 向后翻页游标；取值来自上一页响应的 `next_page`；不传即从第一页开始。 |

### `Skills.Versions.ListAutoPaging`

```go
func (r *SkillVersionService) ListAutoPaging(ctx context.Context, id string, params SkillVersionListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[SkillVersion]
```

[实现](../forward/skillversion.go#L50)。自动遍历 `Skills.Versions.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Skills.Versions.New`

```go
func (r *SkillVersionService) New(ctx context.Context, id string, params SkillVersionNewParams, opts ...option.RequestOption) (res *SkillVersion, err error)
```

**POST `/skills/{id}/versions`** · [实现](../forward/skillversion.go#L67) · 创建资源。

路径参数：`id string`，必填。

返回：`*SkillVersion, error`。[SkillVersion](../forward/skillversion.go#L152)。

请求对象 `params`：[SkillVersionNewParams](../forward/skillversion.go#L78)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Files](../forward/skillversion.go#L80)：`[]io.Reader` | multipart / `files` | 是 | 上传字段，可**重复**出现多次。支持两种形态： • 单个 `.zip` 包； • 裸文件树——每个 part 独立上传一个文件，`filename` 携带相对路径（如 `customer-reply/SKILL.md`、`customer-reply/scripts/run.sh`）。 压缩包本身与解压后总大小均不超过 50 MB。 |

### `Skills.Versions.Get`

```go
func (r *SkillVersionService) Get(ctx context.Context, id string, version string, opts ...option.RequestOption) (res *SkillVersion, err error)
```

**GET `/skills/{id}/versions/{version}`** · [实现](../forward/skillversion.go#L89) · 获取资源详情。

路径参数：`id string`，必填。

路径参数：`version string`，必填。

返回：`*SkillVersion, error`。[SkillVersion](../forward/skillversion.go#L152)。

无请求参数对象；可直接传入请求选项。

### `Skills.Versions.Delete`

```go
func (r *SkillVersionService) Delete(ctx context.Context, id string, version string, opts ...option.RequestOption) (res *DeletedSkillVersion, err error)
```

**DELETE `/skills/{id}/versions/{version}`** · [实现](../forward/skillversion.go#L104) · 删除资源。

路径参数：`id string`，必填。

路径参数：`version string`，必填。

返回：`*DeletedSkillVersion, error`。[DeletedSkillVersion](../forward/skillversion.go#L133)。

无请求参数对象；可直接传入请求选项。

### `Skills.Versions.Download`

```go
func (r *SkillVersionService) Download(ctx context.Context, id string, version string, opts ...option.RequestOption) (res *http.Response, err error)
```

**GET `/skills/{id}/versions/{version}/content`** · [实现](../forward/skillversion.go#L119) · 下载文件内容。

路径参数：`id string`，必填。

路径参数：`version string`，必填。

返回：`*http.Response, error`。HTTP 响应；读取后关闭 Body。

无请求参数对象；可直接传入请求选项。

## Vaults

管理凭据容器及其凭据。

### `Vaults.List`

```go
func (r *VaultService) List(ctx context.Context, params VaultListParams, opts ...option.RequestOption) (res *pagination.PageCursor[Vault], err error)
```

**GET `/vaults`** · [实现](../forward/vault.go#L31) · 列出资源。

返回：`*pagination.PageCursor[Vault], error`。[Vault](../forward/vault.go#L120)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[VaultListParams](../forward/vault.go#L51)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Limit](../forward/vault.go#L53)：`param.Opt[int64]` | query / `limit` | 否 | 分页大小，最大 100。 |
| [Page](../forward/vault.go#L55)：`param.Opt[string]` | query / `page` | 否 | 分页游标（推荐使用），取值来自上一页响应的 `next_page`；与 `after_id`、`before_id` 互斥。 |
| [AfterID](../forward/vault.go#L57)：`param.Opt[string]` | query / `after_id` | 否 | 向后翻页游标；与 `page`、`before_id` 互斥。 |
| [BeforeID](../forward/vault.go#L59)：`param.Opt[string]` | query / `before_id` | 否 | 向前翻页游标；与 `page`、`after_id` 互斥。 |
| [Name](../forward/vault.go#L61)：`param.Opt[string]` | query / `name` | 否 | 按 `display_name` 搜索。 |

### `Vaults.ListAutoPaging`

```go
func (r *VaultService) ListAutoPaging(ctx context.Context, params VaultListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[Vault]
```

[实现](../forward/vault.go#L47)。自动遍历 `Vaults.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Vaults.New`

```go
func (r *VaultService) New(ctx context.Context, params VaultNewParams, opts ...option.RequestOption) (res *Vault, err error)
```

**POST `/vaults`** · [实现](../forward/vault.go#L70) · 创建资源。

返回：`*Vault, error`。[Vault](../forward/vault.go#L120)。

请求对象 `params`：[VaultNewParams](../forward/vault.go#L80)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [DisplayName](../forward/vault.go#L82)：`string` | json / `display_name` | 是 | Vault 展示名。 |
| [Metadata](../forward/vault.go#L84)：`map[string]any` | json / `metadata` | 否 | 元数据对象；`created_by` 为保留字段，不可传入（传入返回 400）。 |
| [IdempotencyKey](../forward/vault.go#L86)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 可选创建请求幂等键。传入时相同 key 只能用于相同请求；不传时不提供本地幂等重放保护。 |

### `Vaults.Get`

```go
func (r *VaultService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *Vault, err error)
```

**GET `/vaults/{id}`** · [实现](../forward/vault.go#L97) · 获取资源详情。

路径参数：`id string`，必填。

返回：`*Vault, error`。[Vault](../forward/vault.go#L120)。

无请求参数对象；可直接传入请求选项。

### `Vaults.Delete`

```go
func (r *VaultService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error)
```

**DELETE `/vaults/{id}`** · [实现](../forward/vault.go#L109) · 删除资源。

路径参数：`id string`，必填。

返回：`error`。成功时 error 为 nil。

无请求参数对象；可直接传入请求选项。

### `Vaults.Credentials.List`

```go
func (r *VaultCredentialService) List(ctx context.Context, id string, params VaultCredentialListParams, opts ...option.RequestOption) (res *pagination.PageCursor[VaultCredential], err error)
```

**GET `/vaults/{id}/credentials`** · [实现](../forward/vaultcredential.go#L30) · 列出资源。

路径参数：`id string`，必填。

返回：`*pagination.PageCursor[VaultCredential], error`。[VaultCredential](../forward/vaultcredential.go#L135)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[VaultCredentialListParams](../forward/vaultcredential.go#L53)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Limit](../forward/vaultcredential.go#L55)：`param.Opt[int64]` | query / `limit` | 否 | 分页大小，最大 100。 |
| [Page](../forward/vaultcredential.go#L57)：`param.Opt[string]` | query / `page` | 否 | 分页游标（推荐使用），取值来自上一页响应的 `next_page`；与 `after_id`、`before_id` 互斥。 |
| [AfterID](../forward/vaultcredential.go#L59)：`param.Opt[string]` | query / `after_id` | 否 | 向后翻页游标；与 `page`、`before_id` 互斥。 |
| [BeforeID](../forward/vaultcredential.go#L61)：`param.Opt[string]` | query / `before_id` | 否 | 向前翻页游标；与 `page`、`after_id` 互斥。 |
| [Name](../forward/vaultcredential.go#L63)：`param.Opt[string]` | query / `name` | 否 | 按 `mcp_server_url` 搜索。 |

### `Vaults.Credentials.ListAutoPaging`

```go
func (r *VaultCredentialService) ListAutoPaging(ctx context.Context, id string, params VaultCredentialListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[VaultCredential]
```

[实现](../forward/vaultcredential.go#L49)。自动遍历 `Vaults.Credentials.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `Vaults.Credentials.New`

```go
func (r *VaultCredentialService) New(ctx context.Context, id string, params VaultCredentialNewParams, opts ...option.RequestOption) (res *VaultCredential, err error)
```

**POST `/vaults/{id}/credentials`** · [实现](../forward/vaultcredential.go#L72) · 创建资源。

路径参数：`id string`，必填。

返回：`*VaultCredential, error`。[VaultCredential](../forward/vaultcredential.go#L135)。

请求对象 `params`：[VaultCredentialNewParams](../forward/vaultcredential.go#L85)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Auth](../forward/vaultcredential.go#L87)：`map[string]any` | json / `auth` | 是 | Credential 认证信息，支持 `static_bearer`、`mcp_oauth`；响应只返回脱敏后的非密文字段。 |
| [DisplayName](../forward/vaultcredential.go#L89)：`param.Opt[string]` | json / `display_name` | 否 | 兼容字段；当前不持久化，Forward 响应固定为空字符串。 |
| [Metadata](../forward/vaultcredential.go#L91)：`map[string]any` | json / `metadata` | 否 | 元数据对象；`created_by` 为保留字段，不可传入（传入返回 400）。 |
| [IdempotencyKey](../forward/vaultcredential.go#L93)：`param.Opt[string]` | header / `Idempotency-Key` | 否 | 可选创建请求幂等键。相同 key 只能用于相同请求。 |

### `Vaults.Credentials.Get`

```go
func (r *VaultCredentialService) Get(ctx context.Context, id string, credID string, opts ...option.RequestOption) (res *VaultCredential, err error)
```

**GET `/vaults/{id}/credentials/{cred_id}`** · [实现](../forward/vaultcredential.go#L106) · 获取资源详情。

路径参数：`id string`，必填。

路径参数：`credID string`，必填。

返回：`*VaultCredential, error`。[VaultCredential](../forward/vaultcredential.go#L135)。

无请求参数对象；可直接传入请求选项。

### `Vaults.Credentials.Delete`

```go
func (r *VaultCredentialService) Delete(ctx context.Context, id string, credID string, opts ...option.RequestOption) (err error)
```

**DELETE `/vaults/{id}/credentials/{cred_id}`** · [实现](../forward/vaultcredential.go#L121) · 删除资源。

路径参数：`id string`，必填。

路径参数：`credID string`，必填。

返回：`error`。成功时 error 为 nil。

无请求参数对象；可直接传入请求选项。

## MemoryStores

管理记忆存储、记忆条目及版本。删除条目、删除存储、归档存储和版本脱敏是不同操作。

### `MemoryStores.List`

```go
func (r *MemoryStoreService) List(ctx context.Context, params MemoryStoreListParams, opts ...option.RequestOption) (res *pagination.Page[MemoryStore], err error)
```

**GET `/memory_stores`** · [实现](../forward/memorystore.go#L32) · 列出资源。

返回：`*pagination.Page[MemoryStore], error`。[MemoryStore](../forward/memorystore.go#L182)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[MemoryStoreListParams](../forward/memorystore.go#L52)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Limit](../forward/memorystore.go#L54)：`param.Opt[int64]` | query / `limit` | 否 | 每页返回数量上限，1..100，默认 20。 |
| [BeforeID](../forward/memorystore.go#L56)：`param.Opt[string]` | query / `before_id` | 否 | 向前翻页游标，与 `after_id` 互斥。 |
| [AfterID](../forward/memorystore.go#L58)：`param.Opt[string]` | query / `after_id` | 否 | 向后翻页游标，与 `before_id` 互斥。 |
| [SystemManaged](../forward/memorystore.go#L60)：`param.Opt[bool]` | query / `system_managed` | 否 | 三态过滤：`true` 只返回系统默认库；`false` 只返回用户创建的库；不传不过滤。 |

### `MemoryStores.ListAutoPaging`

```go
func (r *MemoryStoreService) ListAutoPaging(ctx context.Context, params MemoryStoreListParams, opts ...option.RequestOption) *pagination.PageAutoPager[MemoryStore]
```

[实现](../forward/memorystore.go#L48)。自动遍历 `MemoryStores.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `MemoryStores.New`

```go
func (r *MemoryStoreService) New(ctx context.Context, params MemoryStoreNewParams, opts ...option.RequestOption) (res *MemoryStore, err error)
```

**POST `/memory_stores`** · [实现](../forward/memorystore.go#L69) · 创建资源。

返回：`*MemoryStore, error`。[MemoryStore](../forward/memorystore.go#L182)。

请求对象 `params`：[MemoryStoreNewParams](../forward/memorystore.go#L77)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Name](../forward/memorystore.go#L79)：`string` | json / `name` | 是 | Store 展示名，非空。不允许非打印控制字符（`U+0000`–`U+001F`、`U+007F`），换行 `\n`、回车 `\r`、制表 `\t` 除外。 |
| [Description](../forward/memorystore.go#L81)：`param.Opt[string]` | json / `description` | 否 | 自由文本描述。不允许非打印控制字符。 |
| [Metadata](../forward/memorystore.go#L83)：`map[string]any` | json / `metadata` | 否 | 键值元数据，值必须为字符串。最多 **15** 个键；键 1..64 字符；值 ≤512 字符。`created_by` 是 Forward 保留键，服务端自动写入 `"forward"`；调用方传入 `created_by` 会返回 `400 invalid_request_error`。详见 Store metadata 约束。 |
| [IdempotencyKey](../forward/memorystore.go#L85)：`string` | header / `Idempotency-Key` | 是 | 创建请求幂等键。相同 key 只能用于相同请求体；不传返回 `400`。 |

### `MemoryStores.Get`

```go
func (r *MemoryStoreService) Get(ctx context.Context, memoryStoreID string, opts ...option.RequestOption) (res *MemoryStore, err error)
```

**GET `/memory_stores/{memory_store_id}`** · [实现](../forward/memorystore.go#L98) · 获取资源详情。

路径参数：`memoryStoreID string`，必填。

返回：`*MemoryStore, error`。[MemoryStore](../forward/memorystore.go#L182)。

无请求参数对象；可直接传入请求选项。

### `MemoryStores.Update`

```go
func (r *MemoryStoreService) Update(ctx context.Context, memoryStoreID string, params MemoryStoreUpdateParams, opts ...option.RequestOption) (res *MemoryStore, err error)
```

**POST `/memory_stores/{memory_store_id}`** · [实现](../forward/memorystore.go#L110) · 更新资源。

路径参数：`memoryStoreID string`，必填。

返回：`*MemoryStore, error`。[MemoryStore](../forward/memorystore.go#L182)。

请求对象 `params`：[MemoryStoreUpdateParams](../forward/memorystore.go#L121)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Name](../forward/memorystore.go#L123)：`param.Opt[string]` | json / `name` | 否 | 新名称。传入时非空且不含非打印控制字符。 |
| [Description](../forward/memorystore.go#L125)：`param.Opt[string]` | json / `description` | 否 | 新描述。传入时不含非打印控制字符。 |
| [Metadata](../forward/memorystore.go#L127)：`map[string]any` | json / `metadata` | 否 | 新元数据，**整体替换**当前 metadata（非合并）。约束详见 Store metadata 约束。 |

### `MemoryStores.Delete`

```go
func (r *MemoryStoreService) Delete(ctx context.Context, memoryStoreID string, opts ...option.RequestOption) (res *DeletedMemoryStore, err error)
```

**DELETE `/memory_stores/{memory_store_id}`** · [实现](../forward/memorystore.go#L140) · 删除资源。

路径参数：`memoryStoreID string`，必填。

返回：`*DeletedMemoryStore, error`。[DeletedMemoryStore](../forward/memorystore.go#L163)。

无请求参数对象；可直接传入请求选项。

### `MemoryStores.Archive`

```go
func (r *MemoryStoreService) Archive(ctx context.Context, memoryStoreID string, opts ...option.RequestOption) (res *MemoryStore, err error)
```

**POST `/memory_stores/{memory_store_id}/archive`** · [实现](../forward/memorystore.go#L152) · 归档资源。

路径参数：`memoryStoreID string`，必填。

返回：`*MemoryStore, error`。[MemoryStore](../forward/memorystore.go#L182)。

无请求参数对象；可直接传入请求选项。

### `MemoryStores.Memories.List`

```go
func (r *MemoryStoreMemoryService) List(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryListParams, opts ...option.RequestOption) (res *pagination.Page[Memory], err error)
```

**GET `/memory_stores/{memory_store_id}/memories`** · [实现](../forward/memorystorememory.go#L30) · 列出资源。

路径参数：`memoryStoreID string`，必填。

返回：`*pagination.Page[Memory], error`。[Memory](../forward/memorystorememory.go#L181)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[MemoryStoreMemoryListParams](../forward/memorystorememory.go#L53)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Limit](../forward/memorystorememory.go#L55)：`param.Opt[int64]` | query / `limit` | 否 | 每页返回数量上限，1..100，默认 20。 |
| [BeforeID](../forward/memorystorememory.go#L57)：`param.Opt[string]` | query / `before_id` | 否 | 向前翻页游标，与 `after_id` 互斥。 |
| [AfterID](../forward/memorystorememory.go#L59)：`param.Opt[string]` | query / `after_id` | 否 | 向后翻页游标，与 `before_id` 互斥。 |
| [PathPrefix](../forward/memorystorememory.go#L61)：`param.Opt[string]` | query / `path_prefix` | 否 | 按 `path` 前缀过滤。**这是纯字符串前缀匹配，不是目录语义** —— `path_prefix=a/b` 也会命中 `a/bc.md`。 |

### `MemoryStores.Memories.ListAutoPaging`

```go
func (r *MemoryStoreMemoryService) ListAutoPaging(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Memory]
```

[实现](../forward/memorystorememory.go#L49)。自动遍历 `MemoryStores.Memories.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `MemoryStores.Memories.New`

```go
func (r *MemoryStoreMemoryService) New(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryNewParams, opts ...option.RequestOption) (res *Memory, err error)
```

**POST `/memory_stores/{memory_store_id}/memories`** · [实现](../forward/memorystorememory.go#L70) · 创建资源。

路径参数：`memoryStoreID string`，必填。

返回：`*Memory, error`。[Memory](../forward/memorystorememory.go#L181)。

请求对象 `params`：[MemoryStoreMemoryNewParams](../forward/memorystorememory.go#L81)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Path](../forward/memorystorememory.go#L83)：`string` | json / `path` | 是 | 库内相对路径，大小写敏感。约束详见 path 规则。 |
| [Content](../forward/memorystorememory.go#L85)：`string` | json / `content` | 是 | UTF-8 明文内容，非 base64；原始字节 ≤100 KiB。约束详见 content 约束。 |
| [Metadata](../forward/memorystorememory.go#L87)：`map[string]any` | json / `metadata` | 否 | 键值元数据，值必须为字符串。最多 **16** 个键。约束详见 Memory metadata 约束。 |

### `MemoryStores.Memories.Get`

```go
func (r *MemoryStoreMemoryService) Get(ctx context.Context, memoryStoreID string, memoryID string, opts ...option.RequestOption) (res *Memory, err error)
```

**GET `/memory_stores/{memory_store_id}/memories/{memory_id}`** · [实现](../forward/memorystorememory.go#L100) · 获取资源详情。

路径参数：`memoryStoreID string`，必填。

路径参数：`memoryID string`，必填。

返回：`*Memory, error`。[Memory](../forward/memorystorememory.go#L181)。

无请求参数对象；可直接传入请求选项。

### `MemoryStores.Memories.Update`

```go
func (r *MemoryStoreMemoryService) Update(ctx context.Context, memoryStoreID string, memoryID string, params MemoryStoreMemoryUpdateParams, opts ...option.RequestOption) (res *Memory, err error)
```

**POST `/memory_stores/{memory_store_id}/memories/{memory_id}`** · [实现](../forward/memorystorememory.go#L115) · 更新资源。

路径参数：`memoryStoreID string`，必填。

路径参数：`memoryID string`，必填。

返回：`*Memory, error`。[Memory](../forward/memorystorememory.go#L181)。

请求对象 `params`：[MemoryStoreMemoryUpdateParams](../forward/memorystorememory.go#L129)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Content](../forward/memorystorememory.go#L131)：`string` | json / `content` | 是 | 新内容，UTF-8 明文；原始字节 ≤100 KiB。约束详见 content 约束。 |
| [ContentSHA256](../forward/memorystorememory.go#L133)：`param.Opt[string]` | json / `content_sha256` | 否 | 期望的当前内容 SHA-256，用于乐观并发控制。不一致时返回 `409`。 |
| [Metadata](../forward/memorystorememory.go#L135)：`map[string]any` | json / `metadata` | 否 | 新元数据，**整体替换**当前 metadata（非合并）。未传入时保持原 metadata 不变。约束详见 Memory metadata 约束。 |

### `MemoryStores.Memories.Delete`

```go
func (r *MemoryStoreMemoryService) Delete(ctx context.Context, memoryStoreID string, memoryID string, opts ...option.RequestOption) (res *DeletedMemory, err error)
```

**DELETE `/memory_stores/{memory_store_id}/memories/{memory_id}`** · [实现](../forward/memorystorememory.go#L148) · 删除资源。

路径参数：`memoryStoreID string`，必填。

路径参数：`memoryID string`，必填。

返回：`*DeletedMemory, error`。[DeletedMemory](../forward/memorystorememory.go#L162)。

无请求参数对象；可直接传入请求选项。

### `MemoryStores.MemoryVersions.List`

```go
func (r *MemoryStoreMemoryVersionService) List(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryVersionListParams, opts ...option.RequestOption) (res *pagination.Page[MemoryVersion], err error)
```

**GET `/memory_stores/{memory_store_id}/memory_versions`** · [实现](../forward/memorystorememoryversion.go#L31) · 列出资源。

路径参数：`memoryStoreID string`，必填。

返回：`*pagination.Page[MemoryVersion], error`。[MemoryVersion](../forward/memorystorememoryversion.go#L100)、[分页类型](../convention/pagination/pagination.go)。

请求对象 `params`：[MemoryStoreMemoryVersionListParams](../forward/memorystorememoryversion.go#L54)。

| Go 字段与类型 | 协议位置 / 名称 | 必填 | 说明 |
| --- | --- | --- | --- |
| [Limit](../forward/memorystorememoryversion.go#L56)：`param.Opt[int64]` | query / `limit` | 否 | 每页返回数量上限，1..100，默认 20。 |
| [BeforeID](../forward/memorystorememoryversion.go#L58)：`param.Opt[string]` | query / `before_id` | 否 | 向前翻页游标，与 `after_id` 互斥。 |
| [AfterID](../forward/memorystorememoryversion.go#L60)：`param.Opt[string]` | query / `after_id` | 否 | 向后翻页游标，与 `before_id` 互斥。 |
| [MemoryID](../forward/memorystorememoryversion.go#L62)：`param.Opt[string]` | query / `memory_id` | 否 | 只返回该 memory（`mem_...`）的版本，用于查看单条记忆的变更历史。 |

### `MemoryStores.MemoryVersions.ListAutoPaging`

```go
func (r *MemoryStoreMemoryVersionService) ListAutoPaging(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryVersionListParams, opts ...option.RequestOption) *pagination.PageAutoPager[MemoryVersion]
```

[实现](../forward/memorystorememoryversion.go#L50)。自动遍历 `MemoryStores.MemoryVersions.List` 的分页；参数与该 List 方法一致，不新增 HTTP 路由。调用后检查 `pager.Err()`。

### `MemoryStores.MemoryVersions.Get`

```go
func (r *MemoryStoreMemoryVersionService) Get(ctx context.Context, memoryStoreID string, memoryVersionID string, opts ...option.RequestOption) (res *MemoryVersion, err error)
```

**GET `/memory_stores/{memory_store_id}/memory_versions/{memory_version_id}`** · [实现](../forward/memorystorememoryversion.go#L71) · 获取资源详情。

路径参数：`memoryStoreID string`，必填。

路径参数：`memoryVersionID string`，必填。

返回：`*MemoryVersion, error`。[MemoryVersion](../forward/memorystorememoryversion.go#L100)。

无请求参数对象；可直接传入请求选项。

### `MemoryStores.MemoryVersions.Redact`

```go
func (r *MemoryStoreMemoryVersionService) Redact(ctx context.Context, memoryStoreID string, memoryVersionID string, opts ...option.RequestOption) (res *MemoryVersion, err error)
```

**POST `/memory_stores/{memory_store_id}/memory_versions/{memory_version_id}/redact`** · [实现](../forward/memorystorememoryversion.go#L86) · 对记忆版本脱敏。

路径参数：`memoryStoreID string`，必填。

路径参数：`memoryVersionID string`，必填。

返回：`*MemoryVersion, error`。[MemoryVersion](../forward/memorystorememoryversion.go#L100)。

无请求参数对象；可直接传入请求选项。

## Models

查询当前账号可用模型和能力；创建模板或 Agent 时应使用列表中启用的模型 ID。

### `Models.List`

```go
func (r *ModelService) List(ctx context.Context, opts ...option.RequestOption) (res *ModelListResponse, err error)
```

**GET `/models`** · [实现](../forward/model.go#L24) · 列出资源。

返回：`*ModelListResponse, error`。[ModelListResponse](../forward/model.go#L68)。

无请求参数对象；可直接传入请求选项。

## 生命周期与完整场景

Identity / Template 配置影响新会话；会话资源和记忆挂载应在需要的作用域内设置。Schedule 触发后查询 `ScheduleRuns`，Batch 创建后查询 `Batches.Get`、`Batches.Tasks.List` 与输出，不能把创建成功或 queued 视为完成。清理时使用对应 Cancel、Archive、Detach、Delete，解除挂载不会删除 Memory Store。

可运行的资源准备、对话验证与清理流程见 [Forward 场景目录](../examples/forward)。
