package common

import (
	"net/http"

	"github.com/juju/errors"
)

const (
	// 200
	Finish      = "Finish"
	StopResolve = "StopResolve"

	// 404
	NotFoundError          = "NotFoundError"
	CacheFileNotFoundError = "CacheFileNotFoundError"

	// 400
	ParameterMissingError             = "ParameterMissingError"
	AccountError                      = "AccountError"
	NotSuperAdministratorError        = "NotSuperAdministratorError"
	NotOrganizationAdministratorError = "NotOrganizationAdministratorError"

	// 500
	UnknownError           = "UnknownError"
	ServerError            = "ServerError"
	NetworkError           = "NetworkError"
	LoginCookieExpireError = "LoginCookieExpireError"
	ONESVersionError       = "ONESVersionError"
)

var ErrorCodeMap = map[string]int{
	Finish:                            http.StatusOK,
	StopResolve:                       http.StatusOK,
	NotFoundError:                     http.StatusNotFound,
	CacheFileNotFoundError:            http.StatusNotFound,
	ParameterMissingError:             http.StatusBadRequest,
	NotSuperAdministratorError:        http.StatusBadRequest,
	NotOrganizationAdministratorError: http.StatusBadRequest,
	AccountError:                      http.StatusBadRequest,
	ServerError:                       http.StatusInternalServerError,
	NetworkError:                      http.StatusInternalServerError,
	ONESVersionError:                  http.StatusInternalServerError,
}

type Err struct {
	*errors.Err
	Code    int
	ErrCode string
	Body    interface{}
}

func (e *Err) String() string {
	return e.ErrCode
}

func Errors(ErrCode string, body interface{}) error {
	e := new(Err)
	e.ErrCode = ErrCode
	e.Body = body
	e.Code = ErrorCodeMap[ErrCode]
	return e
}
