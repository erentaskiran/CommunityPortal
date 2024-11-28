package api

import (
	"net/http"

	"github.com/erentaskiran/project123123123/pkg/utils"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, "")
}
