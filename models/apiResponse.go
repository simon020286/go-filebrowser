package models

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

// APIResponse struct definition.
type APIResponse struct {
	Data         interface{} `json:"data"`
	Success      bool        `json:"success"`
	ErrorMessage string      `json:"errorMsg"`
	Error        error       `json:"error"`
	Status       uint16      `json:"-"`
}

// Send response.
func (apiResponse *APIResponse) Send(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	str, err := json.Marshal(apiResponse)
	if err == nil {
		ctx.Write(str)
	}
}

// NewAPIResponse for constructor APIResponse.
func NewAPIResponse(options ...func(*APIResponse)) *APIResponse {
	apiResponse := &APIResponse{}

	// Default values
	apiResponse.Success = true
	apiResponse.ErrorMessage = ""
	apiResponse.Error = nil
	apiResponse.Status = 200

	// Option parameters values:
	for _, op := range options {
		op(apiResponse)
	}

	return apiResponse
}

// OptionData default.
func OptionData(data interface{}) func(apiResponse *APIResponse) {
	return func(apiResponse *APIResponse) {
		apiResponse.Data = data
	}
}

// OptionError default.
func OptionError(errorMsg string, err error) func(apiResponse *APIResponse) {
	return func(apiResponse *APIResponse) {
		apiResponse.Error = err
		apiResponse.ErrorMessage = errorMsg
		apiResponse.Status = 500
	}
}
