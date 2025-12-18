package longreminder

import (
	"encoding/json"
	"log"
	"murky_api/internal/constants"
	"murky_api/internal/context"
	"murky_api/internal/model"
	"murky_api/internal/routing"
	"net/http"
	"time"
)

func Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := context.GetDb(r)
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}
		user := context.GetCurrentUser(r)

		var req CreateRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			routing.WriteInvalidJsonResponse(w)
			return
		}

		validationResult := req.Validate(db)
		if validationResult != nil {
			routing.WriteValidationErrorResponse(w, *validationResult)
			return
		}

		reminder := model.LongReminder{
			Id:           req.Id,
			Title:        req.Title,
			IntervalDays: req.IntervalDays,
		}

		reminder, err = model.CreateLongReminder(db, reminder, user.Id)
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}

		routing.WriteJsonResponse(w, http.StatusNoContent, nil)
	}
}

func GetList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := context.GetDb(r)
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}
		user := context.GetCurrentUser(r)

		reminders, err := model.GetLongRemindersByUserId(db, user.Id)
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}

		resp := ListResponse{
			Data: make([]ListItem, 0),
		}
		for _, reminder := range reminders {
			listItem := ListItem{
				Id:           reminder.Id,
				Title:        reminder.Title,
				IntervalDays: reminder.IntervalDays,
				CreatedAt:    reminder.CreatedAt.Format(constants.ApiDateTimeFormat),
				IsEnabled:    reminder.IsEnabled,
			}

			if reminder.UpdatedAt != nil {
				s := reminder.UpdatedAt.Format(constants.ApiDateTimeFormat)
				listItem.UpdatedAt = &s
			}

			if reminder.LastRemindedAt != nil {
				s := reminder.LastRemindedAt.Format(constants.ApiDateTimeFormat)
				listItem.LastRemindedAt = &s
			}

			if reminder.MarkedDoneAt != nil {
				s := reminder.MarkedDoneAt.Format(constants.ApiDateTimeFormat)
				listItem.MarkedDoneAt = &s
			}

			resp.Data = append(resp.Data, listItem)
		}

		routing.WriteJsonResponse(w, http.StatusOK, resp)
	}
}

func Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := context.GetDb(r)
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}
		user := context.GetCurrentUser(r)

		rowsDeleted, err := model.DeleteLongReminderByUserIdAndId(db, user.Id, r.PathValue("id"))
		if err != nil || rowsDeleted == 0 {
			routing.WriteNotFoundResponse(w)
			return
		}

		routing.WriteJsonResponse(w, http.StatusNoContent, nil)
	}
}

func Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := context.GetDb(r)
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}
		user := context.GetCurrentUser(r)
		id := r.PathValue("id")

		reminder, err := model.GetLongReminderByUserIdAndId(db, user.Id, id)
		if err != nil {
			routing.WriteNotFoundResponse(w)
			return
		}

		var req UpdateRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			routing.WriteInvalidJsonResponse(w)
			return
		}

		validationResult := req.Validate(db)
		if validationResult != nil {
			routing.WriteValidationErrorResponse(w, *validationResult)
			return
		}

		if req.Title != nil {
			reminder.Title = *req.Title
		}
		if req.IntervalDays != nil {
			reminder.IntervalDays = *req.IntervalDays
		}
		if req.MarkedDoneAt != nil {
			markedDoneAt, err := time.Parse(constants.ApiDateTimeFormat, *req.MarkedDoneAt)
			if err != nil {
				log.Println(err)
				routing.WriteInternalServerErrorResponse(w)
			}
			reminder.MarkedDoneAt = &markedDoneAt
		}
		if req.IsEnabled != nil {
			reminder.IsEnabled = *req.IsEnabled
		}

		err = model.UpdateLongReminder(db, id, reminder)
		if err != nil {
			log.Println(err)
			routing.WriteInternalServerErrorResponse(w)
			return
		}

		routing.WriteJsonResponse(w, http.StatusNoContent, nil)
	}
}
