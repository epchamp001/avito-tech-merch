package models

type UserInfoResponse struct {
	Username     string         `json:"username"`
	Balance      int            `json:"balance"`
	Purchases    []*Purchase    `json:"purchases"`
	Transactions []*Transaction `json:"transactions"`
}
