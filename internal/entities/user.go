package entities

type Admin struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

type AdminWithPassword struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u *AdminWithPassword) WithoutPassword() *Admin {
	return &Admin{
		ID:    u.ID,
		Login: u.Login,
	}
}

type AdminClaims struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}
