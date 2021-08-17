package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"restgo/api/models"
	"restgo/api/responses"
)

func (a *App) CreateTopUp(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "success"}
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
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	} else {
		wallet := models.Wallet{}
		_, err := wallet.UpdateWalletFromTopUpByUserId(userId, topUp.Amount, a.DB)
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

func (a *App) GetTopUpHistoryByUserId(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "success"}
	userId := r.Context().Value("UserID").(float64)

	topUps, err := models.GetTopUpHistoryByUserId(int(userId), a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	} else {
		resp["data"] = topUps
		responses.JSON(w, http.StatusOK, resp)
		return
	}
}
