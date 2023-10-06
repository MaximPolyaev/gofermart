package entities

type AuthPayload struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
