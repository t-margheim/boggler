package boggle

type errorResponse struct {
	Error   string `json:"err"`
	Message string `json:"msg"`
}
type response struct {
	Words []string `json:"words"`
}
