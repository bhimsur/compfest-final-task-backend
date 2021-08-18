package controllers

import (
	"net/http"
	"restgo/api/models"
	"restgo/api/responses"
)

func (a *App) GetWalletByUserId(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "true", "message": "success"}
	user := r.Context().Value("UserID").(float64)
	userId := int(user)

	wallet, err := models.GetWalletByUserId(userId, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	} else {
		resp["data"] = wallet.Amount
		responses.JSON(w, http.StatusOK, resp)
		return
	}
}
