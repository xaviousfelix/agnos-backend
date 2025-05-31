package models

type Patient struct {
	ID             uint   `json:"id"`
	HospitalID     uint   `json:"hospital_id"`
	FirstNameTh    string `json:"first_name_th"`
	MiddleNameTh   string `json:"middle_name_th"`
	LastNameTh     string `json:"last_name_th"`
	FirstNameEn    string `json:"first_name_en"`
	MiddleNameEn   string `json:"middle_name_en"`
	LastNameEn     string `json:"last_name_en"`
	DateOfBirth    string `json:"date_of_birth"`
	PatientHN      string `json:"patient_hn"`
	NationalID     string `json:"national_id"`
	PassportID     string `json:"passport_id"`
	PhoneNumber    string `json:"phone_number"`
	Email          string `json:"email"`
	Gender         string `json:"gender"`
	Address        string `json:"address"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}


func (Patient) TableName() string {
	return "hospital.patients"
}
