package taildir

import (
	"fmt"
	"testing"
)

func Test_demo(t *testing.T) {
	m := make(map[string]string)
	m[""] = "sss"
	for k, v := range m {
		fmt.Println(k, ":", v)
	}
}
