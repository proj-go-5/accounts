package entities

type User struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

type UserWithPassword struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u *UserWithPassword) WithoutPassword() *User {
	return &User{
		ID:    u.ID,
		Login: u.Login,
	}
}

type UserClaims struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}
