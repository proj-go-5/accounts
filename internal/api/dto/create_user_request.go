package dto

type CreateUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
