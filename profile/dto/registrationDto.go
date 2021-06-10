package dto

type RegistrationDto struct {
	Username     string   `json:"username" validate:"required,bad_username"`
	Password     string   `json:"password" validate:"required,common_pass,weak_pass"`
	Name         string   `json:"name" validate:"required"`
	Surname      string   `json:"surname" validate:"required"`
	Email        string   `json:"email" validate:"required,email"`
	Telephone    string   `json:"telephone"`
	Gender       string   `json:"gender"`
	BirthDate    string   `json:"birthDate"`
	IsPrivate    bool     `json:"isPrivate"`
	Biography    string   `json:"biography"`
	WebSite      string   `json:"website"`
	InterestedIn []string `json:"interestedIn"`
}
