package models

type Phone struct {
	Id          string `json:"id"`
	UserID      string `json:"user_id"`
	Phone       string `json:"phone"`
	Description string `json:"description"`
	IsFax       bool   `json:"is_fax"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type PhonePrimaryKey struct {
	Id     string `json:"id"`
}

type CreatePhone struct {
	UserID      string `json:"user_id"`
	Phone       string `json:"phone"`
	Description string `json:"description"`
	IsFax       bool   `json:"is_fax"`
}

type UpdatePhone struct {
	Id          string `json:"id"`
	UserID      string `json:"user_id"`
	Phone       string `json:"phone"`
	Description string `json:"description"`
	IsFax       bool   `json:"is_fax"`
}

type GetListPhoneRequest struct {
	UserID string `json:"user_id"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type GetListPhoneResponse struct {
	Count  int      `json:"count"`
	Phones []*Phone `json:"Phones"`
}
