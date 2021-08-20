package controllers

import (
	"net/http"
	"restgo/api/models"
	"restgo/api/responses"
	"sort"
	"time"
)

type Unverified struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Type      string    `json:"type"`
}

func (a *App) UnverifiedList(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": true, "message": "Successfully Retrieve Data"}

	userRole := r.Context().Value("Role").(string)

	if userRole != "admin" {
		resp["status"] = false
		resp["message"] = "Only admin can retrieve data"
		responses.JSON(w, http.StatusUnauthorized, resp)
		return
	}

	unverifieds := []Unverified{}

	users, err := models.GetUnverifiedUser(a.DB)
	programs, err := models.GetUnverifiedDonationProgram(a.DB)
	withdraws, err := models.GetUnverifiedWithdrawal(a.DB)

	for _, v := range *users {
		unverified := Unverified{
			ID:        v.ID,
			Name:      v.Name,
			CreatedAt: v.CreatedAt,
			Type:      "fundraiser",
		}
		unverifieds = append(unverifieds, unverified)
	}

	for _, v := range *programs {
		unverified := Unverified{
			ID:        v.ID,
			Name:      models.GetName(v.UserID, a.DB),
			CreatedAt: v.CreatedAt,
			Type:      "program",
		}
		unverifieds = append(unverifieds, unverified)
	}

	for _, v := range *withdraws {
		unverified := Unverified{
			ID:        v.ID,
			Name:      models.GetName(v.UserID, a.DB),
			CreatedAt: v.CreatedAt,
			Type:      "withdraw",
		}
		unverifieds = append(unverifieds, unverified)
	}

	sort.Slice(unverifieds, func(i, j int) bool {
		return unverifieds[i].CreatedAt.Before(unverifieds[j].CreatedAt)
	})

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	} else {
		resp["data"] = unverifieds
		responses.JSON(w, http.StatusCreated, resp)
		return
	}
}
