package memory

import (
	"zl2501-final-project/raftcluster/store/backendstore"
)

func init() {
	var memUserStorage backendstore.UserStorageInterface
	var memPostStorage backendstore.TweetStorageInterface
	memUserStorage = &MemUserStore{
		userMap: make(map[uint]*backendstore.UserEntity),
		//users:       list.New(),
		userNameSet: make(map[string]bool),
		pkCounter:   100, // Start from 100
	}
	memPostStorage = &MemTweetStore{
		tweetMap:  make(map[uint]*backendstore.TweetEntity),
		pkCounter: 1000, // Start from 100
	}
	memModels := backendstore.Manager{
		UserStorage:  memUserStorage,
		TweetStorage: memPostStorage,
	}
	// Register the implementation of storage manager.
	backendstore.RegisterDriver("memory", &memModels)
}
