package dao

type User struct {
	Id        string `gorm:"primaryKey"`
	Email     string `gorm:"unique;size:256"`
	Password  string
	Nickname  string
	Bio       string
	Gender    int
	Birthday  int
	CreatedAt int64 `gorm:"autoCreateTime:milli"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli"`
	DeletedAt int64
}
