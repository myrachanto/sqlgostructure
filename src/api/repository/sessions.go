package repository

import (
	"strconv"
	"time"

	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/sqlgostructure/src/api/models"
	"github.com/myrachanto/sqlgostructure/src/support"
)

var Sessionsrepo sessionsrepo

type sessionsrepo struct{}

// /Dealing with sessions
func (r *sessionsrepo) CreateSession(session *models.Session) (*models.Session, httperrors.HttpErr) {

	code, err1 := r.GeneSessioncode(session.Usercode)
	if err1 != nil {
		return nil, err1
	}
	session.Code = code
	gorm, err := IndexRepo.Getconnected()
	if err != nil {
		return nil, err
	}
	defer IndexRepo.DbClose(gorm)
	errs := gorm.Create(session).Error
	if errs != nil {
		return nil, httperrors.NewNotFoundError("Error updating!")
	}

	return session, nil
}
func (r sessionsrepo) GeneSessioncode(usercode string) (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[1:5]
	count, err := r.Count()
	co := count + 1
	if err != nil {
		return "", httperrors.NewNotFoundError("no results found")
	}
	cod := "sessions-" + strconv.FormatUint(uint64(co), 10) + "-" + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r sessionsrepo) GeneTokencode(usercode string) (string, httperrors.HttpErr) {

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp
	cod := "token-" + usercode + special
	code := support.Hasher(cod)
	if code == "" {
		return "", httperrors.NewNotFoundError("THe string is empty")
	}
	return code, nil
}
func (r *sessionsrepo) GetOne(code string) (session *models.Session, errors httperrors.HttpErr) {
	stringresults := httperrors.ValidStringNotEmpty(code)
	if stringresults.Noerror() {
		return nil, stringresults
	}
	gorm, err := IndexRepo.Getconnected()
	if err != nil {
		return nil, err
	}
	defer IndexRepo.DbClose(gorm)
	errs := gorm.Where("code = ?", code).First(&session).Error
	if errs != nil {
		return nil, httperrors.NewNotFoundError("no results found")
	}
	return session, nil
}
func (r sessionsrepo) Count() (int64, httperrors.HttpErr) {
	gorm, err := IndexRepo.Getconnected()
	if err != nil {
		return 0, err
	}
	defer IndexRepo.DbClose(gorm)
	result := &models.Session{}
	var count int64
	errs := gorm.Find(&result).Count(&count).Error
	if errs != nil {
		return 0, httperrors.NewNotFoundError("Couldnt count the results")
	}
	return count, nil
}
