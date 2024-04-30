package rdbrepository

import (
	"context"
	"encoding/json"
	"fmt"
	"meet/internal/pkg/app/model"
	"meet/internal/pkg/app/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	userKeyLogin = "user:login:%s"
	userKeyTgID  = "user:tgID:%d"
)

var userCacheDuration = time.Hour * 1

type userRedisRepository struct {
	rdb *redis.Client
}

func NewUserRedisRepository(rdb *redis.Client) repository.UserRepository {
	return &userRedisRepository{rdb}
}

func (ur *userRedisRepository) GetByLogin(login string) (*model.User, error) {
	result, err := ur.rdb.Get(context.Background(), fmt.Sprintf(userKeyLogin, login)).Result()
	if err == redis.Nil {
		return nil, repository.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	u := new(model.User)
	err = json.Unmarshal([]byte(result), u)

	return u, err
}

func (ur *userRedisRepository) HasByLogin(login string) (bool, error) {
	count, err := ur.rdb.Exists(context.Background(), fmt.Sprintf(userKeyLogin, login)).Result()

	return count > 0, err
}

func (ur *userRedisRepository) GetByTgID(id int64) (*model.User, error) {
	result, err := ur.rdb.Get(context.Background(), fmt.Sprintf(userKeyTgID, id)).Result()
	if err == redis.Nil {
		return nil, repository.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	u := new(model.User)
	err = json.Unmarshal([]byte(result), u)

	return u, err
}

func (ur *userRedisRepository) HasByTgID(id int64) (bool, error) {
	count, err := ur.rdb.Exists(context.Background(), fmt.Sprintf(userKeyTgID, id)).Result()

	return count > 0, err
}

func (ur *userRedisRepository) Add(ctx context.Context, user *model.User) error {
	user.BeforeAdd()

	return ur.set(ctx, user)
}

func (ur *userRedisRepository) Update(user *model.User) error {
	user.BeforeUpdate()

	return ur.set(context.Background(), user)
}

func (ur *userRedisRepository) set(ctx context.Context, user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if err := ur.rdb.Set(ctx, fmt.Sprintf(userKeyLogin, user.Login), string(data), userCacheDuration).Err(); err != nil {
		return err
	}

	if !user.TgID.Valid {
		return nil
	}

	return ur.rdb.Set(ctx, fmt.Sprintf(userKeyTgID, user.TgID.Int64), string(data), userCacheDuration).Err()
}

func (ur *userRedisRepository) Remove(user *model.User) error {
	return ur.rdb.Del(context.Background(), fmt.Sprintf(userKeyLogin, user.Login), fmt.Sprintf(userKeyTgID, user.TgID.Int64)).Err()
}
