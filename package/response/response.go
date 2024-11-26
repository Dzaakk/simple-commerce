package response

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
		return Response(500, "Internal Server Error", message)
	}
	return Response(500, "Internal Server Error", nil)
}

func BadRequest(message string) ApiResponse {
	if message != "" {
		return Response(400, "Bad Request", message)
	}
	return Response(400, "Bad Request", nil)
}

func Success(data interface{}) ApiResponse {
	return Response(200, "Success", data)
}

func NotFound(message any) ApiResponse {
	if message != "" {
		return Response(404, "Not Found", message)
	}
	return Response(404, "Not Found", nil)
}
