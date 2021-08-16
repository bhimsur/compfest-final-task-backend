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
	userIdTmp := r.Context().Value("UserID").(float64)
	userId := int(userIdTmp)
	amount := r.Context().Value("Amount").(float64)
	user, err := models.GetUserById(userId)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return 
	}

	topUp := &models.TopUp{
		Amount: amount,
		User: user,
		UserID, userId,
	}

	// NOTES: Bagian sini ngapain sih?
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	if err = json.Unmarshal(body, &topUp), err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	if err = topUp.Validate(); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	topUpCreated, err := topUp.CreateTopUp(a.DB)
	// NOTES: Udah dipindah ke models.topUp
	// walletUpdate := models.Wallet{}
	// walletUpdate.Amount = topUp.Amount + walletUpdate.Amount
	// _, _ = walletUpdate.UpdateWalletFromTopUpByUserId(userId, walletUpdate.Amount, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	resp["data"] = topUpCreated
	responses.JSON(w, http.StatusCreated, resp)
	return
}
