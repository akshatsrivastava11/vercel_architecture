package utils

import "math/rand"

const MAX_LEN = 5

func Generate_random() string {
	ans := ""
	var randomString = "1234567890qwertyuiopasdfghjklzxcvbnm"
	for i := 0; i < MAX_LEN; i++ {
		index := rand.Intn(len(randomString))
		ans += string(randomString[index])
	}
	return ans
}
