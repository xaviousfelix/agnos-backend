package tests

import (
	"backend/config"
	"backend/controllers"
	"fmt"
	"time"

	// "backend/models"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectTestDB() *gorm.DB {
	dsn := "host=db user=postgres password=232546 dbname=hospital_test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect test database:", err)
	}
	return db
}

func init() {
	config.DB = ConnectTestDB()
}

func setupStaffRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/staff/login", controllers.LoginStaff)
	return r
}

func setupStaffCreateRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/staff/create", controllers.CreateStaff) // ใช้ controller ของคุณ
	return r
}

func seedTestStaff() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	config.DB.Exec("DELETE FROM hospital.staff WHERE username = ?", "teststaff")
	config.DB.Exec(`
    INSERT INTO hospital.staff (
        username,
        password_hash,
        hospital_id,
        first_name,
        last_name,
        email,
        role,
        is_active,
        created_at,
        updated_at
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		"teststaff",             // username
		hashedPassword,          // password_hash
		1,                       // hospital_id
		"Tester",                // first_name
		"Mock",                  // last_name
		"teststaff@example.com", // email
		"staff",                 // role
		true)                    // is_active

}
func TestStaffLogin_Success(t *testing.T) {
	seedTestStaff()

	r := setupStaffRouter() // ✅ ใช้อันนี้
	input := map[string]interface{}{
		"username":    "teststaff",
		"password":    "12345678",
		"hospital_id": 1,
	}
	bodyBytes, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/staff/login", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 but got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "token") {
		t.Errorf("Expected token in response, got: %v", w.Body.String())
	}
}

func TestStaffLogin_Fail_WrongPassword(t *testing.T) {
	r := setupStaffRouter()

	body := map[string]interface{}{
		"username":    "testuser",
		"password":    "wrongpass",
		"hospital_id": 1,
	}
	jsonValue, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/staff/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 but got %d", w.Code)
	}
}

func TestStaffLogin_Fail_EmptyField(t *testing.T) {
	r := setupStaffRouter()

	body := map[string]interface{}{
		"username":    "",
		"password":    "12345678",
		"hospital_id": 1,
	}
	jsonValue, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/staff/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %d", w.Code)
	}

}

func TestStaffLogin_Fail_NoUser(t *testing.T) {
	r := setupStaffRouter()

	body := map[string]interface{}{
		"username":    "nouser",
		"password":    "password123",
		"hospital_id": 1,
	}
	jsonValue, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/staff/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 but got %d", w.Code)
	}
}

func TestStaffCreate_Success(t *testing.T) {
	r := setupStaffCreateRouter()

	// ✅ สร้างชื่อไม่ซ้ำกันในทุกครั้งที่รัน
	uniqueUsername := fmt.Sprintf("newstaff_%d", time.Now().UnixNano())

	body := map[string]interface{}{
		"username":      uniqueUsername,
		"password_hash": "12345678",
		"hospital_id":   1,
		"first_name":    "New",
		"last_name":     "Staff",
		"email":         fmt.Sprintf("%s@example.com", uniqueUsername),
		"role":          "staff",
	}
	jsonValue, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/staff/create", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Errorf("Expected status 201 or 200, but got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), uniqueUsername) {
		t.Errorf("Expected response to contain %s, got: %s", uniqueUsername, w.Body.String())
	}
}

func TestStaffCreate_Fail_MissingFields(t *testing.T) {
	r := setupStaffCreateRouter()

	// ขาด username และ password
	body := map[string]interface{}{
		"hospital_id": 1,
		"first_name":  "Missing",
		"last_name":   "Fields",
		"email":       "fail@example.com",
		"role":        "staff",
	}
	jsonValue, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/staff/create", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 Bad Request, but got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "required") && !strings.Contains(w.Body.String(), "missing") {
		t.Errorf("Expected error message about missing fields, got: %v", w.Body.String())
	}
}
