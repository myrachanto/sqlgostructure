package models

import (
	"regexp"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	httperrors "github.com/myrachanto/erroring"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ExpiresAt = time.Now().Add(time.Hour * 24).Unix()

type User struct {
	FName      string     `json:"f_name"`
	LName      string     `json:"l_name"`
	UName      string     `json:"u_name"`
	Phone      string     `json:"phone"`
	Address    string     `json:"address"`
	Dob        *time.Time `json:"dob"`
	Picture    string     `json:"picture"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	Role       string     `json:"role"`
	Admin      string     `json:"admin"`
	Supervisor string     `json:"supervisor"`
	Employee   string     `json:"employee"`
	Usercode   string     `json:"usercode"`
	Shopalias  string     `json:"shopalias"`
	gorm.Model
}
type UserDTO struct {
	FName     string `json:"f_name"`
	LName     string `json:"l_name"`
	UName     string `json:"u_name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Picture   string `json:"picture"`
	Email     string `json:"email"`
	Usercode  string `json:"usercode"`
	Shopalias string `json:"shopalias"`
	gorm.Model
}
type Verify struct {
	Question string `json:"question,omitempty"`
	Answer   string `json:"answer,omitempty"`
	Hint     string `json:"hint,omitempty"`
}
type Session struct {
	Code         string    `json:"code,omitempty"`
	Username     string    `json:"username,omitempty"`
	Usercode     string    `json:"usercode,omitempty"`
	RefleshToken string    `json:"reflesh_token,omitempty"`
	TokenId      string    `json:"token_id,omitempty"`
	Shopalias    string    `json:"shopalias,omitempty"`
	Bizname      string    `json:"bizname,omitempty"`
	UserAgent    string    `json:"user_agent,omitempty"`
	ClientIp     string    `json:"client_ip,omitempty"`
	IsBlocked    bool      `json:"is_blocked,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
	gorm.Model
}
type Auth struct {
	//User User `gorm:"foreignKey:UserID; not null"`
	Admin               string    `json:"admin,omitempty"`
	Supervisor          string    `json:"supervisor,omitempty"`
	Employee            string    `json:"employee,omitempty"`
	Usercode            string    `json:"usercode,omitempty"`
	UName               string    `json:"uname,omitempty"`
	Picture             string    `json:"picture,omitempty"`
	Token               string    `bson:"token" json:"token,omitempty"`
	TokenExpires        time.Time `json:"token_expires,omitempty"`
	RefleshToken        string    `json:"reflesh_token,omitempty"`
	RefleshTokenExpires time.Time `json:"reflesh_token_expires,omitempty"`
	SessionCode         string    `json:"session_code,omitempty"`
	Role                string    `json:"role,omitempty"`
	Shopalias           string    `json:"shopalias,omitempty"`
	Business            string    `json:"business,omitempty"`
	Bizname             string    `json:"bizname,omitempty"`
	// Base                support.Base `json:"base,omitempty"`
}
type LoginUser struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserAgent string `json:"useragent"`
}

// Token struct declaration
type Token struct {
	Usercode   string `json:"usercode,omitempty"`
	UName      string `json:"uname,omitempty"`
	Email      string `json:"email,omitempty"`
	Role       string `json:"role,omitempty"`
	Admin      string `json:"admin,omitempty"`
	Supervisor string `json:"supervisor,omitempty"`
	Employee   string `json:"employee,omitempty"`
	*jwt.StandardClaims
}

func (user *User) ConvertUserToDTO() *UserDTO {
	return &UserDTO{
		// ID:            user.ID,
		FName:    user.FName,
		LName:    user.LName,
		UName:    user.UName,
		Phone:    user.Phone,
		Address:  user.Address,
		Picture:  user.Picture,
		Email:    user.Email,
		Usercode: user.Usercode,
	}
}

func (user User) ValidateEmail(email string) (matchedString bool) {
	stringresults := httperrors.ValidStringNotEmpty(email)
	if stringresults.Noerror() {
		return false
	}
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&amp;'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	matchedString = re.MatchString(email)
	return matchedString
}
func ValidateEmail(email string) (matchedString bool) {
	stringresults := httperrors.ValidStringNotEmpty(email)
	if stringresults.Noerror() {
		return false
	}
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&amp;'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	matchedString = re.MatchString(email)
	return matchedString
}
func (user User) ValidatePassword(password string) (bool, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(password)
	if stringresults.Noerror() {
		return false, stringresults
	}
	if len(password) < 5 {
		return false, httperrors.NewBadRequestError("your password need more characters!")
	} else if len(password) > 32 {
		return false, httperrors.NewBadRequestError("your password is way too long!")
	}
	return true, nil
}
func HashPassword(password string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(password)
	if stringresults.Noerror() {
		return "", httperrors.NewBadRequestError("your password Must not be empty!")
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", httperrors.NewNotFoundError("soemthing went wrong!")
	}
	return string(pass), nil

}
func (user LoginUser) Compare(p1, p2 string) bool {
	stringresults := httperrors.ValidStringNotEmpty(p1)
	if stringresults.Noerror() {
		return false
	}
	stringresults2 := httperrors.ValidStringNotEmpty(p2)
	if stringresults2.Noerror() {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(p2), []byte(p1))
	return err == nil
}
func (user User) Compare(p1, p2 string) bool {
	stringresults := httperrors.ValidStringNotEmpty(p1)
	if stringresults.Noerror() {
		return false
	}
	stringresults2 := httperrors.ValidStringNotEmpty(p2)
	if stringresults2.Noerror() {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(p2), []byte(p1))
	return err == nil
}
func (loginuser LoginUser) Validate() httperrors.HttpErr {
	if loginuser.Email == "" {
		return httperrors.NewNotFoundError("Invalid Email")
	}
	if loginuser.Password == "" {
		return httperrors.NewNotFoundError("Invalid password")
	}
	return nil
}
func (user User) Validate() httperrors.HttpErr {
	if user.FName == "" {
		return httperrors.NewNotFoundError("Invalid first Name")
	}
	if user.LName == "" {
		return httperrors.NewNotFoundError("Invalid last name")
	}
	if user.UName == "" {
		return httperrors.NewNotFoundError("Invalid username")
	}
	if user.Phone == "" {
		return httperrors.NewNotFoundError("Invalid phone number")
	}
	if user.Email == "" {
		return httperrors.NewNotFoundError("Invalid Email")
	}
	// if user.Address == "" {
	// 	return httperrors.NewNotFoundError("Invalid Address")
	// }
	if user.Password == "" {
		return httperrors.NewNotFoundError("Invalid password")
	}
	// if user.Picture == "" {
	// 	return httperrors.NewNotFoundError("Invalid picture")
	// }
	return nil
}

func (verify Verify) Validate() httperrors.HttpErr {
	if verify.Question == "" {
		return httperrors.NewNotFoundError("Invalid question")
	}
	if verify.Answer == "" {
		return httperrors.NewNotFoundError("Invalid aswer")
	}
	if verify.Hint == "" && verify.Hint == verify.Answer {
		return httperrors.NewNotFoundError("Invalid hint")
	}
	return nil
}
func (verify Verify) HashAwnser(p string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(p)
	if stringresults.Noerror() {
		return "", httperrors.NewBadRequestError("your password Must not be empty!")
	}
	passAnswer, err := bcrypt.GenerateFromPassword([]byte(verify.Answer), 10)
	if err != nil {
		return "", httperrors.NewNotFoundError("type a stronger password!")
	}
	return string(passAnswer), nil

}
