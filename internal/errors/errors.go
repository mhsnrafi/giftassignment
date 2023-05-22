package errors

import "errors"

var (
	ErrGiftAlreadyAssigned = errors.New("gift is already assigned to this employee")
	ErrNoAppropriateGift   = errors.New("no appropriate gift found")
	ErrInvalidEmployeeID   = errors.New("invalid employee id")
	ErrGiftNotReturnable   = errors.New("gift is not returnable")
	ErrCouldNotObtainLock  = errors.New("could not obtain lock")
)
