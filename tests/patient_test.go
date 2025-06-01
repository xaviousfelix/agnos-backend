package tests

import (
	"backend/config"
	"backend/controllers"
	"backend/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
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

// ✅ เชื่อมต่อ DB สำหรับ test ก่อนทุกเทสต์
func init() {
	config.DB = ConnectTestDB()
}
func fakeAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// mock staff info
		c.Set("staff_id", uint(1))
		c.Set("hospital_id", uint(1))
		c.Set("role", "staff")
		c.Next()
	}
}

func setupRouterWithAuth() *gin.Engine {
	r := gin.Default()
	r.Use(fakeAuthMiddleware()) // mock token auth
	r.GET("/patient/search/:id", controllers.SearchPatientByNationalOrPassportID)
	return r
}

func getTestToken() string {
	token, _ := utils.GenerateToken(1, 1, "staff")
	return token
}

func TestSearchPatient_Success(t *testing.T) {

	router := setupRouterWithAuth()

	req := httptest.NewRequest("GET", "/patient/search/1234567890123", nil)
	req.Header.Set("Authorization", "Bearer "+getTestToken())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 but got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "1234567890123") {
		t.Errorf("Expected patient data, got: %v", w.Body.String())
	}
}

func TestSearchPatient_NotFound(t *testing.T) {
	r := setupRouterWithAuth()
	req := httptest.NewRequest("GET", "/patient/search/0000000000000", nil)
	req.Header.Set("Authorization", "Bearer "+getTestToken())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404 but got %d", w.Code)
	}
}
