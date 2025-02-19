package dtos

type UserResponse struct {
	UserId    int    `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `josn:"lastName"`
	Email     string `json:"email"`
}
