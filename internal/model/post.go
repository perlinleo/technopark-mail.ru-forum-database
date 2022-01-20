package model

import "time"

type Post struct {
	ID            int32     `json:"id,omitempty"`
	Parent        int32     `json:"parent"`
	Thread        int32     `json:"thread,omitempty"`
	Forum         string    `json:"forum,omitempty"`
	Author        string    `json:"author"`
	Created       time.Time `json:"created,omitempty"`
	IsEdited      bool      `json:"isEdited"`
	Message       string    `json:"message"`
	Path          []int32   `json:"-"`
	Childs        []Post    `json:"childs,omitempty"`
	ParentPointer *Post     `json:"-"`
}

type PostFull struct {
	Author *User   `json:"author,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
	Post   *Post   `json:"post,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
}

type PostUpdate struct {
	Message string `json:"message,omitempty"`
}
