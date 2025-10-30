package entity

type User struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phonenumber"`
	Name        string `json:"name"`
	Password    string `json:"-"`
}
