package binding

type Login struct {
	Username string `json:"username" form:"username" mod:"trim" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Remember bool   `json:"remember" form:"remember"`
}
