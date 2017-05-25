package AccessIdValidator

import "appengine"

type AccessIdValidator interface {
	ValidateAccessId(ctx appengine.Context, accessId []byte) (bool, error)
}
