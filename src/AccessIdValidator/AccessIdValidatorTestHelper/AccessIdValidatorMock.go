package AccessIdValidatorTestHelper

import (
	"golang.org/x/net/context"
)

type CallArgs struct {
	ctx      context.Context
	accessId []byte
}

type AccessIdValidatorMock struct {
	validAccessIds []string
	err            error

	CallCount int
	CallArgs  []CallArgs
}

func NewAccessIdValidatorMock(validAccessIds []string, err error) *AccessIdValidatorMock {
	mock := &AccessIdValidatorMock{
		validAccessIds: validAccessIds,
		err:            err,
		CallArgs:       make([]CallArgs, 0, 10),
	}
	return mock
}

func (mock *AccessIdValidatorMock) ValidateAccessId(ctx context.Context, accessId []byte) (bool, error) {
	mock.CallCount++
	mock.CallArgs = append(mock.CallArgs, CallArgs{ctx, accessId})

	if mock.err != nil {
		return false, mock.err
	}

	idStr := string(accessId)
	for _, id := range mock.validAccessIds {
		if id == idStr {
			return true, nil
		}
	}

	return false, nil
}
