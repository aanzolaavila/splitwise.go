package splitwise

// import (
// "fmt"
// "net/http"
// "net/url"
// "time"
// )
//
// const (
// DiscouragedDefaultBaseUrl = "https://secure.splitwise.com"
// )
//
// type OldHttpClient struct {
// *http.Client
// BaseUrl string
// }
//
// func (c *OldHttpClient) Do(req *http.Request) (*http.Response, error) {
// if req.URL == nil {
// url, err := url.Parse(c.BaseUrl)
// if err != nil {
// return nil, err
// }
//
// req.URL = url
// } else {
// req.URL.Host = DefaultBaseUrl
// }
//
// return c.Client.Do(req)
// }
//
// func defaultClient() *OldHttpClient {
// httpClient := &http.Client{
// Timeout: 10 * time.Second,
// }
//
// return &OldHttpClient{
// Client:  httpClient,
// BaseUrl: DefaultBaseUrl,
// }
// }
//
// type TokenClient struct {
// *OldHttpClient
// Token string
// }
//
// func NewHttpClient(baseClient *http.Client) (*OldHttpClient, error) {
// return NewClientWithBaseUrl(baseClient, DiscouragedDefaultBaseUrl)
// }
//
// func NewClientWithBaseUrl(baseClient *http.Client, baseUrl string) (*OldHttpClient, error) {
// if baseClient == nil {
// return nil, fmt.Errorf("base client should not be nil")
// }
//
// if baseUrl == "" {
// return nil, fmt.Errorf("base url should not be zero-valued")
// }
//
// return &OldHttpClient{
// Client:  baseClient,
// BaseUrl: baseUrl,
// }, nil
// }
//
// func NewTokenClient(token string, baseClient *OldHttpClient) (*TokenClient, error) {
// if token == "" {
// return nil, fmt.Errorf("api key should not be zero-valued")
// }
//
// if baseClient == nil {
// baseClient = defaultClient()
// }
//
// return &TokenClient{
// HttpClient: baseClient,
// Token:      token,
// }, nil
// }
//
// func (c *TokenClient) Do(req *http.Request) (*http.Response, error) {
// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
//
// return c.Client.Do(req)
// }
