package serializer

type ResponseUserLogin struct {
	UserId int64  `json:"user_id"`
	Atoken string `json:"atoken"`
	Rtoken string `json:"rtoken"`
}

type ResponseUserFollow struct {
	UserId   int64  `json:"user_id" db:"user_id"`
	Username string `json:"username" db:"username"`
}
