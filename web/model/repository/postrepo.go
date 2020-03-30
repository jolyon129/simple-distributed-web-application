package repository

import (
	"log"
	"sync"
	"time"
	"zl2501-final-project/web/model/storage"
)


type PostRepo struct {
	sync.Mutex
	Storage storage.PostStorageInterface
}

type PostInfo struct {
	UserID  uint
	Content string
}

func (postRepo *PostRepo) CreateNewPost(p PostInfo) (uint, error) {
	postRepo.Lock()
	defer postRepo.Unlock()
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
	postRepo.Lock()
	postRepo.Unlock()
	p, err := postRepo.Storage.Read(pId)
	if err != nil {
		log.Println(err)
		return nil
	} else {
		return &p
	}
}
