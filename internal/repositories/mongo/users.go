package repositories

type User struct {
	FirstName string `json:"first_name"`

	LastName string `json:"last_name"`

	Nickname string `json:"nickname"`

	Password string `json:"password"`

	Email string `json:"email"`

	Country string `json:"country"`

	//TODO: created_at, updated_at
}
