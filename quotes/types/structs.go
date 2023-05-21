package types

type Quote struct {
	Id     int    `json:"id"`
	Quote  string `json:"quote"`
	Author string `json:"author"`
	Sub_by string `json:"sub_by"`
	Date   string `json:"date"`
	Active int    `json:"active"`
}

type PostData struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
	Sub_by string `json:"sub_by"`
}