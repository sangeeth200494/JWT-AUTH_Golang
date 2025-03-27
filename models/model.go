package models

// type User struct {
// 	UserName string
// 	Password string
// }

type User struct {
	ID       int64  `gorm:"primaryKey" json:"user_id"`
	Username string `gorm:"size:100;not null" json:"username"`
	Password string `gorm:"unique;not null" json:"password"`
}

// type Article struct {
// 	ID       string `json:”id”`
// 	Title    string `json:”title”`
// 	Body     string `json:”body”`
// 	AuthorID string `json:”author_id”`
// }
