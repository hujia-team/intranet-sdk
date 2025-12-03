package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateToken(ak string, sk string) string {
	today := time.Now().Format("2006-1-2") // Go的time格式化布局与常见的YYYY-MM-DD不同，这里是2006-01-02对应YYYY-MM-DD
	currentHour := time.Now().Hour()
	txt := fmt.Sprintf("%s_%s_%s_%d", today, ak, sk, currentHour)

	hasher := md5.New()
	hasher.Write([]byte(txt))
	return hex.EncodeToString(hasher.Sum(nil))
}
