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
		if len(*donationHistories) == 0 {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		} else {
			resp["data"] = donationHistories
			responses.JSON(w, http.StatusOK, resp)
			return
		}
	}
}

func (a *App) DonateNow(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "success"}
	userId := r.Context().Value("UserID").(float64)
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	donationProgram, err := models.GetDonationProgramById(id, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	} else {
		if donationProgram.Status != "verified" {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
	}

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

	userWallet, err := models.GetWalletByUserId(int(userId), a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	} else {
		if userWallet.Amount <= 0 || userWallet.Amount < donation.Amount {
			resp["status"] = false
			resp["message"] = "insufficient balance, please top up"
			responses.JSON(w, http.StatusBadRequest, resp)
			return
		}
	}

	donation.UserID = uint(userId)
	donationCreated, err := donation.SaveDonation(a.DB)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	} else {
		program := models.DonationProgram{}
		_, err := program.UpdateDonationProgramAmountById(id, donation.Amount, a.DB)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		} else {
			wallet := models.Wallet{}
			_, err := wallet.UpdateWalletById(int(userWallet.UserID), donation.Amount, a.DB)
			if err != nil {
				responses.ERROR(w, http.StatusBadRequest, err)
				return
			} else {
				resp["data"] = donationCreated
				responses.JSON(w, http.StatusCreated, resp)
				return
			}
		}
	}
}
