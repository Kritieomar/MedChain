package blockchain

type Block struct {
	Timestamp             int64                     `json:"timestamp"`
	PreviousHash          []byte                    `json:"prevBlockHash"`
	MyBlockHash           []byte                    `json:"hash"`
	AllData               []byte                    `json:"medicalData"`
	IPFSHash              string                    `json:"ipfsHash"`
	TransactionID         string                    `json:"recordId"`
	Owner                 string                    `json:"owner"`
	DoctorsWithPermission []string                  `json:"doctorsWithPermission"`
	PatientProfiles       map[string]PatientProfile `json:"patientProfiles"`
	Interactions          []string                  `json:"interactions"`
}

type PatientProfile struct {
	PatientID          string   `json:"patientId"`
	AcceptedRecords    []string `json:"acceptedRecords"`
	RejectedRecords    []string `json:"rejectedRecords"`
	InteractionHistory []string `json:"interactionHistory"`
}

type DoctorProfile struct {
	DoctorID           string   `json:"doctorId"`
	RecordsAdded       []string `json:"recordsAdded"`
	GrantedPermissions []string `json:"grantedPermissions"`
}

type Blockchain struct {
	Blocks []*Block `json:"blocks"`
}
