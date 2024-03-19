package mymd5

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func computeMD5(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func Show() {
	a := computeMD5("admin")
	fmt.Println("admin = ", a)
	b := computeMD5("666")
	fmt.Println("666 = ", b)
	c := computeMD5("ZAQ!2wsx")
	fmt.Println("ZAQ!2wsx = ", c)
}

// admin =  21232f297a57a5a743894a0e4a801fc3
// 666 =  fae0b27c451c728867a567e8c1bb4e53
// ZAQ!2wsx =  e3bc38a4faa625d074664d572d810c1e
