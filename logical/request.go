package logical

import (
	"errors"
)

// Request is a struct that stores the parameters and context
// of a request being made to Vault. It is used to abstract
// the details of the higher level request protocol from the handlers.
type Request struct {
	// Operation is the requested operation type
	Operation Operation

	// Path is the part of the request path not consumed by the
	// routing. As an example, if the original request path is "prod/aws/foo"
	// and the AWS logical backend is mounted at "prod/aws/", then the
	// final path is "foo" since the mount prefix is trimmed.
	Path string

	// Request data is an opaque map that must have string keys.
	Data map[string]interface{}

	// Storage can be used to durably store and retrieve state.
	Storage Storage

	// Secret will be non-nil only for Revoke and Renew operations
	// to represent the secret that was returned prior.
	Secret *Secret
}

// Get returns a data field and guards for nil Data
func (r *Request) Get(key string) interface{} {
	if r.Data == nil {
		return nil
	}
	return r.Data[key]
}

// GetString returns a data field as a string
func (r *Request) GetString(key string) string {
	raw := r.Get(key)
	s, _ := raw.(string)
	return s
}

// RenewRequest creates the structure of the renew request.
func RenewRequest(
	path string, secret *Secret, data map[string]interface{}) *Request {
	return &Request{
		Operation: RenewOperation,
		Path:      path,
		Data:      data,
		Secret:    secret,
	}
}

// RevokeRequest creates the structure of the revoke request.
func RevokeRequest(
	path string, secret *Secret, data map[string]interface{}) *Request {
	return &Request{
		Operation: RevokeOperation,
		Path:      path,
		Data:      data,
		Secret:    secret,
	}
}

// RollbackRequest creates the structure of the revoke request.
func RollbackRequest(path string) *Request {
	return &Request{
		Operation: RollbackOperation,
		Path:      path,
	}
}

// Operation is an enum that is used to specify the type
// of request being made
type Operation string

const (
	// The operations below are called per path
	ReadOperation   Operation = "read"
	WriteOperation            = "write"
	DeleteOperation           = "delete"
	ListOperation             = "list"
	HelpOperation             = "help"

	// The operations below are called globally, the path is less relevant.
	RevokeOperation   Operation = "revoke"
	RenewOperation              = "renew"
	RollbackOperation           = "rollback"
)

var (
	// ErrUnsupportedOperation is returned if the operation is not supported
	// by the logical backend.
	ErrUnsupportedOperation = errors.New("unsupported operation")

	// ErrUnsupportedPath is returned if the path is not supported
	// by the logical backend.
	ErrUnsupportedPath = errors.New("unsupported path")

	// ErrInvalidRequest is returned if the request is invalid
	ErrInvalidRequest = errors.New("invalid request")
)
