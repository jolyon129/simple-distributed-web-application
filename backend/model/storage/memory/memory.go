package memory

import (
	"zl2501-final-project/backend/model/storage"
)

func init() {
	var memUserStorage storage.UserStorageInterface
	var memPostStorage storage.TweetStorageInterface
	memUserStorage = &MemUserStore{
		userMap: make(map[uint]*storage.UserEntity),
		//users:       list.New(),
		userNameSet: make(map[string]bool),
		pkCounter:   100, // Start from 100
	}
	memPostStorage = &MemTweetStore{
		tweetMap:  make(map[uint]*storage.TweetEntity),
		pkCounter: 1000, // Start from 100
	}
	memModels := storage.Manager{
		UserStorage:  memUserStorage,
		TweetStorage: memPostStorage,
	}
	// Register the implementation of storage manager.
	storage.RegisterDriver("memory", &memModels)
}
