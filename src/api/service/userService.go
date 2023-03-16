package service

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/sqlgostructure/src/api/models"
	"github.com/myrachanto/sqlgostructure/src/api/repository"
	"github.com/myrachanto/sqlgostructure/src/emails"
)

var (
	UserService UserServiceInterface = &userService{}
)

type UserServiceInterface interface {
	Create(user *models.User) (string, httperrors.HttpErr)
	Login(auser *models.LoginUser) (*models.Auth, httperrors.HttpErr)
	RenewAccessToken(renew string) (*models.Auth, httperrors.HttpErr)
	Logout(token string) (string, httperrors.HttpErr)
	GetOne(id string) (*models.User, httperrors.HttpErr)
	GetAll(string) ([]*models.User, httperrors.HttpErr)
	// Forgot(email string) httperrors.HttpErr
	Update(id string, user *models.User) httperrors.HttpErr
	Delete(id int) (string, httperrors.HttpErr)
	PasswordUpdate(oldpassword, email, newpassword string) httperrors.HttpErr
	PasswordReset(email, password string) httperrors.HttpErr
}
type userService struct {
	repo repository.UserrepoInterface
}

func NewUserService(repository repository.UserrepoInterface) UserServiceInterface {
	return &userService{
		repository,
	}
}
func (service *userService) Create(user *models.User) (string, httperrors.HttpErr) {
	return service.repo.Create(user)
}

func (service *userService) Login(auser *models.LoginUser) (*models.Auth, httperrors.HttpErr) {
	// fmt.Println("====================frffqf")
	return service.repo.Login(auser)
}

func (service *userService) RenewAccessToken(renew string) (*models.Auth, httperrors.HttpErr) {
	return service.repo.RenewAccessToken(renew)
}
func (service *userService) Logout(token string) (string, httperrors.HttpErr) {
	return service.repo.Logout(token)
}
func (service *userService) GetOne(code string) (*models.User, httperrors.HttpErr) {
	return service.repo.GetOne(code)
}

func (service *userService) GetAll(search string) ([]*models.User, httperrors.HttpErr) {
	return service.repo.GetAll(search)
}

func (service *userService) Update(code string, user *models.User) httperrors.HttpErr {
	return service.repo.Update(code, user)
}

// func (service *userService) Forgot(email string) httperrors.HttpErr {
// 	e, p, err1 := service.repo.Forgot(email)
// 	go emails.ToEmailPassword(e, p)
// 	return err1
// }

func (service *userService) PasswordReset(email, password string) httperrors.HttpErr {
	_, err1 := service.repo.PasswordReset(email, password)
	// fmt.Println("====================frffqf")
	go emails.ToEmailPassword(password, email)
	return err1
}
func (service *userService) PasswordUpdate(oldpassword, email, newpassword string) httperrors.HttpErr {
	pass, email, err1 := service.repo.PasswordUpdate(oldpassword, email, newpassword)
	go emails.ToEmailPassword(pass, email)
	return err1
}
func (service *userService) Delete(id int) (string, httperrors.HttpErr) {
	return service.repo.DeleteById(id)
}
