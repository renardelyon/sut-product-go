package model

type Status string

const (
	Pending  Status = "pending"
	Rejected Status = "rejected"
	Accepted Status = "Accepted"
)

type UserProduct struct {
	AdminId     string
	Fullname    string
	Username    string
	UserId      string
	ProductId   string
	RequestDate string
	Status      Status
	Product     Product `gorm:"foreignKey:ProductId"`
}

type Product struct {
	Id   string `gorm:"primaryKey"`
	Name string
	Url  string
}
