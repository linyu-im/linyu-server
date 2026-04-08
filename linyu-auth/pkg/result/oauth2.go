package result

type Oauth2UrlResult struct {
	AuthUrl     string `json:"authUrl"`
	RedirectUrl string `json:"redirectUrl"`
}
