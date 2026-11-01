package errors

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
)

// ICCError Interface implements an error interface.
// It contains the return http status, error message and function
// It also has a function to convert a response to a peer struct
type ICCError interface {
	Status() int32
	Message() string
	GetErrorResponse() *peer.Response
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
	return fmt.Sprintf("%d %s", c.status, c.err.Error())
}

// GetErrorResponse converts an Httperror instance to a peer.Response
func (c *CCError) GetErrorResponse() *peer.Response {
	return &peer.Response{
		Status:  errorStatus(c.status),
		Message: c.err.Error(),
	}
}

// errorStatus replaces statuses below 400 with 500, since Fabric peers endorse
// any response with status below 400 as a successful transaction.
func errorStatus(status int32) int32 {
	if status < 400 {
		return 500
	}

	return status
}

// NewCCError creates a new CCError instance.
// Statuses below 400 are replaced by 500.
func NewCCError(errMsg string, status int32) *CCError {
	return &CCError{
		status: errorStatus(status),
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

// WrapErrorWithStatus wraps an existing error and adds a status to it.
// Statuses below 400 are replaced by 500.
func WrapErrorWithStatus(err error, errMsg string, status int32) *CCError {
	newErr := WrapError(err, errMsg)
	newErr.status = errorStatus(status)

	return newErr
}
