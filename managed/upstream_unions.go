package managed

import (
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
)

func init() {
	apijson.RegisterUnion[BrowserStateChangeUnionParam](
		"type",
		apijson.Discriminator[BrowserStateChangeTabOpenedParam]("tab_opened"),
		apijson.Discriminator[BrowserStateChangeDownloadStartedParam]("download_started"),
		apijson.Discriminator[BrowserStateChangeDownloadCompletedParam]("download_completed"),
		apijson.Discriminator[BrowserStateChangeDownloadFailedParam]("download_failed"),
	)
}
func init() {
	apijson.RegisterUnion[ContentBlockParamUnion](
		"type",
		apijson.Discriminator[TextBlockParam]("text"),
		apijson.Discriminator[ImageBlockParam]("image"),
		apijson.Discriminator[RequestDocumentBlockParam]("document"),
		apijson.Discriminator[SearchResultBlockParam]("search_result"),
		apijson.Discriminator[ThinkingBlockParam]("thinking"),
		apijson.Discriminator[RedactedThinkingBlockParam]("redacted_thinking"),
		apijson.Discriminator[ToolUseBlockParam]("tool_use"),
		apijson.Discriminator[ToolResultBlockParam]("tool_result"),
		apijson.Discriminator[ServerToolUseBlockParam]("server_tool_use"),
		apijson.Discriminator[WebSearchToolResultBlockParam]("web_search_tool_result"),
		apijson.Discriminator[WebFetchToolResultBlockParam]("web_fetch_tool_result"),
		apijson.Discriminator[AdvisorToolResultBlockParam]("advisor_tool_result"),
		apijson.Discriminator[CodeExecutionToolResultBlockParam]("code_execution_tool_result"),
		apijson.Discriminator[BashCodeExecutionToolResultBlockParam]("bash_code_execution_tool_result"),
		apijson.Discriminator[TextEditorCodeExecutionToolResultBlockParam]("text_editor_code_execution_tool_result"),
		apijson.Discriminator[ToolSearchToolResultBlockParam]("tool_search_tool_result"),
		apijson.Discriminator[MCPToolUseBlockParam]("mcp_tool_use"),
		apijson.Discriminator[RequestMCPToolResultBlockParam]("mcp_tool_result"),
		apijson.Discriminator[ContainerUploadBlockParam]("container_upload"),
		apijson.Discriminator[CompactionBlockParam]("compaction"),
		apijson.Discriminator[RequestToolAdditionBlockParam]("tool_addition"),
		apijson.Discriminator[RequestToolRemovalBlockParam]("tool_removal"),
		apijson.Discriminator[FallbackBlockParam]("fallback"),
	)
}
func init() {
	apijson.RegisterUnion[ImageBlockParamSourceUnion](
		"type",
		apijson.Discriminator[Base64ImageSourceParam]("base64"),
		apijson.Discriminator[URLImageSourceParam]("url"),
		apijson.Discriminator[FileImageSourceParam]("file"),
	)
}
func init() {
	apijson.RegisterUnion[RequestDocumentBlockSourceUnionParam](
		"type",
		apijson.Discriminator[Base64PDFSourceParam]("base64"),
		apijson.Discriminator[PlainTextSourceParam]("text"),
		apijson.Discriminator[ContentBlockSourceParam]("content"),
		apijson.Discriminator[URLPDFSourceParam]("url"),
		apijson.Discriminator[FileDocumentSourceParam]("file"),
	)
}
func init() {
	apijson.RegisterUnion[RequestToolAdditionBlockToolUnionParam](
		"type",
		apijson.Discriminator[ToolChangeToolReferenceParam]("tool_reference"),
		apijson.Discriminator[ToolChangeMCPToolReferenceParam]("mcp_tool_reference"),
		apijson.Discriminator[ToolChangeMCPToolsetReferenceParam]("mcp_toolset_reference"),
	)
}
func init() {
	apijson.RegisterUnion[RequestToolRemovalBlockToolUnionParam](
		"type",
		apijson.Discriminator[ToolChangeToolReferenceParam]("tool_reference"),
		apijson.Discriminator[ToolChangeMCPToolReferenceParam]("mcp_tool_reference"),
		apijson.Discriminator[ToolChangeMCPToolsetReferenceParam]("mcp_toolset_reference"),
	)
}
func init() {
	apijson.RegisterUnion[ServerToolUseBlockParamCallerUnion](
		"type",
		apijson.Discriminator[DirectCallerParam]("direct"),
		apijson.Discriminator[ServerToolCallerParam]("code_execution_20250825"),
		apijson.Discriminator[ServerToolCaller20260120Param]("code_execution_20260120"),
	)
}
func init() {
	apijson.RegisterUnion[TextCitationParamUnion](
		"type",
		apijson.Discriminator[CitationCharLocationParam]("char_location"),
		apijson.Discriminator[CitationPageLocationParam]("page_location"),
		apijson.Discriminator[CitationContentBlockLocationParam]("content_block_location"),
		apijson.Discriminator[CitationWebSearchResultLocationParam]("web_search_result_location"),
		apijson.Discriminator[CitationSearchResultLocationParam]("search_result_location"),
	)
}
func init() {
	apijson.RegisterUnion[ToolResultBlockParamContentUnion](
		"type",
		apijson.Discriminator[TextBlockParam]("text"),
		apijson.Discriminator[ImageBlockParam]("image"),
		apijson.Discriminator[SearchResultBlockParam]("search_result"),
		apijson.Discriminator[RequestDocumentBlockParam]("document"),
		apijson.Discriminator[ToolReferenceBlockParam]("tool_reference"),
		apijson.Discriminator[BrowserStateBlockParam]("browser_state"),
	)
}
func init() {
	apijson.RegisterUnion[ToolUseBlockParamCallerUnion](
		"type",
		apijson.Discriminator[DirectCallerParam]("direct"),
		apijson.Discriminator[ServerToolCallerParam]("code_execution_20250825"),
		apijson.Discriminator[ServerToolCaller20260120Param]("code_execution_20260120"),
	)
}
func init() {
	apijson.RegisterUnion[WebFetchToolResultBlockParamCallerUnion](
		"type",
		apijson.Discriminator[DirectCallerParam]("direct"),
		apijson.Discriminator[ServerToolCallerParam]("code_execution_20250825"),
		apijson.Discriminator[ServerToolCaller20260120Param]("code_execution_20260120"),
	)
}
func init() {
	apijson.RegisterUnion[WebSearchToolResultBlockParamCallerUnion](
		"type",
		apijson.Discriminator[DirectCallerParam]("direct"),
		apijson.Discriminator[ServerToolCallerParam]("code_execution_20250825"),
		apijson.Discriminator[ServerToolCaller20260120Param]("code_execution_20260120"),
	)
}
