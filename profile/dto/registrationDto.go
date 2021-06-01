package dto

type RegistrationDto struct {
	Username string `json:"username" validate:"required"`
	Name string `json:"name" validate:"required"`
	Surname string `json:"surname" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Telephone string `json:"telephone" validate:"required"`
	Gender string `json:"gender" validate:"required"`
	BirthDate string `json:"birthDate" validate:"required"`
	IsPrivate bool `json:"isPrivate"`
	Biography string `json:"biography" validate:"required"`
	WebSite string `json:"website"`
	InterestedIn []string `json:"interestedIn"`
}
