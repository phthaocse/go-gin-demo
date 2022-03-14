package schema

type DefaultQuery struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}
