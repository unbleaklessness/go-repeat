package main

type ierrori interface {
	Unwrap() error
	Error() string
	Message() string
	Code() int
}

type ierror struct {
	e error
	m string
	c int
}

func (e ierror) Unwrap() error {
	return e.e
}

func (e ierror) Error() string {
	if e.e != nil {
		return e.e.Error()
	}
	return "empty error string"
}

func (e ierror) Message() string {
	return e.m
}

func (e ierror) Code() int {
	return e.c
}
