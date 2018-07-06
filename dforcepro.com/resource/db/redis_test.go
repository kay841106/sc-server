package db

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func _Test_Redis(t *testing.T) {
	redis := Redis{"localhost:6379", "", 1}
	err := redis.GetClient().Set("key", "aaaa", 5*time.Minute).Err()
	if err != nil {
		fmt.Println(err.Error())
	}

	val, err := redis.GetClient().Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
	redis.GetClient().FlushAll()
	assert.Equal(t, 1, 1, "should be equal.")
}

func _Test_private(t *testing.T) {
	serverPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		// _di.Log.Err(err.Error())
		// w.WriteHeader(http.StatusInternalServerError)
	}

	// err := redis.GetClient().Set(key, , 5*time.Minute).Err()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	serverPublicKey := &serverPrivateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(serverPublicKey)
	fmt.Println(string(derPkix))
	// fmt.Println(serverPrivateKey)
	// fmt.Println(serverPublicKey)
}

func Test_RedisExist(t *testing.T) {
	redis := Redis{"localhost:6379", "", RedisDB_Alert}
	client := redis.GetClient()
	key := "test"
	value := client.Incr(key).Val()
	fmt.Println(value)
	value = client.Incr(key).Val()
	fmt.Println(value)
	result := client.Get(key).String()
	fmt.Println(result)

}
