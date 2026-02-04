package response

import "net/http"

type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ApiResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

func Response(code int, message string, data interface{}) ApiResponse {
	Response := ApiResponse{
		Meta: Meta{
			Code:    code,
			Message: message,
		},
		Data: data,
	}
	return Response
}

func InternalServerError(message string) ApiResponse {
	if message != "" {
		return Response(http.StatusInternalServerError, "Internal Server Error", message)
	}
	return Response(http.StatusInternalServerError, "Internal Server Error", nil)
}

func BadRequest(message string) ApiResponse {
	if message != "" {
		return Response(http.StatusBadRequest, "Bad Request", message)
	}
	return Response(http.StatusBadRequest, "Bad Request", nil)
}

func Success(data interface{}) ApiResponse {
	return Response(http.StatusOK, "Success", data)
}

func NotFound(message string) ApiResponse {
	if message != "" {
		return Response(http.StatusNotFound, "Not Found", message)
	}
	return Response(http.StatusNotFound, "Not Found", nil)
}

func Unauthorized(message string) ApiResponse {
	if message != "" {
		return Response(http.StatusUnauthorized, "Unauthorized", message)
	}
	return Response(http.StatusUnauthorized, "Unauthorized", nil)
}

func InvalidRequestData() ApiResponse {
	return Response(http.StatusBadRequest, "Invalid Request Data", nil)
}

func InvalidEmailOrPassword() ApiResponse {
	return Response(http.StatusBadRequest, "Invalid Email or Password", nil)
}
