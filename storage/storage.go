package storage

import (
	"time"
	"log"
	"math/rand"
	"github.com/garyburd/redigo/redis"
	"fmt"
)

const (
	uidLength = 5
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type Storage struct {
	Address    string
	Pool 	   *redis.Pool
}

type Url struct {
	UID	   string `json:"uid"`
	Url 	   string `json:"url"`
	CreatedAt  string `json:"createdAt"` // RFC3339
}

func InitStorage(adr string) *Storage {
	s := &Storage{
		Address: adr,
		Pool: &redis.Pool{
			MaxIdle:     1000,
			MaxActive:   1000,
			IdleTimeout: 600 * time.Second,
			Wait:        true,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", adr)
				if err != nil {
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
		},
	}
	if s.Pool == nil {
		return nil
	}

	return s
}

func CreateUrl(storage *Storage, urlAddress string) Url {
	url := Url{
		UID: "",
		Url: urlAddress,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	rc := storage.Pool.Get()
	defer rc.Close()

	for {
		url.UID = RandStringBytes(uidLength)
		key := fmt.Sprintf("urlshort:url:%s", url.UID)
		exists, err := rc.Do("HMGET", key, "uid")
		if err != nil {
			log.Println(err)
			continue
		}
		if exists != nil {
			break
		}
	}

	key := fmt.Sprintf("urlshort:url:%s", url.UID)
	if _, err := rc.Do("HMSET", key, "uid", url.UID, "url", url.Url, "createdAt", url.CreatedAt); err != nil {
		log.Fatal(err)
	}

	rc.Flush()

	return url;
}

func GetUrl(storage *Storage, uid string) *Url {
	rc := storage.Pool.Get()
	defer rc.Close()

	key := fmt.Sprintf("urlshort:url:%s", uid)
	u, err := redis.Values(rc.Do("HMGET", key, "uid", "url", "createdAt"))
	if err != nil {
		log.Fatal(err)
	}
	url := new(Url)
	_, err = redis.Scan(u, &url.UID, &url.Url, &url.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}

	rc.Flush()

	return url;
}

func UpdateStatistics(storage *Storage, uid string) {
	rc := storage.Pool.Get()
	defer rc.Close()

	key := fmt.Sprintf("urlshort:statistics:%s", uid)
	if _, err := rc.Do("INCR", key); err != nil {
		log.Fatal(err)
	}
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}