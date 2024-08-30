package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	ErrUserExists = errors.New("user already exists")
)

type KVS struct {
	Cli *redis.Client
}

func NewKVS() *KVS {
	return nil
}

func (k *KVS) CreateUser(ctx context.Context, userName, pw string) error {
	if ok, err := k.Cli.Exists(ctx, userName).Result(); err != nil {
		return fmt.Errorf("failed to check existence of key: %v", err)
	} else if ok != 0 {
		return ErrUserExists
	}
	return k.Cli.Set(ctx, userName, pw, 0).Err()
}

func (k *KVS) UserExists(ctx context.Context, userName string) (bool, error) {
	ok, err := k.Cli.Exists(ctx, userName).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key: %v", err)
	}
	return ok != 0, nil
}
