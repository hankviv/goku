package uuid

import (
	"fmt"
	"math/rand"
	"strings"
)

func UUIDv4() (uid string) {
	u := make([]byte, 16)
	_, _ = rand.Read(u)
	u[6] = (u[6] & 0x0f) | 0x40
	u[8] = (u[8] & 0x3f) | 0x80
	uid = strings.ToUpper(fmt.Sprintf("%x%x%x%x%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:]))
	return
}

