package forward

import (
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

// OutcomeRubricParam defines acceptance criteria as inline text or an uploaded file.
// Set Type to "text" with Content, or "file" with FileID.
type OutcomeRubricParam struct {
	Type    string            `json:"type" api:"required"`
	Content param.Opt[string] `json:"content,omitzero"`
	FileID  param.Opt[string] `json:"file_id,omitzero"`
	paramObj
}

func (r OutcomeRubricParam) MarshalJSON() ([]byte, error) {
	type shadow OutcomeRubricParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OutcomeRubricParam) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// OutcomeRubric contains the text or file reference used to evaluate an Outcome.
type OutcomeRubric struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	FileID  string `json:"file_id"`
	JSON    struct {
		Type        respjson.Field
		Content     respjson.Field
		FileID      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r OutcomeRubric) RawJSON() string                  { return r.JSON.raw }
func (r *OutcomeRubric) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
