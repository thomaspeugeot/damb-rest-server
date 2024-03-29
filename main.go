package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Unit is Military Element in the Mockup
type Unit struct {
	ID          string    `json:"ID"`
	Coordinates []float64 `json:"coordinates"`
}

type allUnits []Unit

var units = allUnits{
	{
		ID: "1",
		Coordinates: []float64{
			3.0609209,
			50.2362764},
	},
	{
		ID: "2",
		Coordinates: []float64{
			3.1587604,
			50.7170335},
	},
	{
		ID: "3",
		Coordinates: []float64{
			3.4587604,
			50.7170335},
	},
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func createUnit(w http.ResponseWriter, r *http.Request) {
	var newUnit Unit
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the unit title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newUnit)
	units = append(units, newUnit)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newUnit)
}

func getOneUnit(w http.ResponseWriter, r *http.Request) {
	unitID := mux.Vars(r)["id"]

	for _, singleUnit := range units {
		if singleUnit.ID == unitID {
			json.NewEncoder(w).Encode(singleUnit)
		}
	}
}

func getAllUnits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST")
	w.Header().Add("Access-Control-Allow-Methods", "OPTION")
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(units)
}

func updateUnit(w http.ResponseWriter, r *http.Request) {
	unitID := mux.Vars(r)["id"]
	var updatedUnit Unit

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the unit title and description only in order to update")
	}
	json.Unmarshal(reqBody, &updatedUnit)

	for i, singleUnit := range units {
		if singleUnit.ID == unitID {
			/* 			singleUnit.Title = updatedUnit.Title
			   			singleUnit.Description = updatedUnit.Description */
			units = append(units[:i], singleUnit)
			json.NewEncoder(w).Encode(singleUnit)
		}
	}
}

func deleteUnit(w http.ResponseWriter, r *http.Request) {
	unitID := mux.Vars(r)["id"]

	for i, singleUnit := range units {
		if singleUnit.ID == unitID {
			units = append(units[:i], units[i+1:]...)
			fmt.Fprintf(w, "The unit with ID %v has been deleted successfully", unitID)
		}
	}
}

func main() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/unit", createUnit).Methods("POST")
	router.HandleFunc("/units", getAllUnits).Methods("GET")
	router.HandleFunc("/units/{id}", getOneUnit).Methods("GET")
	router.HandleFunc("/units/{id}", updateUnit).Methods("PATCH")
	router.HandleFunc("/units/{id}", deleteUnit).Methods("DELETE")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.Fatal(http.ListenAndServe(":8080", router))
}
