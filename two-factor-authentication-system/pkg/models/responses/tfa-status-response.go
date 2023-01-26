package responses

import "github.com/montaser55/two-factor-authentication-service/pkg/utils/enums"

type TfaStatusResponse struct {
	IsEnabled      bool                 `json:"isEnabled"`
	TfaChannelType enums.TfaChannelType `json:"tfaChannelType"`
}
