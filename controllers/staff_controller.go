package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginStaff(c *gin.Context) {
	var input struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		HospitalID uint   `json:"hospital_id"`
	}

	// ⬇️ Bind และ log input เพื่อ debug
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	log.Printf("🟨 input.Username = %s | input.HospitalID = %d", input.Username, input.HospitalID)

	if input.Username == "" || input.Password == "" || input.HospitalID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, password, and hospital_id are required"})
		return
	}

	var staff models.Staff
	err := config.DB.Table("hospital.staff").
		Where("username = ? AND hospital_id = ?", input.Username, input.HospitalID).
		First(&staff).Error
	if err != nil {
		log.Printf("❌ GORM error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or hospital"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	token, err := utils.GenerateToken(staff.ID, staff.HospitalID, staff.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Login successful",
		"staff_id":    staff.ID,
		"username":    staff.Username,
		"hospital_id": staff.HospitalID,
		"role":        staff.Role,
		"is_active":   staff.IsActive,
		"token":       token,
	})
}

func CreateStaff(c *gin.Context) {
	var input models.Staff
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("📥 Input: Username=%s, HospitalID=%d, Email=%s\n", input.Username, input.HospitalID, input.Email)

	if input.Username == "" || input.PasswordHash == "" || input.Email == "" || input.FirstName == "" || input.LastName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	var existing models.Staff
	if err := config.DB.Table("hospital.staff").
		Where("username = ? AND hospital_id = ?", input.Username, input.HospitalID).
		First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists in this hospital"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("❌ GORM check username error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing username"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	staff := models.Staff{
		Username:     input.Username,
		PasswordHash: string(hashedPassword),
		HospitalID:   input.HospitalID,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Email:        input.Email,
		Role:         "staff",
		IsActive:     true,
	}

	if err := config.DB.Table("hospital.staff").Create(&staff).Error; err != nil {
		log.Printf("❌ Create staff failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create staff"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Staff created successfully",
		"staff":   staff,
	})
}
