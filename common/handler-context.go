package common

type HandleFunc func(ctx HContext) error

type HContext interface {
	// Header

	// Returns the HTTP request header specified by the field.
	Get(key string, defaultValue ...string) string
	// Sets the response’s HTTP header field to the specified key, value.
	Set(key string, value string)
	// Returns the HTTP request header "Authorization" field.
	Authorization() string
	// Returns the HTTP request header "RequestId" field.
	RequestId() string

	// Request

	// Return HTTP method of the request.
	Method() string
	// Return path part of the request URL.
	Path() string
	// Binds the request body to a struct.
	BodyParser(out interface{}) error
	// Binds the query string parameter to a struct.
	QueryParser(out interface{}) error
	// Get value of query string parameter, If there is no query string, it returns an empty string.
	Query(key string, defaultValue ...string) string
	// Get value of route parameters
	Params(key string, defaultValue ...string) string

	// stores variables scoped to the request
	// or get variables scoped from the request
	Locals(key string, value ...interface{}) interface{}

	// middleware
	Next() error

	// response
	SendStatus(status int) error
	SendJSON(status int, data interface{}) error
}
