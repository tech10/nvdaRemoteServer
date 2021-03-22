package server

import (
	"math/rand"
	"strconv"
	"time"
)

func gen_key() string {
	min := 1000000
	max := 10000000
	var key string
	var c *ClientChannel
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 20; i++ {
		key = strconv.Itoa(rand.Intn(max-min) + min)
		c = FindChannel(key)
		if c == nil {
			break
		}
	}
	return key
}
