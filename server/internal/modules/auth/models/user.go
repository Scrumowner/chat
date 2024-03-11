package models

type Tabler interface {
	TableName() string
}

type User struct {
	ID       string `json:"id" db:"id" db_type:"text"`
	Username string `json:"username" db:"username" db_type:"text"`
	Password string `json:"password" db:"password" db_type:"text"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) GetID() string {
	return u.ID
}
func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) GetPassword() string {
	return u.Password
}

func (u *User) SetID(id string) {
	u.ID = id
}
func (u *User) SetUsername(username string) {
	u.Username = username
}
func (u *User) SetPassword(password string) {
	u.Password = password
}
