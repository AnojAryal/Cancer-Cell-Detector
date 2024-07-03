package initializers

import "github.com/anojaryal/Cancer-Cell-Detector/models"

func SyncDatabase() {
 
    DB.AutoMigrate(
        &models.User{},
        &models.PasswordResetToken{},
        &models.Hospital{},
        &models.Patient{},
        &models.Address{},
        &models.CellTest{},
        &models.Result{},
        &models.CellTestImage{},
        &models.ResultImage{},
    )
}
