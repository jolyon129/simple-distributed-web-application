package repository

import (
	"log"
	"time"
	"zl2501-final-project/web/model/storage"
)

//TODO:
// Add lock!
type PostRepo struct {
	Storage storage.PostStorageInterface
}

type PostInfo struct {
	UserID  uint
	Content string
}

func (postRepo *PostRepo) CreateNewPost(p PostInfo) (uint, error) {
	ID, err := postRepo.Storage.Create(&storage.PostEntity{
		ID:          0,
		UserID:      p.UserID,
		Content:     p.Content,
		CreatedTime: time.Time{},
	})
	if err != nil {
		log.Println(err)
		return 0, err
	} else {
		return ID, nil
	}
}

// Return a post Entity according the post id
func (postRepo *PostRepo) SelectById(pId uint) *storage.PostEntity {
	p, err := postRepo.Storage.Read(pId)
	if err != nil {
		log.Println(err)
		return nil
	} else {
		return &p
	}
}
