package notify

import (
	"context"
	"fmt"
)

type Notify struct{}

func (n *Notify) Send(ctx context.Context, to string, msg string) {
	fmt.Printf("To %s, msg: %s", to, msg)
}
