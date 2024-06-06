package queue

import "errors"

type AlreadyActiveError struct {
}

func (e *AlreadyActiveError) Error() string {
	return "queue already active"
}

func (e *AlreadyActiveError) Is(target error) bool {
	var val *AlreadyActiveError
	if !errors.As(target, &val) {
		return false
	}

	return true
}

type InActiveError struct {
}

func (e *InActiveError) Error() string {
	return "queue inactive"
}

func (e *InActiveError) Is(target error) bool {
	var val *InActiveError
	if !errors.As(target, &val) {
		return false
	}

	return true
}
