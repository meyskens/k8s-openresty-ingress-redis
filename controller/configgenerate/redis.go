package configgenerate

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

// UpdateRedis sets the DomainConfigValues in redis
func UpdateRedis(in []DomainConfigValues) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	for _, domain := range in {
		value, _ := json.Marshal(domain.Values)
		err := client.Set(domain.Domain, value, 0).Err()
		if err != nil {
			fmt.Println(err)
		}
	}
}
