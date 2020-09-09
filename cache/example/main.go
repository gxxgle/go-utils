package main

import (
	"log"

	"github.com/gxxgle/go-utils/cache"
)

type person struct {
	Name string
	Age  int
}

func main() {
	c, err := cache.NewRedisCacher(&cache.RedisConfig{
		Addr:       "192.168.31.61:6379",
		Password:   "KgqvdOdYV5",
		MaxRetries: 10,
	})
	if err != nil {
		log.Fatalln(err)
	}

	k := "xxx:yyy"
	a := &person{"xiaoming", 12}
	b := &person{}

	log.Println("before:", b)
	log.Println(c.Set(k, a, 0))
	log.Println(c.Get(k, b))
	log.Println("after:", b)
}
