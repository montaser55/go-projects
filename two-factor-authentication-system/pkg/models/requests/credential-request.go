package requests

type CredentialRequest struct {
	UserId  int64  `json:"userId"`
	PinHash string `json:"pinHash"`
}
