package domain

type User struct {
	UserId   string
	Email    string
	Password string
	Nickname string
	Bio      string
	Gender   int // 0 - male, 1 - female
	Birthday int
}
