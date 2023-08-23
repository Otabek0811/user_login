package models

type User struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	Age       int    `json:"age"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserPrimaryKey struct {
	Id    string `json:"id"`
	Login string `json:"login"`
	Name  string  `json:"name"`
	UserID string 	`json:"user_id"`
}

type CreateUser struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Age      int    `json:"age"`
}

type UpdateUser struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Age      int    `json:"age"`
}

type GetListUserRequest struct {
	UserID string `json:"user_id"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type GetListUserResponse struct {
	Count int     `json:"count"`
	Users []*User `json:"users"`
}
