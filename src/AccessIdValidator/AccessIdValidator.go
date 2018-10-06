package AccessIdValidator

import "golang.org/x/net/context"

type AccessIdValidator interface {
	ValidateAccessId(ctx context.Context, accessId []byte) (bool, error)

	EnsureUpdatedIds(ctx context.Context) error
}
