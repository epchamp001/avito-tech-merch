package models

type UserInfo struct {
	UserID       int            `json:"user_id"`
	Username     string         `json:"username"`
	Balance      int            `json:"balance"`
	Purchases    []*Purchase    `json:"purchases"`
	Transactions []*Transaction `json:"transactions"`
}
