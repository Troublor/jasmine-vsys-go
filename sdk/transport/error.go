package transport

import (
	"encoding/json"
	"fmt"
	sdkErr "github.com/Troublor/jasmine-vsys-go/sdk/error"
	"io/ioutil"
	"net/http"
)

type VsysErr struct {
	sdkErr.Err
	Code int
	Data map[string]interface{}
	Raw  string
}

func NewVsysError(httpResp *http.Response) VsysErr {
	if httpResp == nil {
		return VsysErr{
			Err: sdkErr.Err{
				Msg: "no http response",
			},
			Code: -1,
		}
	}

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return VsysErr{
			Err: sdkErr.Err{
				Msg: "failed to read http response body, " + err.Error(),
			},
			Code: -1,
		}
	}

	var errResp ErrorResponse
	err = json.Unmarshal(body, &errResp)
	if err != nil {
		return VsysErr{
			Err: sdkErr.Err{
				Msg: "failed to unmarshal http response body '" + string(body) + "', " + err.Error(),
			},
			Code: -1,
		}
	}

	return VsysErr{
		Err: sdkErr.Err{
			Msg: errResp.Message,
		},
		Code: errResp.Code,
		Data: errResp.Data,
		Raw:  string(body),
	}
}

func (e VsysErr) Error() string {
	return fmt.Sprintf("Vsys error %d: %s, %s", e.Code, e.Msg, e.Raw)
}
