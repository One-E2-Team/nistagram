package model

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Privilege struct {
	Name string `json:"name"`
}

type Role struct {
	Name       string      `json:"name"`
	Privileges []Privilege `json:"privileges"`
}

type User struct {
	Credentials Credentials `json:"credentials"`
	Salt        string      `json:"salt"`
	APIToken    string      `json:"apiToken"`
	IsDeleted   bool        `json:"isDeleted"`
	Roles       []Role      `json:"roles"`
}
