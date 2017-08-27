package Parser

import (
	"github.com/go-redis/redis"
	"fmt"

	"strings"
)

var client *redis.Client =  redis.NewClient(&redis.Options{
Addr:     "localhost:6379",
Password: "", // no password set
DB:       0,  // use default DB
})


func PingRedis() {

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	if err != nil{
		panic(err)
	}
	// Output: PONG <nil>
}

func PushData(kv *HostKV){
	key := kv.host+":"+kv.app
	err := client.HSet(key,kv.key,kv.val)
	err2 :=client.HSet(key,"host",kv.host)
	client.HSet(key,"topic",kv.app)
	if err != nil || err2 != nil{
		fmt.Println(err,err2)
	}
}

func AllKeys() ([]string,error){
	res := client.Keys("*")
	keys ,err := res.Result()
	return keys,err
}

func GetHMap(key string) map[string]string{
	fields := strings.Split(key,":")
	res := client.HLen(key)
	num,_ := res.Result()
	if num < int64(QuerySectionNumber(fields[1])){
		fmt.Println("not ready now...")
		return nil
	}else{
		hmap := client.HGetAll(key)
		m := hmap.Val()
		return m
	}
}

/*
func ExampleClient() {
	err := client.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exists")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exists
}
*/