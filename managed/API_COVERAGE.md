# 95 个 Managed API

基线：固定上游提交 `de6914c`。范围仅限官网 Conventions 与 Managed Mode。下表的服务均为具体 struct，方法均执行 HTTP 请求。完整路径前缀为 `/api/v1/cloud`。

排除 Forward Mode / Webhooks、8 个 Search、Deployment 子路径下的 2 个 Runs API、Session Cancel、已废弃的 Skill Update、OAuth start，以及官网此范围没有定义的 上游 Models.Get / Vaults.Update。Conventions 是共享协议，不计入 HTTP 操作数。

| 客户端资源入口 | API 数量 |
|---|---:|
| `Agents` | 6 |
| `Deployments` | 8 |
| `Models` | 1 |
| `Sessions` | 19 |
| `DeploymentRuns` | 2 |
| `Vaults` | 12 |
| `MemoryStores` | 14 |
| `Skills` | 9 |
| `Environments` | 14 |
| `Dreams` | 5 |
| `Files` | 5 |
| **合计** | **95** |

| SDK 方法 | HTTP 路由 | 官网定义 |
|---|---|---|
| `client.Agents.List` | `GET /agents` | [文档](https://docs.qoder.com/cloud-agents/api/agents/list) |
| `client.Agents.New` | `POST /agents` | [文档](https://docs.qoder.com/cloud-agents/api/agents/create) |
| `client.Agents.Get` | `GET /agents/{agent_id}` | [文档](https://docs.qoder.com/cloud-agents/api/agents/get) |
| `client.Agents.Update` | `POST /agents/{agent_id}` | [文档](https://docs.qoder.com/cloud-agents/api/agents/update) |
| `client.Agents.Archive` | `POST /agents/{agent_id}/archive` | [文档](https://docs.qoder.com/cloud-agents/api/agents/archive) |
| `client.Agents.Versions.List` | `GET /agents/{agent_id}/versions` | [文档](https://docs.qoder.com/cloud-agents/api/agents/list-versions) |
| `client.DeploymentRuns.List` | `GET /deployment_runs` | [文档](https://docs.qoder.com/cloud-agents/api/deployments/list-all-runs) |
| `client.DeploymentRuns.Get` | `GET /deployment_runs/{run_id}` | [文档](https://docs.qoder.com/cloud-agents/api/deployments/get-run-global) |
| `client.Deployments.List` | `GET /deployments` | [文档](https://docs.qoder.com/cloud-agents/api/deployments/list) |
| `client.Deployments.New` | `POST /deployments` | [文档](https://docs.qoder.com/cloud-agents/api/deployments/create) |
| `client.Deployments.Get` | `GET /deployments/{id}` | [文档](https://docs.qoder.com/cloud-agents/api/deployments/get) |
| `client.Deployments.Update` | `POST /deployments/{id}` | [文档](https://docs.qoder.com/cloud-agents/api/deployments/update) |
| `client.Deployments.Archive` | `POST /deployments/{id}/archive` | [文档](https://docs.qoder.com/cloud-agents/api/deployments/archive) |
| `client.Deployments.Pause` | `POST /deployments/{id}/pause` | [文档](https://docs.qoder.com/cloud-agents/api/deployments/pause) |
| `client.Deployments.Run` | `POST /deployments/{id}/run` | [文档](https://docs.qoder.com/cloud-agents/api/deployments/run) |
| `client.Deployments.Unpause` | `POST /deployments/{id}/unpause` | [文档](https://docs.qoder.com/cloud-agents/api/deployments/unpause) |
| `client.Dreams.List` | `GET /dreams` | [文档](https://docs.qoder.com/cloud-agents/api/dreams/list) |
| `client.Dreams.New` | `POST /dreams` | [文档](https://docs.qoder.com/cloud-agents/api/dreams/create) |
| `client.Dreams.Get` | `GET /dreams/{id}` | [文档](https://docs.qoder.com/cloud-agents/api/dreams/get) |
| `client.Dreams.Archive` | `POST /dreams/{id}/archive` | [文档](https://docs.qoder.com/cloud-agents/api/dreams/archive) |
| `client.Dreams.Cancel` | `POST /dreams/{id}/cancel` | [文档](https://docs.qoder.com/cloud-agents/api/dreams/cancel) |
| `client.Environments.List` | `GET /environments` | [文档](https://docs.qoder.com/cloud-agents/api/environments/list) |
| `client.Environments.New` | `POST /environments` | [文档](https://docs.qoder.com/cloud-agents/api/environments/create) |
| `client.Environments.Delete` | `DELETE /environments/{environment_id}` | [文档](https://docs.qoder.com/cloud-agents/api/environments/delete) |
| `client.Environments.Get` | `GET /environments/{environment_id}` | [文档](https://docs.qoder.com/cloud-agents/api/environments/get) |
| `client.Environments.Update` | `POST /environments/{environment_id}` | [文档](https://docs.qoder.com/cloud-agents/api/environments/update) |
| `client.Environments.Archive` | `POST /environments/{environment_id}/archive` | [文档](https://docs.qoder.com/cloud-agents/api/environments/archive) |
| `client.Environments.Work.List` | `GET /environments/{environment_id}/work` | [文档](https://docs.qoder.com/cloud-agents/api/environments/work/list) |
| `client.Environments.Work.Poll` | `GET /environments/{environment_id}/work/poll` | [文档](https://docs.qoder.com/cloud-agents/api/environments/work/poll) |
| `client.Environments.Work.Stats` | `GET /environments/{environment_id}/work/stats` | [文档](https://docs.qoder.com/cloud-agents/api/environments/work/stats) |
| `client.Environments.Work.Get` | `GET /environments/{environment_id}/work/{work_id}` | [文档](https://docs.qoder.com/cloud-agents/api/environments/work/get) |
| `client.Environments.Work.Update` | `POST /environments/{environment_id}/work/{work_id}` | [文档](https://docs.qoder.com/cloud-agents/api/environments/work/update-metadata) |
| `client.Environments.Work.Ack` | `POST /environments/{environment_id}/work/{work_id}/ack` | [文档](https://docs.qoder.com/cloud-agents/api/environments/work/ack) |
| `client.Environments.Work.Heartbeat` | `POST /environments/{environment_id}/work/{work_id}/heartbeat` | [文档](https://docs.qoder.com/cloud-agents/api/environments/work/heartbeat) |
| `client.Environments.Work.Stop` | `POST /environments/{environment_id}/work/{work_id}/stop` | [文档](https://docs.qoder.com/cloud-agents/api/environments/work/stop) |
| `client.Files.List` | `GET /files` | [文档](https://docs.qoder.com/cloud-agents/api/files/list) |
| `client.Files.Upload` | `POST /files` | [文档](https://docs.qoder.com/cloud-agents/api/files/upload) |
| `client.Files.Delete` | `DELETE /files/{file_id}` | [文档](https://docs.qoder.com/cloud-agents/api/files/delete) |
| `client.Files.GetMetadata` | `GET /files/{file_id}` | [文档](https://docs.qoder.com/cloud-agents/api/files/get) |
| `client.Files.Download` | `GET /files/{file_id}/content` | [文档](https://docs.qoder.com/cloud-agents/api/files/download) |
| `client.MemoryStores.List` | `GET /memory_stores` | [文档](https://docs.qoder.com/cloud-agents/api/memory-stores/list) |
| `client.MemoryStores.New` | `POST /memory_stores` | [文档](https://docs.qoder.com/cloud-agents/api/memory-stores/create) |
| `client.MemoryStores.Delete` | `DELETE /memory_stores/{memory_store_id}` | [文档](https://docs.qoder.com/cloud-agents/api/memory-stores/delete) |
| `client.MemoryStores.Get` | `GET /memory_stores/{memory_store_id}` | [文档](https://docs.qoder.com/cloud-agents/api/memory-stores/get) |
| `client.MemoryStores.Update` | `POST /memory_stores/{memory_store_id}` | [文档](https://docs.qoder.com/cloud-agents/api/memory-stores/update) |
| `client.MemoryStores.Archive` | `POST /memory_stores/{memory_store_id}/archive` | [文档](https://docs.qoder.com/cloud-agents/api/memory-stores/archive) |
| `client.MemoryStores.Memories.List` | `GET /memory_stores/{memory_store_id}/memories` | [文档](https://docs.qoder.com/cloud-agents/api/memory-stores/list-memories) |
| `client.MemoryStores.Memories.New` | `POST /memory_stores/{memory_store_id}/memories` | [文档](https://docs.qoder.com/cloud-agents/api/memory-stores/create-memory) |
| `client.MemoryStores.Memories.Delete` | `DELETE /memory_stores/{memory_store_id}/memories/{memory_id}` | [文档](https://docs.qoder.com/cloud-agents/api/memory-stores/delete-memory) |
| `client.MemoryStores.Memories.Get` | `GET /memory_stores/{memory_store_id}/memories/{memory_id}` | [文档](https://docs.qoder.com/cloud-agents/api/memory-stores/get-memory) |
| `client.MemoryStores.Memories.Update` | `POST /memory_stores/{memory_store_id}/memories/{memory_id}` | [文档](https://docs.qoder.com/cloud-agents/api/memory-stores/update-memory) |
| `client.MemoryStores.MemoryVersions.List` | `GET /memory_stores/{memory_store_id}/memory_versions` | [文档](https://docs.qoder.com/cloud-agents/api/memory-stores/list-versions) |
| `client.MemoryStores.MemoryVersions.Get` | `GET /memory_stores/{memory_store_id}/memory_versions/{memory_version_id}` | [文档](https://docs.qoder.com/cloud-agents/api/memory-stores/get-version) |
| `client.MemoryStores.MemoryVersions.Redact` | `POST /memory_stores/{memory_store_id}/memory_versions/{memory_version_id}/redact` | [文档](https://docs.qoder.com/cloud-agents/api/memory-stores/redact-version) |
| `client.Models.List` | `GET /models` | [文档](https://docs.qoder.com/cloud-agents/api/models/list) |
| `client.Sessions.List` | `GET /sessions` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/list) |
| `client.Sessions.New` | `POST /sessions` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/create) |
| `client.Sessions.Delete` | `DELETE /sessions/{session_id}` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/delete) |
| `client.Sessions.Get` | `GET /sessions/{session_id}` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/get) |
| `client.Sessions.Update` | `POST /sessions/{session_id}` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/update) |
| `client.Sessions.Archive` | `POST /sessions/{session_id}/archive` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/archive) |
| `client.Sessions.Events.List` | `GET /sessions/{session_id}/events` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/list-events) |
| `client.Sessions.Events.Send` | `POST /sessions/{session_id}/events` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/send-event) |
| `client.Sessions.Events.StreamEvents` | `GET /sessions/{session_id}/events/stream` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/stream-events) |
| `client.Sessions.Resources.List` | `GET /sessions/{session_id}/resources` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/list-resources) |
| `client.Sessions.Resources.Add` | `POST /sessions/{session_id}/resources` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/add-resource) |
| `client.Sessions.Resources.Delete` | `DELETE /sessions/{session_id}/resources/{resource_id}` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/delete-resource) |
| `client.Sessions.Resources.Get` | `GET /sessions/{session_id}/resources/{resource_id}` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/get-resource) |
| `client.Sessions.Resources.Update` | `POST /sessions/{session_id}/resources/{resource_id}` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/update-resource) |
| `client.Sessions.Threads.List` | `GET /sessions/{session_id}/threads` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/list-threads) |
| `client.Sessions.Threads.Get` | `GET /sessions/{session_id}/threads/{thread_id}` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/get-thread) |
| `client.Sessions.Threads.Archive` | `POST /sessions/{session_id}/threads/{thread_id}/archive` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/archive-thread) |
| `client.Sessions.Threads.Events.List` | `GET /sessions/{session_id}/threads/{thread_id}/events` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/list-thread-events) |
| `client.Sessions.Threads.Events.StreamEvents` | `GET /sessions/{session_id}/threads/{thread_id}/stream` | [文档](https://docs.qoder.com/cloud-agents/api/sessions/stream-thread-events) |
| `client.Skills.List` | `GET /skills` | [文档](https://docs.qoder.com/cloud-agents/api/skills/list) |
| `client.Skills.New` | `POST /skills` | [文档](https://docs.qoder.com/cloud-agents/api/skills/create) |
| `client.Skills.Delete` | `DELETE /skills/{skill_id}` | [文档](https://docs.qoder.com/cloud-agents/api/skills/delete) |
| `client.Skills.Get` | `GET /skills/{skill_id}` | [文档](https://docs.qoder.com/cloud-agents/api/skills/get) |
| `client.Skills.Versions.List` | `GET /skills/{skill_id}/versions` | [文档](https://docs.qoder.com/cloud-agents/api/skills/list-versions) |
| `client.Skills.Versions.New` | `POST /skills/{skill_id}/versions` | [文档](https://docs.qoder.com/cloud-agents/api/skills/create-version) |
| `client.Skills.Versions.Delete` | `DELETE /skills/{skill_id}/versions/{version}` | [文档](https://docs.qoder.com/cloud-agents/api/skills/delete-version) |
| `client.Skills.Versions.Get` | `GET /skills/{skill_id}/versions/{version}` | [文档](https://docs.qoder.com/cloud-agents/api/skills/get-version) |
| `client.Skills.Versions.Download` | `GET /skills/{skill_id}/versions/{version}/content` | [文档](https://docs.qoder.com/cloud-agents/api/skills/get-version-content) |
| `client.Vaults.List` | `GET /vaults` | [文档](https://docs.qoder.com/cloud-agents/api/vaults/list) |
| `client.Vaults.New` | `POST /vaults` | [文档](https://docs.qoder.com/cloud-agents/api/vaults/create) |
| `client.Vaults.Delete` | `DELETE /vaults/{vault_id}` | [文档](https://docs.qoder.com/cloud-agents/api/vaults/delete) |
| `client.Vaults.Get` | `GET /vaults/{vault_id}` | [文档](https://docs.qoder.com/cloud-agents/api/vaults/get) |
| `client.Vaults.Archive` | `POST /vaults/{vault_id}/archive` | [文档](https://docs.qoder.com/cloud-agents/api/vaults/archive) |
| `client.Vaults.Credentials.List` | `GET /vaults/{vault_id}/credentials` | [文档](https://docs.qoder.com/cloud-agents/api/vaults/list-credentials) |
| `client.Vaults.Credentials.New` | `POST /vaults/{vault_id}/credentials` | [文档](https://docs.qoder.com/cloud-agents/api/vaults/create-credential) |
| `client.Vaults.Credentials.Delete` | `DELETE /vaults/{vault_id}/credentials/{credential_id}` | [文档](https://docs.qoder.com/cloud-agents/api/vaults/delete-credential) |
| `client.Vaults.Credentials.Get` | `GET /vaults/{vault_id}/credentials/{credential_id}` | [文档](https://docs.qoder.com/cloud-agents/api/vaults/get-credential) |
| `client.Vaults.Credentials.Update` | `POST /vaults/{vault_id}/credentials/{credential_id}` | [文档](https://docs.qoder.com/cloud-agents/api/vaults/update-credential) |
| `client.Vaults.Credentials.Archive` | `POST /vaults/{vault_id}/credentials/{credential_id}/archive` | [文档](https://docs.qoder.com/cloud-agents/api/vaults/archive-credential) |
| `client.Vaults.Credentials.MCPOAuthValidate` | `POST /vaults/{vault_id}/credentials/{credential_id}/mcp_oauth_validate` | [文档](https://docs.qoder.com/cloud-agents/api/vaults/validate-credential) |

类型字段快照来自固定上游提交，并按 Qoder 命名统一转换。`testdata/api-contracts.json` 保存 95 条路由及官网响应示例，`api_contract_test.go` 同时检查不存在超出清单的服务方法。
