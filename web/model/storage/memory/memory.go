package memory

import (
	"container/list"
	"zl2501-final-project/web/model/storage"
)

func init() {
	var memUserStorage storage.UserStorageInterface
	var memPostStorage storage.PostStorageInterface
	memUserStorage = &MemUserStore{
		userMap: make(map[uint]*storage.UserEntity),
		//users:       list.New(),
		userNameSet: make(map[string]bool),
		pkCounter:   100, // Start from 100
	}
	memPostStorage = &MemPostStore{
		postMap:   make(map[uint]*storage.PostEntity),
		posts:     list.New(),
		pkCounter: 1000, // Start from 100
	}
	memModels := storage.Manager{
		UserStorage: memUserStorage,
		PostStorage: memPostStorage,
	}
	// Register the implementation of storage manager.
	storage.RegisterDriver("memory", &memModels)
}
