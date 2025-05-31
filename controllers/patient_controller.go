package controllers

import (
	"backend/models"
	"backend/config"
	"net/http"
	// "strconv"

	"github.com/gin-gonic/gin"
)

func SearchPatient(c *gin.Context) {

    // ดึง hospital_id จาก context
	hospitalID, exists := c.Get("hospital_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// รับ query param (optional)
	nationalID := c.Query("national_id")
	passportID := c.Query("passport_id")
	firstName := c.Query("first_name")
    middleName := c.Query("middle_name")
	lastName := c.Query("last_name")
	dob := c.Query("date_of_birth") // ควรเป็น format YYYY-MM-DD
	phone := c.Query("phone_number")
    Gender := c.Query("gender")
	email := c.Query("email")

    // สร้าง query
	query := config.DB.Table("hospital.patients").Where("hospital_id = ?", hospitalID)

    if firstName != "" {
		query = query.Where("first_name_th ILIKE ?", "%"+firstName+"%")
	}
    if middleName != "" {
		query = query.Where("middle_name_th ILIKE ?", "%"+middleName+"%")
	}
	if lastName != "" {
		query = query.Where("last_name_th ILIKE ?", "%"+lastName+"%")
	}
	if nationalID != "" {
		query = query.Where("national_id = ?", nationalID)
	}
	if passportID != "" {
		query = query.Where("passport_id = ?", passportID)
	}
	if Gender != "" {
		query = query.Where("gender = ?", Gender)
	}
    if dob != "" {
		query = query.Where("date_of_birth = ?", dob)
	}
	if phone != "" {
		query = query.Where("phone_number = ?", phone)
	}
	if email != "" {
		query = query.Where("email ILIKE ?", "%"+email+"%")
	}

	// ✅ ถ้าไม่มีการส่งค่าใดๆ มาเลย ให้แจ้งว่าไม่สามารถค้นหาได้
if nationalID == "" && passportID == "" &&
	firstName == "" && middleName == "" &&
	lastName == "" && dob == "" &&
	phone == "" && email == "" {
	c.JSON(http.StatusBadRequest, gin.H{"error": "At least one search parameter is required"})
	return
}


	var patients []models.Patient
	if err := query.Find(&patients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ดึงข้อมูลผู้ป่วย
	if err := query.Find(&patients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ถ้าไม่พบข้อมูลผู้ป่วย
	if len(patients) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No patients found"})
		return
	}

	// ส่งผลลัพธ์กลับ
	c.JSON(http.StatusOK, gin.H{"patients": patients})
}

func SearchPatientByNationalOrPassportID(c *gin.Context) {
	id := c.Param("id")

	// ✅ ดึง hospital_id จาก context
	hospitalID, exists := c.Get("hospital_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var patient models.Patient
	err := config.DB.Table("hospital.patients").
		Where("(national_id = ? OR passport_id = ?) AND hospital_id = ?", id, id, hospitalID).
		First(&patient).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	c.JSON(http.StatusOK, patient)
}


func GetAllPatients(c *gin.Context) {
	var patients []models.Patient

	err := config.DB.Table("hospital.patients").Find(&patients).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"patients": patients})
}
