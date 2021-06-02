package dto

type PersonalDataDTO struct {
	Username  string `json:"username"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birthDate"`
	Biography string `json:"biography"`
	Website   string `json:"website"`
}
