package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
)

func NewRedisClientForTest(t *testing.T) *redis.Client {
	t.Helper()

	host := "127.0.0.1"
	port := 36379

	cli := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       0,
	})
	if err := cli.Ping(context.Background()).Err(); err != nil {
		t.Fatalf("couldn't connect to Redis: %v", err)
	}
	return cli
}

func TestKVS_Save(t *testing.T) {
	t.Parallel()

	cli := NewRedisClientForTest(t)
	kvs := &KVS{Cli: cli}
	userName := "TestKVS_Save"
	pw := "username"
	ctx := context.Background()
	t.Cleanup(func() {
		cli.Del(ctx, userName)
	})
	if err := kvs.CreateUser(ctx, userName, pw); err != nil {
		t.Fatalf("SaveUser failed: %v", err)
	}
	if err := kvs.CreateUser(ctx, userName, pw); err != ErrUserExists {
		t.Fatalf("expected %v, got %v", ErrUserExists, err)
	}
}

func TestKVS_Load(t *testing.T) {
	t.Parallel()

	cli := NewRedisClientForTest(t)
	kvs := &KVS{Cli: cli}
	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		userName := "test"
		pw := "password"
		ctx := context.Background()
		t.Cleanup(func() {
			cli.Del(ctx, userName)
		})
		err := kvs.CreateUser(ctx, userName, pw)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}
		ok, err := kvs.UserExists(ctx, userName)
		if err != nil {
			t.Fatalf("failed to check user: %v", err)
		}
		if !ok {
			t.Errorf("user doesn't exist")
		}
	})
	t.Run("ng", func(t *testing.T) {
		t.Parallel()

		userName := "test"
		pw := "password"
		ctx := context.Background()
		t.Cleanup(func() {
			cli.Del(ctx, userName)
		})
		err := kvs.CreateUser(ctx, userName, pw)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}
		ok, err := kvs.UserExists(ctx, "test1")
		if err != nil {
			t.Fatalf("failed to check user: %v", err)
		}
		if ok {
			t.Errorf("user exists, but should not")
		}
	})
}
