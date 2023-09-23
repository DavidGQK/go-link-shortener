package models

type RequestShortenLink struct {
	URL string `json:"url"`
}

type ResponseShortenLink struct {
	Result string `json:"result"`
}
