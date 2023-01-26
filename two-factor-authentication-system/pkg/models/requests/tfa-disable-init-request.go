package requests

import "github.com/montaser55/two-factor-authentication-service/pkg/utils/enums"

type TfaDisableInitRequest struct {
	UserId         int64                `json:"userId"`
	TfaChannelType enums.TfaChannelType `json:"tfaChannelType"`
}
