package utilities

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kwesidev/authserver/internal/models"
)

func GetJsonInput(input interface{}, req *http.Request) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	err = json.Unmarshal(body, &input)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// JSONError send json error messages
func JSONError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	generalErrorResponse := models.GeneralErrorResponse{}
	generalErrorResponse.Status = code
	generalErrorResponse.Success = false
	generalErrorResponse.ErrorMessage = error
	json.NewEncoder(w).Encode(generalErrorResponse)
}

// JSONResponse Send Json response messages
func JSONResponse(w http.ResponseWriter, result interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
