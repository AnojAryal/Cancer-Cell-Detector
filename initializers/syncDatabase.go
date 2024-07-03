package initializers

import "github.com/anojaryal/Cancer-Cell-Detector/models"

func SyncDatabase() {
    // AutoMigrate all models from the models package
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
