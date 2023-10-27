package app

import (
	"context"
	"emtest/internal/models"
	"emtest/pkg/logger"
	"emtest/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (a *App) StorePerson(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	req := models.HttpPostRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		r.Body.Close()
		return
	}
	defer r.Body.Close()
	a.Log.Debug("decoded person", logger.Val("person", fmt.Sprintf("%v", req)))
	if err := a.PS.StorePerson(ctx, req); err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, "person added")
}

func (a *App) FilterPersons(w http.ResponseWriter, r *http.Request) {
	a.Log.Info("filterPerson request")
	ctx := context.Background()
	filter := models.NewHttpFilterRequest()
	ages := r.URL.Query()["age"]
	if len(ages) > 0 {
		for _, v := range ages {
			age, err := strconv.Atoi(v)
			if err != nil && age < 1 {
				utils.ErrorJSON(w, err, http.StatusBadRequest)
				return
			}
			filter.Age = append(filter.Age, age)
		}
	}

	ageGtstring := r.URL.Query().Get("age_gt")
	if ageGtstring != "" {
		ageGt, err := strconv.Atoi(ageGtstring)
		if err != nil {
			utils.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}
		filter.AgeGT = ageGt
	}
	ageLtstring := r.URL.Query().Get("age_lt")
	if ageLtstring != "" {
		ageLt, err := strconv.Atoi(ageLtstring)
		if err != nil {
			utils.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}
		filter.AgeLT = ageLt
	}
	pageString := r.URL.Query().Get("page")
	if pageString != "" {
		page, err := strconv.Atoi(pageString)
		if err != nil {
			utils.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}
		filter.Page = page
	}
	perpageString := r.URL.Query().Get("per_page")
	if perpageString != "" {
		perpage, err := strconv.Atoi(perpageString)
		if err != nil {
			utils.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}
		filter.PerPage = perpage
	}

	if filter.AgeGT != 0 && filter.AgeLT != 0 && filter.AgeGT > filter.AgeLT {
		utils.ErrorJSON(w, errors.New("wrong values of age_lt and age_gt"), http.StatusBadRequest)
		return
	}

	filter.Name = r.URL.Query()["name"]
	filter.Surname = r.URL.Query()["surname"]
	filter.Patronymic = r.URL.Query()["patronymic"]
	filter.Nationality = r.URL.Query()["nationality"]
	filter.Gender = r.URL.Query().Get("gender")
	filter.OrderBy = r.URL.Query().Get("order_by")
	filterresult, err := a.PS.FilterPersons(ctx, filter)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, filterresult)
}

func (a *App) DeletePersonById(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	if err := a.PS.DeletePersonByID(ctx, id); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, "person deleted")
}

func (a *App) UpdatePersonById(w http.ResponseWriter, r *http.Request) {
	p := models.PersonFull{}
	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	p.ID = id
	ctx := context.Background()
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	newperson, err := a.PS.DB.UpdatePersonById(ctx, p)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusConflict)
	}
	utils.WriteJSON(w, newperson)
}
