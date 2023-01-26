package responses

import "github.com/montaser55/two-factor-authentication-service/pkg/utils/enums"

type GenerateSecretResponse struct {
	ReferenceId         string               `json:"referenceId"`
	QrCodeMessage       string               `json:"qrCodeMessage"`
	PlainSecret         string               `json:"plainSecret"`
	TfaChannelType      enums.TfaChannelType `json:"tfaChannelType"`
	ExpiryTimeInSeconds int                  `json:"expiryTimeInSeconds"`
}
