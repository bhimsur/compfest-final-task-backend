package controllers

import (
	"net/http"
	"restgo/api/models"
	"restgo/api/responses"
	"sort"
	"time"
)

type Notification struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
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

	notifications := []Notification{}

	users, err := models.GetUnverifiedUser(a.DB)
	programs, err := models.GetUnverifiedDonationProgram(a.DB)
	withdraws, err := models.GetUnverifiedWithdrawal(a.DB)

	for _, v := range *users {
		notification := Notification{
			ID:        v.ID,
			Name:      v.Name,
			CreatedAt: v.CreatedAt,
			Type:      "fundraiser",
		}
		notifications = append(notifications, notification)
	}

	for _, v := range *programs {
		notification := Notification{
			ID:        v.ID,
			Name:      v.Name,
			CreatedAt: v.CreatedAt,
			Title:     v.Title,
			Type:      "program",
		}
		notifications = append(notifications, notification)
	}

	for _, v := range *withdraws {
		notification := Notification{
			ID:        v.ID,
			Name:      v.Name,
			CreatedAt: v.CreatedAt,
			Title:     v.Title,
			Type:      "withdraw",
		}
		notifications = append(notifications, notification)
	}

	sort.Slice(notifications, func(i, j int) bool {
		return notifications[i].CreatedAt.Before(notifications[j].CreatedAt)
	})

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	} else {
		resp["data"] = notifications
		responses.JSON(w, http.StatusCreated, resp)
		return
	}
}
