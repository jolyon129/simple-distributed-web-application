package storage

import "time"

type User struct {
	ID        int // The DB will fill this field
	UserName  string
	Password  string
	followers []int
	Posts     []int
}

type Post struct {
	ID         int
	Content    string
	createTime time.Time
}

type UserStoreInterface interface {
	Create(user *User) int
	Delete(user *User)
	Get(ID int) *User
	Update(ID int, user *User) int
}

type PostStoreInterface interface {
	Create(post *Post)
	Delete(post *Post)
	Get(ID int) *Post
	Update(ID int, post *Post) int
}
