package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ttacon/libphonenumber"
)

// PhoneNumber defines JSON structure of response
type PhoneNumber struct {
	CountryCode int32                `json:"country_code"`
	PhoneNumber uint64               `json:"phone_number"`
	Extension   string               `json:"extension"`
	Formatted   PhoneNumberFormatted `json:"formatted"`
}

// PhoneNumberFormatted defines JSON structure of formatted key of response
type PhoneNumberFormatted struct {
	National      string `json:"national"`
	International string `json:"international"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/{number}", Parse)

	log.Fatal(http.ListenAndServe(":8080", router))
}

// Index welcomes anyone to provide instructions
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

// Parse Parses given phone number and return fragments as JSON
func Parse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	rawNumber := vars["number"]
	parsedNum, parseErr := libphonenumber.Parse(fmt.Sprintf("+%s", rawNumber), "US")
	if parseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(parseErr.Error())
		return
	}

	number := PhoneNumber{
		CountryCode: parsedNum.GetCountryCode(),
		PhoneNumber: parsedNum.GetNationalNumber(),
		Extension:   parsedNum.GetExtension(),
		Formatted: PhoneNumberFormatted{
			National:      libphonenumber.Format(parsedNum, libphonenumber.NATIONAL),
			International: libphonenumber.Format(parsedNum, libphonenumber.INTERNATIONAL),
		},
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(number); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
	}
}
