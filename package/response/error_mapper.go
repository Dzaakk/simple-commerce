package response

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
)

// ErrorResponse converts any error into a consistent API response.
// AppError is preferred; other errors fall back to keyword-based mapping.
func ErrorResponse(err error) (int, ApiResponse) {
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
	case errors.Is(err, ErrEmailNotVerified):
		return http.StatusForbidden, Response(http.StatusForbidden, http.StatusText(http.StatusForbidden), err.Error())
	case errors.Is(err, ErrInvalidRefreshToken):
		return http.StatusUnauthorized, Response(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), err.Error())
	case errors.Is(err, ErrEmailAlreadyExist):
		return http.StatusConflict, Response(http.StatusConflict, http.StatusText(http.StatusConflict), err.Error())
	case errors.Is(err, ErrInvalidActivationCode):
		return http.StatusBadRequest, Response(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
	case errors.Is(err, ErrUserNotFound):
		return http.StatusNotFound, Response(http.StatusNotFound, http.StatusText(http.StatusNotFound), err.Error())
	case errors.Is(err, sql.ErrNoRows):
		return http.StatusNotFound, Response(http.StatusNotFound, http.StatusText(http.StatusNotFound), "resource not found")
	}

	msg := strings.ToLower(err.Error())
	code := http.StatusInternalServerError

	switch {
	case strings.Contains(msg, "unauthorized"):
		code = http.StatusUnauthorized
	case strings.Contains(msg, "forbidden"):
		code = http.StatusForbidden
	case strings.Contains(msg, "not found"):
		code = http.StatusNotFound
	case strings.Contains(msg, "no rows updated"):
		code = http.StatusNotFound
	case strings.Contains(msg, "already exists"), strings.Contains(msg, "duplicate"), strings.Contains(msg, "conflict"):
		code = http.StatusConflict
	case strings.Contains(msg, "status is not"):
		code = http.StatusConflict
	case strings.Contains(msg, "invalid parameter"),
		strings.Contains(msg, "invalid request"),
		strings.Contains(msg, "invalid signature"),
		strings.Contains(msg, "invalid"):
		code = http.StatusBadRequest
	case strings.Contains(msg, "stock"):
		code = http.StatusConflict
	}

	if code >= 500 {
		return code, Response(code, http.StatusText(code), nil)
	}

	return code, Response(code, http.StatusText(code), err.Error())
}
