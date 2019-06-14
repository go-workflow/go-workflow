package model_test

import (
	"fmt"
	"testing"

	"github.com/go-workflow/go-workflow/workflow-engine/model"
)

func TestClient(t *testing.T) {
	err := model.RedisClient.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}
	val, err := model.RedisClient.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key:", val)
}
