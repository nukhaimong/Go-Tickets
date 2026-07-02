package dto

type Response struct {
	ID        uint   `json:"Id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Token     string `json:"token,omitempty"`
	CreatedAt string `json:"created_at"`
}
