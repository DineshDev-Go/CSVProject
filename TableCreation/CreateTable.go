package tablecreation

import (
	db "CSVPROJECT/Db"
	modelstruct "CSVPROJECT/ModelStruct"
	"log"
)

func ExecuteSchemaUpdates() error {
	log.Println("ExecuteSchemaUpdates (+)")

	rawDB, lErr := db.LocalDBConnect()
	if lErr != nil {
		log.Println("executeDataPurge :", lErr.Error())

		return lErr

	}

	tx := rawDB.Begin()

	if tx.Error != nil {
		log.Println("Migration Start Error:", tx.Error)
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
			log.Println("Panic caught during migration. Transaction rolled back.")
		}
	}()

	modelsToMigrate := []interface{}{
		&modelstruct.AllDataStruct{},
		&modelstruct.CustomerData{},
		&modelstruct.OrderData{},
		&modelstruct.ProductData{},
	}

	for idx, m := range modelsToMigrate {
		if err := tx.AutoMigrate(m); err != nil {
			log.Printf("Migration Error MX-%03d: %v\n", idx+1, err)
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Migration Commit Error:", err)
		return err
	}

	log.Println("MySQL schema migrations completed successfully.")

	log.Println("ExecuteSchemaUpdates (-)")

	return nil
}
