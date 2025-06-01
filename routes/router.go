package routes

import (
	"net/http"

	"backend/controllers"
	"backend/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// Route สำหรับ root path
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API is running",
		})
	})
	// Route สำหรับจัดการข้อมูลแพทย์
	r.POST("/staff/login", controllers.LoginStaff)
	// Route สำหรับสร้าง staff ใหม่
	r.POST("/staff/create", controllers.CreateStaff)

	//กลุ่มที่ต้อง login ด้วย JWT
	auth := r.Group("/patient", middleware.AuthMiddleware())
	{
		auth.GET("/search/:id", controllers.SearchPatientByNationalOrPassportID) // path param
		auth.GET("/search", controllers.SearchPatient)                           //query param
	}
}
