package models

type Merch struct {
	ID    int    `json:"id"   gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name" gorm:"uniqueIndex;not null"`
	Price int    `json:"price" gorm:"not null;check:price > 0"`
}
