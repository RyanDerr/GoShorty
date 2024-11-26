package response

type EntryListResponse struct {
	Entries []EntryResponse `json:"entries"`
}

type EntryResponse struct {
	Id      uint   `json:"id"`
	Content string `json:"content"`
	UserId  uint   `json:"user_id"`
}
