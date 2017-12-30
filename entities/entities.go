package entities

// TweeterStats blablabla
type TweeterStats struct {
	FullName string `json:"fullName"`
	Username string `json:"username"`

	TweetsCount uint `json:"tweetsCount"`
}

// Tweeter blablabla
type Tweeter struct {
	FullName string
	Username string
}
