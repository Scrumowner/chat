package service

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"server/internal/infrastructure/db"
	"server/internal/infrastructure/tool"
	"server/internal/modules/auth/models"
)

type Auther interface {
	Register(u models.User) error
	Login(u models.User) (string, string, string, error)
	Verify(t string) bool
}

type AuthService struct {
	adapter *db.SqlAdapter
	crypt   *tool.Bcrypt
	jwt     *tool.TokenManager
}

func NewAuthService(sql *sqlx.DB, token []byte) *AuthService {
	return &AuthService{
		adapter: db.NewSqlAdapter(sql),
		crypt:   tool.NewBcrypt(),
		jwt:     tool.NewTokenManager(token),
	}
}

func (a *AuthService) Register(u models.User) error {
	hash, err := a.crypt.GeneratePassword(u.Password)
	if err != nil {
		log.Println("error wgen generate password", err)
		return err
	}
	user := models.User{
		ID:       uuid.NewString(),
		Username: u.Username,
		Password: hash,
	}
	err = a.adapter.Insert(&user)
	if err != nil {
		log.Println("error when insert into db", err)
		return err

	}
	return nil

}

func (a *AuthService) Login(u models.User) (string, string, string, error) {
	users := []*models.User{}
	eq := make(map[string]interface{})
	eq["username"] = u.Username
	err := a.adapter.Select(&u, db.Condition{Eq: eq}, &users)
	if err != nil {
		log.Println("error when select from db", err)
		return "", "", "", err
	}
	if len(users) == 0 {
		return "", "", "", err
	}
	user := users[0]
	log.Println(user.ID, user.Username, user.Password)
	isCompare := a.crypt.CompareHashWithPassowrd(user.Password, u.Password)
	if !isCompare {
		log.Println("error when compare hash and password", err)
		return "", "", "", err
	}
	token, err := a.jwt.Sign(user.ID, user.Username)
	if err != nil {
		log.Println("error when sign user info to jwt ", err)
		return "", "", "", err
	}
	return token, user.ID, user.Username, nil

}

func (a *AuthService) Verify(t string) bool {
	claims, err := a.jwt.Unsign(t)
	if err != nil {
		log.Println("errr when undign token to claims", err)
		return false
	}
	log.Println(claims)
	eq := make(map[string]interface{})
	eq["id"] = claims.ID
	eq["username"] = claims.Username
	user := models.User{
		ID:       claims.ID,
		Username: claims.Username,
	}
	users := []*models.User{}
	err = a.adapter.Select(&user, db.Condition{Eq: eq}, &users)
	if err != nil {
		log.Println("errr when select user from db", err)
		return false
	}
	u := users[0]
	if u.ID != claims.ID && u.Username != claims.Username {
		return false
	}
	return true
}
