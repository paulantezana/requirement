package controller

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/requirement/config"
	"github.com/paulantezana/requirement/models"
	"github.com/paulantezana/requirement/utilities"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
)

type loginDataResponse struct {
	User  interface{} `json:"user"`
	Token interface{} `json:"token"`
}

func Login(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Hash password
	cc := sha256.Sum256([]byte(user.Password))
	pwd := fmt.Sprintf("%x", cc)

	// Validate user and email
	if db.Where("user_name = ? and password = ?", user.UserName, pwd).First(&user).RecordNotFound() {
		if db.Debug().Where("email = ? and password = ?", user.UserName, pwd).First(&user).RecordNotFound() {
			return c.JSON(http.StatusOK, utilities.Response{
				Success: false,
				Message: "El nombre de usuario o contraseña es incorecta",
			})
		}
	}

	// Check state user
	if !user.State {
		return c.NoContent(http.StatusForbidden)
	}

	// Prepare response data
	user.Password = ""

	// get token key
	token := utilities.GenerateJWT(user)

	// Login success
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Bienvenido al sistema %s", user.UserName),
		Data: loginDataResponse{
			User:  user,
			Token: token,
		},
	})
}

func ForgotSearch(c echo.Context) error {
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Validations
	if err := db.Where("email = ?", user.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: "Tu búsqueda no arrojó ningún resultado. Vuelve a intentarlo con otros datos.",
			Success: false,
		})
	}

	// Generate key validation
	key := (int)(rand.Float32() * 10000000)
	user.Key = fmt.Sprint(key)

	// Update database
	if err := db.Model(&user).Update(user).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// SEND EMAIL get html template
	t, err := template.ParseFiles("templates/email.html")
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// SEND EMAIL new buffer
	buf := new(bytes.Buffer)
	err = t.Execute(buf, user)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// SEND EMAIL
	err = config.SendEmail(user.Email, fmt.Sprint(key)+" es el código de recuperación de tu cuenta en SINST", buf.String())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Response success api service
	user.Key = ""
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    user.ID,
	})
}

func ForgotValidate(c echo.Context) error {
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validations
	if err := db.Where("id = ? AND key = ?", user.ID, user.Key).First(&user).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: "El número que ingresaste no coincide con tu código. Vuelve a intentarlo",
			Success: false,
		})
	}

	user.Password = ""

	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    user.ID,
	})
}

func ForgotChange(c echo.Context) error {
	xUser := models.User{}
	user := models.User{}
	if err := c.Bind(&xUser); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	if err := db.Where("id = ?", xUser.ID).First(&user).Error; err != nil {
		return err
	}

	// Encrypted old password
	tt := sha256.Sum256([]byte(xUser.Password))
	ttt := fmt.Sprintf("%x", tt)

	// Update
	user.Password = ttt
	user.Key = ""
	if err := db.Model(&user).Update(user).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    user.ID,
		Message: fmt.Sprintf("La contraseña del usuario %s se cambio exitosamente", user.UserName),
	})
}

func GetUsers(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Pagination calculate
	if request.CurrentPage == 0 {
		request.CurrentPage = 1
	}
	offset := request.Limit*request.CurrentPage - request.Limit

	// Check the number of matches
	var total uint
	users := make([]models.User, 0)

	// Find users
	if err := db.Where("lower(user_name) LIKE lower(?)", "%"+request.Search+"%").
		Or("lower(dni) LIKE lower(?)", "%"+request.Search+"%").
		Or("lower(last_name) LIKE lower(?)", "%"+request.Search+"%").
		Or("lower(first_name) LIKE lower(?)", "%"+request.Search+"%").
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&users).
		Offset(-1).Limit(-1).Count(&total).
		Error; err != nil {
		return err
	}

	// Type response
	// 0 = all data
	// 1 = minimal data
	if request.Type == 1 {
		customUsers := make([]models.User, 0)
		for _, user := range users {
			customUsers = append(customUsers, models.User{
				ID:       user.ID,
				UserName: user.UserName,
			})
		}
		return c.JSON(http.StatusCreated, utilities.Response{
			Success:     true,
			Data:        customUsers,
			Total:       total,
			CurrentPage: request.CurrentPage,
		})
	}
	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success:     true,
		Data:        users,
		Total:       total,
		CurrentPage: request.CurrentPage,
	})
}

func GetUserByID(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&user, user.ID).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    user,
	})
}

func CreateUser(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// Default empty values
	if len(user.Profile) == 0 {
		user.Profile = "user"
	}
	user.Avatar = "static/logo.png"

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Hash password
	cc := sha256.Sum256([]byte(user.Password))
	pwd := fmt.Sprintf("%x", cc)
	user.Password = pwd

	// Insert user in database
	if err := db.Create(&user).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    user.ID,
		Message: fmt.Sprintf("El usuario %s se registro exitosamente", user.UserName),
	})
}

func UpdateUser(c echo.Context) error {
	// Get data request
	newUser := models.User{}
	if err := c.Bind(&newUser); err != nil {
		return err
	}
	oldUser := models.User{
		ID: newUser.ID,
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	if db.First(&oldUser).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", oldUser.ID),
		})
	}

	// Update user in database
	if err := db.Model(&newUser).Update(newUser).Error; err != nil {
		return err
	}
	if !newUser.State {
		if err := db.Model(newUser).UpdateColumn("state", false).Error; err != nil {
			return err
		}
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    newUser.ID,
	})
}

func DeleteUser(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	if db.First(&user).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", user.ID),
		})
	}

	// Delete user in database
	if err := db.Delete(&user).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    user.ID,
	})
}

func UploadAvatarUser(c echo.Context) error {
	// Read form fields
	idUser := c.FormValue("id")
	user := models.User{}

	fmt.Println(c.FormFile("picture"))

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	if db.First(&user, "id = ?", idUser).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", user.ID),
		})
	}

	// Source
	file, err := c.FormFile("picture")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	ccc := sha256.Sum256([]byte(string(user.ID)))
	name := fmt.Sprintf("%x%s", ccc, filepath.Ext(file.Filename))
	avatarSRC := "static/profiles/" + name
	dst, err := os.Create(avatarSRC)
	if err != nil {
		return err
	}
	defer dst.Close()
	user.Avatar = avatarSRC

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Update database user
	if err := db.Model(&user).Update(user).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    user.ID,
	})
}

func ResetPasswordUser(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	if db.First(&user, "id = ?", user.ID).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", user.ID),
		})
	}

	// Set new password
	cc := sha256.Sum256([]byte(user.DNI + user.UserName))
	pwd := fmt.Sprintf("%x", cc)
	user.Password = pwd

	// Update user in database
	if err := db.Model(&user).Update(user).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("La contraseña del usuario se reseto extosamente. ahora su numevacontraseña es %s", user.DNI+user.UserName),
	})
}

func ChangePasswordUser(c echo.Context) error {
	// Get data request
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	if db.First(&user, "id = ?", user.ID).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", user.ID),
		})
	}

	// Change password
	if len(user.Password) > 0 {
		if len(user.OldPassword) == 0 {
			return c.JSON(http.StatusOK, utilities.Response{
				Success: false,
				Message: "Ingrese la contraseña antigua",
			})
		}

		// Validate old password
		ccc := sha256.Sum256([]byte(user.OldPassword))
		old := fmt.Sprintf("%x", ccc)

		if db.Where("password = ?", old).First(&user).RecordNotFound() {
			return c.JSON(http.StatusOK, utilities.Response{
				Success: false,
				Message: "La contraseña antigua es incorrecta",
			})
		}

		// Set new password
		cc := sha256.Sum256([]byte(user.Password))
		pwd := fmt.Sprintf("%x", cc)
		user.Password = pwd
	}

	// Update user in database
	if err := db.Model(&user).Update(user).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("La contraseña del usuario %s se cambio exitosamente", user.UserName),
	})
}