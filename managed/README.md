# Managed Cloud Agents SDK

此包实现 Qoder 官网 Managed Mode 范围内选定的 **95 个 HTTP API**。结构和辅助工具以固定的上游提交 `de6914c` 为基线。完整清单见 [API_COVERAGE.md](API_COVERAGE.md)。

服务采用具体 struct，完整保留 API 所需的参数、领域实体、联合类型、序列化与辅助方法。公开类型和文件去掉上游的 Beta 前缀，例如 `AgentService`、`AgentNewParams`、`agent.go`、`session.go`；提供方相关命名统一使用 Qoder。`Client` 专用于 Managed Mode，直接通过 `client.Agents`、`client.Sessions` 访问。自动分页方法不计为新的 HTTP API。

## 调用

```go
import (
    "context"
    managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
    "github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
)

client := managed.NewClient(option.WithAccessToken("<PAT or SAT>"))
agent, err := client.Agents.New(context.Background(), managed.AgentNewParams{
    Name: "code-reviewer",
    Model: managed.ManagedAgentsModelConfigParams{ID: "<model ID from Models.List>"},
    System: managed.String("Review code and explain your findings."),
    Metadata: map[string]string{"team": "sdk"},
})
```

默认地址为 `https://api.qoder.com/api/v1/cloud/`。`NewClient` 读取 `QODER_ACCESS_TOKEN` 和 `QODER_BASE_URL`；显式 option 优先。鉴权使用 `Authorization: Bearer ...`。`Betas` 会发送为 `x-qoder-beta`，不附加未在 Qoder 协议中定义的 beta/version 参数或请求头。SAT 可作为访问令牌传入；此包不增加令牌交换 API。

| 客户端入口 | 子服务 |
|---|---|
| `Agents` | `Versions` |
| `Sessions` | `Events`、`Resources`、`Threads`、`Threads.Events` |
| `Deployments`、`DeploymentRuns` | DeploymentRuns 是客户端直接入口 |
| `Dreams` | — |
| `Environments` | `Work` |
| `Skills` | `Versions` |
| `Vaults` | `Credentials` |
| `Files` | — |
| `MemoryStores` | `Memories`、`MemoryVersions` |
| `Models` | — |

## Conventions

公共能力集中在 `convention` 及其子包。以下工具是 API 依赖的一部分，保留其完整定义与实现，不用简化的替代实现裁剪字段或辅助方法：

- `convention`：HTTP 执行、重试、错误别名、下载链接和 `UploadFile`。
- `convention/option`：地址、访问令牌、HTTP 客户端、超时、重试、请求头、查询参数和 JSON 覆盖。
- `convention/param`：`Opt`、`Null`、`NullStruct`、`NullMap`、`NullSlice`、联合类型及 `SetExtraFields`。
- `convention/pagination`：游标分页及 `ListAutoPaging`；保留 `HasMore`、`FirstID`、`LastID`、`NextPage`，后续请求使用 `next_page`，清除兼容游标。
- `convention/respjson`：字段是否出现、是否为 null，以及未知字段和原始 JSON。
- `convention/ssestream`：增量 SSE 的 `Next`、`Current`、`Err`、`Close` 和 `LastEventID`。

可选字段的零值代表省略；`managed.String("")` 表示发送空字符串；`param.Null[string]()` 表示发送 null。对象使用 `param.NullStruct[T]()`，map 使用 `param.NullMap[map[string]string]()`。删除单个 metadata 键可用 `params.SetExtraFields(map[string]any{"metadata": map[string]any{"obsolete": nil}})`。

默认最多重试两次，遵循 `Retry-After` / `Retry-After-Ms` 和指数退避。409 不会自动重试；没有 `Idempotency-Key` 的写请求仅对 429 重试，避免在结果不确定时重复执行。可用 `WithMaxRetries(0)` 关闭重试。上下文取消会停止请求与退避。API 错误可通过 `errors.As(err, &apiError)` 提取 `*convention.Error` 的 `StatusCode`、`RequestID`、`Code`、`Message` 和 `Type()`。

SSE 返回正在读取的流；调用者关闭它并检查 `Err()`。同 ID 的 `event_start` / `event_delta` 会全部保留，也会保留未知事件的原始 JSON。需要重连时，调用者使用 `option.WithHeader("Last-Event-ID", stream.LastEventID())` 发起新连接。

## Qoder 协议适配

结构、工具方法与命名替换前一致，保留现有协议适配；HTTP 协议以 [Qoder API Reference](https://docs.qoder.com/cloud-agents/api/conventions/schemas) 为准：

| 差异 | 实现 |
|---|---|
| Credential 更新 | Qoder 仅支持更新 `auth` 和 `metadata`。`VaultCredentialUpdateParams.DisplayName` 保留以兼容现有代码，但已弃用；设置非空、空字符串或 null 都会在发送请求前报错，不能用于修改名称 |
| 文件下载返回临时链接 | `Files.Download` 获取链接后下载内容，返回 `*http.Response`；PAT 和 API middleware 不会传给存储请求。调用者关闭 Body |
| Skill 上传 | 使用重复的 `files` 表单字段；`convention.UploadFile` 支持 zip 文件名和文件树相对路径；metadata 作为 JSON 字符串表单字段发送 |
| Skill 字段 | Go 的 `DisplayName`、`LatestVersionID` 映射 Qoder 的 `display_title`、`latest_version`；字符串 source 解码到 `Source.Type` |
| Skill 版本下载 | 直接返回 ZIP 响应；Qoder 路径使用版本号，应传 `SkillVersion.Version`，而非 `SkillVersion.ID` |
| Agent 模型 | 兼容字符串和对象响应；增加 `ContextWindow`，保留 effort；显式设置 Qoder 不支持的 `Speed` 会返回错误 |
| Session / Deployment | 增加 `EnvironmentVariables`；Session 使用字符串 map，Deployment 使用可空字符串 |
| Memory 更新 | `Precondition.ContentSha256` 转换为 Qoder 顶层 `content_sha256`；保留 metadata |
| 模型目录 | 增加 Qoder 的 `Source`、`IsEnabled`、`Efforts`、上下文窗口等字段 |
| 其他 | 补齐 Qoder 的兼容分页游标、名称筛选、文件 metadata 和 Git resource 的 `Password` |

上游的 `WorkspaceID` 字段保留用于类型对齐，Qoder 的账户/空间范围由访问令牌决定，该字段不发送 workspace 请求头。上游保留的其他字段能否使用，仍由 Qoder 服务端的 API 定义决定。Qoder 已废弃的 Skills 参数可通过请求 option/extra fields 表达，未增加专用字段。

Go 参数中的 `Betas` 和 `QoderBeta` 用来指定 `x-qoder-beta` 协议版本，不是业务类型的 Beta 前缀。`Environments.Work.Poll` 的可选 `QoderWorkerID` 对应官网定义的 `Worker-ID` 请求头。

当前依赖范围没有云平台适配、配置文件/身份联合、本地工具执行器、Worker 框架或 MCP 执行适配。所需的 JSON、multipart、查询、分页、错误、SSE 工具及其内部依赖完整保留；具体边界见 [PACKAGE_SCOPE.md](PACKAGE_SCOPE.md)。

## 迁移和验证

这是一次公开 SDK 调用方式变更：旧的 `ManagedClient`、`AgentsAPI` interface、`CreateAgent(...).Execute()`、旧参数/响应类型和旧 SSE/PageIterator helper 被新的具体服务替代。调用方需按上面的服务树迁移。

业务源码和测试位于 `managed/`；共享结构和工具统一位于根目录 `convention/`，与 `managed/`、`forward/` 齐平。根模块依赖保持不变。Forward 也采用具体服务，复用公共 HTTP、参数、分页、SSE、错误和下载工具，见 [Forward README](../forward/README.md)。原生成器仍生成旧版结构，不应用它覆盖本分支的适配实现。

```sh
go test -race ./...
```

根模块验证包含 Managed 的 95 条官网路由契约、领域字段对齐快照，以及 Forward 的 110 条 API 契约和协议测试（请求/null/联合类型、multipart、临时下载链接、分页、SSE、错误、重试、取消、自定义 HTTP 客户端）。离线测试使用本地 HTTP transport 或 httptest，不访问真实账户或创建远端资源。

95 个 API 使用对应的 `agent_test.go`、`session_test.go` 等文件直接调用方法。原独立 `test/` 模块已移除；两种模式的测试与源码同目录，沿用上游的外部测试 package 写法。真实环境用例位于对应 API 的 `*_live_test.go`，需要 `-tags live` 和显式测试配置。SDK 通用用法见 [根目录指南](../README.md)。
