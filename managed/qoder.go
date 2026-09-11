// Qoder managed API definitions.
package managed

import (
	"encoding/json"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

type QoderBeta = string

type Currency string

// A monetary amount in a specific currency.
type MonetaryAmount struct {
	// Amount in minor units of the currency, as an integer decimal string with no
	// leading zeros: "2500" is $25.00 and "50" is fifty cents. A string rather than a
	// number so no float rounding is ever applied.
	Amount string `json:"amount" api:"required"`
	// Uppercase ISO-4217 currency code. `USD` is the only currency currently
	// supported; the accepted set is closed and grows only when a new currency is
	// priced.
	//
	// Any of "USD".
	Currency Currency `json:"currency" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Amount      respjson.Field
		Currency    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r MonetaryAmount) RawJSON() string { return r.JSON.raw }
func (r *MonetaryAmount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this MonetaryAmount to a MonetaryAmountParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// MonetaryAmountParam.Overrides()
func (r MonetaryAmount) ToParam() MonetaryAmountParam {
	return param.Override[MonetaryAmountParam](json.RawMessage(r.RawJSON()))
}

// A monetary amount in a specific currency.
//
// The properties Amount, Currency are required.
type MonetaryAmountParam struct {
	// Amount in minor units of the currency, as an integer decimal string with no
	// leading zeros: "2500" is $25.00 and "50" is fifty cents. A string rather than a
	// number so no float rounding is ever applied.
	Amount string `json:"amount" api:"required"`
	// Uppercase ISO-4217 currency code. `USD` is the only currency currently
	// supported; the accepted set is closed and grows only when a new currency is
	// priced.
	//
	// Any of "USD".
	Currency Currency `json:"currency,omitzero" api:"required"`
	paramObj
}

func (r MonetaryAmountParam) MarshalJSON() (data []byte, err error) {
	type shadow MonetaryAmountParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MonetaryAmountParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

const CurrencyUSD Currency = "USD"
