package listeners

import "context"

type ListenerServer interface {
	Start(ctx context.Context)
	Stop()
}
