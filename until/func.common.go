package until

import (
	"crypto/md5"
	"fmt"
	"time"
)

const (
	defaultKey = "HSJnpYbnLamBhu"
)

func GenerateToken(offset int64) string {
	var time = time.Now().Unix()
	if offset != 0 {
		time += offset
	}
	str := []byte(fmt.Sprintf("%s%d%s", defaultKey, time, defaultKey))
	newToken := md5.Sum(str)
	return fmt.Sprintf("%X", newToken)
}
