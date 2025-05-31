package models

type Staff struct {
	ID           uint   `json:"id"`
	Username     string `json:"username" binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"` // ← สำคัญ!
	HospitalID   uint   `json:"hospital_id" binding:"required"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	IsActive     bool   `json:"is_active"`
}

func (Staff) TableName() string {
	return "hospital.staff"
}
