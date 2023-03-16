package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/imagery"
	model "github.com/myrachanto/sqlgostructure/src/api/models"
	"github.com/myrachanto/sqlgostructure/src/api/service"
)

// UserController ...
var (
	UserController UserControllerInterface = &userController{}
)

type UserControllerInterface interface {
	Create(c echo.Context) error
	Login(c echo.Context) error
	RenewAccessToken(c echo.Context) error
	Logout(c echo.Context) error
	GetOne(c echo.Context) error
	// Forgot(c echo.Context) error
	GetAll(c echo.Context) error
	Update(c echo.Context) error
	PasswordUpdate(c echo.Context) error
	PasswordReset(c echo.Context) error
	Delete(c echo.Context) error
}

type userController struct {
	service service.UserServiceInterface
}

func NewUserController(ser service.UserServiceInterface) UserControllerInterface {
	return &userController{
		ser,
	}
}

// ///////controllers/////////////////

// Create godoc
// @Summary Create a user
// @Description Create a new user item
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/users [post]
func (controller userController) Create(c echo.Context) error {

	user := &model.User{}
	user.FName = c.FormValue("fname")
	user.LName = c.FormValue("lname")
	user.UName = c.FormValue("uname")
	user.Phone = c.FormValue("phone")
	user.Address = c.FormValue("address")
	user.Email = c.FormValue("email")
	user.Password = c.FormValue("password")
	// fmt.Println("----------------------------step1")
	// user.Business = c.FormValue("business")

	pic, err2 := c.FormFile("picture")
	if pic != nil {
		//    fmt.Println(pic.Filename)
		if err2 != nil {
			httperror := httperrors.NewBadRequestError("Invalid picture")
			return c.JSON(httperror.Code(), err2.Error())
		}
		src, err := pic.Open()
		if err != nil {
			httperror := httperrors.NewBadRequestError("the picture is corrupted")
			return c.JSON(httperror.Code(), err.Error())
		}
		defer src.Close()
		// filePath := "./public/imgs/users/"
		filePath := "./src/public/imgs/users/" + user.FName + pic.Filename
		filePath1 := "/imgs/users/" + user.FName + pic.Filename
		// Destination
		dst, err4 := os.Create(filePath)
		if err4 != nil {
			httperror := httperrors.NewBadRequestError("the Directory mess")
			return c.JSON(httperror.Code(), err4.Error())
		}
		defer dst.Close()
		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			if err2 != nil {
				httperror := httperrors.NewBadRequestError("error filling")
				return c.JSON(httperror.Code(), httperror.Message())
			}

		}
		// fmt.Println("----------------------------step2")
		//resize the image and replace the old one
		imagery.Imageryrepository.Imagetype(filePath, filePath, 400, 500)

		// fmt.Println("----------------------------step3")
		user.Picture = filePath1
		_, err1 := controller.service.Create(user)
		// fmt.Println("----------------------------step4", err1)
		if err1 != nil {
			return c.JSON(err1.Code(), err1.Message())
		}
		if _, err = io.Copy(dst, src); err != nil {
			if err2 != nil {
				httperror := httperrors.NewBadRequestError("error filling")
				return c.JSON(httperror.Code(), httperror.Message())
			}
		}
		return c.JSON(http.StatusCreated, "user created succesifully")
	} else {
		_, err1 := controller.service.Create(user)
		if err1 != nil {
			return c.JSON(err1.Code(), err1.Message())
		}
		return c.JSON(http.StatusCreated, "user created succesifully")
	}
}

// Login godoc
// @Summary Login a user
// @Description Login user
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} User
// @Failure 400 {object} support.HttpError
// @Router /front/login [post]
func (controller userController) Login(c echo.Context) error {
	user := &model.LoginUser{}
	user.Email = c.FormValue("email")
	user.Password = c.FormValue("password")
	user.UserAgent = c.Request().UserAgent()
	auth, problem := controller.service.Login(user)
	if problem != nil {
		fmt.Println(problem)
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, auth)
}

// renewAccesstoken godoc
// @Summary renewAccesstoken a user
// @Description renewAccesstoken user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Failure 400 {object} support.HttpError
// @Router /front/renewAccesstoken [post]
func (controller userController) RenewAccessToken(c echo.Context) error {
	renewaccessToken := c.FormValue("renewaccessToken")
	// fmt.Println(">>>>>>>>>>>>>>>>>>>", renewaccessToken)
	auth, problem := controller.service.RenewAccessToken(renewaccessToken)
	if problem != nil {
		fmt.Println(problem)
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, auth)
}

// logout godoc
// @Summary logout a user
// @Description logout user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/logout [post]
func (controller userController) Logout(c echo.Context) error {
	token := string(c.Param("token"))
	_, problem := controller.service.Logout(token)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "succeessifully logged out")
}

// GetAll godoc
// @Summary GetAll a user
// @Description Getall users
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/users [get]
func (controller userController) GetAll(c echo.Context) error {
	search := c.QueryParam("search")
	users, err3 := controller.service.GetAll(search)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, users)
}

// @Summary Get a user
// @Description Get item
// @Tags users
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/users [get]
func (controller userController) GetOne(c echo.Context) error {
	code := c.Param("code")
	user, problem := controller.service.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, user)
}

// Forgot godoc
// @Summary Forgot a user
// @Description Forgot user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Failure 400 {object} support.HttpError
// @Router /front/forgot [post]
// func (controller userController) Forgot(c echo.Context) error {
// 	email := c.FormValue("email")
// 	problem := controller.service.Forgot(email)
// 	if problem != nil {
// 		return c.JSON(problem.Code(), problem.Message())
// 	}
// 	return c.JSON(http.StatusOK, "updated succesifully")
// }

// Forgot godoc
// @Summary Forgot a user
// @Description Forgot user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Failure 400 {object} support.HttpError
// @Router /front/forgot [post]
func (controller userController) PasswordUpdate(c echo.Context) error {
	// fmt.Println("-----------------0")
	oldpassword := c.FormValue("oldpassword")
	email := c.FormValue("email")
	newpassword := c.FormValue("newpassword")
	problem := controller.service.PasswordUpdate(oldpassword, email, newpassword)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "updated succesifully")
}

// Reset godoc
// @Summary Reset a user
// @Description Reset user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/users/reset [post]
func (controller userController) PasswordReset(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	problem := controller.service.PasswordReset(email, password)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "updated succesifully")
}

// GetOne godoc
// @Summary Update a user
// @Description Update a user item
// @Tags users
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/users [put]
func (controller userController) Update(c echo.Context) error {
	user := &model.User{}
	user.FName = c.FormValue("fname")
	user.LName = c.FormValue("lname")
	user.UName = c.FormValue("uname")
	user.Phone = c.FormValue("phone")
	user.Address = c.FormValue("address")
	user.Email = c.FormValue("email")
	image := c.FormValue("image")
	// user.Business = c.FormValue("business")
	code := c.Param("code")
	pic, err2 := c.FormFile("picture")
	if image == "yes" {
		//    fmt.Println(pic.Filename)
		if err2 != nil {
			httperror := httperrors.NewBadRequestError("Invalid picture")
			return c.JSON(httperror.Code(), err2.Error())
		}
		src, err := pic.Open()
		if err != nil {
			httperror := httperrors.NewBadRequestError("the picture is corrupted")
			return c.JSON(httperror.Code(), err.Error())
		}
		defer src.Close()
		// filePath := "./public/imgs/users/"
		filePath := "./src/public/imgs/users/" + user.FName + pic.Filename
		filePath1 := "/imgs/users/" + user.FName + pic.Filename
		// Destination
		dst, err4 := os.Create(filePath)
		if err4 != nil {
			httperror := httperrors.NewBadRequestError("the Directory mess")
			return c.JSON(httperror.Code(), err4.Error())
		}
		defer dst.Close()
		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			if err2 != nil {
				httperror := httperrors.NewBadRequestError("error filling")
				return c.JSON(httperror.Code(), httperror.Message())
			}
		}

		//resize the image and replace the old one
		imagery.Imageryrepository.Imagetype(filePath, filePath, 400, 500)
		user.Picture = filePath1
		err1 := controller.service.Update(code, user)
		if err1 != nil {
			return c.JSON(err1.Code(), err1.Message())
		}
		if _, err = io.Copy(dst, src); err != nil {
			if err2 != nil {
				httperror := httperrors.NewBadRequestError("error filling")
				return c.JSON(httperror.Code(), httperror.Message())
			}
		}
		return c.JSON(http.StatusCreated, "user created succesifully")
	}
	problem := controller.service.Update(code, user)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusCreated, "Updated successifuly")
}

// Delete godoc
// @Summary Delete a user
// @Description Create a new user item
// @Tags users
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} string
// @Failure 400 {object} support.HttpError
// @Router /api/users [delete]
func (controller userController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperrors.NewBadRequestError("Failed to parse the id")
		return c.JSON(httperror.Code(), err.Error())
	}
	success, failure := controller.service.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
