package requests

import "github.com/montaser55/two-factor-authentication-service/pkg/utils/enums"

type TfaDisableRequest struct {
	UserId         int64                `json:"userId"`
	TfaChannelType enums.TfaChannelType `json:"tfaChannelType"`
	ReferenceId    string               `json:"referenceId"`
	Otp            string               `json:"otp"`
}
