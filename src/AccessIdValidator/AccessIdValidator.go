package AccessIdValidator

import "golang.org/x/net/context"

type AccessIdValidator interface {
	ValidateAccessIdPrimary(ctx context.Context, accessId []byte) (bool, error)
	ValidateAccessIdSecondary(ctx context.Context, accessId []byte) (bool, error)
}
