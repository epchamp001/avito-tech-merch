package models

type UserInfoResponse struct {
	Balance      int           `json:"balance"`
	Purchases    []Purchase    `json:"purchases"`
	Transactions []Transaction `json:"transactions"`
}
