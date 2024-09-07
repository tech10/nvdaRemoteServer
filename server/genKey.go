package server

import (
	"math/rand"
	"strconv"
)

func gen_key() string {
	min_num := 1000000
	max_num := 10000000
	var key string
	var c *ClientChannel
	for i := 0; i < 20; i++ {
		key = strconv.Itoa(rand.Intn(max_num-min_num) + min_num)
		c = FindChannel(key)
		if c == nil {
			break
		}
	}
	return key
}
