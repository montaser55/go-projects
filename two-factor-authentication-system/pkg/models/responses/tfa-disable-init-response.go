package responses

type TfaDisableInitResponse struct {
	ReferenceId         string `json:"referenceId"`
	ExpiryTimeInSeconds int    `json:"expiryTimeInSeconds"`
}
