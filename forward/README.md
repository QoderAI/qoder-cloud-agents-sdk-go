# Forward Go SDK

Forward 包提供 **110 个 HTTP API**，按 26 个具体 service 组织。调用结构与 Managed 一致，保留 Forward 的 Template、Identity、Schedule、Batch 和 Channel 业务模型。

## 创建客户端

```go
import (
    "context"

    "github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
    "github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
    "github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

client := forward.NewClient(option.WithAccessToken("<PAT or SAT>"))
ctx := context.Background()

identity, err := client.Identities.New(ctx, forward.IdentityNewParams{
    ExternalID: "customer-123",
    Name: forward.String("Customer"),
    Metadata: map[string]any{"team": "support"},
})
if err != nil { /* handle error */ }
```

默认地址为 `https://api.qoder.com/api/v1/forward/`。`NewClient` 读取 `QODER_ACCESS_TOKEN` 和 `QODER_FORWARD_BASE_URL`；显式 option 优先。Forward 不读取 Managed 的 `QODER_BASE_URL`。自定义地址传完整 API 根路径，末尾 `/` 可省略。

同一份动态凭据可传给两种模式：`forward.NewClient(option.WithCredential(credential))`。PAT 可通过 `convention.NewPATCredential(token)` 或 `convention.PATCredentialFromEnv("QODER_ACCESS_TOKEN")` 构建。请求发送 `Authorization: Bearer ...`，权限范围由服务端和凭据决定。

| 客户端入口 | 子服务 |
|---|---|
| `Templates` | — |
| `Identities` | `Configs`、`MemoryStores` |
| `Sessions` | `Events`、`Resources`、`Threads`、`Threads.Events` |
| `Schedules`、`ScheduleRuns` | Run 是客户端独立入口 |
| `Batches` | `Tasks` |
| `Channels`、`ChannelPairings` | `Channels.QRSessions`；Pairing 是客户端独立入口 |
| `Environments`、`Files` | — |
| `Skills` | `Versions` |
| `Vaults` | `Credentials` |
| `MemoryStores` | `Memories`、`MemoryVersions` |
| `Models` | — |

## 参数与响应

必填字段使用普通值，可选标量使用 `param.Opt[T]`。零值表示省略；`forward.String("")`、`forward.Bool(false)` 明确发送空字符串和 false；`param.Null[string]()` 发送 JSON null。对象、map、slice 分别支持 `param.NullStruct[T]()`、`param.NullMap[T]()`、`param.NullSlice[T]()`。空数组与省略数组不同，服务端仍会校验具体字段是否允许 null。

模型参数兼容字符串与对象：

```go
template, err := client.Templates.New(ctx, forward.TemplateNewParams{
    Name: "support",
    EnvironmentID: "env_...",
    Model: forward.ModelConfigUnionParam{OfConfig: &forward.ModelConfigParam{
        ID: "<model ID from Models.List>",
        Effort: forward.String("high"),
        ContextWindow: forward.Int(400000),
    }},
    Tools: []forward.ToolParam{{Type: "agent_toolset_20260401"}},
})
```

字符串形式为 `ModelConfigUnionParam{OfString: forward.String("<model ID>")}`。同一联合类型只设置一个分支。响应中的字符串 model 同样可从 `template.Model.ID` 读取。

Template 的工具、MCP、Skill、GitHub、资源绑定和多 Agent 配置具有独立类型。Identity Config 按资源键覆盖，使用 `ToolOverrideParam`、`MCPServerOverrideParam`、`SkillOverrideParam` 等类型，不要求在值中重复 name/id。环境变量保留任意调用方键。开放对象和 metadata 保留 map，不把文档中的示例键固定为字段；可通过 `params.SetExtraFields(...)` 发送后续新增字段。

每个响应对象保留 `RawJSON()`、`JSON.ExtraFields` 和字段元信息，例如 `result.JSON.Description.Valid()` 和 `.Raw()`（`""` 表示未返回，`"null"` 表示 null），可区分未返回、null 和实际值。`SessionEvent.Content` 使用 `json.RawMessage` 保留不同事件的数组、字符串或对象；内容块数组可通过 `ContentBlocks()` 解析。未知事件和字段不会丢失。

批量归档 Schedule 必须指定 `Scope: "by_schedule_ids"`，去重后的 ID 数量为 1～50：

```go
archived, err := client.Schedules.ArchiveMany(ctx, forward.ScheduleArchiveManyParams{
    Scope: "by_schedule_ids",
    ScheduleIDs: []string{"sched_..."},
})
if err != nil { /* handle error */ }
```

## 分页、流和文件

```go
pager := client.Templates.ListAutoPaging(ctx, forward.TemplateListParams{
    Limit: forward.Int(20), Status: forward.String("active"),
})
for pager.Next() {
    template := pager.Current()
    _ = template.ID
}
if err := pager.Err(); err != nil { /* handle error */ }
```

`List` 返回单页，可手动调用 `GetNextPage()`。Template、Identity、Session、Schedule、Batch、Channel、Memory Store 等使用 `after_id` / `before_id`；Environment、File、Skill、Vault 等使用 `next_page` → `page`，后续请求清除兼容游标。列表筛选和 limit 会保留。服务端返回不前进的游标时返回错误，避免无限翻页。

```go
stream := client.Sessions.Events.StreamEvents(ctx, "sess_...", forward.SessionEventStreamParams{
    EventDeltas: []string{"agent.message"},
})
defer stream.Close()
for stream.Next() {
    event := stream.Current()
    _ = event.RawJSON()
}
if err := stream.Err(); err != nil { /* handle error */ }
```

流使用调用方的 context；必须关闭并检查 `Err()`。同 ID 的 start/delta 事件全部保留。需要重连时，由调用者再次调用 `StreamEvents`，将 `stream.LastEventID()` 填入 `SessionEventStreamParams.LastEventID`，或使用 `option.WithHeader("Last-Event-ID", ...)`；不再自动重连。

`Files.Upload` 返回 `*FileMetadata`；`Files.Download` 获取临时链接后返回文件的 `*http.Response`，调用者关闭 Body。临时下载不会携带 API 的 PAT 或 middleware。可用 `io.Copy(destination, response.Body)` 写入目标 writer。已获得 Batch 等临时链接时，可调用 `convention.DownloadTo(ctx, nil, url, destination)`。

`Skills.New`、`Skills.Versions.New` 接收 `[]io.Reader`，以重复 `files` 字段上传，并返回 `*Skill` / `*SkillVersion`。`convention.UploadFile{Reader: reader, Name: "my-skill/SKILL.md"}` 保留文件树相对路径；metadata 编码成 JSON 字符串表单字段。`Skills.Versions.Download` 直接返回 ZIP 的 HTTP 响应；version 参数使用 `SkillVersion.Version`，而非 ID。

## 错误与配置

通过 `errors.As(err, &apiErr)` 提取 `*convention.Error`（`forward.Error` 为别名）的 `StatusCode`、`RequestID`、`Code` 和 `Message`。调用方可使用 `option.WithResponseInto(&rawResponse)` 获取 API 原始响应。

客户端、服务、方法均可设置 `option.RequestOption`，后者覆盖前者。默认最多重试两次，遵循服务端重试时间；409 不重试，无幂等键的写请求仅重试 429。需要关闭重试可传 `option.WithMaxRetries(0)`。超时、自定义 HTTP 客户端、请求头、查询和 JSON 覆盖均使用共享 option；共享客户端前先完成 Options 配置。

## 从旧 Forward SDK 迁移

这是公开 Go API 的不兼容改造，当前保留 110 个 HTTP 操作，暂不提供 Service Account Token 管理接口，旧 Go 符号不保留兼容层。

| 旧调用 | 新调用 |
|---|---|
| `NewForwardClient(credential, ...)` | `NewClient(option.WithCredential(credential), ...)`，直接返回 Client |
| `Templates.CreateTemplate(ctx).CreateTemplateRequest(body).Execute()` | `Templates.New(ctx, TemplateNewParams{...})` |
| `Sessions.GetSession(ctx, id).Execute()` | `Sessions.Get(ctx, id)` |
| `Sessions.SendSessionEvents(...).Execute()` | `Sessions.Events.Send(ctx, id, SessionEventSendParams{...})` |
| `Identities.GetIdentityConfig(...).Execute()` | `Identities.Configs.Get(ctx, identityID, templateID)` |
| `Schedules.ListScheduleRuns(...).Execute()` | `ScheduleRuns.List(ctx, ScheduleRunListParams{IdentityID: ...})` |
| `Channels.PairChannel(...).Execute()` | `ChannelPairings.New(ctx, ChannelPairingNewParams{...})` |
| `GetId()` / `GetData()` | `.ID` / `.Data`；字段存在性用 `.JSON` |
| `CreateXRequest` / `GetXResponse` 等重复类型 | `XNewParams` / 统一实体 `X` |
| `ForwardAPIError`、独立 retry / SSE / PageIterator helper | `convention.Error`、共享 option / ssestream / pagination |

方法通常返回 `(*T, error)`；无响应体的删除方法只返回 error。幂等键等请求头放入 Params 或 option，不再以 builder setter 设置。自动分页方法不计入 HTTP API 数量。Identity 的 Template 列表方法命名为 `ListTemplates`，仍请求原有 `/identities/{identity_id}/agents` 路径。

当前相邻 codegen 仍生成旧布局，不能直接覆盖本实现。SDK 通用用法见 [根目录指南](../README.md)。
