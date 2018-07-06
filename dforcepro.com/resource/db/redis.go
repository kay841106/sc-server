package db

import "github.com/go-redis/redis"

var (
	redisConnPool = make(map[int]*redis.Client)
)

type Redis struct {
	Host string `yaml:"host"`
	Pwd  string `yaml:"pass"`
	db   int    `yaml:"db"`
}

const (
	RedisDB_Default  = 0
	RedisDB_LoginKey = 1
	RedisDB_Token    = 2
	RedisDB_Alert    = 3
)

func (r *Redis) DB(db int) *Redis {
	r.db = db
	return r
}

func (r *Redis) GetClient() *redis.Client {
	redisClient, ok := redisConnPool[r.db]
	if ok {
		return redisClient
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     r.Host,
		Password: r.Pwd, // no password set
		DB:       r.db,  // use default DB
	})
	redisConnPool[r.db] = redisClient
	return redisClient
}

func (r *Redis) Close() {
	for _, val := range redisConnPool {
		val.Close()
	}
}
