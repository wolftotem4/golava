package binding

type Register struct {
	Username string `json:"username" form:"username" mod:"trim" binding:"required"`
	Password string `json:"password" form:"password" binding:"required,min=8"`
}
