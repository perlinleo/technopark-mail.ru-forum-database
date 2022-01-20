package model

type Forum struct {
	ID      int32  `json:"-"`
	Posts   int32  `json:"posts"`
	Slug    string `json:"slug"`
	Threads int32  `json:"threads"`
	Title   string `json:"title"`
	User    string `json:"user"`
}
