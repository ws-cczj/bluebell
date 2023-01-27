package serializer

type ResponseUserLogin struct {
	UserId string `json:"user_id"`
	Atoken string `json:"atoken"`
	Rtoken string `json:"rtoken"`
}

type ResponseUserFollow struct {
	UserId   string `json:"user_id" db:"user_id"`
	Username string `json:"username" db:"username"`
}
