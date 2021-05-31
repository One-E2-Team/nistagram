package dto

type RegistrationDto struct {
	Username string `json:"username"`
	Name string `json:"name"`
	Surname string `json:"surname"`
	Email string `json:"email"`
	Telephone string `json:"telephone"`
	Gender string `json:"gender"`
	BirthDate string `json:"birthDate"`
	IsPrivate bool `json:"isPrivate"`
	Biography string `json:"biography"`
	WebSite string `json:"website"`
	InterestedIn []string `json:"interestedIn"`
}
