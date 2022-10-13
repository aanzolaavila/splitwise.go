package resources

import "time"

type User struct {
	Entity
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	Email              string `json:"email"`
	RegistrationStatus string `json:"registration_status"`
	Picture            struct {
		Small  string `json:"small"`
		Medium string `json:"medium"`
		Large  string `json:"large"`
	} `json:"picture"`
	NotificationsRead  time.Time `json:"notifications_read"`
	NotificationsCount int       `json:"notifications_count"`
	Notifications      struct {
		AddedAsFriend bool `json:"added_as_friend"`
	} `json:"notifications"`
	DefaultCurrency string `json:"default_currency"`
	Locale          string `json:"locale"`
}

// type CurrentUserRequiredParams struct {
// }
//
// func (c CurrentUserRequiredParams) Set(key string, value string) error {
// return fmt.Errorf("there is no required field for /get_current_user")
// }
//
// type CurrentUserOp struct {
// }
//
// func (o CurrentUserOp) Do(client ApiClient, params CurrentUserRequiredParams) (Response[User], error) {
// const path = "/get_current_user"
//
// res, err := client.Get(path, url.Values{}, url.Values{})
// if err != nil {
// return Response[User]{}, err
// }
//
// body, err := io.ReadAll(res.Body)
// if err != nil {
// return Response[User]{}, err
// }
// defer res.Body.Close()
//
// if res.StatusCode != http.StatusOK {
// var msgMap ErrorResponse
// var msg string
// err = json.Unmarshal(body, &msgMap)
// if err == nil {
// msg = msgMap.Message
// }
//
// return Response[User]{
// Error: ErrorResponse{
// ErrCode: int16(res.StatusCode),
// Message: msg,
// },
// }, nil
// }
//
// var user User
// err = json.Unmarshal(body, &user)
// if err != nil {
// return Response[User]{}, err
// }
//
// return Response[User]{
// Result: user,
// }, nil
// }
