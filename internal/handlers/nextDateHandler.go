package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/azoma13/go_final_project/configs"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("incorrect date format")
	}
	parseDate, err := time.Parse(configs.DateLayout, date)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err)
	}

	repeats := strings.Split(repeat, " ")
	lengthRepeats := len(repeats)
	zeroRepeat := repeats[0]
	switch zeroRepeat {
	case "d":

		if lengthRepeats != 2 {
			return "", fmt.Errorf("invailed repetition format")
		}

		dayToRepeat, err := strconv.Atoi(repeats[1])
		if err != nil {
			return "", err
		}

		if dayToRepeat > 400 {
			return "", fmt.Errorf("maximum allowed interval of 400 days has been exceeded")
		}

		newDate := parseDate.AddDate(0, 0, dayToRepeat)
		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, dayToRepeat)
		}

		return newDate.Format(configs.DateLayout), nil
	case "y":

		if lengthRepeats != 1 {
			return "", fmt.Errorf("invailed repetition format")
		}

		newDate := parseDate.AddDate(1, 0, 0)
		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}
		return newDate.Format(configs.DateLayout), nil
	default:
		return "invalid character", nil
	}
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	nowTime, err := time.Parse(configs.DateLayout, now)
	if err != nil {
		sendErrorResponse(w, "date is in an incorrect format", http.StatusBadRequest)
		return
	}
	nextDate, err := NextDate(nowTime, date, repeat)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(nextDate))
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
}
