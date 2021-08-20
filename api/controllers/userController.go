package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"restgo/api/models"
	"restgo/api/responses"
	"restgo/api/utils"
	"strconv"

	"github.com/gorilla/mux"
)

func (a *App) UserSignUp(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "Succesfully registered"}

	user := &models.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	usr, _ := user.GetUser(a.DB)
	if usr != nil {
		resp["status"] = false
		resp["message"] = "User already registered, please login"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	user.Prepare()

	err = user.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if user.Role != "fundraiser" {
		user.Status = "verified"
	}

	if user.Role == "fundraiser" {
		user.Status = "pending"
	}

	userCreated, err := user.SaveUser(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	} else {
		wallet := models.Wallet{}
		wallet.UserID = userCreated.ID
		_, err := wallet.InitWallet(a.DB)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		} else {
			resp["user"] = userCreated
			responses.JSON(w, http.StatusCreated, resp)
			return
		}
	}
}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "logged in"}

	user := &models.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user.Prepare()

	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	usr, err := user.GetUser(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if usr == nil {
		resp["status"] = false
		resp["message"] = "login failed, please sign up"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	err = models.CheckPasswordHash(user.Password, usr.Password)
	if err != nil {
		resp["status"] = false
		resp["message"] = "login failed, wrong password"
		responses.JSON(w, http.StatusForbidden, resp)
		return
	}

	token, err := utils.EncodeAuthToken(usr.ID, usr.Email, string(usr.Role), string(usr.Status))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["token"] = token
	responses.JSON(w, http.StatusOK, resp)
}

func (a *App) VerifyFundraiser(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "fundraiser successfully verified"}

	vars := mux.Vars(r)

	userRole := r.Context().Value("Role").(string)

	userId, _ := strconv.Atoi(vars["id"])

	user, err := models.GetUserById(userId, a.DB)

	if userRole != "admin" {
		resp["status"] = false
		resp["message"] = "Only admin can verify fundraiser"
		responses.JSON(w, http.StatusUnauthorized, resp)
		return
	} else if user.Status == "verified" {
		resp["message"] = "Fundraiser already verified"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	verifyUser := models.User{}

	_, err = verifyUser.VerifyFundraiser(userId, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, resp)
}

func (a *App) GetUserById(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "success"}
	user := r.Context().Value("UserID").(float64)
	userID := int(user)
	userDetail, err := models.GetUserById(userID, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	} else {
		resp["data"] = userDetail
		responses.JSON(w, http.StatusOK, resp)
	}
}

func (a *App) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "user updated successfully"}
	user := r.Context().Value("UserID").(float64)
	userID := int(user)
	_, err := models.GetUserById(userID, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	userUpdate := models.User{}
	if err = json.Unmarshal(body, &userUpdate); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userUpdate.Prepare()
	_, err = userUpdate.UpdateUser(userID, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	} else {
		responses.JSON(w, http.StatusOK, resp)
		return
	}
}

func (a *App) GetUnverifiedUser(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "User successfully retrieved"}

	userRole := r.Context().Value("Role").(string)

	if userRole != "admin" {
		resp["status"] = false
		resp["message"] = "You don't have authorities"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	u, err := models.GetUnverifiedUser(a.DB)

	if err != nil {
		resp["data"] = make([]string, 0)
		responses.JSON(w, http.StatusOK, resp)
		return
	}

	resp["data"] = u
	responses.JSON(w, http.StatusOK, resp)
	return
}

type PasswordField struct {
	CurrentPassword string `json:"currentPassword"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func (a *App) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "Password Changed Successfully"}
	user := r.Context().Value("UserID").(float64)
	userID := int(user)

	current, _ := models.GetUserByID(userID, a.DB)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	passwordBody := &PasswordField{}

	err = json.Unmarshal(body, &passwordBody)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = models.CheckPasswordHash(passwordBody.CurrentPassword, current.Password)

	if err != nil {
		resp["status"] = false
		resp["message"] = "Wrong Password"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	if passwordBody.ConfirmPassword != passwordBody.Password {
		resp["status"] = false
		resp["message"] = "Password and Confirm Password must same"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	userUpdate := models.User{Password: passwordBody.Password}
	_, err = userUpdate.ChangePassword(userID, a.DB)

	if err != nil {
		resp["message"] = "Failed To Change Password"
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, resp)
}
