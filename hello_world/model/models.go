package model

type Bill struct {
	ID       string `gorm:"primarykey"`
	State    bool   `gorm:"default:true"`
	Title    string
	Sum      float32 `gorm:"default:0"`
	Currency uint32  `gorm:"default:1"`
}
