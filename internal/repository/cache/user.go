package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/tinywaves/gohub/internal/domain"
)

type UserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func InitUserCache(client redis.Cmdable, expiration time.Duration) *UserCache {
	return &UserCache{
		client:     client,
		expiration: expiration,
	}
}

func (uCache *UserCache) Get(ctx context.Context, id string) (domain.User, error) {
	unmarshaled, err := uCache.client.Get(ctx, uCache.domainKey(id)).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err = json.Unmarshal(unmarshaled, &user)
	return user, err
}

func (uCache *UserCache) Set(ctx context.Context, uDomain domain.User) error {
	marshaled, err := json.Marshal(uDomain)
	if err != nil {
		return err
	}
	return uCache.client.Set(ctx, uCache.domainKey(uDomain.UserId), marshaled, uCache.expiration).Err()
}

func (uCache *UserCache) domainKey(id string) string {
	return fmt.Sprintf("gohub:cache:domain:%s", id)
}
