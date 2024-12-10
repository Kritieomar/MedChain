package api

import (
	"encoding/json"
	"log"
	"my-blockchain/blockchain"
	"net/http"

	"github.com/gorilla/mux"
)

func GetDoctorProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	doctorID := vars["doctor_id"]

	profile, exists := blockchain.DoctorProfiles[doctorID]
	if !exists {
		http.Error(w, "Doctor profile not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

func GetPatientProfile(w http.ResponseWriter, r *http.Request) {
	// Extract patient ID from the request
	vars := mux.Vars(r)
	patientID := vars["patient_id"]

	// Retrieve the patient profile from the blockchain's PatientProfiles map
	profile, exists := blockchain.PatientProfiles[patientID]
	if !exists {
		http.Error(w, "Patient profile not found", http.StatusNotFound)
		return
	}

	// Return the profile as a JSON response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

// AddRecord handles POST requests to add medical records
func AddRecord(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Data                  string   `json:"data"`
		DoctorID              string   `json:"doctor_id"`
		RecordID              string   `json:"record_id"`
		DoctorsWithPermission []string `json:"doctors_with_permission"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Upload to IPFS
	ipfsHash, err := blockchain.AddFileToIPFS(request.Data)
	if err != nil {
		http.Error(w, "Failed to upload to IPFS", http.StatusInternalServerError)
		return
	}

	// Add block to the blockchain with permissions
	err = blockchain.BC.AddBlockWithMetadata(request.Data, ipfsHash, request.DoctorID, request.RecordID, request.DoctorsWithPermission, []string{}, "")
	if err != nil {
		log.Printf("Error adding block to blockchain: %v", err)
		http.Error(w, "Failed to add block to blockchain", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message":  "Record successfully added by doctor",
		"cid":      ipfsHash,
		"recordID": request.RecordID,
	}
	json.NewEncoder(w).Encode(response)
}

// AcceptRecord handles POST requests for the patient to accept the record
func AcceptRecord(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RecordID  string `json:"record_id"`
		PatientID string `json:"patient_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Search for the record by RecordID
	for _, block := range blockchain.BC.Blocks {
		if block.TransactionID == request.RecordID {
			// Retrieve the patient's profile
			profile, exists := block.PatientProfiles[request.PatientID]
			if !exists {
				// If the profile does not exist, create a new one
				profile = blockchain.PatientProfile{
					PatientID:          request.PatientID,
					AcceptedRecords:    []string{},
					RejectedRecords:    []string{},
					InteractionHistory: []string{},
				}
			}

			// Modify the profile
			profile.AcceptedRecords = append(profile.AcceptedRecords, block.IPFSHash)
			profile.InteractionHistory = append(profile.InteractionHistory, "Record accepted by patient "+request.PatientID)

			// Save the updated profile back to the map
			block.PatientProfiles[request.PatientID] = profile

			// Respond with success
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message":   "Record accepted successfully",
				"recordID":  request.RecordID,
				"patientID": request.PatientID,
			})
			return
		}
	}

	// If record not found, return an error
	http.Error(w, "Record not found", http.StatusNotFound)
}

func RejectRecord(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RecordID  string `json:"record_id"`
		PatientID string `json:"patient_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Search for the record by RecordID
	for _, block := range blockchain.BC.Blocks {
		if block.TransactionID == request.RecordID {
			// Retrieve the patient's profile
			profile, exists := block.PatientProfiles[request.PatientID]
			if !exists {
				// If the profile does not exist, create a new one
				profile = blockchain.PatientProfile{
					PatientID:          request.PatientID,
					AcceptedRecords:    []string{},
					RejectedRecords:    []string{},
					InteractionHistory: []string{},
				}
			}

			// Modify the profile
			profile.RejectedRecords = append(profile.RejectedRecords, block.IPFSHash)
			profile.InteractionHistory = append(profile.InteractionHistory, "Record rejected by patient "+request.PatientID)

			// Save the updated profile back to the map
			block.PatientProfiles[request.PatientID] = profile

			// Respond with success
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message":   "Record rejected successfully",
				"recordID":  request.RecordID,
				"patientID": request.PatientID,
			})
			return
		}
	}

	// If record not found, return an error
	http.Error(w, "Record not found", http.StatusNotFound)
}

// GetFileFromIPFS handles GET requests to retrieve a file from IPFS using its CID
func GetFileFromIPFS(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid := vars["cid"]

	if cid == "" {
		http.Error(w, "Missing CID parameter", http.StatusBadRequest)
		return
	}

	// Retrieve data from IPFS
	data, err := blockchain.GetFileFromIPFS(cid)
	if err != nil {
		log.Printf("Error retrieving file from IPFS for CID %s: %v", cid, err)
		http.Error(w, "Failed to retrieve file from IPFS", http.StatusInternalServerError)
		return
	}

	// Return data and CID as a response
	response := map[string]string{"cid": cid, "data": data}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetBlockchain handles GET requests to view the current state of the blockchain
func GetBlockchain(w http.ResponseWriter, r *http.Request) {
	// Return the blockchain
	json.NewEncoder(w).Encode(blockchain.BC)
}

// GrantPermission handles POST requests to grant a doctor access to a medical record
func GrantPermission(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RecordID string `json:"record_id"`
		DoctorID string `json:"doctor_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Find the block by RecordID and grant permission
	for _, block := range blockchain.BC.Blocks {
		if block.TransactionID == request.RecordID {
			block.AddDoctorPermission(request.DoctorID)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message":  "Permission granted successfully",
				"recordID": request.RecordID,
				"doctorID": request.DoctorID,
			})
			return
		}
	}

	http.Error(w, "Record not found", http.StatusNotFound)
}

// LogInteraction handles POST requests to log interactions for a record
func LogInteraction(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RecordID    string `json:"record_id"`
		Interaction string `json:"interaction"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Find the block by RecordID and log the interaction
	for _, block := range blockchain.BC.Blocks {
		if block.TransactionID == request.RecordID {
			block.LogInteraction(request.Interaction)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message":     "Interaction logged successfully",
				"recordID":    request.RecordID,
				"interaction": request.Interaction,
			})
			return
		}
	}

	http.Error(w, "Record not found", http.StatusNotFound)
}
