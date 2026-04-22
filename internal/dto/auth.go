package dto

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Nickname    string `json:"nickname"`
	Affiliation string `json:"affiliation"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
