package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"backend/utils" // ✅ ปรับให้ตรงกับ module name

	// ...
)

func GetExternalPatient(c *gin.Context) {
	id := c.Param("id")

	patient, err := utils.FetchPatientFromAPI(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, patient)
}
