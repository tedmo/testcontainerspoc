package app

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type CreateUserReq struct {
	Name string `json:"name"`
}
