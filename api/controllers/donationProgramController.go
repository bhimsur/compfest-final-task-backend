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
	"github.com/jinzhu/gorm"
)

type _DonationProgram struct {
	ID             uint          `json:"id"`
	Title          string        `json:"title"`
	Detail         string        `json:"detail"`
	Amount         float64       `json:"amount"`
	Deadline       string        `json:"deadline"`
	UserID         uint          `json:"user_id"`
	FundraiserName string        `json:"fundraiser_name"`
	Status         models.Status `json:"status"`
	Withdrawn      float64       `json:"withdrawn"`
	Collected      float64       `json:"collected"`
}

type _ProgramInput struct {
	Title    string  `json:"title"`
	Detail   string  `json:"detail"`
	Amount   float64 `json:"amount"`
	Deadline string  `json:"deadline"`
}

func (a *App) CreateDonationProgram(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "Donation program successfully created"}

	userId := r.Context().Value("UserID").(float64)
	userStatus := r.Context().Value("Status").(string)
	userRole := r.Context().Value("Role").(string)

	if userStatus != "verified" || userRole != "fundraiser" {
		resp["status"] = false
		resp["message"] = "Only fundraiser can create donation program"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	} else {
		if userStatus != "verified" {
			resp["message"] = "Only verified fundraiser can create donation program"
			responses.JSON(w, http.StatusBadRequest, resp)
			return
		}
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	_input := &_ProgramInput{}
	err = json.Unmarshal(body, &_input)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	donationProgram := &models.DonationProgram{
		UserID:   uint(userId),
		Deadline: utils.DateToRFC(_input.Deadline),
		Title:    _input.Title,
		Amount:   _input.Amount,
		Detail:   _input.Detail,
	}

	donationProgram.Prepare()

	if err = donationProgram.Validate(); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

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
	resp["data"] = convertToDetail(donationPrograms, a.DB)
	responses.JSON(w, http.StatusOK, resp)
}

type _DonationProgramDetail struct {
	ID             uint          `json:"id"`
	Title          string        `json:"title"`
	Detail         string        `json:"detail"`
	Amount         float64       `json:"amount"`
	Deadline       string        `json:"deadline"`
	UserID         uint          `json:"user_id"`
	FundraiserName string        `json:"fundraiser_name"`
	Status         models.Status `json:"status"`
	Withdrawn      float64       `json:"withdrawn"`
	Collected      float64       `json:"collected"`
	Donation       []_Donation   `json:"Donation"`
}

type _Donation struct {
	Amount float64 `json:"amount"`
	Name   string  `json:"name"`
}

func (a *App) GetDonationProgramById(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "success"}
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	donationProgram, err := models.GetDonationProgramById(id, a.DB)

	detail := _DonationProgramDetail{
		ID:             donationProgram.ID,
		Title:          donationProgram.Title,
		Detail:         donationProgram.Detail,
		Amount:         donationProgram.Amount,
		FundraiserName: models.GetName(donationProgram.UserID, a.DB),
		Deadline:       utils.RFCToDate(donationProgram.Deadline),
		UserID:         donationProgram.UserID,
		Status:         donationProgram.Status,
		Withdrawn:      donationProgram.GetWithdrawedAmount(a.DB),
		Collected:      donationProgram.GetAvailableAmount(a.DB),
		Donation:       *convertDonation(&donationProgram.Donation, a.DB),
	}

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	resp["data"] = detail
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

func (a *App) GetUnverifiedDonationProgram(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "Donation Program successfully retrieved"}

	userRole := r.Context().Value("Role").(string)

	if userRole != "admin" {
		resp["status"] = false
		resp["message"] = "You don't have authorities"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	dp, err := models.GetUnverifiedDonationProgram(a.DB)

	if err != nil {
		resp["data"] = make([]string, 0)
		responses.JSON(w, http.StatusOK, resp)
		return
	}

	resp["data"] = dp
	responses.JSON(w, http.StatusOK, resp)
	return
}

func (a *App) SearchDonationProgram(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "Donation Program successfully retrieved"}

	keyword := r.URL.Query().Get("keyword")

	donationPrograms, err := models.SearchDonationProgram(keyword, a.DB)

	donationDetails := convertToDetail(donationPrograms, a.DB)
	if err != nil {
		resp["data"] = make([]string, 0)
		responses.JSON(w, http.StatusOK, resp)
		return
	}

	resp["data"] = donationDetails
	responses.JSON(w, http.StatusOK, resp)
	return
}

func convertDonation(donations *[]models.Donation, db *gorm.DB) *[]_Donation {
	_donations := []_Donation{}

	for _, d := range *donations {
		dn := _Donation{
			Amount: d.Amount,
			Name:   models.GetName(d.UserID, db),
		}
		_donations = append(_donations, dn)
	}

	return &_donations
}

func convertToDetail(donationPrograms *[]models.DonationProgram, db *gorm.DB) *[]_DonationProgram {

	donationDetails := []_DonationProgram{}
	for _, elem := range *donationPrograms {
		detail := _DonationProgram{
			ID:             elem.ID,
			Title:          elem.Title,
			Detail:         elem.Detail,
			Amount:         elem.Amount,
			FundraiserName: models.GetName(elem.UserID, db),
			Deadline:       utils.RFCToDate(elem.Deadline),
			UserID:         elem.UserID,
			Status:         elem.Status,
			Withdrawn:      elem.GetWithdrawedAmount(db),
			Collected:      elem.GetAvailableAmount(db),
		}
		donationDetails = append(donationDetails, detail)
	}

	return &donationDetails
}
