package services

import (
    "backend/config"
    "backend/models"
)

func FindStaffByUsername(username string) (*models.Staff, error) {
    var staff models.Staff
    result := config.DB.Where("username = ?", username).First(&staff)
    if result.Error != nil {
        return nil, result.Error
    }
    return &staff, nil
}
