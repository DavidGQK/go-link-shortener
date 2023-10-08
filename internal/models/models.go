package models

type RequestShortenLink struct {
	URL string `json:"url"`
}

type ResponseShortenLink struct {
	Result string `json:"result"`
}

type RequestLinks struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type RequestBatchLinks []RequestLinks

type ResponseLinks struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type ResponseBatchLinks []ResponseLinks
