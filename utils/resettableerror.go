package utils

import "sync"

type ResettableError struct {
	err  error
	once sync.Once
}

func NewResettaleError() *ResettableError {
	return &ResettableError{
		once: sync.Once{},
		err:  nil,
	}
}

func (r *ResettableError) Error() string {
	if r.err != nil {
		return r.err.Error()
	}

	return ""
}

func (r *ResettableError) Set(err error) {
	if err == nil {
		return
	}

	r.once.Do(func() {
		r.err = err
	})
}

func (r *ResettableError) Reset() {
	r.once = sync.Once{}
	r.err = nil
}
