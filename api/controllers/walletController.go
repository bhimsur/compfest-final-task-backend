package controllers

import (
	"net/http"
	"restgo/api/models"
	"restgo/api/responses"
	"sort"
	"time"
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

type _WalletHistory struct {
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
	Action string    `json:"action"`
}

func (a *App) WalletHistory(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "successs"}
	user := r.Context().Value("UserID").(float64)
	userId := int(user)

	_walletHistory := []_WalletHistory{}

	topups, err := models.TopupHistory(userId, a.DB)
	donates, err := models.GetDonationHistoryFromUser(userId, a.DB)

	for _, v := range *topups {
		wallet := _WalletHistory{
			Date:   v.CreatedAt,
			Amount: v.Amount,
			Action: "topup",
		}
		_walletHistory = append(_walletHistory, wallet)
	}

	for _, v := range *donates {
		wallet := _WalletHistory{
			Date:   v.CreatedAt,
			Amount: v.Amount,
			Action: "donate",
		}
		_walletHistory = append(_walletHistory, wallet)
	}

	sort.Slice(_walletHistory, func(i, j int) bool {
		return _walletHistory[i].Date.After(_walletHistory[j].Date)
	})

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	} else {
		resp["data"] = _walletHistory
		responses.JSON(w, http.StatusCreated, resp)
		return
	}
}
