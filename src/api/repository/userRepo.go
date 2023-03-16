package repository

import (
	"fmt"
	"strconv"
	"time"

	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/sqlgostructure/src/api/models"
	"github.com/myrachanto/sqlgostructure/src/pasetos"
	"github.com/myrachanto/sqlgostructure/src/support"
)

// Userrepository repository
var (
	Userrepository UserrepoInterface = &userrepository{}
	Userrepo                         = userrepository{}
)

type Key struct {
	EncryptionKey string `mapstructure:"EncryptionKey"`
}

type UserrepoInterface interface {
	Create(user *models.User) (string, httperrors.HttpErr)
	Login(auser *models.LoginUser) (*models.Auth, httperrors.HttpErr)
	RenewAccessToken(renewAccesstoken string) (*models.Auth, httperrors.HttpErr)
	Logout(token string) (string, httperrors.HttpErr)
	GetOne(id string) (*models.User, httperrors.HttpErr)
	GetAll(string) ([]*models.User, httperrors.HttpErr)
	// Forgot(email string) (string, string, httperrors.HttpErr)
	DeleteById(id int) (string, httperrors.HttpErr)
	Update(string, *models.User) httperrors.HttpErr
	PasswordUpdate(oldpassword, email, newpassword string) (string, string, httperrors.HttpErr)
	PasswordReset(email, newpassword string) (string, httperrors.HttpErr)
	Count() (int64, httperrors.HttpErr)
}
type userrepository struct{}

func NewUserRepo() UserrepoInterface {
	return &userrepository{}
}

func (r *userrepository) Create(user *models.User) (string, httperrors.HttpErr) {
	if err1 := user.Validate(); err1 != nil {
		return "", err1
	}
	ok, err1 := user.ValidatePassword(user.Password)
	if !ok {
		return "", err1
	}
	ok = user.ValidateEmail(user.Email)
	if !ok {
		return "", httperrors.NewNotFoundError("Your email format is wrong!")
	}
	code, errs := r.genecode()
	if errs != nil {
		return "", errs
	}
	user.Usercode = code

	//get user

	ok = r.emailexist(user.Email)
	if ok {
		return "", httperrors.NewNotFoundError("Email eist")
	}
	count := 0
	gorm, err := IndexRepo.Getconnected()
	if err != nil {
		return "", err
	}
	defer IndexRepo.DbClose(gorm)
	if count >= 1 {
		user.Role = "level3"
		user.Admin = "notadmin"
		user.Supervisor = "notsupervisor"
		user.Employee = "employee"

		hashpassword, err2 := models.HashPassword(user.Password)
		if err2 != nil {
			return "", err2
		}
		user.Password = hashpassword
		errs := gorm.Create(user).Error
		if errs != nil {
			return "", httperrors.NewBadRequestError(fmt.Sprintf("Create user Failed, %d", err))
		}
		return "user succesifully created", nil

	} else {
		user.Role = "level3"
		user.Admin = "admin"
		user.Supervisor = "supervisor"
		user.Employee = "employee"

		hashpassword, err2 := models.HashPassword(user.Password)
		if err2 != nil {
			return "", err2
		}
		user.Password = hashpassword

		//insert user
		errs := gorm.Create(user).Error
		if errs != nil {
			return "", httperrors.NewBadRequestError(fmt.Sprintf("Create user Failed, %d", err))
		}
		return "user successifully created", nil
	}
}

func (r *userrepository) Login(user *models.LoginUser) (*models.Auth, httperrors.HttpErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	u, err := r.getOneByEmail(user.Email)
	if err != nil {
		return nil, httperrors.NewNotFoundError("wrong email password combo!")
	}

	ok := user.Compare(user.Password, u.Password)
	if !ok {
		return nil, httperrors.NewNotFoundError("wrong email password combo!")
	}
	maker, errs := pasetos.NewPasetoMaker()
	if err != nil {
		return nil, errs
	}
	tokencode, errs := Sessionsrepo.GeneTokencode(u.Usercode)
	if errs != nil {
		return nil, errs
	}
	renewtokencode, errs := Sessionsrepo.GeneSessioncode(u.Usercode)
	if errs != nil {
		return nil, errs
	}
	data := &pasetos.Data{
		Code:     tokencode,
		Usercode: u.Usercode,
		Username: u.UName,
		Email:    u.Email,
		Admin:    u.Admin,
	}
	tokenString, payload, errs := maker.CreateToken(data, time.Hour*3)
	if errs != nil {
		return nil, errs
	}
	data.Code = renewtokencode
	RefleshToken, refleshtok, errs := maker.CreateToken(data, time.Hour*24)
	if errs != nil {
		return nil, errs
	}
	sessiond, errs := Sessionsrepo.CreateSession(&models.Session{
		Code: renewtokencode,
		// TokenId:      tokencode,
		Username:     u.UName,
		Usercode:     u.Usercode,
		RefleshToken: RefleshToken,
		UserAgent:    user.UserAgent,
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refleshtok.ExpiredAt,
	})
	if errs != nil {
		return nil, errs
	}
	auths := &models.Auth{Usercode: u.Usercode, Role: u.Role, Admin: u.Admin, Supervisor: u.Supervisor, Employee: u.Employee, Picture: u.Picture, UName: u.UName, Token: tokenString, RefleshToken: RefleshToken, SessionCode: sessiond.Code, TokenExpires: payload.ExpiredAt, RefleshTokenExpires: sessiond.ExpiresAt}
	return auths, nil
}

func (r *userrepository) RenewAccessToken(renewAccesstoken string) (*models.Auth, httperrors.HttpErr) {
	maker, err := pasetos.NewPasetoMaker()
	if err != nil {
		return nil, err
	}
	refleshpayload, err := maker.VerifyToken(renewAccesstoken)
	if err != nil {
		return nil, err
	}
	sessions, err := Sessionsrepo.GetOne(refleshpayload.Code)
	if err != nil {
		return nil, err
	}
	if sessions.IsBlocked {
		if err != nil {
			return nil, httperrors.NewAnuthorizedError("your Session is blocked")
		}
	}
	if sessions.Username != refleshpayload.Username {
		if err != nil {
			return nil, httperrors.NewAnuthorizedError("your Session is blocked -u")
		}
	}

	tokencode, errs := Sessionsrepo.GeneTokencode(sessions.Usercode)
	if errs != nil {
		return nil, errs
	}
	tokenString, payload, errs := maker.CreateToken(&pasetos.Data{
		Username: refleshpayload.Username,
		Code:     tokencode,
		Usercode: sessions.Usercode,
		Email:    refleshpayload.Email,
		Admin:    refleshpayload.Admin,
		Bizname:  refleshpayload.Bizname,
	}, time.Hour*1)
	if errs != nil {
		return nil, errs
	}
	auths := &models.Auth{Usercode: sessions.Usercode, Admin: refleshpayload.Admin, UName: sessions.Username, Token: tokenString, TokenExpires: payload.ExpiredAt, RefleshTokenExpires: sessions.ExpiresAt}
	return auths, nil
}
func (r *userrepository) Logout(token string) (string, httperrors.HttpErr) {

	stringresults := httperrors.ValidStringNotEmpty(token)
	if stringresults.Noerror() {
		return "", stringresults
	}
	// collection := db.Mongodb.Collection("auth")
	// filter1 := bson.M{"token": token}
	// _, err3 := collection.DeleteOne(ctx, filter1)

	//inser user
	var err3 any
	if err3 != nil {
		return "", httperrors.NewBadRequestError("something went wrong login out!")
	}
	return "something went wrong login out!", nil
}
func (r *userrepository) GetOne(code string) (user *models.User, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	gorm, err := IndexRepo.Getconnected()
	if err != nil {
		return nil, err
	}
	defer IndexRepo.DbClose(gorm)
	errs := gorm.Where("usercode = ?", code).First(&user).Error
	if errs != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return user, nil
}
func (r *userrepository) GetAll(search string) ([]*models.User, httperrors.HttpErr) {
	results := []*models.User{}
	gorm, err := IndexRepo.Getconnected()
	if err != nil {
		return nil, err
	}
	defer IndexRepo.DbClose(gorm)
	errs := gorm.Find(&results).Error
	if errs != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return results, nil

}

func (r *userrepository) PasswordUpdate(oldpassword, email, newpassword string) (string, string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(oldpassword)
	if stringresults.Noerror() {
		return "", "", stringresults
	}
	stringresults2 := httperrors.ValidStringNotEmpty(email)
	if stringresults2.Noerror() {
		return "", "", stringresults2
	}
	stringresults3 := httperrors.ValidStringNotEmpty(newpassword)
	if stringresults3.Noerror() {
		return "", "", stringresults3
	}
	gorm, err := IndexRepo.Getconnected()
	if err != nil {
		return "", "", err
	}
	defer IndexRepo.DbClose(gorm)
	result := &models.User{}
	errs := gorm.Where("email = ? AND password = ?", email, oldpassword).First(result).Error
	if errs != nil {
		return "", "", httperrors.NewNotFoundError("Those credentials doesnt match!")
	}
	hashpassword, err2 := models.HashPassword(newpassword)
	if err2 != nil {
		return "", "", err2
	}
	errs = gorm.Model(result).Where("email = ?", email).Update("password", hashpassword).Error
	if errs != nil {
		return "", "", httperrors.NewNotFoundError("Error updating!")
	}
	return newpassword, email, nil
}
func (r *userrepository) PasswordReset(email, password string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(password)
	if stringresults.Noerror() {
		return "", stringresults
	}
	stringresults2 := httperrors.ValidStringNotEmpty(email)
	if stringresults2.Noerror() {
		return "", stringresults2
	}
	hashpassword, err2 := models.HashPassword(password)
	if err2 != nil {
		return "", err2
	}
	gorm, err := IndexRepo.Getconnected()
	if err != nil {
		return "", err
	}
	defer IndexRepo.DbClose(gorm)
	result := &models.User{}
	errs := gorm.Where("email = ? AND password = ?", email, password).First(result).Error
	if errs != nil {
		return "", httperrors.NewNotFoundError("Those credentials doesnt match!")
	}
	errs = gorm.Model(result).Where("email = ?", email).Update("password", hashpassword).Error
	if errs != nil {
		return "", httperrors.NewNotFoundError("Error updating!")
	}
	return "password uppdated successifully", nil
}

func (r *userrepository) Update(id string, user *models.User) httperrors.HttpErr {
	stringresults := httperrors.ValidStringNotEmpty(id)
	if stringresults.Noerror() {
		return stringresults
	}
	ok, err1 := user.ValidatePassword(user.Password)
	if !ok {
		return err1
	}
	ok = user.ValidateEmail(user.Email)
	if !ok {
		return httperrors.NewNotFoundError("Your email format is wrong!")
	}
	hashpassword, err2 := models.HashPassword(user.Password)
	if err2 != nil {
		return err2
	}
	user.Password = hashpassword
	uuser, err := r.getOneByCode(user.Usercode)
	if err != nil {
		return err
	}
	if user.FName == "" {
		user.FName = uuser.FName
	}
	if user.LName == "" {
		user.LName = uuser.LName
	}
	if user.UName == "" {
		user.UName = uuser.UName
	}
	if user.Phone == "" {
		user.Phone = uuser.Phone
	}
	if user.Address == "" {
		user.Address = uuser.Address
	}
	if user.Picture == "" {
		user.Picture = uuser.Picture
	}
	if user.Email == "" {
		user.Email = uuser.Email
	}
	if hashpassword == "" {
		user.Password = uuser.Password
	}
	gorm, err := IndexRepo.Getconnected()
	if err != nil {
		return err
	}
	defer IndexRepo.DbClose(gorm)
	errs := gorm.Save(user).Error
	if errs != nil {
		return httperrors.NewBadRequestError("Failed to update the resource")
	}
	return nil
}
func (r userrepository) DeleteById(id int) (string, httperrors.HttpErr) {
	gorm, err := IndexRepo.Getconnected()
	if err != nil {
		return "", err
	}
	defer IndexRepo.DbClose(gorm)
	result := &models.User{}
	errs := gorm.Where("id = ?", id).Delete(result).Error
	if errs != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	return "deleted successfully", nil

}
func (r userrepository) DeleteByCode(code string) (string, httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return "", stringresults
	}
	gorm, err := IndexRepo.Getconnected()
	if err != nil {
		return "", err
	}
	defer IndexRepo.DbClose(gorm)
	result := &models.User{}
	errs := gorm.Where("usercode = ?", code).Delete(result).Error
	if errs != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	return "deleted successfully", nil

}
func (r userrepository) genecode() (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	count, err := r.Count()
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	cod := "UserCode-" + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r userrepository) getOneByCode(code string) (result *models.User, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	gorm, err := IndexRepo.Getconnected()
	if err != nil {
		return nil, err
	}
	defer IndexRepo.DbClose(gorm)
	errs := gorm.Where("usercode = ?", code).First(&result).Error
	if errs != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return result, nil
}
func (r *userrepository) getOneByEmail(email string) (result *models.User, err httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(email)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	ok := models.ValidateEmail(email)
	if !ok {
		return nil, httperrors.NewNotFoundError("Email validation failed")
	}
	gorm, err := IndexRepo.Getconnected()
	if err != nil {
		return nil, err
	}
	defer IndexRepo.DbClose(gorm)
	errs := gorm.Where("email = ?", email).First(&result).Error
	if errs != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return result, nil
}
func (r userrepository) emailexist(email string) bool {
	stringresults := httperrors.ValidStringNotEmpty(email)
	if stringresults.Noerror() {
		return stringresults.Noerror()
	}
	gorm, err := IndexRepo.Getconnected()
	if err != nil {
		return false
	}
	defer IndexRepo.DbClose(gorm)
	result := &models.User{}
	errs := gorm.Where("email = ?", email).First(&result).Error
	return errs == nil
}

func (r userrepository) Count() (int64, httperrors.HttpErr) {
	gorm, err := IndexRepo.Getconnected()
	if err != nil {
		return 0, err
	}
	defer IndexRepo.DbClose(gorm)
	result := &models.User{}
	var count int64
	errs := gorm.Find(&result).Count(&count).Error
	if errs != nil {
		return 0, httperrors.NewNotFoundError("Couldnt count the results")
	}
	return count, nil
}
