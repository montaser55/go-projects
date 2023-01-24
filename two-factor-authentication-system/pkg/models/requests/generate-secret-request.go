package requests

type GenerateSecretRequest struct {
	UserId         int64  `json:"userId"`
	TfaChannelType string `json:"tfaChannelType"`
}
