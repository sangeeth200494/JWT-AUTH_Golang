package models

import "fmt"

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

type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

func (e *APIResponse) Error() string {
	return fmt.Sprintf("Error%d:%s", e.Code, e.Message)
}
