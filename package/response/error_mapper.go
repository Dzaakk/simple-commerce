package response

import (
	"database/sql"
	"errors"
	"net/http"
)

// ErrorResponse maps typed and sentinel errors to a stable public response.
// Unknown errors remain opaque and are logged by the HTTP middleware.
func ErrorResponse(err error) (int, APIResponse) {
	if err == nil {
		return http.StatusOK, Success(nil)
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		code := appErr.Code
		if code == 0 {
			code = http.StatusInternalServerError
		}
		data := appErr.Details
		if data == nil {
			data = appErr.Message
		}
		return code, Response(code, http.StatusText(code), data)
	}

	switch {
	case errors.Is(err, ErrInvalidCredentials):
		return http.StatusUnauthorized, Response(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), err.Error())
	case errors.Is(err, ErrInvalidRefreshToken):
		return http.StatusUnauthorized, Response(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), err.Error())
	case errors.Is(err, ErrEmailAlreadyExist):
		return http.StatusConflict, Response(http.StatusConflict, http.StatusText(http.StatusConflict), err.Error())
	case errors.Is(err, ErrUserNotFound):
		return http.StatusNotFound, Response(http.StatusNotFound, http.StatusText(http.StatusNotFound), err.Error())
	case errors.Is(err, sql.ErrNoRows):
		return http.StatusNotFound, Response(http.StatusNotFound, http.StatusText(http.StatusNotFound), "resource not found")
	}

	return http.StatusInternalServerError, Response(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil)
}
