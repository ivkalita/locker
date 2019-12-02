package locker

import "context"

type Interface interface {
	Lock(context.Context) error
	Unlock(context.Context) error
}
