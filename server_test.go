package ji

import (
	"context"
	"fmt"
	"testing"
)

func TestStartServerByName(t *testing.T) {
	err := StartServerByName(context.Background(), "tangerine/csdn", nil, nil)
	fmt.Println(err)
}
