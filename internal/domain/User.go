package domain

type User struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
	Code      string `json:"code"`
	Expiry    string `json:"expiry"`
	Verified  bool   `json:"verified"`
	UserType  string `json:"user_type"`
}
