package schema

type Ticket struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content"`
	Assignee *int   `json:"assignee"`
}
