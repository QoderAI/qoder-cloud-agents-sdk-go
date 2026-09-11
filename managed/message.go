// Qoder managed API definitions.
package managed

import (
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/paramutil"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/constant"
)

// The properties EncryptedContent, Type are required.
type AdvisorRedactedResultBlockParam struct {
	// Opaque blob produced by a prior response; must be round-tripped verbatim.
	EncryptedContent string            `json:"encrypted_content" api:"required"`
	StopReason       param.Opt[string] `json:"stop_reason,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "advisor_redacted_result".
	Type constant.AdvisorRedactedResult `json:"type" default:"advisor_redacted_result"`
	paramObj
}

func (r AdvisorRedactedResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow AdvisorRedactedResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdvisorRedactedResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Text, Type are required.
type AdvisorResultBlockParam struct {
	Text       string            `json:"text" api:"required"`
	StopReason param.Opt[string] `json:"stop_reason,omitzero"`
	// This field can be elided, and will marshal its zero value as "advisor_result".
	Type constant.AdvisorResult `json:"type" default:"advisor_result"`
	paramObj
}

func (r AdvisorResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow AdvisorResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdvisorResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Content, ToolUseID, Type are required.
type AdvisorToolResultBlockParam struct {
	Content   AdvisorToolResultBlockParamContentUnion `json:"content,omitzero" api:"required"`
	ToolUseID string                                  `json:"tool_use_id" api:"required"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "advisor_tool_result".
	Type constant.AdvisorToolResult `json:"type" default:"advisor_tool_result"`
	paramObj
}

func (r AdvisorToolResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow AdvisorToolResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdvisorToolResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AdvisorToolResultBlockParamContentUnion struct {
	OfRequestAdvisorToolResultError     *AdvisorToolResultErrorParam     `json:",omitzero,inline"`
	OfRequestAdvisorResultBlock         *AdvisorResultBlockParam         `json:",omitzero,inline"`
	OfRequestAdvisorRedactedResultBlock *AdvisorRedactedResultBlockParam `json:",omitzero,inline"`
	paramUnion
}

func (u AdvisorToolResultBlockParamContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfRequestAdvisorToolResultError, u.OfRequestAdvisorResultBlock, u.OfRequestAdvisorRedactedResultBlock)
}
func (u *AdvisorToolResultBlockParamContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *AdvisorToolResultBlockParamContentUnion) asAny() any {
	if !param.IsOmitted(u.OfRequestAdvisorToolResultError) {
		return u.OfRequestAdvisorToolResultError
	} else if !param.IsOmitted(u.OfRequestAdvisorResultBlock) {
		return u.OfRequestAdvisorResultBlock
	} else if !param.IsOmitted(u.OfRequestAdvisorRedactedResultBlock) {
		return u.OfRequestAdvisorRedactedResultBlock
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdvisorToolResultBlockParamContentUnion) GetErrorCode() *string {
	if vt := u.OfRequestAdvisorToolResultError; vt != nil {
		return (*string)(&vt.ErrorCode)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdvisorToolResultBlockParamContentUnion) GetText() *string {
	if vt := u.OfRequestAdvisorResultBlock; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdvisorToolResultBlockParamContentUnion) GetEncryptedContent() *string {
	if vt := u.OfRequestAdvisorRedactedResultBlock; vt != nil {
		return &vt.EncryptedContent
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdvisorToolResultBlockParamContentUnion) GetType() *string {
	if vt := u.OfRequestAdvisorToolResultError; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestAdvisorResultBlock; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestAdvisorRedactedResultBlock; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AdvisorToolResultBlockParamContentUnion) GetStopReason() *string {
	if vt := u.OfRequestAdvisorResultBlock; vt != nil && vt.StopReason.Valid() {
		return &vt.StopReason.Value
	} else if vt := u.OfRequestAdvisorRedactedResultBlock; vt != nil && vt.StopReason.Valid() {
		return &vt.StopReason.Value
	}
	return nil
}

// The properties ErrorCode, Type are required.
type AdvisorToolResultErrorParam struct {
	// Any of "max_uses_exceeded", "prompt_too_long", "too_many_requests",
	// "overloaded", "unavailable", "execution_time_exceeded", "model_not_found".
	ErrorCode AdvisorToolResultErrorParamErrorCode `json:"error_code,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "advisor_tool_result_error".
	Type constant.AdvisorToolResultError `json:"type" default:"advisor_tool_result_error"`
	paramObj
}

func (r AdvisorToolResultErrorParam) MarshalJSON() (data []byte, err error) {
	type shadow AdvisorToolResultErrorParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AdvisorToolResultErrorParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AdvisorToolResultErrorParamErrorCode string

// The properties Data, MediaType, Type are required.
type Base64ImageSourceParam struct {
	Data string `json:"data" api:"required" format:"byte"`
	// Any of "image/jpeg", "image/png", "image/gif", "image/webp".
	MediaType Base64ImageSourceMediaType `json:"media_type,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "base64".
	Type constant.Base64 `json:"type" default:"base64"`
	paramObj
}

func (r Base64ImageSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow Base64ImageSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *Base64ImageSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type Base64ImageSourceMediaType string

// The properties Data, MediaType, Type are required.
type Base64PDFSourceParam struct {
	Data string `json:"data" api:"required" format:"byte"`
	// This field can be elided, and will marshal its zero value as "application/pdf".
	MediaType constant.ApplicationPDF `json:"media_type" default:"application/pdf"`
	// This field can be elided, and will marshal its zero value as "base64".
	Type constant.Base64 `json:"type" default:"base64"`
	paramObj
}

func (r Base64PDFSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow Base64PDFSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *Base64PDFSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties FileID, Type are required.
type BashCodeExecutionOutputBlockParam struct {
	FileID string `json:"file_id" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "bash_code_execution_output".
	Type constant.BashCodeExecutionOutput `json:"type" default:"bash_code_execution_output"`
	paramObj
}

func (r BashCodeExecutionOutputBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow BashCodeExecutionOutputBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BashCodeExecutionOutputBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Content, ReturnCode, Stderr, Stdout, Type are required.
type BashCodeExecutionResultBlockParam struct {
	Content    []BashCodeExecutionOutputBlockParam `json:"content,omitzero" api:"required"`
	ReturnCode int64                               `json:"return_code" api:"required"`
	Stderr     string                              `json:"stderr" api:"required"`
	Stdout     string                              `json:"stdout" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "bash_code_execution_result".
	Type constant.BashCodeExecutionResult `json:"type" default:"bash_code_execution_result"`
	paramObj
}

func (r BashCodeExecutionResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow BashCodeExecutionResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BashCodeExecutionResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Content, ToolUseID, Type are required.
type BashCodeExecutionToolResultBlockParam struct {
	Content   BashCodeExecutionToolResultBlockParamContentUnion `json:"content,omitzero" api:"required"`
	ToolUseID string                                            `json:"tool_use_id" api:"required"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "bash_code_execution_tool_result".
	Type constant.BashCodeExecutionToolResult `json:"type" default:"bash_code_execution_tool_result"`
	paramObj
}

func (r BashCodeExecutionToolResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow BashCodeExecutionToolResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BashCodeExecutionToolResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BashCodeExecutionToolResultBlockParamContentUnion struct {
	OfRequestBashCodeExecutionToolResultError *BashCodeExecutionToolResultErrorParam `json:",omitzero,inline"`
	OfRequestBashCodeExecutionResultBlock     *BashCodeExecutionResultBlockParam     `json:",omitzero,inline"`
	paramUnion
}

func (u BashCodeExecutionToolResultBlockParamContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfRequestBashCodeExecutionToolResultError, u.OfRequestBashCodeExecutionResultBlock)
}
func (u *BashCodeExecutionToolResultBlockParamContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BashCodeExecutionToolResultBlockParamContentUnion) asAny() any {
	if !param.IsOmitted(u.OfRequestBashCodeExecutionToolResultError) {
		return u.OfRequestBashCodeExecutionToolResultError
	} else if !param.IsOmitted(u.OfRequestBashCodeExecutionResultBlock) {
		return u.OfRequestBashCodeExecutionResultBlock
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BashCodeExecutionToolResultBlockParamContentUnion) GetErrorCode() *string {
	if vt := u.OfRequestBashCodeExecutionToolResultError; vt != nil {
		return (*string)(&vt.ErrorCode)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BashCodeExecutionToolResultBlockParamContentUnion) GetContent() []BashCodeExecutionOutputBlockParam {
	if vt := u.OfRequestBashCodeExecutionResultBlock; vt != nil {
		return vt.Content
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BashCodeExecutionToolResultBlockParamContentUnion) GetReturnCode() *int64 {
	if vt := u.OfRequestBashCodeExecutionResultBlock; vt != nil {
		return &vt.ReturnCode
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BashCodeExecutionToolResultBlockParamContentUnion) GetStderr() *string {
	if vt := u.OfRequestBashCodeExecutionResultBlock; vt != nil {
		return &vt.Stderr
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BashCodeExecutionToolResultBlockParamContentUnion) GetStdout() *string {
	if vt := u.OfRequestBashCodeExecutionResultBlock; vt != nil {
		return &vt.Stdout
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BashCodeExecutionToolResultBlockParamContentUnion) GetType() *string {
	if vt := u.OfRequestBashCodeExecutionToolResultError; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestBashCodeExecutionResultBlock; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// The properties ErrorCode, Type are required.
type BashCodeExecutionToolResultErrorParam struct {
	// Any of "invalid_tool_input", "unavailable", "too_many_requests",
	// "execution_time_exceeded", "output_file_too_large".
	ErrorCode BashCodeExecutionToolResultErrorParamErrorCode `json:"error_code,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "bash_code_execution_tool_result_error".
	Type constant.BashCodeExecutionToolResultError `json:"type" default:"bash_code_execution_tool_result_error"`
	paramObj
}

func (r BashCodeExecutionToolResultErrorParam) MarshalJSON() (data []byte, err error) {
	type shadow BashCodeExecutionToolResultErrorParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BashCodeExecutionToolResultErrorParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BashCodeExecutionToolResultErrorParamErrorCode string

// The caller's browser state after a browser toolset member call — the full
// inventory of open tabs, which tab is active, and any side effects (tabs opened,
// download state changes) the call produced.
//
// At most one per `tool_result`, only on a non-error result answering a browser
// toolset member `tool_use`. The server renders the model-visible text from it;
// the model never sees the raw fields.
//
// The properties Tabs, Type are required.
type BrowserStateBlockParam struct {
	// All tabs open in the browser after this call — the full inventory, not a delta.
	// May be empty. Whenever non-empty, exactly one entry carries `active: true`.
	Tabs []BrowserStateTabEntryParam `json:"tabs,omitzero" api:"required"`
	// Tabs opened and download state changes during this call. "Nothing to report" is
	// expressed by omitting the field, never by an empty list.
	StateChanges []BrowserStateChangeUnionParam `json:"state_changes,omitzero"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as "browser_state".
	Type constant.BrowserState `json:"type" default:"browser_state"`
	paramObj
}

func (r BrowserStateBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow BrowserStateBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserStateBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BrowserStateChangeUnionParam struct {
	OfTabOpened         *BrowserStateChangeTabOpenedParam         `json:",omitzero,inline"`
	OfDownloadStarted   *BrowserStateChangeDownloadStartedParam   `json:",omitzero,inline"`
	OfDownloadCompleted *BrowserStateChangeDownloadCompletedParam `json:",omitzero,inline"`
	OfDownloadFailed    *BrowserStateChangeDownloadFailedParam    `json:",omitzero,inline"`
	paramUnion
}

func (u BrowserStateChangeUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfTabOpened, u.OfDownloadStarted, u.OfDownloadCompleted, u.OfDownloadFailed)
}
func (u *BrowserStateChangeUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BrowserStateChangeUnionParam) asAny() any {
	if !param.IsOmitted(u.OfTabOpened) {
		return u.OfTabOpened
	} else if !param.IsOmitted(u.OfDownloadStarted) {
		return u.OfDownloadStarted
	} else if !param.IsOmitted(u.OfDownloadCompleted) {
		return u.OfDownloadCompleted
	} else if !param.IsOmitted(u.OfDownloadFailed) {
		return u.OfDownloadFailed
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BrowserStateChangeUnionParam) GetTabID() *string {
	if vt := u.OfTabOpened; vt != nil {
		return &vt.TabID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BrowserStateChangeUnionParam) GetPath() *string {
	if vt := u.OfDownloadCompleted; vt != nil && vt.Path.Valid() {
		return &vt.Path.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BrowserStateChangeUnionParam) GetSizeBytes() *int64 {
	if vt := u.OfDownloadCompleted; vt != nil && vt.SizeBytes.Valid() {
		return &vt.SizeBytes.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BrowserStateChangeUnionParam) GetError() *string {
	if vt := u.OfDownloadFailed; vt != nil && vt.Error.Valid() {
		return &vt.Error.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BrowserStateChangeUnionParam) GetType() *string {
	if vt := u.OfTabOpened; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDownloadStarted; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDownloadCompleted; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDownloadFailed; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BrowserStateChangeUnionParam) GetDownloadID() *string {
	if vt := u.OfDownloadStarted; vt != nil {
		return (*string)(&vt.DownloadID)
	} else if vt := u.OfDownloadCompleted; vt != nil {
		return (*string)(&vt.DownloadID)
	} else if vt := u.OfDownloadFailed; vt != nil {
		return (*string)(&vt.DownloadID)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BrowserStateChangeUnionParam) GetURL() *string {
	if vt := u.OfDownloadStarted; vt != nil {
		return (*string)(&vt.URL)
	} else if vt := u.OfDownloadCompleted; vt != nil {
		return (*string)(&vt.URL)
	} else if vt := u.OfDownloadFailed; vt != nil {
		return (*string)(&vt.URL)
	}
	return nil
}

// A file download that finished during this call, reported with the same
// `download_id` as its `download_started` — or without a prior `download_started`,
// when the download finished during the call that started it (at most one state
// change per `download_id` per result).
//
// The properties DownloadID, Type, URL are required.
type BrowserStateChangeDownloadCompletedParam struct {
	// The caller-assigned identifier for this download, stable across the state
	// changes reporting it.
	DownloadID string `json:"download_id" api:"required"`
	// The final post-redirect URL the download was served from.
	URL string `json:"url" api:"required"`
	// Where the executor saved the file, on the executor's filesystem. Only included
	// when another tool in the same environment can read the file at that path.
	Path param.Opt[string] `json:"path,omitzero"`
	// The completed download's size.
	SizeBytes param.Opt[int64] `json:"size_bytes,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "download_completed".
	Type constant.DownloadCompleted `json:"type" default:"download_completed"`
	paramObj
}

func (r BrowserStateChangeDownloadCompletedParam) MarshalJSON() (data []byte, err error) {
	type shadow BrowserStateChangeDownloadCompletedParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserStateChangeDownloadCompletedParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A file download that failed — or was cancelled — during this call.
//
// The properties DownloadID, Type, URL are required.
type BrowserStateChangeDownloadFailedParam struct {
	// The caller-assigned identifier for this download, stable across the state
	// changes reporting it.
	DownloadID string `json:"download_id" api:"required"`
	// The final post-redirect URL the download was served from.
	URL string `json:"url" api:"required"`
	// The failure or cancellation detail, when known.
	Error param.Opt[string] `json:"error,omitzero"`
	// This field can be elided, and will marshal its zero value as "download_failed".
	Type constant.DownloadFailed `json:"type" default:"download_failed"`
	paramObj
}

func (r BrowserStateChangeDownloadFailedParam) MarshalJSON() (data []byte, err error) {
	type shadow BrowserStateChangeDownloadFailedParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserStateChangeDownloadFailedParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A file download that started during this call.
//
// The properties DownloadID, Type, URL are required.
type BrowserStateChangeDownloadStartedParam struct {
	// The caller-assigned identifier for this download, stable across the state
	// changes reporting it.
	DownloadID string `json:"download_id" api:"required"`
	// The final post-redirect URL the download was served from.
	URL string `json:"url" api:"required"`
	// This field can be elided, and will marshal its zero value as "download_started".
	Type constant.DownloadStarted `json:"type" default:"download_started"`
	paramObj
}

func (r BrowserStateChangeDownloadStartedParam) MarshalJSON() (data []byte, err error) {
	type shadow BrowserStateChangeDownloadStartedParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserStateChangeDownloadStartedParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A tab this call's execution opened that remains open at its end — the creation
// delta of the `tabs` inventory, not an event log.
//
// Carries only the `tab_id`; the tab's `title` and `url` live on its `tabs` entry,
// which must include the same `tab_id`. A tab opened during a failed call gets no
// deferred `tab_opened`; it simply appears in the next result's `tabs` inventory.
//
// The properties TabID, Type are required.
type BrowserStateChangeTabOpenedParam struct {
	// The `tab_id` of the opened tab, present in `tabs`.
	TabID string `json:"tab_id" api:"required"`
	// This field can be elided, and will marshal its zero value as "tab_opened".
	Type constant.TabOpened `json:"type" default:"tab_opened"`
	paramObj
}

func (r BrowserStateChangeTabOpenedParam) MarshalJSON() (data []byte, err error) {
	type shadow BrowserStateChangeTabOpenedParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserStateChangeTabOpenedParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// One open browser tab reported in a `browser_state` block's `tabs` inventory.
//
// `tab_id` is the caller-assigned identifier for the tab; `title` and `url`
// describe the page the tab is currently showing and may be empty strings (a blank
// tab legitimately has both empty). `active` marks the tab that is active after
// this call; whenever `tabs` is non-empty, exactly one entry is marked.
//
// The properties TabID, Title, URL are required.
type BrowserStateTabEntryParam struct {
	// The caller-assigned identifier for this tab, unique within the inventory.
	TabID string `json:"tab_id" api:"required"`
	// The title of the page the tab is showing. May be empty.
	Title string `json:"title" api:"required"`
	// The URL of the page the tab is showing. May be empty.
	URL string `json:"url" api:"required"`
	// Whether this tab is the active tab after this call. Whenever `tabs` is
	// non-empty, exactly one entry is marked `active: true`.
	Active param.Opt[bool] `json:"active,omitzero"`
	paramObj
}

func (r BrowserStateTabEntryParam) MarshalJSON() (data []byte, err error) {
	type shadow BrowserStateTabEntryParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BrowserStateTabEntryParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// This struct has a constant value, construct it with
// [NewCacheControlEphemeralParam].
type CacheControlEphemeralParam struct {
	// The time-to-live for the cache control breakpoint.
	//
	// This may be one the following values:
	//
	// - `5m`: 5 minutes
	// - `1h`: 1 hour
	//
	// Defaults to `5m`.
	//
	// Any of "5m", "1h".
	TTL  CacheControlEphemeralTTL `json:"ttl,omitzero"`
	Type constant.Ephemeral       `json:"type" default:"ephemeral"`
	paramObj
}

func (r CacheControlEphemeralParam) MarshalJSON() (data []byte, err error) {
	type shadow CacheControlEphemeralParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CacheControlEphemeralParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The time-to-live for the cache control breakpoint.
//
// This may be one the following values:
//
// - `5m`: 5 minutes
// - `1h`: 1 hour
//
// Defaults to `5m`.
type CacheControlEphemeralTTL string

// The properties CitedText, DocumentIndex, DocumentTitle, EndCharIndex,
// StartCharIndex, Type are required.
type CitationCharLocationParam struct {
	DocumentTitle  param.Opt[string] `json:"document_title,omitzero" api:"required"`
	CitedText      string            `json:"cited_text" api:"required"`
	DocumentIndex  int64             `json:"document_index" api:"required"`
	EndCharIndex   int64             `json:"end_char_index" api:"required"`
	StartCharIndex int64             `json:"start_char_index" api:"required"`
	// This field can be elided, and will marshal its zero value as "char_location".
	Type constant.CharLocation `json:"type" default:"char_location"`
	paramObj
}

func (r CitationCharLocationParam) MarshalJSON() (data []byte, err error) {
	type shadow CitationCharLocationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CitationCharLocationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties CitedText, DocumentIndex, DocumentTitle, EndBlockIndex,
// StartBlockIndex, Type are required.
type CitationContentBlockLocationParam struct {
	DocumentTitle param.Opt[string] `json:"document_title,omitzero" api:"required"`
	// The full text of the cited block range, concatenated.
	//
	// Always equals the contents of `content[start_block_index:end_block_index]`
	// joined together. The text block is the minimal citable unit; this field is never
	// a substring of a single block. Not counted toward output tokens, and not counted
	// toward input tokens when sent back in subsequent turns.
	CitedText     string `json:"cited_text" api:"required"`
	DocumentIndex int64  `json:"document_index" api:"required"`
	// Exclusive 0-based end index of the cited block range in the source's `content`
	// array.
	//
	// Always greater than `start_block_index`; a single-block citation has
	// `end_block_index = start_block_index + 1`.
	EndBlockIndex int64 `json:"end_block_index" api:"required"`
	// 0-based index of the first cited block in the source's `content` array.
	StartBlockIndex int64 `json:"start_block_index" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "content_block_location".
	Type constant.ContentBlockLocation `json:"type" default:"content_block_location"`
	paramObj
}

func (r CitationContentBlockLocationParam) MarshalJSON() (data []byte, err error) {
	type shadow CitationContentBlockLocationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CitationContentBlockLocationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties CitedText, DocumentIndex, DocumentTitle, EndPageNumber,
// StartPageNumber, Type are required.
type CitationPageLocationParam struct {
	DocumentTitle   param.Opt[string] `json:"document_title,omitzero" api:"required"`
	CitedText       string            `json:"cited_text" api:"required"`
	DocumentIndex   int64             `json:"document_index" api:"required"`
	EndPageNumber   int64             `json:"end_page_number" api:"required"`
	StartPageNumber int64             `json:"start_page_number" api:"required"`
	// This field can be elided, and will marshal its zero value as "page_location".
	Type constant.PageLocation `json:"type" default:"page_location"`
	paramObj
}

func (r CitationPageLocationParam) MarshalJSON() (data []byte, err error) {
	type shadow CitationPageLocationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CitationPageLocationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties CitedText, EndBlockIndex, SearchResultIndex, Source,
// StartBlockIndex, Title, Type are required.
type CitationSearchResultLocationParam struct {
	Title param.Opt[string] `json:"title,omitzero" api:"required"`
	// The full text of the cited block range, concatenated.
	//
	// Always equals the contents of `content[start_block_index:end_block_index]`
	// joined together. The text block is the minimal citable unit; this field is never
	// a substring of a single block. Not counted toward output tokens, and not counted
	// toward input tokens when sent back in subsequent turns.
	CitedText string `json:"cited_text" api:"required"`
	// Exclusive 0-based end index of the cited block range in the source's `content`
	// array.
	//
	// Always greater than `start_block_index`; a single-block citation has
	// `end_block_index = start_block_index + 1`.
	EndBlockIndex int64 `json:"end_block_index" api:"required"`
	// 0-based index of the cited search result among all `search_result` content
	// blocks in the request, in the order they appear across messages and tool
	// results.
	//
	// Counted separately from `document_index`; server-side web search results are not
	// included in this count.
	SearchResultIndex int64  `json:"search_result_index" api:"required"`
	Source            string `json:"source" api:"required"`
	// 0-based index of the first cited block in the source's `content` array.
	StartBlockIndex int64 `json:"start_block_index" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "search_result_location".
	Type constant.SearchResultLocation `json:"type" default:"search_result_location"`
	paramObj
}

func (r CitationSearchResultLocationParam) MarshalJSON() (data []byte, err error) {
	type shadow CitationSearchResultLocationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CitationSearchResultLocationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties CitedText, EncryptedIndex, Title, Type, URL are required.
type CitationWebSearchResultLocationParam struct {
	Title          param.Opt[string] `json:"title,omitzero" api:"required"`
	CitedText      string            `json:"cited_text" api:"required"`
	EncryptedIndex string            `json:"encrypted_index" api:"required"`
	URL            string            `json:"url" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "web_search_result_location".
	Type constant.WebSearchResultLocation `json:"type" default:"web_search_result_location"`
	paramObj
}

func (r CitationWebSearchResultLocationParam) MarshalJSON() (data []byte, err error) {
	type shadow CitationWebSearchResultLocationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CitationWebSearchResultLocationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CitationsConfigParam struct {
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	paramObj
}

func (r CitationsConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow CitationsConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CitationsConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties FileID, Type are required.
type CodeExecutionOutputBlockParam struct {
	FileID string `json:"file_id" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "code_execution_output".
	Type constant.CodeExecutionOutput `json:"type" default:"code_execution_output"`
	paramObj
}

func (r CodeExecutionOutputBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow CodeExecutionOutputBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CodeExecutionOutputBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Content, ReturnCode, Stderr, Stdout, Type are required.
type CodeExecutionResultBlockParam struct {
	Content    []CodeExecutionOutputBlockParam `json:"content,omitzero" api:"required"`
	ReturnCode int64                           `json:"return_code" api:"required"`
	Stderr     string                          `json:"stderr" api:"required"`
	Stdout     string                          `json:"stdout" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "code_execution_result".
	Type constant.CodeExecutionResult `json:"type" default:"code_execution_result"`
	paramObj
}

func (r CodeExecutionResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow CodeExecutionResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CodeExecutionResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Content, ToolUseID, Type are required.
type CodeExecutionToolResultBlockParam struct {
	// Code execution result with encrypted stdout for PFC + web_search results.
	Content   CodeExecutionToolResultBlockParamContentUnion `json:"content,omitzero" api:"required"`
	ToolUseID string                                        `json:"tool_use_id" api:"required"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "code_execution_tool_result".
	Type constant.CodeExecutionToolResult `json:"type" default:"code_execution_tool_result"`
	paramObj
}

func (r CodeExecutionToolResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow CodeExecutionToolResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CodeExecutionToolResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type CodeExecutionToolResultBlockParamContentUnion struct {
	OfError                                    *CodeExecutionToolResultErrorParam      `json:",omitzero,inline"`
	OfResultBlock                              *CodeExecutionResultBlockParam          `json:",omitzero,inline"`
	OfRequestEncryptedCodeExecutionResultBlock *EncryptedCodeExecutionResultBlockParam `json:",omitzero,inline"`
	paramUnion
}

func (u CodeExecutionToolResultBlockParamContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfError, u.OfResultBlock, u.OfRequestEncryptedCodeExecutionResultBlock)
}
func (u *CodeExecutionToolResultBlockParamContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *CodeExecutionToolResultBlockParamContentUnion) asAny() any {
	if !param.IsOmitted(u.OfError) {
		return u.OfError
	} else if !param.IsOmitted(u.OfResultBlock) {
		return u.OfResultBlock
	} else if !param.IsOmitted(u.OfRequestEncryptedCodeExecutionResultBlock) {
		return u.OfRequestEncryptedCodeExecutionResultBlock
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CodeExecutionToolResultBlockParamContentUnion) GetErrorCode() *string {
	if vt := u.OfError; vt != nil {
		return (*string)(&vt.ErrorCode)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CodeExecutionToolResultBlockParamContentUnion) GetStdout() *string {
	if vt := u.OfResultBlock; vt != nil {
		return &vt.Stdout
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CodeExecutionToolResultBlockParamContentUnion) GetEncryptedStdout() *string {
	if vt := u.OfRequestEncryptedCodeExecutionResultBlock; vt != nil {
		return &vt.EncryptedStdout
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CodeExecutionToolResultBlockParamContentUnion) GetType() *string {
	if vt := u.OfError; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfResultBlock; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestEncryptedCodeExecutionResultBlock; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CodeExecutionToolResultBlockParamContentUnion) GetReturnCode() *int64 {
	if vt := u.OfResultBlock; vt != nil {
		return (*int64)(&vt.ReturnCode)
	} else if vt := u.OfRequestEncryptedCodeExecutionResultBlock; vt != nil {
		return (*int64)(&vt.ReturnCode)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CodeExecutionToolResultBlockParamContentUnion) GetStderr() *string {
	if vt := u.OfResultBlock; vt != nil {
		return (*string)(&vt.Stderr)
	} else if vt := u.OfRequestEncryptedCodeExecutionResultBlock; vt != nil {
		return (*string)(&vt.Stderr)
	}
	return nil
}

// Returns a pointer to the underlying variant's Content property, if present.
func (u CodeExecutionToolResultBlockParamContentUnion) GetContent() []CodeExecutionOutputBlockParam {
	if vt := u.OfResultBlock; vt != nil {
		return vt.Content
	} else if vt := u.OfRequestEncryptedCodeExecutionResultBlock; vt != nil {
		return vt.Content
	}
	return nil
}

type CodeExecutionToolResultErrorCode string

// The properties ErrorCode, Type are required.
type CodeExecutionToolResultErrorParam struct {
	// Any of "invalid_tool_input", "unavailable", "too_many_requests",
	// "execution_time_exceeded".
	ErrorCode CodeExecutionToolResultErrorCode `json:"error_code,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "code_execution_tool_result_error".
	Type constant.CodeExecutionToolResultError `json:"type" default:"code_execution_tool_result_error"`
	paramObj
}

func (r CodeExecutionToolResultErrorParam) MarshalJSON() (data []byte, err error) {
	type shadow CodeExecutionToolResultErrorParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CodeExecutionToolResultErrorParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A compaction block containing summary of previous context.
//
// Users should round-trip these blocks from responses to subsequent requests to
// maintain context across compaction boundaries.
//
// When content is None, the block represents a failed compaction. The server
// treats these as no-ops. Empty string content is not allowed.
//
// The property Type is required.
type CompactionBlockParam struct {
	// Summary of previously compacted content, or null if compaction failed
	Content param.Opt[string] `json:"content,omitzero"`
	// Opaque metadata from prior compaction, to be round-tripped verbatim
	EncryptedContent param.Opt[string] `json:"encrypted_content,omitzero"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as "compaction".
	Type constant.Compaction `json:"type" default:"compaction"`
	paramObj
}

func (r CompactionBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow CompactionBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CompactionBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A content block that represents a file to be uploaded to the container Files
// uploaded via this block will be available in the container's input directory.
//
// The properties FileID, Type are required.
type ContainerUploadBlockParam struct {
	FileID string `json:"file_id" api:"required"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as "container_upload".
	Type constant.ContainerUpload `json:"type" default:"container_upload"`
	paramObj
}

func (r ContainerUploadBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ContainerUploadBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ContainerUploadBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ContentBlockParamUnion struct {
	OfText                              *TextBlockParam                              `json:",omitzero,inline"`
	OfImage                             *ImageBlockParam                             `json:",omitzero,inline"`
	OfDocument                          *RequestDocumentBlockParam                   `json:",omitzero,inline"`
	OfSearchResult                      *SearchResultBlockParam                      `json:",omitzero,inline"`
	OfThinking                          *ThinkingBlockParam                          `json:",omitzero,inline"`
	OfRedactedThinking                  *RedactedThinkingBlockParam                  `json:",omitzero,inline"`
	OfToolUse                           *ToolUseBlockParam                           `json:",omitzero,inline"`
	OfToolResult                        *ToolResultBlockParam                        `json:",omitzero,inline"`
	OfServerToolUse                     *ServerToolUseBlockParam                     `json:",omitzero,inline"`
	OfWebSearchToolResult               *WebSearchToolResultBlockParam               `json:",omitzero,inline"`
	OfWebFetchToolResult                *WebFetchToolResultBlockParam                `json:",omitzero,inline"`
	OfAdvisorToolResult                 *AdvisorToolResultBlockParam                 `json:",omitzero,inline"`
	OfCodeExecutionToolResult           *CodeExecutionToolResultBlockParam           `json:",omitzero,inline"`
	OfBashCodeExecutionToolResult       *BashCodeExecutionToolResultBlockParam       `json:",omitzero,inline"`
	OfTextEditorCodeExecutionToolResult *TextEditorCodeExecutionToolResultBlockParam `json:",omitzero,inline"`
	OfToolSearchToolResult              *ToolSearchToolResultBlockParam              `json:",omitzero,inline"`
	OfMCPToolUse                        *MCPToolUseBlockParam                        `json:",omitzero,inline"`
	OfMCPToolResult                     *RequestMCPToolResultBlockParam              `json:",omitzero,inline"`
	OfContainerUpload                   *ContainerUploadBlockParam                   `json:",omitzero,inline"`
	OfCompaction                        *CompactionBlockParam                        `json:",omitzero,inline"`
	OfToolAddition                      *RequestToolAdditionBlockParam               `json:",omitzero,inline"`
	OfToolRemoval                       *RequestToolRemovalBlockParam                `json:",omitzero,inline"`
	OfFallback                          *FallbackBlockParam                          `json:",omitzero,inline"`
	paramUnion
}

func (u ContentBlockParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfText,
		u.OfImage,
		u.OfDocument,
		u.OfSearchResult,
		u.OfThinking,
		u.OfRedactedThinking,
		u.OfToolUse,
		u.OfToolResult,
		u.OfServerToolUse,
		u.OfWebSearchToolResult,
		u.OfWebFetchToolResult,
		u.OfAdvisorToolResult,
		u.OfCodeExecutionToolResult,
		u.OfBashCodeExecutionToolResult,
		u.OfTextEditorCodeExecutionToolResult,
		u.OfToolSearchToolResult,
		u.OfMCPToolUse,
		u.OfMCPToolResult,
		u.OfContainerUpload,
		u.OfCompaction,
		u.OfToolAddition,
		u.OfToolRemoval,
		u.OfFallback)
}
func (u *ContentBlockParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ContentBlockParamUnion) asAny() any {
	if !param.IsOmitted(u.OfText) {
		return u.OfText
	} else if !param.IsOmitted(u.OfImage) {
		return u.OfImage
	} else if !param.IsOmitted(u.OfDocument) {
		return u.OfDocument
	} else if !param.IsOmitted(u.OfSearchResult) {
		return u.OfSearchResult
	} else if !param.IsOmitted(u.OfThinking) {
		return u.OfThinking
	} else if !param.IsOmitted(u.OfRedactedThinking) {
		return u.OfRedactedThinking
	} else if !param.IsOmitted(u.OfToolUse) {
		return u.OfToolUse
	} else if !param.IsOmitted(u.OfToolResult) {
		return u.OfToolResult
	} else if !param.IsOmitted(u.OfServerToolUse) {
		return u.OfServerToolUse
	} else if !param.IsOmitted(u.OfWebSearchToolResult) {
		return u.OfWebSearchToolResult
	} else if !param.IsOmitted(u.OfWebFetchToolResult) {
		return u.OfWebFetchToolResult
	} else if !param.IsOmitted(u.OfAdvisorToolResult) {
		return u.OfAdvisorToolResult
	} else if !param.IsOmitted(u.OfCodeExecutionToolResult) {
		return u.OfCodeExecutionToolResult
	} else if !param.IsOmitted(u.OfBashCodeExecutionToolResult) {
		return u.OfBashCodeExecutionToolResult
	} else if !param.IsOmitted(u.OfTextEditorCodeExecutionToolResult) {
		return u.OfTextEditorCodeExecutionToolResult
	} else if !param.IsOmitted(u.OfToolSearchToolResult) {
		return u.OfToolSearchToolResult
	} else if !param.IsOmitted(u.OfMCPToolUse) {
		return u.OfMCPToolUse
	} else if !param.IsOmitted(u.OfMCPToolResult) {
		return u.OfMCPToolResult
	} else if !param.IsOmitted(u.OfContainerUpload) {
		return u.OfContainerUpload
	} else if !param.IsOmitted(u.OfCompaction) {
		return u.OfCompaction
	} else if !param.IsOmitted(u.OfToolAddition) {
		return u.OfToolAddition
	} else if !param.IsOmitted(u.OfToolRemoval) {
		return u.OfToolRemoval
	} else if !param.IsOmitted(u.OfFallback) {
		return u.OfFallback
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetText() *string {
	if vt := u.OfText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetTransformations() *ImageTransformationsParam {
	if vt := u.OfImage; vt != nil {
		return &vt.Transformations
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetContext() *string {
	if vt := u.OfDocument; vt != nil && vt.Context.Valid() {
		return &vt.Context.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetSignature() *string {
	if vt := u.OfThinking; vt != nil {
		return &vt.Signature
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetThinking() *string {
	if vt := u.OfThinking; vt != nil {
		return &vt.Thinking
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetData() *string {
	if vt := u.OfRedactedThinking; vt != nil {
		return &vt.Data
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetServerName() *string {
	if vt := u.OfMCPToolUse; vt != nil {
		return &vt.ServerName
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetFileID() *string {
	if vt := u.OfContainerUpload; vt != nil {
		return &vt.FileID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetEncryptedContent() *string {
	if vt := u.OfCompaction; vt != nil && vt.EncryptedContent.Valid() {
		return &vt.EncryptedContent.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetFrom() *FallbackInfoParam {
	if vt := u.OfFallback; vt != nil {
		return &vt.From
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetTo() *FallbackInfoParam {
	if vt := u.OfFallback; vt != nil {
		return &vt.To
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetTrigger() *any {
	if vt := u.OfFallback; vt != nil {
		return &vt.Trigger
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetType() *string {
	if vt := u.OfText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfImage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDocument; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSearchResult; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfThinking; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRedactedThinking; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfToolUse; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfToolResult; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfServerToolUse; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfWebSearchToolResult; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfWebFetchToolResult; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAdvisorToolResult; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCodeExecutionToolResult; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfBashCodeExecutionToolResult; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfTextEditorCodeExecutionToolResult; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfToolSearchToolResult; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMCPToolUse; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMCPToolResult; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfContainerUpload; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCompaction; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfToolAddition; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfToolRemoval; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFallback; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetTitle() *string {
	if vt := u.OfDocument; vt != nil && vt.Title.Valid() {
		return &vt.Title.Value
	} else if vt := u.OfSearchResult; vt != nil {
		return (*string)(&vt.Title)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetID() *string {
	if vt := u.OfToolUse; vt != nil {
		return (*string)(&vt.ID)
	} else if vt := u.OfServerToolUse; vt != nil {
		return (*string)(&vt.ID)
	} else if vt := u.OfMCPToolUse; vt != nil {
		return (*string)(&vt.ID)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetName() *string {
	if vt := u.OfToolUse; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfServerToolUse; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfMCPToolUse; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetToolsetName() *string {
	if vt := u.OfToolUse; vt != nil && vt.ToolsetName.Valid() {
		return &vt.ToolsetName.Value
	} else if vt := u.OfToolResult; vt != nil && vt.ToolsetName.Valid() {
		return &vt.ToolsetName.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetToolUseID() *string {
	if vt := u.OfToolResult; vt != nil {
		return (*string)(&vt.ToolUseID)
	} else if vt := u.OfWebSearchToolResult; vt != nil {
		return (*string)(&vt.ToolUseID)
	} else if vt := u.OfWebFetchToolResult; vt != nil {
		return (*string)(&vt.ToolUseID)
	} else if vt := u.OfAdvisorToolResult; vt != nil {
		return (*string)(&vt.ToolUseID)
	} else if vt := u.OfCodeExecutionToolResult; vt != nil {
		return (*string)(&vt.ToolUseID)
	} else if vt := u.OfBashCodeExecutionToolResult; vt != nil {
		return (*string)(&vt.ToolUseID)
	} else if vt := u.OfTextEditorCodeExecutionToolResult; vt != nil {
		return (*string)(&vt.ToolUseID)
	} else if vt := u.OfToolSearchToolResult; vt != nil {
		return (*string)(&vt.ToolUseID)
	} else if vt := u.OfMCPToolResult; vt != nil {
		return (*string)(&vt.ToolUseID)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetIsError() *bool {
	if vt := u.OfToolResult; vt != nil && vt.IsError.Valid() {
		return &vt.IsError.Value
	} else if vt := u.OfMCPToolResult; vt != nil && vt.IsError.Valid() {
		return &vt.IsError.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's CacheControl property, if present.
func (u ContentBlockParamUnion) GetCacheControl() *CacheControlEphemeralParam {
	if vt := u.OfText; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfImage; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfDocument; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfSearchResult; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfToolUse; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfToolResult; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfServerToolUse; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfWebSearchToolResult; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfWebFetchToolResult; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfAdvisorToolResult; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfCodeExecutionToolResult; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfBashCodeExecutionToolResult; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfTextEditorCodeExecutionToolResult; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfToolSearchToolResult; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfMCPToolUse; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfMCPToolResult; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfContainerUpload; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfCompaction; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfToolAddition; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfToolRemoval; vt != nil {
		return &vt.CacheControl
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ContentBlockParamUnion) GetCitations() (res contentBlockParamUnionCitations) {
	if vt := u.OfText; vt != nil {
		res.any = &vt.Citations
	} else if vt := u.OfDocument; vt != nil {
		res.any = &vt.Citations
	} else if vt := u.OfSearchResult; vt != nil {
		res.any = &vt.Citations
	}
	return
}

// Can have the runtime types [*[]TextCitationParamUnion],
// [*CitationsConfigParam]
type contentBlockParamUnionCitations struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]qoder.TextCitationParamUnion:
//	case *qoder.CitationsConfigParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u contentBlockParamUnionCitations) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionCitations) GetEnabled() *bool {
	switch vt := u.any.(type) {
	case *CitationsConfigParam:
		return paramutil.AddrIfPresent(vt.Enabled)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ContentBlockParamUnion) GetSource() (res contentBlockParamUnionSource) {
	if vt := u.OfImage; vt != nil {
		res.any = vt.Source.asAny()
	} else if vt := u.OfDocument; vt != nil {
		res.any = vt.Source.asAny()
	} else if vt := u.OfSearchResult; vt != nil {
		res.any = &vt.Source
	}
	return
}

// Can have the runtime types [*Base64ImageSourceParam],
// [*URLImageSourceParam], [*FileImageSourceParam],
// [*Base64PDFSourceParam], [*PlainTextSourceParam],
// [*ContentBlockSourceParam], [*URLPDFSourceParam],
// [*FileDocumentSourceParam], [*string]
type contentBlockParamUnionSource struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *qoder.Base64ImageSourceParam:
//	case *qoder.URLImageSourceParam:
//	case *qoder.FileImageSourceParam:
//	case *qoder.Base64PDFSourceParam:
//	case *qoder.PlainTextSourceParam:
//	case *qoder.ContentBlockSourceParam:
//	case *qoder.URLPDFSourceParam:
//	case *qoder.FileDocumentSourceParam:
//	case *string:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u contentBlockParamUnionSource) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionSource) GetContent() *ContentBlockSourceContentUnionParam {
	switch vt := u.any.(type) {
	case *RequestDocumentBlockSourceUnionParam:
		return vt.GetContent()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionSource) GetData() *string {
	switch vt := u.any.(type) {
	case *ImageBlockParamSourceUnion:
		return vt.GetData()
	case *RequestDocumentBlockSourceUnionParam:
		return vt.GetData()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionSource) GetMediaType() *string {
	switch vt := u.any.(type) {
	case *ImageBlockParamSourceUnion:
		return vt.GetMediaType()
	case *RequestDocumentBlockSourceUnionParam:
		return vt.GetMediaType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionSource) GetType() *string {
	switch vt := u.any.(type) {
	case *ImageBlockParamSourceUnion:
		return vt.GetType()
	case *RequestDocumentBlockSourceUnionParam:
		return vt.GetType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionSource) GetURL() *string {
	switch vt := u.any.(type) {
	case *ImageBlockParamSourceUnion:
		return vt.GetURL()
	case *RequestDocumentBlockSourceUnionParam:
		return vt.GetURL()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionSource) GetFileID() *string {
	switch vt := u.any.(type) {
	case *ImageBlockParamSourceUnion:
		return vt.GetFileID()
	case *RequestDocumentBlockSourceUnionParam:
		return vt.GetFileID()
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ContentBlockParamUnion) GetContent() (res contentBlockParamUnionContent) {
	if vt := u.OfSearchResult; vt != nil {
		res.any = &vt.Content
	} else if vt := u.OfToolResult; vt != nil {
		res.any = &vt.Content
	} else if vt := u.OfWebSearchToolResult; vt != nil {
		res.any = vt.Content.asAny()
	} else if vt := u.OfWebFetchToolResult; vt != nil {
		res.any = vt.Content.asAny()
	} else if vt := u.OfAdvisorToolResult; vt != nil {
		res.any = vt.Content.asAny()
	} else if vt := u.OfCodeExecutionToolResult; vt != nil {
		res.any = vt.Content.asAny()
	} else if vt := u.OfBashCodeExecutionToolResult; vt != nil {
		res.any = vt.Content.asAny()
	} else if vt := u.OfTextEditorCodeExecutionToolResult; vt != nil {
		res.any = vt.Content.asAny()
	} else if vt := u.OfToolSearchToolResult; vt != nil {
		res.any = vt.Content.asAny()
	} else if vt := u.OfMCPToolResult; vt != nil {
		res.any = vt.Content.asAny()
	} else if vt := u.OfCompaction; vt != nil && vt.Content.Valid() {
		res.any = &vt.Content.Value
	}
	return
}

// Can have the runtime types [_[]TextBlockParam],
// [_[]ToolResultBlockParamContentUnion], [*[]WebSearchResultBlockParam],
// [*WebFetchToolResultErrorBlockParam], [*WebFetchBlockParam],
// [*AdvisorToolResultErrorParam], [*AdvisorResultBlockParam],
// [*AdvisorRedactedResultBlockParam],
// [*CodeExecutionToolResultErrorParam], [*CodeExecutionResultBlockParam],
// [*EncryptedCodeExecutionResultBlockParam],
// [*BashCodeExecutionToolResultErrorParam],
// [*BashCodeExecutionResultBlockParam],
// [*TextEditorCodeExecutionToolResultErrorParam],
// [*TextEditorCodeExecutionViewResultBlockParam],
// [*TextEditorCodeExecutionCreateResultBlockParam],
// [*TextEditorCodeExecutionStrReplaceResultBlockParam],
// [*ToolSearchToolResultErrorParam],
// [*ToolSearchToolSearchResultBlockParam], [*string]
type contentBlockParamUnionContent struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]qoder.TextBlockParam:
//	case *[]qoder.ToolResultBlockParamContentUnion:
//	case *[]qoder.WebSearchResultBlockParam:
//	case *qoder.WebFetchToolResultErrorBlockParam:
//	case *qoder.WebFetchBlockParam:
//	case *qoder.AdvisorToolResultErrorParam:
//	case *qoder.AdvisorResultBlockParam:
//	case *qoder.AdvisorRedactedResultBlockParam:
//	case *qoder.CodeExecutionToolResultErrorParam:
//	case *qoder.CodeExecutionResultBlockParam:
//	case *qoder.EncryptedCodeExecutionResultBlockParam:
//	case *qoder.BashCodeExecutionToolResultErrorParam:
//	case *qoder.BashCodeExecutionResultBlockParam:
//	case *qoder.TextEditorCodeExecutionToolResultErrorParam:
//	case *qoder.TextEditorCodeExecutionViewResultBlockParam:
//	case *qoder.TextEditorCodeExecutionCreateResultBlockParam:
//	case *qoder.TextEditorCodeExecutionStrReplaceResultBlockParam:
//	case *qoder.ToolSearchToolResultErrorParam:
//	case *qoder.ToolSearchToolSearchResultBlockParam:
//	case *string:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u contentBlockParamUnionContent) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetURL() *string {
	switch vt := u.any.(type) {
	case *WebFetchToolResultBlockParamContentUnion:
		return vt.GetURL()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetRetrievedAt() *string {
	switch vt := u.any.(type) {
	case *WebFetchToolResultBlockParamContentUnion:
		return vt.GetRetrievedAt()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetText() *string {
	switch vt := u.any.(type) {
	case *AdvisorToolResultBlockParamContentUnion:
		return vt.GetText()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetEncryptedContent() *string {
	switch vt := u.any.(type) {
	case *AdvisorToolResultBlockParamContentUnion:
		return vt.GetEncryptedContent()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetEncryptedStdout() *string {
	switch vt := u.any.(type) {
	case *CodeExecutionToolResultBlockParamContentUnion:
		return vt.GetEncryptedStdout()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetFileType() *string {
	switch vt := u.any.(type) {
	case *TextEditorCodeExecutionToolResultBlockParamContentUnion:
		return vt.GetFileType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetNumLines() *int64 {
	switch vt := u.any.(type) {
	case *TextEditorCodeExecutionToolResultBlockParamContentUnion:
		return vt.GetNumLines()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetStartLine() *int64 {
	switch vt := u.any.(type) {
	case *TextEditorCodeExecutionToolResultBlockParamContentUnion:
		return vt.GetStartLine()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetTotalLines() *int64 {
	switch vt := u.any.(type) {
	case *TextEditorCodeExecutionToolResultBlockParamContentUnion:
		return vt.GetTotalLines()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetIsFileUpdate() *bool {
	switch vt := u.any.(type) {
	case *TextEditorCodeExecutionToolResultBlockParamContentUnion:
		return vt.GetIsFileUpdate()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetLines() []string {
	switch vt := u.any.(type) {
	case *TextEditorCodeExecutionToolResultBlockParamContentUnion:
		return vt.GetLines()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetNewLines() *int64 {
	switch vt := u.any.(type) {
	case *TextEditorCodeExecutionToolResultBlockParamContentUnion:
		return vt.GetNewLines()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetNewStart() *int64 {
	switch vt := u.any.(type) {
	case *TextEditorCodeExecutionToolResultBlockParamContentUnion:
		return vt.GetNewStart()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetOldLines() *int64 {
	switch vt := u.any.(type) {
	case *TextEditorCodeExecutionToolResultBlockParamContentUnion:
		return vt.GetOldLines()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetOldStart() *int64 {
	switch vt := u.any.(type) {
	case *TextEditorCodeExecutionToolResultBlockParamContentUnion:
		return vt.GetOldStart()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetToolReferences() []ToolReferenceBlockParam {
	switch vt := u.any.(type) {
	case *ToolSearchToolResultBlockParamContentUnion:
		return vt.GetToolReferences()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetErrorCode() *string {
	switch vt := u.any.(type) {
	case *WebSearchToolResultBlockParamContentUnion:
		if vt.OfError != nil {
			return (*string)(&vt.OfError.ErrorCode)
		}
	case *WebFetchToolResultBlockParamContentUnion:
		return vt.GetErrorCode()
	case *AdvisorToolResultBlockParamContentUnion:
		return vt.GetErrorCode()
	case *CodeExecutionToolResultBlockParamContentUnion:
		return vt.GetErrorCode()
	case *BashCodeExecutionToolResultBlockParamContentUnion:
		return vt.GetErrorCode()
	case *TextEditorCodeExecutionToolResultBlockParamContentUnion:
		return vt.GetErrorCode()
	case *ToolSearchToolResultBlockParamContentUnion:
		return vt.GetErrorCode()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetType() *string {
	switch vt := u.any.(type) {
	case *WebSearchToolResultBlockParamContentUnion:
		if vt.OfError != nil {
			return (*string)(&vt.OfError.Type)
		}
	case *WebFetchToolResultBlockParamContentUnion:
		return vt.GetType()
	case *AdvisorToolResultBlockParamContentUnion:
		return vt.GetType()
	case *CodeExecutionToolResultBlockParamContentUnion:
		return vt.GetType()
	case *BashCodeExecutionToolResultBlockParamContentUnion:
		return vt.GetType()
	case *TextEditorCodeExecutionToolResultBlockParamContentUnion:
		return vt.GetType()
	case *ToolSearchToolResultBlockParamContentUnion:
		return vt.GetType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetStopReason() *string {
	switch vt := u.any.(type) {
	case *AdvisorToolResultBlockParamContentUnion:
		return vt.GetStopReason()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetReturnCode() *int64 {
	switch vt := u.any.(type) {
	case *CodeExecutionToolResultBlockParamContentUnion:
		return vt.GetReturnCode()
	case *BashCodeExecutionToolResultBlockParamContentUnion:
		return vt.GetReturnCode()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetStderr() *string {
	switch vt := u.any.(type) {
	case *CodeExecutionToolResultBlockParamContentUnion:
		return vt.GetStderr()
	case *BashCodeExecutionToolResultBlockParamContentUnion:
		return vt.GetStderr()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetStdout() *string {
	switch vt := u.any.(type) {
	case *CodeExecutionToolResultBlockParamContentUnion:
		return vt.GetStdout()
	case *BashCodeExecutionToolResultBlockParamContentUnion:
		return vt.GetStdout()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionContent) GetErrorMessage() *string {
	switch vt := u.any.(type) {
	case *TextEditorCodeExecutionToolResultBlockParamContentUnion:
		return vt.GetErrorMessage()
	case *ToolSearchToolResultBlockParamContentUnion:
		return vt.GetErrorMessage()
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u contentBlockParamUnionContent) GetContent() (res contentBlockParamUnionContentContent) {
	switch vt := u.any.(type) {
	case *WebFetchToolResultBlockParamContentUnion:
		res.any = vt.GetContent()
	case *CodeExecutionToolResultBlockParamContentUnion:
		res.any = vt.GetContent()
	case *BashCodeExecutionToolResultBlockParamContentUnion:
		res.any = vt.GetContent()
	case *TextEditorCodeExecutionToolResultBlockParamContentUnion:
		res.any = vt.GetContent()
	}
	return res
}

// Can have the runtime types [*RequestDocumentBlockParam],
// [_[]CodeExecutionOutputBlockParam],
// [_[]BashCodeExecutionOutputBlockParam], [*string]
type contentBlockParamUnionContentContent struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *qoder.RequestDocumentBlockParam:
//	case *[]qoder.CodeExecutionOutputBlockParam:
//	case *[]qoder.BashCodeExecutionOutputBlockParam:
//	case *string:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u contentBlockParamUnionContentContent) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's Input property, if present.
func (u ContentBlockParamUnion) GetInput() *any {
	if vt := u.OfToolUse; vt != nil {
		return &vt.Input
	} else if vt := u.OfServerToolUse; vt != nil {
		return &vt.Input
	} else if vt := u.OfMCPToolUse; vt != nil {
		return &vt.Input
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ContentBlockParamUnion) GetCaller() (res contentBlockParamUnionCaller) {
	if vt := u.OfToolUse; vt != nil {
		res.any = vt.Caller.asAny()
	} else if vt := u.OfServerToolUse; vt != nil {
		res.any = vt.Caller.asAny()
	} else if vt := u.OfWebSearchToolResult; vt != nil {
		res.any = vt.Caller.asAny()
	} else if vt := u.OfWebFetchToolResult; vt != nil {
		res.any = vt.Caller.asAny()
	}
	return
}

// Can have the runtime types [*DirectCallerParam],
// [*ServerToolCallerParam], [*ServerToolCaller20260120Param]
type contentBlockParamUnionCaller struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *qoder.DirectCallerParam:
//	case *qoder.ServerToolCallerParam:
//	case *qoder.ServerToolCaller20260120Param:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u contentBlockParamUnionCaller) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionCaller) GetType() *string {
	switch vt := u.any.(type) {
	case *ToolUseBlockParamCallerUnion:
		return vt.GetType()
	case *ServerToolUseBlockParamCallerUnion:
		return vt.GetType()
	case *WebSearchToolResultBlockParamCallerUnion:
		return vt.GetType()
	case *WebFetchToolResultBlockParamCallerUnion:
		return vt.GetType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionCaller) GetToolID() *string {
	switch vt := u.any.(type) {
	case *ToolUseBlockParamCallerUnion:
		return vt.GetToolID()
	case *ServerToolUseBlockParamCallerUnion:
		return vt.GetToolID()
	case *WebSearchToolResultBlockParamCallerUnion:
		return vt.GetToolID()
	case *WebFetchToolResultBlockParamCallerUnion:
		return vt.GetToolID()
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ContentBlockParamUnion) GetTool() (res contentBlockParamUnionTool) {
	if vt := u.OfToolAddition; vt != nil {
		res.any = vt.Tool.asAny()
	} else if vt := u.OfToolRemoval; vt != nil {
		res.any = vt.Tool.asAny()
	}
	return
}

// Can have the runtime types [*ToolChangeToolReferenceParam],
// [*ToolChangeMCPToolReferenceParam],
// [*ToolChangeMCPToolsetReferenceParam]
type contentBlockParamUnionTool struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *qoder.ToolChangeToolReferenceParam:
//	case *qoder.ToolChangeMCPToolReferenceParam:
//	case *qoder.ToolChangeMCPToolsetReferenceParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u contentBlockParamUnionTool) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionTool) GetName() *string {
	switch vt := u.any.(type) {
	case *RequestToolAdditionBlockToolUnionParam:
		return vt.GetName()
	case *RequestToolRemovalBlockToolUnionParam:
		return vt.GetName()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionTool) GetType() *string {
	switch vt := u.any.(type) {
	case *RequestToolAdditionBlockToolUnionParam:
		return vt.GetType()
	case *RequestToolRemovalBlockToolUnionParam:
		return vt.GetType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionTool) GetServerName() *string {
	switch vt := u.any.(type) {
	case *RequestToolAdditionBlockToolUnionParam:
		return vt.GetServerName()
	case *RequestToolRemovalBlockToolUnionParam:
		return vt.GetServerName()
	}
	return nil
}

// The properties Content, Type are required.
type ContentBlockSourceParam struct {
	Content ContentBlockSourceContentUnionParam `json:"content,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "content".
	Type constant.Content `json:"type" default:"content"`
	paramObj
}

func (r ContentBlockSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow ContentBlockSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ContentBlockSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ContentBlockSourceContentUnionParam struct {
	OfString                    param.Opt[string]                     `json:",omitzero,inline"`
	OfContentBlockSourceContent []ContentBlockSourceContentUnionParam `json:",omitzero,inline"`
	paramUnion
}

func (u ContentBlockSourceContentUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfContentBlockSourceContent)
}
func (u *ContentBlockSourceContentUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ContentBlockSourceContentUnionParam) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfContentBlockSourceContent) {
		return &u.OfContentBlockSourceContent
	}
	return nil
}

// Tool invocation directly from the model.
//
// This struct has a constant value, construct it with [NewDirectCallerParam].
type DirectCallerParam struct {
	Type constant.Direct `json:"type" default:"direct"`
	paramObj
}

func (r DirectCallerParam) MarshalJSON() (data []byte, err error) {
	type shadow DirectCallerParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DirectCallerParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Code execution result with encrypted stdout for PFC + web_search results.
//
// The properties Content, EncryptedStdout, ReturnCode, Stderr, Type are required.
type EncryptedCodeExecutionResultBlockParam struct {
	Content         []CodeExecutionOutputBlockParam `json:"content,omitzero" api:"required"`
	EncryptedStdout string                          `json:"encrypted_stdout" api:"required"`
	ReturnCode      int64                           `json:"return_code" api:"required"`
	Stderr          string                          `json:"stderr" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "encrypted_code_execution_result".
	Type constant.EncryptedCodeExecutionResult `json:"type" default:"encrypted_code_execution_result"`
	paramObj
}

func (r EncryptedCodeExecutionResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow EncryptedCodeExecutionResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *EncryptedCodeExecutionResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A `fallback` block echoed back from a prior response.
//
// Accepted in `messages[].content` and not rendered into the prompt; not validated
// against the request's `fallbacks` chain or top-level `model`.
//
// Echo the assistant turn back verbatim, including this block in its original
// position. The block marks the boundary between content produced before and after
// a fallback hop, and the server relies on that boundary to validate the turn:
// when thinking runs flank the boundary, omitting the block merges them into one
// span the server cannot validate (the request is rejected), and moving it into
// the middle of a single run is likewise rejected; between non-thinking blocks the
// block's placement has no validation effect.
//
// The properties From, To, Type are required.
type FallbackBlockParam struct {
	// Identifies one hop of a fallback transition.
	From FallbackInfoParam `json:"from,omitzero" api:"required"`
	// Identifies one hop of a fallback transition.
	To FallbackInfoParam `json:"to,omitzero" api:"required"`
	// The response block's `trigger`, echoed verbatim. Accepted and ignored by the
	// server; any object or `null` is allowed.
	Trigger any `json:"trigger,omitzero"`
	// This field can be elided, and will marshal its zero value as "fallback".
	Type constant.Fallback `json:"type" default:"fallback"`
	paramObj
}

func (r FallbackBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow FallbackBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *FallbackBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Identifies one hop of a fallback transition.
//
// The property Model is required.
type FallbackInfoParam struct {
	// The model that will complete your prompt.
	//
	// See [models](https://docs.qoder.com/cloud-agents/api/models/list) for additional
	// details and options.
	Model Model `json:"model,omitzero" api:"required"`
	paramObj
}

func (r FallbackInfoParam) MarshalJSON() (data []byte, err error) {
	type shadow FallbackInfoParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *FallbackInfoParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties FileID, Type are required.
type FileDocumentSourceParam struct {
	FileID string `json:"file_id" api:"required"`
	// This field can be elided, and will marshal its zero value as "file".
	Type constant.File `json:"type" default:"file"`
	paramObj
}

func (r FileDocumentSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow FileDocumentSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *FileDocumentSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties FileID, Type are required.
type FileImageSourceParam struct {
	FileID string `json:"file_id" api:"required"`
	// This field can be elided, and will marshal its zero value as "file".
	Type constant.File `json:"type" default:"file"`
	paramObj
}

func (r FileImageSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow FileImageSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *FileImageSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Source, Type are required.
type ImageBlockParam struct {
	Source ImageBlockParamSourceUnion `json:"source,omitzero" api:"required"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// Configures the transformations the server applies to this image before the model
	// observes it. Each key names a condition the server transforms images for; its
	// value selects the transformation applied. Omitted keys keep their default
	// behavior, and an empty object is equivalent to omitting the field.
	Transformations ImageTransformationsParam `json:"transformations,omitzero"`
	// This field can be elided, and will marshal its zero value as "image".
	Type constant.Image `json:"type" default:"image"`
	paramObj
}

func (r ImageBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ImageBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ImageBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ImageBlockParamSourceUnion struct {
	OfBase64 *Base64ImageSourceParam `json:",omitzero,inline"`
	OfURL    *URLImageSourceParam    `json:",omitzero,inline"`
	OfFile   *FileImageSourceParam   `json:",omitzero,inline"`
	paramUnion
}

func (u ImageBlockParamSourceUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfBase64, u.OfURL, u.OfFile)
}
func (u *ImageBlockParamSourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ImageBlockParamSourceUnion) asAny() any {
	if !param.IsOmitted(u.OfBase64) {
		return u.OfBase64
	} else if !param.IsOmitted(u.OfURL) {
		return u.OfURL
	} else if !param.IsOmitted(u.OfFile) {
		return u.OfFile
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ImageBlockParamSourceUnion) GetData() *string {
	if vt := u.OfBase64; vt != nil {
		return &vt.Data
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ImageBlockParamSourceUnion) GetMediaType() *string {
	if vt := u.OfBase64; vt != nil {
		return (*string)(&vt.MediaType)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ImageBlockParamSourceUnion) GetURL() *string {
	if vt := u.OfURL; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ImageBlockParamSourceUnion) GetFileID() *string {
	if vt := u.OfFile; vt != nil {
		return &vt.FileID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ImageBlockParamSourceUnion) GetType() *string {
	if vt := u.OfBase64; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfURL; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFile; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Configures the transformations the server applies to this image before the model
// observes it. Each key names a condition the server transforms images for; its
// value selects the transformation applied. Omitted keys keep their default
// behavior, and an empty object is equivalent to omitting the field.
type ImageTransformationsParam struct {
	// What the server does when this image exceeds the model's maximum image size.
	// `"downsize"` (the default) scales the image down to fit, which changes the
	// dimensions the model observes without telling you. `"error"` instead rejects the
	// request with a 400 error naming the image's dimensions and the largest
	// dimensions that fit, so you can scale the image deliberately — your image is
	// never silently scaled down.
	//
	// Any of "downsize", "error".
	OversizedImage ImageTransformationsParamOversizedImage `json:"oversized_image,omitzero"`
	paramObj
}

func (r ImageTransformationsParam) MarshalJSON() (data []byte, err error) {
	type shadow ImageTransformationsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ImageTransformationsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// What the server does when this image exceeds the model's maximum image size.
// `"downsize"` (the default) scales the image down to fit, which changes the
// dimensions the model observes without telling you. `"error"` instead rejects the
// request with a 400 error naming the image's dimensions and the largest
// dimensions that fit, so you can scale the image deliberately — your image is
// never silently scaled down.
type ImageTransformationsParamOversizedImage string

// The properties ID, Input, Name, ServerName, Type are required.
type MCPToolUseBlockParam struct {
	ID    string `json:"id" api:"required"`
	Input any    `json:"input,omitzero" api:"required"`
	Name  string `json:"name" api:"required"`
	// The name of the MCP server
	ServerName string `json:"server_name" api:"required"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as "mcp_tool_use".
	Type constant.MCPToolUse `json:"type" default:"mcp_tool_use"`
	paramObj
}

func (r MCPToolUseBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow MCPToolUseBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MCPToolUseBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Data, MediaType, Type are required.
type PlainTextSourceParam struct {
	Data string `json:"data" api:"required"`
	// This field can be elided, and will marshal its zero value as "text/plain".
	MediaType constant.TextPlain `json:"media_type" default:"text/plain"`
	// This field can be elided, and will marshal its zero value as "text".
	Type constant.Text `json:"type" default:"text"`
	paramObj
}

func (r PlainTextSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow PlainTextSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *PlainTextSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Data, Type are required.
type RedactedThinkingBlockParam struct {
	// The `data` value of this redacted thinking block, exactly as returned by the API
	// in a previous response. Opaque and encrypted; pass it back unchanged.
	Data string `json:"data" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "redacted_thinking".
	Type constant.RedactedThinking `json:"type" default:"redacted_thinking"`
	paramObj
}

func (r RedactedThinkingBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow RedactedThinkingBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RedactedThinkingBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Source, Type are required.
type RequestDocumentBlockParam struct {
	Source  RequestDocumentBlockSourceUnionParam `json:"source,omitzero" api:"required"`
	Context param.Opt[string]                    `json:"context,omitzero"`
	Title   param.Opt[string]                    `json:"title,omitzero"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	Citations    CitationsConfigParam       `json:"citations,omitzero"`
	// This field can be elided, and will marshal its zero value as "document".
	Type constant.Document `json:"type" default:"document"`
	paramObj
}

func (r RequestDocumentBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow RequestDocumentBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RequestDocumentBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RequestDocumentBlockSourceUnionParam struct {
	OfBase64  *Base64PDFSourceParam    `json:",omitzero,inline"`
	OfText    *PlainTextSourceParam    `json:",omitzero,inline"`
	OfContent *ContentBlockSourceParam `json:",omitzero,inline"`
	OfURL     *URLPDFSourceParam       `json:",omitzero,inline"`
	OfFile    *FileDocumentSourceParam `json:",omitzero,inline"`
	paramUnion
}

func (u RequestDocumentBlockSourceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfBase64,
		u.OfText,
		u.OfContent,
		u.OfURL,
		u.OfFile)
}
func (u *RequestDocumentBlockSourceUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *RequestDocumentBlockSourceUnionParam) asAny() any {
	if !param.IsOmitted(u.OfBase64) {
		return u.OfBase64
	} else if !param.IsOmitted(u.OfText) {
		return u.OfText
	} else if !param.IsOmitted(u.OfContent) {
		return u.OfContent
	} else if !param.IsOmitted(u.OfURL) {
		return u.OfURL
	} else if !param.IsOmitted(u.OfFile) {
		return u.OfFile
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RequestDocumentBlockSourceUnionParam) GetContent() *ContentBlockSourceContentUnionParam {
	if vt := u.OfContent; vt != nil {
		return &vt.Content
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RequestDocumentBlockSourceUnionParam) GetURL() *string {
	if vt := u.OfURL; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RequestDocumentBlockSourceUnionParam) GetFileID() *string {
	if vt := u.OfFile; vt != nil {
		return &vt.FileID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RequestDocumentBlockSourceUnionParam) GetData() *string {
	if vt := u.OfBase64; vt != nil {
		return (*string)(&vt.Data)
	} else if vt := u.OfText; vt != nil {
		return (*string)(&vt.Data)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RequestDocumentBlockSourceUnionParam) GetMediaType() *string {
	if vt := u.OfBase64; vt != nil {
		return (*string)(&vt.MediaType)
	} else if vt := u.OfText; vt != nil {
		return (*string)(&vt.MediaType)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RequestDocumentBlockSourceUnionParam) GetType() *string {
	if vt := u.OfBase64; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfContent; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfURL; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFile; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// The properties ToolUseID, Type are required.
type RequestMCPToolResultBlockParam struct {
	ToolUseID string          `json:"tool_use_id" api:"required"`
	IsError   param.Opt[bool] `json:"is_error,omitzero"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam                 `json:"cache_control,omitzero"`
	Content      RequestMCPToolResultBlockParamContentUnion `json:"content,omitzero"`
	// This field can be elided, and will marshal its zero value as "mcp_tool_result".
	Type constant.MCPToolResult `json:"type" default:"mcp_tool_result"`
	paramObj
}

func (r RequestMCPToolResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow RequestMCPToolResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RequestMCPToolResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RequestMCPToolResultBlockParamContentUnion struct {
	OfString                    param.Opt[string] `json:",omitzero,inline"`
	OfMCPToolResultBlockContent []TextBlockParam  `json:",omitzero,inline"`
	paramUnion
}

func (u RequestMCPToolResultBlockParamContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfMCPToolResultBlockContent)
}
func (u *RequestMCPToolResultBlockParamContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *RequestMCPToolResultBlockParamContentUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfMCPToolResultBlockContent) {
		return &u.OfMCPToolResultBlockContent
	}
	return nil
}

// Mid-conversation directive to surface a declared tool.
//
// `tool` references a tool (or MCP toolset) by name from the request's `tools`; it
// is offered to the model from this point in the conversation onward.
//
// The properties Tool, Type are required.
type RequestToolAdditionBlockParam struct {
	// Reference to a single tool the caller declared directly in `tools[]`. Does not
	// accept the composed `{server}_{name}` form the server assigns to MCP-resolved
	// tools — use `mcp_tool_reference` or `mcp_toolset_reference` for those.
	Tool RequestToolAdditionBlockToolUnionParam `json:"tool,omitzero" api:"required"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as "tool_addition".
	Type constant.ToolAddition `json:"type" default:"tool_addition"`
	paramObj
}

func (r RequestToolAdditionBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow RequestToolAdditionBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RequestToolAdditionBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RequestToolAdditionBlockToolUnionParam struct {
	OfToolReference       *ToolChangeToolReferenceParam       `json:",omitzero,inline"`
	OfMCPToolReference    *ToolChangeMCPToolReferenceParam    `json:",omitzero,inline"`
	OfMCPToolsetReference *ToolChangeMCPToolsetReferenceParam `json:",omitzero,inline"`
	paramUnion
}

func (u RequestToolAdditionBlockToolUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfToolReference, u.OfMCPToolReference, u.OfMCPToolsetReference)
}
func (u *RequestToolAdditionBlockToolUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *RequestToolAdditionBlockToolUnionParam) asAny() any {
	if !param.IsOmitted(u.OfToolReference) {
		return u.OfToolReference
	} else if !param.IsOmitted(u.OfMCPToolReference) {
		return u.OfMCPToolReference
	} else if !param.IsOmitted(u.OfMCPToolsetReference) {
		return u.OfMCPToolsetReference
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RequestToolAdditionBlockToolUnionParam) GetName() *string {
	if vt := u.OfToolReference; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfMCPToolReference; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RequestToolAdditionBlockToolUnionParam) GetType() *string {
	if vt := u.OfToolReference; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMCPToolReference; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMCPToolsetReference; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RequestToolAdditionBlockToolUnionParam) GetServerName() *string {
	if vt := u.OfMCPToolReference; vt != nil {
		return (*string)(&vt.ServerName)
	} else if vt := u.OfMCPToolsetReference; vt != nil {
		return (*string)(&vt.ServerName)
	}
	return nil
}

// Mid-conversation directive to withdraw a tool.
//
// `tool` references a tool (or MCP toolset) by name from the request's `tools`; it
// is no longer offered to the model from this point in the conversation onward.
//
// The properties Tool, Type are required.
type RequestToolRemovalBlockParam struct {
	// Reference to a single tool the caller declared directly in `tools[]`. Does not
	// accept the composed `{server}_{name}` form the server assigns to MCP-resolved
	// tools — use `mcp_tool_reference` or `mcp_toolset_reference` for those.
	Tool RequestToolRemovalBlockToolUnionParam `json:"tool,omitzero" api:"required"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as "tool_removal".
	Type constant.ToolRemoval `json:"type" default:"tool_removal"`
	paramObj
}

func (r RequestToolRemovalBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow RequestToolRemovalBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *RequestToolRemovalBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type RequestToolRemovalBlockToolUnionParam struct {
	OfToolReference       *ToolChangeToolReferenceParam       `json:",omitzero,inline"`
	OfMCPToolReference    *ToolChangeMCPToolReferenceParam    `json:",omitzero,inline"`
	OfMCPToolsetReference *ToolChangeMCPToolsetReferenceParam `json:",omitzero,inline"`
	paramUnion
}

func (u RequestToolRemovalBlockToolUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfToolReference, u.OfMCPToolReference, u.OfMCPToolsetReference)
}
func (u *RequestToolRemovalBlockToolUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *RequestToolRemovalBlockToolUnionParam) asAny() any {
	if !param.IsOmitted(u.OfToolReference) {
		return u.OfToolReference
	} else if !param.IsOmitted(u.OfMCPToolReference) {
		return u.OfMCPToolReference
	} else if !param.IsOmitted(u.OfMCPToolsetReference) {
		return u.OfMCPToolsetReference
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RequestToolRemovalBlockToolUnionParam) GetName() *string {
	if vt := u.OfToolReference; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfMCPToolReference; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RequestToolRemovalBlockToolUnionParam) GetType() *string {
	if vt := u.OfToolReference; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMCPToolReference; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMCPToolsetReference; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u RequestToolRemovalBlockToolUnionParam) GetServerName() *string {
	if vt := u.OfMCPToolReference; vt != nil {
		return (*string)(&vt.ServerName)
	} else if vt := u.OfMCPToolsetReference; vt != nil {
		return (*string)(&vt.ServerName)
	}
	return nil
}

// The properties Content, Source, Title, Type are required.
type SearchResultBlockParam struct {
	Content []TextBlockParam `json:"content,omitzero" api:"required"`
	Source  string           `json:"source" api:"required"`
	Title   string           `json:"title" api:"required"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	Citations    CitationsConfigParam       `json:"citations,omitzero"`
	// This field can be elided, and will marshal its zero value as "search_result".
	Type constant.SearchResult `json:"type" default:"search_result"`
	paramObj
}

func (r SearchResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow SearchResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SearchResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Tool invocation generated by a server-side tool.
//
// The properties ToolID, Type are required.
type ServerToolCallerParam struct {
	ToolID string `json:"tool_id" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "code_execution_20250825".
	Type constant.CodeExecution20250825 `json:"type" default:"code_execution_20250825"`
	paramObj
}

func (r ServerToolCallerParam) MarshalJSON() (data []byte, err error) {
	type shadow ServerToolCallerParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ServerToolCallerParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties ToolID, Type are required.
type ServerToolCaller20260120Param struct {
	ToolID string `json:"tool_id" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "code_execution_20260120".
	Type constant.CodeExecution20260120 `json:"type" default:"code_execution_20260120"`
	paramObj
}

func (r ServerToolCaller20260120Param) MarshalJSON() (data []byte, err error) {
	type shadow ServerToolCaller20260120Param
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ServerToolCaller20260120Param) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties ID, Input, Name, Type are required.
type ServerToolUseBlockParam struct {
	ID    string `json:"id" api:"required"`
	Input any    `json:"input,omitzero" api:"required"`
	// Any of "advisor", "web_search", "web_fetch", "code_execution",
	// "bash_code_execution", "text_editor_code_execution", "tool_search_tool_regex",
	// "tool_search_tool_bm25".
	Name ServerToolUseBlockParamName `json:"name,omitzero" api:"required"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// Tool invocation directly from the model.
	Caller ServerToolUseBlockParamCallerUnion `json:"caller,omitzero"`
	// This field can be elided, and will marshal its zero value as "server_tool_use".
	Type constant.ServerToolUse `json:"type" default:"server_tool_use"`
	paramObj
}

func (r ServerToolUseBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ServerToolUseBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ServerToolUseBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ServerToolUseBlockParamName string

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ServerToolUseBlockParamCallerUnion struct {
	OfDirect                *DirectCallerParam             `json:",omitzero,inline"`
	OfCodeExecution20250825 *ServerToolCallerParam         `json:",omitzero,inline"`
	OfCodeExecution20260120 *ServerToolCaller20260120Param `json:",omitzero,inline"`
	paramUnion
}

func (u ServerToolUseBlockParamCallerUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDirect, u.OfCodeExecution20250825, u.OfCodeExecution20260120)
}
func (u *ServerToolUseBlockParamCallerUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ServerToolUseBlockParamCallerUnion) asAny() any {
	if !param.IsOmitted(u.OfDirect) {
		return u.OfDirect
	} else if !param.IsOmitted(u.OfCodeExecution20250825) {
		return u.OfCodeExecution20250825
	} else if !param.IsOmitted(u.OfCodeExecution20260120) {
		return u.OfCodeExecution20260120
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ServerToolUseBlockParamCallerUnion) GetType() *string {
	if vt := u.OfDirect; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCodeExecution20250825; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCodeExecution20260120; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ServerToolUseBlockParamCallerUnion) GetToolID() *string {
	if vt := u.OfCodeExecution20250825; vt != nil {
		return (*string)(&vt.ToolID)
	} else if vt := u.OfCodeExecution20260120; vt != nil {
		return (*string)(&vt.ToolID)
	}
	return nil
}

// The properties Text, Type are required.
type TextBlockParam struct {
	Text      string                   `json:"text" api:"required"`
	Citations []TextCitationParamUnion `json:"citations,omitzero"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as "text".
	Type constant.Text `json:"type" default:"text"`
	paramObj
}

func (r TextBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow TextBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TextBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type TextCitationParamUnion struct {
	OfCharLocation            *CitationCharLocationParam            `json:",omitzero,inline"`
	OfPageLocation            *CitationPageLocationParam            `json:",omitzero,inline"`
	OfContentBlockLocation    *CitationContentBlockLocationParam    `json:",omitzero,inline"`
	OfWebSearchResultLocation *CitationWebSearchResultLocationParam `json:",omitzero,inline"`
	OfSearchResultLocation    *CitationSearchResultLocationParam    `json:",omitzero,inline"`
	paramUnion
}

func (u TextCitationParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfCharLocation,
		u.OfPageLocation,
		u.OfContentBlockLocation,
		u.OfWebSearchResultLocation,
		u.OfSearchResultLocation)
}
func (u *TextCitationParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *TextCitationParamUnion) asAny() any {
	if !param.IsOmitted(u.OfCharLocation) {
		return u.OfCharLocation
	} else if !param.IsOmitted(u.OfPageLocation) {
		return u.OfPageLocation
	} else if !param.IsOmitted(u.OfContentBlockLocation) {
		return u.OfContentBlockLocation
	} else if !param.IsOmitted(u.OfWebSearchResultLocation) {
		return u.OfWebSearchResultLocation
	} else if !param.IsOmitted(u.OfSearchResultLocation) {
		return u.OfSearchResultLocation
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetEndCharIndex() *int64 {
	if vt := u.OfCharLocation; vt != nil {
		return &vt.EndCharIndex
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetStartCharIndex() *int64 {
	if vt := u.OfCharLocation; vt != nil {
		return &vt.StartCharIndex
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetEndPageNumber() *int64 {
	if vt := u.OfPageLocation; vt != nil {
		return &vt.EndPageNumber
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetStartPageNumber() *int64 {
	if vt := u.OfPageLocation; vt != nil {
		return &vt.StartPageNumber
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetEncryptedIndex() *string {
	if vt := u.OfWebSearchResultLocation; vt != nil {
		return &vt.EncryptedIndex
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetURL() *string {
	if vt := u.OfWebSearchResultLocation; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetSearchResultIndex() *int64 {
	if vt := u.OfSearchResultLocation; vt != nil {
		return &vt.SearchResultIndex
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetSource() *string {
	if vt := u.OfSearchResultLocation; vt != nil {
		return &vt.Source
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetCitedText() *string {
	if vt := u.OfCharLocation; vt != nil {
		return (*string)(&vt.CitedText)
	} else if vt := u.OfPageLocation; vt != nil {
		return (*string)(&vt.CitedText)
	} else if vt := u.OfContentBlockLocation; vt != nil {
		return (*string)(&vt.CitedText)
	} else if vt := u.OfWebSearchResultLocation; vt != nil {
		return (*string)(&vt.CitedText)
	} else if vt := u.OfSearchResultLocation; vt != nil {
		return (*string)(&vt.CitedText)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetDocumentIndex() *int64 {
	if vt := u.OfCharLocation; vt != nil {
		return (*int64)(&vt.DocumentIndex)
	} else if vt := u.OfPageLocation; vt != nil {
		return (*int64)(&vt.DocumentIndex)
	} else if vt := u.OfContentBlockLocation; vt != nil {
		return (*int64)(&vt.DocumentIndex)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetDocumentTitle() *string {
	if vt := u.OfCharLocation; vt != nil && vt.DocumentTitle.Valid() {
		return &vt.DocumentTitle.Value
	} else if vt := u.OfPageLocation; vt != nil && vt.DocumentTitle.Valid() {
		return &vt.DocumentTitle.Value
	} else if vt := u.OfContentBlockLocation; vt != nil && vt.DocumentTitle.Valid() {
		return &vt.DocumentTitle.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetType() *string {
	if vt := u.OfCharLocation; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfPageLocation; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfContentBlockLocation; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfWebSearchResultLocation; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSearchResultLocation; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetEndBlockIndex() *int64 {
	if vt := u.OfContentBlockLocation; vt != nil {
		return (*int64)(&vt.EndBlockIndex)
	} else if vt := u.OfSearchResultLocation; vt != nil {
		return (*int64)(&vt.EndBlockIndex)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetStartBlockIndex() *int64 {
	if vt := u.OfContentBlockLocation; vt != nil {
		return (*int64)(&vt.StartBlockIndex)
	} else if vt := u.OfSearchResultLocation; vt != nil {
		return (*int64)(&vt.StartBlockIndex)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetTitle() *string {
	if vt := u.OfWebSearchResultLocation; vt != nil && vt.Title.Valid() {
		return &vt.Title.Value
	} else if vt := u.OfSearchResultLocation; vt != nil && vt.Title.Valid() {
		return &vt.Title.Value
	}
	return nil
}

// The properties IsFileUpdate, Type are required.
type TextEditorCodeExecutionCreateResultBlockParam struct {
	IsFileUpdate bool `json:"is_file_update" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "text_editor_code_execution_create_result".
	Type constant.TextEditorCodeExecutionCreateResult `json:"type" default:"text_editor_code_execution_create_result"`
	paramObj
}

func (r TextEditorCodeExecutionCreateResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow TextEditorCodeExecutionCreateResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TextEditorCodeExecutionCreateResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The property Type is required.
type TextEditorCodeExecutionStrReplaceResultBlockParam struct {
	NewLines param.Opt[int64] `json:"new_lines,omitzero"`
	NewStart param.Opt[int64] `json:"new_start,omitzero"`
	OldLines param.Opt[int64] `json:"old_lines,omitzero"`
	OldStart param.Opt[int64] `json:"old_start,omitzero"`
	Lines    []string         `json:"lines,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "text_editor_code_execution_str_replace_result".
	Type constant.TextEditorCodeExecutionStrReplaceResult `json:"type" default:"text_editor_code_execution_str_replace_result"`
	paramObj
}

func (r TextEditorCodeExecutionStrReplaceResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow TextEditorCodeExecutionStrReplaceResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TextEditorCodeExecutionStrReplaceResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Content, ToolUseID, Type are required.
type TextEditorCodeExecutionToolResultBlockParam struct {
	Content   TextEditorCodeExecutionToolResultBlockParamContentUnion `json:"content,omitzero" api:"required"`
	ToolUseID string                                                  `json:"tool_use_id" api:"required"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "text_editor_code_execution_tool_result".
	Type constant.TextEditorCodeExecutionToolResult `json:"type" default:"text_editor_code_execution_tool_result"`
	paramObj
}

func (r TextEditorCodeExecutionToolResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow TextEditorCodeExecutionToolResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TextEditorCodeExecutionToolResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type TextEditorCodeExecutionToolResultBlockParamContentUnion struct {
	OfRequestTextEditorCodeExecutionToolResultError       *TextEditorCodeExecutionToolResultErrorParam       `json:",omitzero,inline"`
	OfRequestTextEditorCodeExecutionViewResultBlock       *TextEditorCodeExecutionViewResultBlockParam       `json:",omitzero,inline"`
	OfRequestTextEditorCodeExecutionCreateResultBlock     *TextEditorCodeExecutionCreateResultBlockParam     `json:",omitzero,inline"`
	OfRequestTextEditorCodeExecutionStrReplaceResultBlock *TextEditorCodeExecutionStrReplaceResultBlockParam `json:",omitzero,inline"`
	paramUnion
}

func (u TextEditorCodeExecutionToolResultBlockParamContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfRequestTextEditorCodeExecutionToolResultError, u.OfRequestTextEditorCodeExecutionViewResultBlock, u.OfRequestTextEditorCodeExecutionCreateResultBlock, u.OfRequestTextEditorCodeExecutionStrReplaceResultBlock)
}
func (u *TextEditorCodeExecutionToolResultBlockParamContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *TextEditorCodeExecutionToolResultBlockParamContentUnion) asAny() any {
	if !param.IsOmitted(u.OfRequestTextEditorCodeExecutionToolResultError) {
		return u.OfRequestTextEditorCodeExecutionToolResultError
	} else if !param.IsOmitted(u.OfRequestTextEditorCodeExecutionViewResultBlock) {
		return u.OfRequestTextEditorCodeExecutionViewResultBlock
	} else if !param.IsOmitted(u.OfRequestTextEditorCodeExecutionCreateResultBlock) {
		return u.OfRequestTextEditorCodeExecutionCreateResultBlock
	} else if !param.IsOmitted(u.OfRequestTextEditorCodeExecutionStrReplaceResultBlock) {
		return u.OfRequestTextEditorCodeExecutionStrReplaceResultBlock
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextEditorCodeExecutionToolResultBlockParamContentUnion) GetErrorCode() *string {
	if vt := u.OfRequestTextEditorCodeExecutionToolResultError; vt != nil {
		return (*string)(&vt.ErrorCode)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextEditorCodeExecutionToolResultBlockParamContentUnion) GetErrorMessage() *string {
	if vt := u.OfRequestTextEditorCodeExecutionToolResultError; vt != nil && vt.ErrorMessage.Valid() {
		return &vt.ErrorMessage.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextEditorCodeExecutionToolResultBlockParamContentUnion) GetContent() *string {
	if vt := u.OfRequestTextEditorCodeExecutionViewResultBlock; vt != nil {
		return &vt.Content
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextEditorCodeExecutionToolResultBlockParamContentUnion) GetFileType() *string {
	if vt := u.OfRequestTextEditorCodeExecutionViewResultBlock; vt != nil {
		return (*string)(&vt.FileType)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextEditorCodeExecutionToolResultBlockParamContentUnion) GetNumLines() *int64 {
	if vt := u.OfRequestTextEditorCodeExecutionViewResultBlock; vt != nil && vt.NumLines.Valid() {
		return &vt.NumLines.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextEditorCodeExecutionToolResultBlockParamContentUnion) GetStartLine() *int64 {
	if vt := u.OfRequestTextEditorCodeExecutionViewResultBlock; vt != nil && vt.StartLine.Valid() {
		return &vt.StartLine.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextEditorCodeExecutionToolResultBlockParamContentUnion) GetTotalLines() *int64 {
	if vt := u.OfRequestTextEditorCodeExecutionViewResultBlock; vt != nil && vt.TotalLines.Valid() {
		return &vt.TotalLines.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextEditorCodeExecutionToolResultBlockParamContentUnion) GetIsFileUpdate() *bool {
	if vt := u.OfRequestTextEditorCodeExecutionCreateResultBlock; vt != nil {
		return &vt.IsFileUpdate
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextEditorCodeExecutionToolResultBlockParamContentUnion) GetLines() []string {
	if vt := u.OfRequestTextEditorCodeExecutionStrReplaceResultBlock; vt != nil {
		return vt.Lines
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextEditorCodeExecutionToolResultBlockParamContentUnion) GetNewLines() *int64 {
	if vt := u.OfRequestTextEditorCodeExecutionStrReplaceResultBlock; vt != nil && vt.NewLines.Valid() {
		return &vt.NewLines.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextEditorCodeExecutionToolResultBlockParamContentUnion) GetNewStart() *int64 {
	if vt := u.OfRequestTextEditorCodeExecutionStrReplaceResultBlock; vt != nil && vt.NewStart.Valid() {
		return &vt.NewStart.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextEditorCodeExecutionToolResultBlockParamContentUnion) GetOldLines() *int64 {
	if vt := u.OfRequestTextEditorCodeExecutionStrReplaceResultBlock; vt != nil && vt.OldLines.Valid() {
		return &vt.OldLines.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextEditorCodeExecutionToolResultBlockParamContentUnion) GetOldStart() *int64 {
	if vt := u.OfRequestTextEditorCodeExecutionStrReplaceResultBlock; vt != nil && vt.OldStart.Valid() {
		return &vt.OldStart.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextEditorCodeExecutionToolResultBlockParamContentUnion) GetType() *string {
	if vt := u.OfRequestTextEditorCodeExecutionToolResultError; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestTextEditorCodeExecutionViewResultBlock; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestTextEditorCodeExecutionCreateResultBlock; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestTextEditorCodeExecutionStrReplaceResultBlock; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// The properties ErrorCode, Type are required.
type TextEditorCodeExecutionToolResultErrorParam struct {
	// Any of "invalid_tool_input", "unavailable", "too_many_requests",
	// "execution_time_exceeded", "file_not_found".
	ErrorCode    TextEditorCodeExecutionToolResultErrorParamErrorCode `json:"error_code,omitzero" api:"required"`
	ErrorMessage param.Opt[string]                                    `json:"error_message,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "text_editor_code_execution_tool_result_error".
	Type constant.TextEditorCodeExecutionToolResultError `json:"type" default:"text_editor_code_execution_tool_result_error"`
	paramObj
}

func (r TextEditorCodeExecutionToolResultErrorParam) MarshalJSON() (data []byte, err error) {
	type shadow TextEditorCodeExecutionToolResultErrorParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TextEditorCodeExecutionToolResultErrorParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type TextEditorCodeExecutionToolResultErrorParamErrorCode string

// The properties Content, FileType, Type are required.
type TextEditorCodeExecutionViewResultBlockParam struct {
	Content string `json:"content" api:"required"`
	// Any of "text", "image", "pdf".
	FileType   TextEditorCodeExecutionViewResultBlockParamFileType `json:"file_type,omitzero" api:"required"`
	NumLines   param.Opt[int64]                                    `json:"num_lines,omitzero"`
	StartLine  param.Opt[int64]                                    `json:"start_line,omitzero"`
	TotalLines param.Opt[int64]                                    `json:"total_lines,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "text_editor_code_execution_view_result".
	Type constant.TextEditorCodeExecutionViewResult `json:"type" default:"text_editor_code_execution_view_result"`
	paramObj
}

func (r TextEditorCodeExecutionViewResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow TextEditorCodeExecutionViewResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TextEditorCodeExecutionViewResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type TextEditorCodeExecutionViewResultBlockParamFileType string

// The properties Signature, Thinking, Type are required.
type ThinkingBlockParam struct {
	// The `signature` value of this thinking block, exactly as returned by the API in
	// a previous response. Used to verify that the block was generated by the model.
	//
	// Thinking blocks must be passed back unmodified and in their original order; a
	// modified block results in a 400 `invalid_request_error`.
	Signature string `json:"signature" api:"required"`
	// The `thinking` text of this block as returned by the API.
	Thinking string `json:"thinking" api:"required"`
	// This field can be elided, and will marshal its zero value as "thinking".
	Type constant.Thinking `json:"type" default:"thinking"`
	paramObj
}

func (r ThinkingBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ThinkingBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ThinkingBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Reference to a single MCP tool by its server and remote name — the same
// `server_name`/`name` pair `mcp_tool_use` carries.
//
// The properties Name, ServerName, Type are required.
type ToolChangeMCPToolReferenceParam struct {
	Name       string `json:"name" api:"required"`
	ServerName string `json:"server_name" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "mcp_tool_reference".
	Type constant.MCPToolReference `json:"type" default:"mcp_tool_reference"`
	paramObj
}

func (r ToolChangeMCPToolReferenceParam) MarshalJSON() (data []byte, err error) {
	type shadow ToolChangeMCPToolReferenceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ToolChangeMCPToolReferenceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Reference to every tool in the named MCP server's toolset.
//
// The properties ServerName, Type are required.
type ToolChangeMCPToolsetReferenceParam struct {
	ServerName string `json:"server_name" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "mcp_toolset_reference".
	Type constant.MCPToolsetReference `json:"type" default:"mcp_toolset_reference"`
	paramObj
}

func (r ToolChangeMCPToolsetReferenceParam) MarshalJSON() (data []byte, err error) {
	type shadow ToolChangeMCPToolsetReferenceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ToolChangeMCPToolsetReferenceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Reference to a single tool the caller declared directly in `tools[]`. Does not
// accept the composed `{server}_{name}` form the server assigns to MCP-resolved
// tools — use `mcp_tool_reference` or `mcp_toolset_reference` for those.
//
// The properties Name, Type are required.
type ToolChangeToolReferenceParam struct {
	Name string `json:"name" api:"required"`
	// This field can be elided, and will marshal its zero value as "tool_reference".
	Type constant.ToolReference `json:"type" default:"tool_reference"`
	paramObj
}

func (r ToolChangeToolReferenceParam) MarshalJSON() (data []byte, err error) {
	type shadow ToolChangeToolReferenceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ToolChangeToolReferenceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Tool reference block that can be included in tool_result content.
//
// The properties ToolName, Type are required.
type ToolReferenceBlockParam struct {
	ToolName string `json:"tool_name" api:"required"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as "tool_reference".
	Type constant.ToolReference `json:"type" default:"tool_reference"`
	paramObj
}

func (r ToolReferenceBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ToolReferenceBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ToolReferenceBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties ToolUseID, Type are required.
type ToolResultBlockParam struct {
	ToolUseID string `json:"tool_use_id" api:"required"`
	// For a toolset member tool_result, the toolset family of the paired tool_use.
	ToolsetName param.Opt[string] `json:"toolset_name,omitzero"`
	IsError     param.Opt[bool]   `json:"is_error,omitzero"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam         `json:"cache_control,omitzero"`
	Content      []ToolResultBlockParamContentUnion `json:"content,omitzero"`
	// This field can be elided, and will marshal its zero value as "tool_result".
	Type constant.ToolResult `json:"type" default:"tool_result"`
	paramObj
}

func (r ToolResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ToolResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ToolResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ToolResultBlockParamContentUnion struct {
	OfText          *TextBlockParam            `json:",omitzero,inline"`
	OfImage         *ImageBlockParam           `json:",omitzero,inline"`
	OfSearchResult  *SearchResultBlockParam    `json:",omitzero,inline"`
	OfDocument      *RequestDocumentBlockParam `json:",omitzero,inline"`
	OfToolReference *ToolReferenceBlockParam   `json:",omitzero,inline"`
	OfBrowserState  *BrowserStateBlockParam    `json:",omitzero,inline"`
	paramUnion
}

func (u ToolResultBlockParamContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfText,
		u.OfImage,
		u.OfSearchResult,
		u.OfDocument,
		u.OfToolReference,
		u.OfBrowserState)
}
func (u *ToolResultBlockParamContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ToolResultBlockParamContentUnion) asAny() any {
	if !param.IsOmitted(u.OfText) {
		return u.OfText
	} else if !param.IsOmitted(u.OfImage) {
		return u.OfImage
	} else if !param.IsOmitted(u.OfSearchResult) {
		return u.OfSearchResult
	} else if !param.IsOmitted(u.OfDocument) {
		return u.OfDocument
	} else if !param.IsOmitted(u.OfToolReference) {
		return u.OfToolReference
	} else if !param.IsOmitted(u.OfBrowserState) {
		return u.OfBrowserState
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolResultBlockParamContentUnion) GetText() *string {
	if vt := u.OfText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolResultBlockParamContentUnion) GetTransformations() *ImageTransformationsParam {
	if vt := u.OfImage; vt != nil {
		return &vt.Transformations
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolResultBlockParamContentUnion) GetContent() []TextBlockParam {
	if vt := u.OfSearchResult; vt != nil {
		return vt.Content
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolResultBlockParamContentUnion) GetContext() *string {
	if vt := u.OfDocument; vt != nil && vt.Context.Valid() {
		return &vt.Context.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolResultBlockParamContentUnion) GetToolName() *string {
	if vt := u.OfToolReference; vt != nil {
		return &vt.ToolName
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolResultBlockParamContentUnion) GetTabs() []BrowserStateTabEntryParam {
	if vt := u.OfBrowserState; vt != nil {
		return vt.Tabs
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolResultBlockParamContentUnion) GetStateChanges() []BrowserStateChangeUnionParam {
	if vt := u.OfBrowserState; vt != nil {
		return vt.StateChanges
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolResultBlockParamContentUnion) GetType() *string {
	if vt := u.OfText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfImage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSearchResult; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDocument; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfToolReference; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfBrowserState; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolResultBlockParamContentUnion) GetTitle() *string {
	if vt := u.OfSearchResult; vt != nil {
		return (*string)(&vt.Title)
	} else if vt := u.OfDocument; vt != nil && vt.Title.Valid() {
		return &vt.Title.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's CacheControl property, if present.
func (u ToolResultBlockParamContentUnion) GetCacheControl() *CacheControlEphemeralParam {
	if vt := u.OfText; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfImage; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfSearchResult; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfDocument; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfToolReference; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfBrowserState; vt != nil {
		return &vt.CacheControl
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ToolResultBlockParamContentUnion) GetCitations() (res toolResultBlockParamContentUnionCitations) {
	if vt := u.OfText; vt != nil {
		res.any = &vt.Citations
	} else if vt := u.OfSearchResult; vt != nil {
		res.any = &vt.Citations
	} else if vt := u.OfDocument; vt != nil {
		res.any = &vt.Citations
	}
	return
}

// Can have the runtime types [*[]TextCitationParamUnion],
// [*CitationsConfigParam]
type toolResultBlockParamContentUnionCitations struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]qoder.TextCitationParamUnion:
//	case *qoder.CitationsConfigParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u toolResultBlockParamContentUnionCitations) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u toolResultBlockParamContentUnionCitations) GetEnabled() *bool {
	switch vt := u.any.(type) {
	case *CitationsConfigParam:
		return paramutil.AddrIfPresent(vt.Enabled)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ToolResultBlockParamContentUnion) GetSource() (res toolResultBlockParamContentUnionSource) {
	if vt := u.OfImage; vt != nil {
		res.any = vt.Source.asAny()
	} else if vt := u.OfSearchResult; vt != nil {
		res.any = &vt.Source
	} else if vt := u.OfDocument; vt != nil {
		res.any = vt.Source.asAny()
	}
	return
}

// Can have the runtime types [*Base64ImageSourceParam],
// [*URLImageSourceParam], [*FileImageSourceParam], [*string],
// [*Base64PDFSourceParam], [*PlainTextSourceParam],
// [*ContentBlockSourceParam], [*URLPDFSourceParam],
// [*FileDocumentSourceParam]
type toolResultBlockParamContentUnionSource struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *qoder.Base64ImageSourceParam:
//	case *qoder.URLImageSourceParam:
//	case *qoder.FileImageSourceParam:
//	case *string:
//	case *qoder.Base64PDFSourceParam:
//	case *qoder.PlainTextSourceParam:
//	case *qoder.ContentBlockSourceParam:
//	case *qoder.URLPDFSourceParam:
//	case *qoder.FileDocumentSourceParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u toolResultBlockParamContentUnionSource) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u toolResultBlockParamContentUnionSource) GetContent() *ContentBlockSourceContentUnionParam {
	switch vt := u.any.(type) {
	case *RequestDocumentBlockSourceUnionParam:
		return vt.GetContent()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u toolResultBlockParamContentUnionSource) GetData() *string {
	switch vt := u.any.(type) {
	case *ImageBlockParamSourceUnion:
		return vt.GetData()
	case *RequestDocumentBlockSourceUnionParam:
		return vt.GetData()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u toolResultBlockParamContentUnionSource) GetMediaType() *string {
	switch vt := u.any.(type) {
	case *ImageBlockParamSourceUnion:
		return vt.GetMediaType()
	case *RequestDocumentBlockSourceUnionParam:
		return vt.GetMediaType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u toolResultBlockParamContentUnionSource) GetType() *string {
	switch vt := u.any.(type) {
	case *ImageBlockParamSourceUnion:
		return vt.GetType()
	case *RequestDocumentBlockSourceUnionParam:
		return vt.GetType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u toolResultBlockParamContentUnionSource) GetURL() *string {
	switch vt := u.any.(type) {
	case *ImageBlockParamSourceUnion:
		return vt.GetURL()
	case *RequestDocumentBlockSourceUnionParam:
		return vt.GetURL()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u toolResultBlockParamContentUnionSource) GetFileID() *string {
	switch vt := u.any.(type) {
	case *ImageBlockParamSourceUnion:
		return vt.GetFileID()
	case *RequestDocumentBlockSourceUnionParam:
		return vt.GetFileID()
	}
	return nil
}

// The properties Content, ToolUseID, Type are required.
type ToolSearchToolResultBlockParam struct {
	Content   ToolSearchToolResultBlockParamContentUnion `json:"content,omitzero" api:"required"`
	ToolUseID string                                     `json:"tool_use_id" api:"required"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "tool_search_tool_result".
	Type constant.ToolSearchToolResult `json:"type" default:"tool_search_tool_result"`
	paramObj
}

func (r ToolSearchToolResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ToolSearchToolResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ToolSearchToolResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ToolSearchToolResultBlockParamContentUnion struct {
	OfRequestToolSearchToolResultError       *ToolSearchToolResultErrorParam       `json:",omitzero,inline"`
	OfRequestToolSearchToolSearchResultBlock *ToolSearchToolSearchResultBlockParam `json:",omitzero,inline"`
	paramUnion
}

func (u ToolSearchToolResultBlockParamContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfRequestToolSearchToolResultError, u.OfRequestToolSearchToolSearchResultBlock)
}
func (u *ToolSearchToolResultBlockParamContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ToolSearchToolResultBlockParamContentUnion) asAny() any {
	if !param.IsOmitted(u.OfRequestToolSearchToolResultError) {
		return u.OfRequestToolSearchToolResultError
	} else if !param.IsOmitted(u.OfRequestToolSearchToolSearchResultBlock) {
		return u.OfRequestToolSearchToolSearchResultBlock
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolSearchToolResultBlockParamContentUnion) GetErrorCode() *string {
	if vt := u.OfRequestToolSearchToolResultError; vt != nil {
		return (*string)(&vt.ErrorCode)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolSearchToolResultBlockParamContentUnion) GetErrorMessage() *string {
	if vt := u.OfRequestToolSearchToolResultError; vt != nil && vt.ErrorMessage.Valid() {
		return &vt.ErrorMessage.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolSearchToolResultBlockParamContentUnion) GetToolReferences() []ToolReferenceBlockParam {
	if vt := u.OfRequestToolSearchToolSearchResultBlock; vt != nil {
		return vt.ToolReferences
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolSearchToolResultBlockParamContentUnion) GetType() *string {
	if vt := u.OfRequestToolSearchToolResultError; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestToolSearchToolSearchResultBlock; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// The properties ErrorCode, Type are required.
type ToolSearchToolResultErrorParam struct {
	// Any of "invalid_tool_input", "unavailable", "too_many_requests",
	// "execution_time_exceeded".
	ErrorCode    ToolSearchToolResultErrorParamErrorCode `json:"error_code,omitzero" api:"required"`
	ErrorMessage param.Opt[string]                       `json:"error_message,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "tool_search_tool_result_error".
	Type constant.ToolSearchToolResultError `json:"type" default:"tool_search_tool_result_error"`
	paramObj
}

func (r ToolSearchToolResultErrorParam) MarshalJSON() (data []byte, err error) {
	type shadow ToolSearchToolResultErrorParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ToolSearchToolResultErrorParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ToolSearchToolResultErrorParamErrorCode string

// The properties ToolReferences, Type are required.
type ToolSearchToolSearchResultBlockParam struct {
	ToolReferences []ToolReferenceBlockParam `json:"tool_references,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "tool_search_tool_search_result".
	Type constant.ToolSearchToolSearchResult `json:"type" default:"tool_search_tool_search_result"`
	paramObj
}

func (r ToolSearchToolSearchResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ToolSearchToolSearchResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ToolSearchToolSearchResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties ID, Input, Name, Type are required.
type ToolUseBlockParam struct {
	ID    string `json:"id" api:"required"`
	Input any    `json:"input,omitzero" api:"required"`
	Name  string `json:"name" api:"required"`
	// For a toolset member tool_use, the toolset family this member belongs to.
	ToolsetName param.Opt[string] `json:"toolset_name,omitzero"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// Tool invocation directly from the model.
	Caller ToolUseBlockParamCallerUnion `json:"caller,omitzero"`
	// This field can be elided, and will marshal its zero value as "tool_use".
	Type constant.ToolUse `json:"type" default:"tool_use"`
	paramObj
}

func (r ToolUseBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ToolUseBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ToolUseBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ToolUseBlockParamCallerUnion struct {
	OfDirect                *DirectCallerParam             `json:",omitzero,inline"`
	OfCodeExecution20250825 *ServerToolCallerParam         `json:",omitzero,inline"`
	OfCodeExecution20260120 *ServerToolCaller20260120Param `json:",omitzero,inline"`
	paramUnion
}

func (u ToolUseBlockParamCallerUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDirect, u.OfCodeExecution20250825, u.OfCodeExecution20260120)
}
func (u *ToolUseBlockParamCallerUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ToolUseBlockParamCallerUnion) asAny() any {
	if !param.IsOmitted(u.OfDirect) {
		return u.OfDirect
	} else if !param.IsOmitted(u.OfCodeExecution20250825) {
		return u.OfCodeExecution20250825
	} else if !param.IsOmitted(u.OfCodeExecution20260120) {
		return u.OfCodeExecution20260120
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolUseBlockParamCallerUnion) GetType() *string {
	if vt := u.OfDirect; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCodeExecution20250825; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCodeExecution20260120; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolUseBlockParamCallerUnion) GetToolID() *string {
	if vt := u.OfCodeExecution20250825; vt != nil {
		return (*string)(&vt.ToolID)
	} else if vt := u.OfCodeExecution20260120; vt != nil {
		return (*string)(&vt.ToolID)
	}
	return nil
}

// The properties Type, URL are required.
type URLImageSourceParam struct {
	URL string `json:"url" api:"required"`
	// This field can be elided, and will marshal its zero value as "url".
	Type constant.URL `json:"type" default:"url"`
	paramObj
}

func (r URLImageSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow URLImageSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *URLImageSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Type, URL are required.
type URLPDFSourceParam struct {
	URL string `json:"url" api:"required"`
	// This field can be elided, and will marshal its zero value as "url".
	Type constant.URL `json:"type" default:"url"`
	paramObj
}

func (r URLPDFSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow URLPDFSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *URLPDFSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Content, Type, URL are required.
type WebFetchBlockParam struct {
	Content RequestDocumentBlockParam `json:"content,omitzero" api:"required"`
	// Fetched content URL
	URL string `json:"url" api:"required"`
	// ISO 8601 timestamp when the content was retrieved
	RetrievedAt param.Opt[string] `json:"retrieved_at,omitzero"`
	// This field can be elided, and will marshal its zero value as "web_fetch_result".
	Type constant.WebFetchResult `json:"type" default:"web_fetch_result"`
	paramObj
}

func (r WebFetchBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow WebFetchBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WebFetchBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Content, ToolUseID, Type are required.
type WebFetchToolResultBlockParam struct {
	Content   WebFetchToolResultBlockParamContentUnion `json:"content,omitzero" api:"required"`
	ToolUseID string                                   `json:"tool_use_id" api:"required"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// Tool invocation directly from the model.
	Caller WebFetchToolResultBlockParamCallerUnion `json:"caller,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "web_fetch_tool_result".
	Type constant.WebFetchToolResult `json:"type" default:"web_fetch_tool_result"`
	paramObj
}

func (r WebFetchToolResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow WebFetchToolResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WebFetchToolResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type WebFetchToolResultBlockParamContentUnion struct {
	OfRequestWebFetchToolResultError *WebFetchToolResultErrorBlockParam `json:",omitzero,inline"`
	OfRequestWebFetchResultBlock     *WebFetchBlockParam                `json:",omitzero,inline"`
	paramUnion
}

func (u WebFetchToolResultBlockParamContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfRequestWebFetchToolResultError, u.OfRequestWebFetchResultBlock)
}
func (u *WebFetchToolResultBlockParamContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *WebFetchToolResultBlockParamContentUnion) asAny() any {
	if !param.IsOmitted(u.OfRequestWebFetchToolResultError) {
		return u.OfRequestWebFetchToolResultError
	} else if !param.IsOmitted(u.OfRequestWebFetchResultBlock) {
		return u.OfRequestWebFetchResultBlock
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u WebFetchToolResultBlockParamContentUnion) GetErrorCode() *string {
	if vt := u.OfRequestWebFetchToolResultError; vt != nil {
		return (*string)(&vt.ErrorCode)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u WebFetchToolResultBlockParamContentUnion) GetContent() *RequestDocumentBlockParam {
	if vt := u.OfRequestWebFetchResultBlock; vt != nil {
		return &vt.Content
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u WebFetchToolResultBlockParamContentUnion) GetURL() *string {
	if vt := u.OfRequestWebFetchResultBlock; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u WebFetchToolResultBlockParamContentUnion) GetRetrievedAt() *string {
	if vt := u.OfRequestWebFetchResultBlock; vt != nil && vt.RetrievedAt.Valid() {
		return &vt.RetrievedAt.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u WebFetchToolResultBlockParamContentUnion) GetType() *string {
	if vt := u.OfRequestWebFetchToolResultError; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestWebFetchResultBlock; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type WebFetchToolResultBlockParamCallerUnion struct {
	OfDirect                *DirectCallerParam             `json:",omitzero,inline"`
	OfCodeExecution20250825 *ServerToolCallerParam         `json:",omitzero,inline"`
	OfCodeExecution20260120 *ServerToolCaller20260120Param `json:",omitzero,inline"`
	paramUnion
}

func (u WebFetchToolResultBlockParamCallerUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDirect, u.OfCodeExecution20250825, u.OfCodeExecution20260120)
}
func (u *WebFetchToolResultBlockParamCallerUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *WebFetchToolResultBlockParamCallerUnion) asAny() any {
	if !param.IsOmitted(u.OfDirect) {
		return u.OfDirect
	} else if !param.IsOmitted(u.OfCodeExecution20250825) {
		return u.OfCodeExecution20250825
	} else if !param.IsOmitted(u.OfCodeExecution20260120) {
		return u.OfCodeExecution20260120
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u WebFetchToolResultBlockParamCallerUnion) GetType() *string {
	if vt := u.OfDirect; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCodeExecution20250825; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCodeExecution20260120; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u WebFetchToolResultBlockParamCallerUnion) GetToolID() *string {
	if vt := u.OfCodeExecution20250825; vt != nil {
		return (*string)(&vt.ToolID)
	} else if vt := u.OfCodeExecution20260120; vt != nil {
		return (*string)(&vt.ToolID)
	}
	return nil
}

// The properties ErrorCode, Type are required.
type WebFetchToolResultErrorBlockParam struct {
	// Any of "invalid_tool_input", "url_too_long", "url_not_allowed",
	// "url_not_in_prior_context", "url_not_accessible", "unsupported_content_type",
	// "too_many_requests", "max_uses_exceeded", "unavailable".
	ErrorCode WebFetchToolResultErrorCode `json:"error_code,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "web_fetch_tool_result_error".
	Type constant.WebFetchToolResultError `json:"type" default:"web_fetch_tool_result_error"`
	paramObj
}

func (r WebFetchToolResultErrorBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow WebFetchToolResultErrorBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WebFetchToolResultErrorBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WebFetchToolResultErrorCode string

// The properties EncryptedContent, Title, Type, URL are required.
type WebSearchResultBlockParam struct {
	EncryptedContent string            `json:"encrypted_content" api:"required"`
	Title            string            `json:"title" api:"required"`
	URL              string            `json:"url" api:"required"`
	PageAge          param.Opt[string] `json:"page_age,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "web_search_result".
	Type constant.WebSearchResult `json:"type" default:"web_search_result"`
	paramObj
}

func (r WebSearchResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow WebSearchResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WebSearchResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties ErrorCode, Type are required.
type WebSearchToolRequestErrorParam struct {
	// Any of "invalid_tool_input", "unavailable", "max_uses_exceeded",
	// "too_many_requests", "query_too_long", "request_too_large".
	ErrorCode WebSearchToolResultErrorCode `json:"error_code,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as
	// "web_search_tool_result_error".
	Type constant.WebSearchToolResultError `json:"type" default:"web_search_tool_result_error"`
	paramObj
}

func (r WebSearchToolRequestErrorParam) MarshalJSON() (data []byte, err error) {
	type shadow WebSearchToolRequestErrorParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WebSearchToolRequestErrorParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Content, ToolUseID, Type are required.
type WebSearchToolResultBlockParam struct {
	Content   WebSearchToolResultBlockParamContentUnion `json:"content,omitzero" api:"required"`
	ToolUseID string                                    `json:"tool_use_id" api:"required"`
	// Create a cache control breakpoint at this content block.
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// Tool invocation directly from the model.
	Caller WebSearchToolResultBlockParamCallerUnion `json:"caller,omitzero"`
	// This field can be elided, and will marshal its zero value as
	// "web_search_tool_result".
	Type constant.WebSearchToolResult `json:"type" default:"web_search_tool_result"`
	paramObj
}

func (r WebSearchToolResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow WebSearchToolResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WebSearchToolResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type WebSearchToolResultBlockParamCallerUnion struct {
	OfDirect                *DirectCallerParam             `json:",omitzero,inline"`
	OfCodeExecution20250825 *ServerToolCallerParam         `json:",omitzero,inline"`
	OfCodeExecution20260120 *ServerToolCaller20260120Param `json:",omitzero,inline"`
	paramUnion
}

func (u WebSearchToolResultBlockParamCallerUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDirect, u.OfCodeExecution20250825, u.OfCodeExecution20260120)
}
func (u *WebSearchToolResultBlockParamCallerUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *WebSearchToolResultBlockParamCallerUnion) asAny() any {
	if !param.IsOmitted(u.OfDirect) {
		return u.OfDirect
	} else if !param.IsOmitted(u.OfCodeExecution20250825) {
		return u.OfCodeExecution20250825
	} else if !param.IsOmitted(u.OfCodeExecution20260120) {
		return u.OfCodeExecution20260120
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u WebSearchToolResultBlockParamCallerUnion) GetType() *string {
	if vt := u.OfDirect; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCodeExecution20250825; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCodeExecution20260120; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u WebSearchToolResultBlockParamCallerUnion) GetToolID() *string {
	if vt := u.OfCodeExecution20250825; vt != nil {
		return (*string)(&vt.ToolID)
	} else if vt := u.OfCodeExecution20260120; vt != nil {
		return (*string)(&vt.ToolID)
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type WebSearchToolResultBlockParamContentUnion struct {
	OfResultBlock []WebSearchResultBlockParam     `json:",omitzero,inline"`
	OfError       *WebSearchToolRequestErrorParam `json:",omitzero,inline"`
	paramUnion
}

func (u WebSearchToolResultBlockParamContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfResultBlock, u.OfError)
}
func (u *WebSearchToolResultBlockParamContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *WebSearchToolResultBlockParamContentUnion) asAny() any {
	if !param.IsOmitted(u.OfResultBlock) {
		return &u.OfResultBlock
	} else if !param.IsOmitted(u.OfError) {
		return u.OfError
	}
	return nil
}

type WebSearchToolResultErrorCode string

// Model identifies a model available to the API.
type Model = string
