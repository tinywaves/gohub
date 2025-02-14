package user

// Entity PO/Model/Entity
type Entity struct {
	Id         string `gorm:"primaryKey"`
	Email      string `gorm:"uniqueIndex;size:256"`
	Password   string
	Nickname   string
	Bio        string
	Gender     string
	Birthday   int64
	CreateTime int64
	UpdateTime int64
}

func (Entity) TableName() string {
	return "users"
}
