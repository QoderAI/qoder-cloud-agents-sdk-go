# SDK 包依赖边界

公共目录 `convention/` 与 `managed/`、`forward/` 齐平。原 `managed/convention` 和 `managed/internal` 合并到此目录。保留编码、参数、分页等子包，以维持原有符号及依赖关系，避免循环引用。

一期 Managed 范围是 95 个 HTTP API，以及这些 API 递归依赖的结构、枚举、联合类型和工具。服务为具体 struct。所需定义和工具完整保留，只进行 Qoder 命名转换及已有 Qoder HTTP 协议适配。

| 保留的包 | API 依赖原因 |
|---|---|
| `managed` | 客户端、95 个方法、参数、响应、领域实体和相关内容块 |
| `forward` | 客户端、109 个方法、参数、领域实体、业务配置及事件类型 |
| `convention` | HTTP、凭据、鉴权 transport、超时、重试、公共错误别名、响应及下载 |
| `convention/option` | 客户端、服务和方法的配置组合 |
| `convention/param` | 省略/null/实际值、联合类型、JSON override、extra fields |
| `convention/respjson` | 响应字段元信息和 RawJSON |
| `convention/constant` | 固定值类型及其默认序列化 |
| `convention/pagination` | 分页响应和自动遍历工具 |
| `convention/ssestream` | Session/Thread 事件流和恢复游标 |
| `convention/apierror` | 结构化错误 |
| `convention/apijson` | 领域类型、联合类型和字段元信息的 JSON 编解码 |
| `convention/apiquery` | 查询参数编码 |
| `convention/apiform` | Files/Skills 的 multipart 编码 |
| `convention/paramutil` | 领域联合类型的参数取值辅助函数 |
| `convention/encoding/json` 及其 `sentinel`、`shims` | 参数编码依赖的 omitzero、显式 null 和兼容工具 |
| `convention/thirdparty/{gjson,sjson,pretty,match}` | 上述完整 JSON 工具的直接或间接依赖，保留原许可证 |

未引入云平台适配、配置文件/身份联合、环境 Worker、模型 fallback、本地工具执行器、MCP 执行适配、文件同步、JSONL 流、上游测试基础设施及其他 API 服务。

`message.go` 只提供被这 95 个 API 引用的内容块等结构和工具，不包含 Messages API。名称转换后相同的结构只保留一个完全等价的声明；所有字段和方法均保留。

命名规则：业务类型及文件移除 Beta 前缀；提供方名称统一为 Qoder。协议用的 `Betas`、`QoderBeta` 和 `x-qoder-beta` 保留。

Forward 采用与 Managed 相同的具体服务和参数约定，覆盖 110 个 HTTP API，暂不提供 Service Account Token 管理接口。两者通过 `option.WithCredential` 使用同一个 Credential，并共享请求、重试、分页、错误、SSE、multipart 和下载实现。Forward 原有 `ForwardClient`、builder/Execute、参数/响应类型和专用 helper 已替换，调用迁移见 [Forward README](../forward/README.md)。Forward 的 Template、Identity、Schedule、Batch、Channel 等业务服务独立保留，不映射为 Managed API。
