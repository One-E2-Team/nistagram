package model

type PersonalData struct {
	Name         string     `json:"name"`
	Surname      string     `json:"surname"`
	Telephone    string     `json:"telephone"`
	Gender       string     `json:"gender"`
	BirthDate    string     `json:"birthDate"`
	ProfileID    uint
}

