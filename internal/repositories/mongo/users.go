package repositories

type User struct {
	Id        string `json:"id" bson:"_id,omitempty"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Nickname  string `json:"nickname" bson:"nickname"`
	Password  string `json:"password" bson:"password"`
	Email     string `json:"email" bson:"email"`
	Country   string `json:"country" bson:"country"`

	//TODO: created_at, updated_at
}
