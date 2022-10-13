package splitwise

// type ErrorResponse struct {
// StatusCode int
// Message    string
// }
//
// type errorMap struct {
// Error  string `json:"error,omitempty"`
// Errors struct {
// Base []string `json:"base"`
// } `json:"errors,omitempty"`
// }
//
// func UnmarshalErrorResponse(code int, data []byte) (resp ErrorResponse, _ error) {
// var msgMap errorMap
// if err := json.Unmarshal(data, &msgMap); err != nil {
// return ErrorResponse{}, err
// }
//
// if msgMap.Error != "" {
// resp.Message = msgMap.Error
// } else {
// if len(msgMap.Errors.Base) > 0 {
// resp.Message = strings.Join(msgMap.Errors.Base, ", ")
// }
// }
//
// return resp, nil
// }
//
// type Response[T any] struct {
// Result T
// Error  ErrorResponse
// }
