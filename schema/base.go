package schema

type DefaultQuery struct {
	Limit  int `form:"limit" binding:"gte=0"`
	Offset int `form:"offset" binding:"gte=0"`
}
