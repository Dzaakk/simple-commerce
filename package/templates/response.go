package templates

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
