package repository

import (
	"context"
	"errors"
	"sync"
	"zl2501-final-project/backend/model/storage"
)

// This is a singleton type
type UserRepo struct {
	//	sync.Mutex // Need lock
	storage storage.UserStorageInterface
	con     *sync.Cond
}

type UserInfo struct {
	UserName string
	Password string
}

func NewUserRepo(storageInterface storage.UserStorageInterface) *UserRepo {
	ret := UserRepo{
		storage: nil,
	}
	ret.storage = storageInterface
	return &ret
}

// Create a new user and return user id
// if the user name is not duplicated.
// Otherwise return error
func (userRepo *UserRepo) CreateNewUser(ctx context.Context, u *UserInfo) (uint, error) {
	result := make(chan uint)
	errorChan := make(chan error)
	go func() {
		userRepo.storage.Create(&storage.UserEntity{
			ID:       0,
			UserName: u.UserName,
			Password: u.Password,
		}, result, errorChan)
	}()

	select {
	case ret := <-result:
		return ret, nil
	case err := <-errorChan:
		return 0, err
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func (userRepo *UserRepo) SelectByName(ctx context.Context, name string) (*storage.UserEntity, error) {
	result := make(chan []*storage.UserEntity)
	errorChan := make(chan error)
	go func() {
		userRepo.storage.FindAll(result, errorChan)
	}()
	select {
	case users := <-result:
		for _, value := range users {
			if value.UserName == name {
				return value, nil
			}
		}
		return nil, errors.New("the name does not exist")
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (userRepo *UserRepo) SelectById(ctx context.Context, uid uint) (*storage.UserEntity, error) {
	result := make(chan *storage.UserEntity)
	errorChan := make(chan error)
	go func() {
		userRepo.storage.Read(uid, result, errorChan)
	}()
	select {
	case ret := <-result:
		return ret, nil
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Add the tweet id into the user
func (u *UserRepo) AddTweetToUser(ctx context.Context, uId uint, pId uint) (bool, error) {
	result := make(chan bool)
	errorChan := make(chan error)
	go func() {
		u.storage.AddTweetToUserDB(uId, pId, result, errorChan)
	}()
	select {
	case <-result:
		return true, nil
	case err := <-errorChan:
		return false, err
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

// Return all users in the database
func (u *UserRepo) FindAllUsers(ctx context.Context) ([]*storage.UserEntity, error) {
	result := make(chan []*storage.UserEntity)
	errorChan := make(chan error)
	go func() {
		u.storage.FindAll(result, errorChan)
	}()
	select {
	case users := <-result:
		return users, nil
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Check whether the user srcId follows the user targetId.
// Take O(#following) time
func (u *UserRepo) CheckWhetherFollowing(ctx context.Context, srcId uint, targetId uint) (bool,
	error) {
	result := make(chan bool)
	errorChan := make(chan error)
	go func() {
		u.storage.CheckWhetherFollowingDB(srcId, targetId, result, errorChan)
	}()
	select {
	case ret := <-result:
		return ret, nil
	case err := <-errorChan:
		return false, err
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

// User srcId starts to follow targetId.
func (u *UserRepo) StartFollowing(ctx context.Context, srcId uint, targetId uint) (bool, error) {
	result := make(chan bool)
	errorChan := make(chan error)
	go func() {
		u.storage.StartFollowingDB(srcId, targetId, result, errorChan)
	}()
	select {
	case err := <-errorChan:
		return false, err
	case ret := <-result:
		return ret, nil
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

// srcId stop following targetId.
// targetId remove the follower srcId.
func (u *UserRepo) StopFollowing(ctx context.Context, srcId uint, targetId uint) (bool, error) {
	result := make(chan bool)
	errorChan := make(chan error)
	go func() {
		u.storage.StopFollowingDB(srcId, targetId, result, errorChan)
	}()
	select {
	case ret := <-result:
		return ret, nil
	case err := <-errorChan:
		return false, err
	case <-ctx.Done():
		return false, nil
	}
}
