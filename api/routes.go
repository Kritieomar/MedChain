package api

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/get-patient-profile/{patient_id}", GetPatientProfile).Methods("GET")

	router.HandleFunc("/api/v1/get-doctor-profile/{doctor_id}", GetDoctorProfile).Methods("GET")
	router.HandleFunc("/api/v1/get-patient-profile/{patient_id}", GetPatientProfile).Methods("GET")

	router.HandleFunc("/api/v1/add-record", AddRecord).Methods("POST")
	router.HandleFunc("/api/v1/get-record/{cid}", GetFileFromIPFS).Methods("GET")
	router.HandleFunc("/api/v1/blockchain", GetBlockchain).Methods("GET")
	router.HandleFunc("/api/v1/grant-permission", GrantPermission).Methods("POST")
	router.HandleFunc("/api/v1/accept-record", AcceptRecord).Methods("POST")
	router.HandleFunc("/api/v1/reject-record", RejectRecord).Methods("POST")

	return router
}
