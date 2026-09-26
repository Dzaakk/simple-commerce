package response

import "net/http"

type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type APIResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

func Response(code int, message string, data interface{}) APIResponse {
	return APIResponse{
		Meta: Meta{
			Code:    code,
			Message: message,
		},
		Data: data,
	}
}

func Success(data interface{}) APIResponse {
	return Response(http.StatusOK, "Success", data)
}

func Unauthorized(message string) APIResponse {
	if message != "" {
		return Response(http.StatusUnauthorized, "Unauthorized", message)
	}
	return Response(http.StatusUnauthorized, "Unauthorized", nil)
}
