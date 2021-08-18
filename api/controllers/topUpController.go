package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"restgo/api/models"
	"restgo/api/responses"
)

func (a *App) CreateTopUp(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "successs"}
	user := r.Context().Value("UserID").(float64)
	userId := int(user)

	topUp := &models.TopUp{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	err = json.Unmarshal(body, &topUp)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if err = topUp.Validate(); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	topUp.UserID = uint(userId)
	topUpCreated, err := topUp.CreateTopUp(a.DB)

	walletUpdate, err := models.GetWalletByUserId(userId, a.DB)
	walletUpdate.Amount += topUp.Amount

	_, err = walletUpdate.UpdateWallet(a.DB)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	} else {
		wallet, err := models.GetWalletByUserId(userId, a.DB)
		wallet.Amount = topUp.Amount
		_, err = wallet.UpdateWallet(a.DB)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		} else {
			resp["data"] = topUpCreated
			responses.JSON(w, http.StatusCreated, resp)
			return
		}
	}
}

func (a *App) TopupHistory(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "successs"}
	user := r.Context().Value("UserID").(float64)
	userId := int(user)

	topUp := &models.TopUp{}

	topups, err := topUp.TopupHistory(userId, a.DB)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	} else {
		resp["data"] = topups
		responses.JSON(w, http.StatusCreated, resp)
		return
	}
}
