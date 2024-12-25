package redisclient

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	GRedIsDB  *redis.Client          // Declare a global variable for the Redis client
	once      sync.Once              // Declare a sync.Once to ensure single initialization
	GRedIsCtx = context.Background() // Create a background context for Redis operations
)

func Init() (lErr error) {
	// Retrieve the Redis host and port from environment variables
	lRedISHost := os.Getenv("REDIS_HOST")
	lRedISPort := os.Getenv("REDIS_PORT")
	// Check if the host or port are empty and return an error if they are
	if strings.EqualFold(lRedISHost, "") {
		return errors.New("redis host or port not present in env set < REDIS_HOST = host_name >")
	}

	// Initialize the Redis client with the host and port
	initialize(lRedISHost, lRedISPort)

	return nil

}

// initialize Redis using the RedIsInit function from the redisclient package
func initialize(pRedisHost, pRedisPort string) {
	if pRedisPort != "" {
		pRedisHost = fmt.Sprintf("%s:%s", pRedisHost, pRedisPort)
	}
	// Ensure the Redis client is initialized only once
	once.Do(func() {
		GRedIsDB = redis.NewClient(&redis.Options{
			// Set the address for the Redis client
			Addr: pRedisHost,
		})
	})
}

// KeyExist checks if a Redis key exists.
func KeyExist(pKey string) (lExist int64, lErr error) {
	lExist, lErr = GRedIsDB.Exists(GRedIsCtx, pKey).Result() // Check if the key exists in Redis.
	if lErr != nil {
		return lExist, lErr // Return false and an error if checking key fails.
	}
	return lExist, nil // Return true if the key exists, otherwise false.
}

func SetKey(pKey, pValue string) (lErr error) {
	lErr = GRedIsDB.Set(GRedIsCtx, pKey, pValue, 0).Err() // Set the value in Redis.
	if lErr != nil {
		return lErr
	}
	return nil
}

func UpdateKey(pKey, pValue string) (lErr error) {

	lExistCount, lErr := KeyExist(pKey)
	if lErr != nil {
		return lErr
	}
	if lExistCount <= 0 {
		return fmt.Errorf("GivenKey %s not exist", pKey)
	}

	lErr = SetKey(pKey, pValue)
	if lErr != nil {
		return lErr
	}
	return nil
}

func GetValue(pKey string) (lRedISValue string, lErr error) {

	lRedISValue, lErr = GRedIsDB.Get(GRedIsCtx, pKey).Result() // Get the value from Redis.
	if lErr != nil {                                           // If there is an error getting the value from Redis.
		return lRedISValue, lErr
	}

	return lRedISValue, nil
}

func DeleteValue(pKey string) (lErr error) {
	lExistCount, lErr := KeyExist(pKey)
	if lErr != nil {
		return lErr
	}
	if lExistCount <= 0 {
		return fmt.Errorf("GivenKey %s not exist", pKey)
	}
	GRedIsDB.Del(GRedIsCtx, pKey) // Get the value from Redis.
	return nil
}

func GetKeys(pPatton string) (lKeys []string, lErr error) {
	lKeys, lErr = GRedIsDB.Keys(GRedIsCtx, "pPatton").Result()
	if lErr != nil {
		return lKeys, lErr
	}
	return lKeys, nil
}
