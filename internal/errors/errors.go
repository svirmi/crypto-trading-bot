package errors

import "fmt"

type ErrorCode string

const (
	// Functional error codes
	NOT_FOUND_ERROR   = "NOT_FOUND_ERROR"
	BAD_REQUEST_ERROR = "BAD_REQUEST_ERROR"
	DUPLICATE_ERROR   = "DUPLICATE_ERROR"

	// Technical error codes
	SERVER_ERROR   = "SERVER_ERROR"
	MONGO_ERROR    = "MONGO_ERROR"
	EXCHANGE_ERROR = "EXCHANGE_ERROR"
)

type CtbError interface {
	Error() string
	Code() ErrorCode
}

type ctb_error struct {
	message string
	code    ErrorCode
}

func (e ctb_error) Error() string {
	return e.message
}

func (e ctb_error) Code() ErrorCode {
	return e.code
}

func WrapNotFound(err error) CtbError {
	if err == nil {
		return nil
	}
	return build_error(NOT_FOUND_ERROR, err.Error())
}

func WrapBadRequest(err error) CtbError {
	if err == nil {
		return nil
	}
	return build_error(BAD_REQUEST_ERROR, err.Error())
}

func WrapDuplicate(err error) CtbError {
	if err == nil {
		return nil
	}
	return build_error(DUPLICATE_ERROR, err.Error())
}

func WrapInternal(err error) CtbError {
	if err == nil {
		return nil
	}
	return build_error(SERVER_ERROR, err.Error())
}

func WrapMongo(err error) CtbError {
	if err == nil {
		return nil
	}
	return build_error(MONGO_ERROR, err.Error())
}

func WrapExchange(err error) CtbError {
	if err == nil {
		return nil
	}
	return build_error(EXCHANGE_ERROR, err.Error())
}

func NotFound(format string, a ...interface{}) CtbError {
	return build_error(NOT_FOUND_ERROR, format, a...)
}

func BadRequest(format string, a ...interface{}) CtbError {
	return build_error(BAD_REQUEST_ERROR, format, a...)
}

func Duplicate(format string, a ...interface{}) CtbError {
	return build_error(DUPLICATE_ERROR, format, a...)
}

func Internal(format string, a ...interface{}) CtbError {
	return build_error(SERVER_ERROR, format, a...)
}

func Mongo(format string, a ...interface{}) CtbError {
	return build_error(MONGO_ERROR, format, a...)
}

func Exchange(format string, a ...interface{}) CtbError {
	return build_error(EXCHANGE_ERROR, format, a...)
}

func build_error(code ErrorCode, f string, a ...interface{}) CtbError {
	return ctb_error{
		message: fmt.Sprintf(f, a...),
		code:    code}
}
