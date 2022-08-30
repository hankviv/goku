package bloom

import (
	"fmt"
	"testing"
	"time"
)

func Test_Add(t *testing.T) {

	timestamp := time.Now().UnixNano() / 1e6

	s := fmt.Sprintf("%d%d", timestamp, 99999)
	fmt.Println(s)
}

func Test_Exists(t *testing.T) {
}
