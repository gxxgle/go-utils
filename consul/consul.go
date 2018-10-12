package consul

import (
	"os"

	"github.com/go-redis/redis"
	"github.com/gxxgle/go-utils/json"
	"github.com/gxxgle/go-utils/log"
	"github.com/hashicorp/consul/api"
)

var Client *api.Client

func init() {
	var (
		err    error
		config = api.DefaultConfig()
		addr   = os.Getenv("CONSUL_HTTP_ADDR")
	)

	if addr != "" {
		config.Address = addr
	}

	Client, err = api.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}
}

func KVPut(key string, value interface{}) error {
	_, err := Client.KV().Put(&api.KVPair{
		Key:   key,
		Value: json.MustMarshal(value),
	}, nil)
	return err
}

func KVGet(key string, value interface{}) error {
	pair, _, err := Client.KV().Get(key, nil)
	if err != nil {
		return err
	}

	if pair == nil || len(pair.Value) <= 0 {
		return redis.Nil
	}

	return json.Unmarshal(pair.Value, value)
}
