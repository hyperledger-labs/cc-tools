package errors

import (
	"errors"

	"github.com/hyperledger/fabric-protos-go/peer"
)

// ICCError Interface implements an error interface.
// It contains the return http status, error message and function
// It also has a function to convert a response to a peer struct
type ICCError interface {
	Status() int32
	Message() string
	GetErrorResponse() peer.Response
	Error() string
}

// CCError struct
type CCError struct {
	status int32
	err    error
}

// Status Returns the http status code
func (c *CCError) Status() int32 {
	return c.status
}

// Message returns the inner error of a CCError instance
func (c *CCError) Message() string {
	return c.err.Error()
}

// Implements the error interface
func (c *CCError) Error() string {
	return string(c.status) + c.err.Error()
}

// GetErrorResponse converts an Httperror instance to a peer.Response
func (c *CCError) GetErrorResponse() peer.Response {
	return peer.Response{
		Status:  c.status,
		Message: c.err.Error(),
	}
}

// NewCCError creates a new CCError instance
func NewCCError(errMsg string, status int32) *CCError {
	return &CCError{
		status: status,
		err:    errors.New(errMsg),
	}
}

// WrapError stacks an error msg on top of the existing one
func WrapError(err error, errMsg string) *CCError {
	if err == nil {
		return NewCCError(errMsg, 500)
	}

	if v, ok := err.(*CCError); ok {
		return NewCCError(
			errMsg+": "+v.Message(),
			v.Status(),
		)
	}
	e := errMsg + ": " + err.Error()
	return NewCCError(e, 500)
}

// WrapErrorWithStatus wraps an existing error and adds a status to it
func WrapErrorWithStatus(err error, errMsg string, status int32) *CCError {
	newErr := WrapError(err, errMsg)
	newErr.status = status

	return newErr
}
