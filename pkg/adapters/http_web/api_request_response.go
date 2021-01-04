package http_web

type ApiResponse struct {
	Status bool `json:"status"`
}

type BadApiResponse struct {
	ApiResponse
	Error string `json:"error"`
}

func NewBadApiResponse(error error) *BadApiResponse {
	return &BadApiResponse{ApiResponse: ApiResponse{Status: false}, Error: error.Error()}
}

type DeleteApiResponse struct {
	ApiResponse
}

type SaveApiResponse struct {
	ApiResponse
	Key string `json:"key"`
}
type LoadApiResponse struct {
	ApiResponse
	Text string `json:"text"`
}
type ListApiResponse struct {
	ApiResponse
	Keys []string `json:"list"`
}

type SaveRequest struct {
	RequiredKey string `json:"required_key"`
	Text        string `json:"text"`
}
