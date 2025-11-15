package dto

type SearchReq struct {
	Query  string `query:"query" validate:"required,min=1,max=100"`
	Limit  int    `query:"limit" validate:"min=1,max=100"`
	Offset int    `query:"offset" validate:"min=0"`
}

type SearchRes struct {
	Results interface{} `json:"results"`
	Count   int64       `json:"count"`
	Query   string      `json:"query"`
}
