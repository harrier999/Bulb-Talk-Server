package models

type User struct {
	user_id       string
	username      string
	password_hash string
	profile_image string
	phone_number  string
	country_code  string
}

type JoinRoom struct {
	room_id    string
	user_id    string
	message_id string
}
