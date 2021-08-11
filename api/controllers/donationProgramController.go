package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"restgo/api/models"
	"restgo/api/responses"
	"strconv"

	"github.com/gorilla/mux"
)

func (a *App) CreateDonationProgram(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "Donation program successfully created"}

	userId := r.Context().Value("UserID").(float64)
	userStatus := r.Context().Value("Status").(string)

	if userStatus != "verified" {
		resp["status"] = false
		resp["message"] = "Only verified fundraiser can create donation program"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	donationProgram := &models.DonationProgram{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &donationProgram)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	donationProgram.Prepare()

	if err = donationProgram.Validate(); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	donationProgram.UserID = uint(userId)

	donationProgramCreated, err := donationProgram.Save(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	resp["data"] = donationProgramCreated
	responses.JSON(w, http.StatusCreated, resp)
}

func (a *App) GetDonationPrograms(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "success"}
	donationPrograms, err := models.GetDonationPrograms(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	resp["data"] = donationPrograms
	responses.JSON(w, http.StatusOK, resp)
}

func (a *App) GetDonationProgramById(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "success"}
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	donationProgram, err := models.GetDonationProgramById(id, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	resp["data"] = donationProgram
	responses.JSON(w, http.StatusOK, resp)
}

func (a *App) GetDonationProgramByFundraiser(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "success"}
	userRole := r.Context().Value("Role").(string)
	userId := r.Context().Value("UserID").(float64)
	userID := int(userId)

	if userRole != "fundraiser" {
		resp["status"] = false
		resp["message"] = "Unauthorized user"
		responses.JSON(w, http.StatusUnauthorized, resp)
		return
	}

	donationPrograms, err := models.GetDonationProgramByFundraiser(userID, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	} else {
		resp["data"] = donationPrograms
		responses.JSON(w, http.StatusOK, resp)
		return
	}
}

func (a *App) VerifyDonationProgram(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "donation program successfully verified"}

	vars := mux.Vars(r)

	userRole := r.Context().Value("Role").(string)

	donationProgramId, _ := strconv.Atoi(vars["id"])

	donationProgram, err := models.GetDonationProgramById(donationProgramId, a.DB)

	if donationProgram.Status == "verified" {
		resp["message"] = "Donation program already verified"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	if userRole != "admin" {
		resp["status"] = false
		resp["message"] = "Only admin can verify fundraiser"
		responses.JSON(w, http.StatusUnauthorized, resp)
		return
	}

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	verifyDonationProgram := models.DonationProgram{}

	_, err = verifyDonationProgram.VerifyDonationProgram(donationProgramId, a.DB)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	} else {
		responses.JSON(w, http.StatusOK, resp)
	}
}
