//go:generate go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config codegen.yaml openapi.yaml
package apispec

import "encoding/json"

// XXX: Workaround for https://github.com/oapi-codegen/oapi-codegen/issues/970
func (r GetTeamPaymentMethod200JSONResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(TeamPaymentMethod(r))
}

// XXX: Workaround for https://github.com/oapi-codegen/oapi-codegen/issues/970
func (r PutTeamPaymentMethod200JSONResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(TeamPaymentMethod(r))
}
