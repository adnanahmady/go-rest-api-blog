package response

type JsonResponse struct {
	Data any `json:"data"`
	Meta any `json:"meta"`
}

func NewJsonResponse(data any, meta any) *JsonResponse {
	return &JsonResponse{
		Data: data,
		Meta: meta,
	}
}
