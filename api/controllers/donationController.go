package controllers

import (
	"net/http"
	"restgo/api/models"
	"restgo/api/responses"
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

func (a *App) DonateNow(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "success"}

	user := r.Context().Value("UserID").(float64)
	userID := int(user)
	amount := r.Context().Value("Amount").(float64)
	donationProgram := r.Context().Value("DonationProgram")	

	newDonation := &models.Donation{
		Amount: amount,
		UserID: userID,
		DonationProgram: donationProgram,
		// TODO: Figure out how to get donationProgramID
		DonationProgramID: donationProgram.ID,
	}

	donationCreated, err := newDonation.SaveDonation(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	resp["data"] = donationCreated
	responses.JSON(w, http.StatusCreated, resp)
	return
}
