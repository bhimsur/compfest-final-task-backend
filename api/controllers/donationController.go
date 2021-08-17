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

func (a *App) GetDonationHistoryFromUser(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "success"}
	user := r.Context().Value("UserID").(float64)
	userID := int(user)

	donationHistories, err := models.GetDonationHistoryFromUser(userID, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	} else {
		resp["data"] = donationHistories
		responses.JSON(w, http.StatusOK, resp)
		return
	}
}

func (a *App) DonateToProgram(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "donation program successfully verified"}

	vars := mux.Vars(r)

	user := r.Context().Value("UserID").(float64)
	userId := int(user)

	donationProgramId, _ := strconv.Atoi(vars["id"])

	donation := &models.Donation{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	err = json.Unmarshal(body, &donation)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if err = donation.Validate(); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	wallet, err := models.GetWalletByUserId(userId, a.DB)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if wallet.Amount < donation.Amount {
		resp["status"] = "fail"
		resp["message"] = "You don't have enough money"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	donation.UserID = uint(userId)
	donation.DonationProgramID = uint(donationProgramId)

	_, err = donation.SaveDonation(a.DB)

	wallet.Amount -= donation.Amount
	_, err = wallet.UpdateWallet(a.DB)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["data"] = donation
	responses.JSON(w, http.StatusOK, resp)
	return
}
