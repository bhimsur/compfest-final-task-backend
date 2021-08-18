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

func (a *App) VerifyWithdrawal(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "Withdrawal successfully verified"}

	vars := mux.Vars(r)

	userRole := r.Context().Value("Role").(string)

	withdrawalId, _ := strconv.Atoi(vars["id"])

	withdrawal, err := models.GetWithdrawalById(withdrawalId, a.DB)

	if withdrawal.Status == "verified" {
		resp["message"] = "Withdrawal already verified"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	if userRole != "admin" {
		resp["status"] = false
		resp["message"] = "Only admin can verify withdrawals"
		responses.JSON(w, http.StatusUnauthorized, resp)
		return
	}

	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	verifyWithdrawal := models.Withdrawal{}
	_, err = verifyWithdrawal.VerifyWithdrawal(withdrawalId, a.DB)

	wallet, err := models.GetWalletByUserId(int(withdrawal.UserID), a.DB)
	wallet.Amount += withdrawal.Amount
	wallet, err = wallet.UpdateWallet(a.DB)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	} else {
		responses.JSON(w, http.StatusOK, resp)
		return
	}
}

func (a *App) CreateWithdrawal(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "Withdrawal successfully created"}

	userId := r.Context().Value("UserID").(float64)
	userStatus := r.Context().Value("Status").(string)
	userRole := r.Context().Value("Role").(string)

	vars := mux.Vars(r)
	donation_id, _ := strconv.Atoi(vars["id"])
	donation, err := models.GetDonationProgramById(donation_id, a.DB)

	if userStatus != "verified" || userRole != "fundraiser" {
		resp["status"] = false
		resp["message"] = "Only verified fundraiser can create withdrawal"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	if uint(userId) != donation.UserID {
		resp["status"] = false
		resp["message"] = "You dont have authorities to withdraw this program money"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	withdrawal := &models.Withdrawal{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &withdrawal)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	withdrawedAmount := donation.GetWithdrawedAmount(a.DB)
	availableAmount := donation.GetAvailableAmount(a.DB)

	if withdrawedAmount+withdrawal.Amount > availableAmount {
		resp["status"] = false
		resp["message"] = "You cannot withdraw amount exceed current amount"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	withdrawal.Prepare()

	if err = withdrawal.Validate(); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	withdrawal.UserID = uint(userId)
	withdrawal.DonationProgramID = uint(donation_id)

	withdrawalCreated, err := withdrawal.CreateWithdrawal(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	resp["data"] = withdrawalCreated
	responses.JSON(w, http.StatusCreated, resp)
	return
}

func (a *App) GetUnverifiedWithdrawal(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "Withdrawal successfully retrieved"}

	userRole := r.Context().Value("Role").(string)

	if userRole != "admin" {
		resp["status"] = false
		resp["message"] = "You don't have authorities"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	withdrawals, err := models.GetUnverifiedWithdrawal(a.DB)

	if err != nil {
		resp["data"] = make([]string, 0)
		responses.JSON(w, http.StatusOK, resp)
		return
	}

	resp["data"] = withdrawals
	responses.JSON(w, http.StatusOK, resp)
	return
}
