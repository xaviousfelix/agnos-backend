package services

import (
    "backend/config"
    "backend/models"
)

func FindPatientByID(id uint) (*models.Patient, error) {
    var patient models.Patient
    result := config.DB.First(&patient, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &patient, nil
}
