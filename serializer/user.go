package serializer

type ResponseUserLogin struct {
	UserId int64  `json:"user_id"`
	Atoken string `json:"atoken,omitempty"`
	Rtoken string `json:"rtoken,omitempty"`
}

type ResponseUserFollow struct {
	UserId   int64  `json:"user_id" db:"user_id"`
	Username string `json:"username" db:"username"`
}
