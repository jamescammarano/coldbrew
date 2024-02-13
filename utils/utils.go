package utils

import (
	"encoding/base64"
	"math/rand"
	"time"
	"unsafe"
)

func Base64() string {
	return base64.StdEncoding.EncodeToString([]byte(randStringBytesMaskImprSrcUnsafe(202)))
}

func MergeMaps(mapArr []map[string]string) map[string]string {
	merged := map[string]string{}

	for _, t := range mapArr {
		for k, v := range t {
			merged[k] = v

		}
	}

	return merged
}

// source: https://stackoverflow.com/a/31832326
func randStringBytesMaskImprSrcUnsafe(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
