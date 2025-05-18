package csvprocess

import (
	db "CSVPROJECT/Db"
	modelstruct "CSVPROJECT/ModelStruct"
	"log"

	"gorm.io/gorm/clause"
)

func InsertCSVIntoDB(pCustomerRecord []modelstruct.CustomerData, pProductRecord []modelstruct.ProductData, pOrderRecord []modelstruct.OrderData, pAllDataRecord []modelstruct.AllDataStruct) error {

	if len(pAllDataRecord) > 0 {

		lErr := HandleInsertRequest(pAllDataRecord)
		if lErr != nil {
			log.Println("InsertCSVIntoDB01 :", lErr.Error())
			return lErr
		}

	}

	if len(pCustomerRecord) > 0 {

		lErr := HandleInsertRequest(pCustomerRecord)
		if lErr != nil {
			log.Println("InsertCSVIntoDB01 :", lErr.Error())
			return lErr

		}

	}
	if len(pProductRecord) > 0 {

		lErr := HandleInsertRequest(pProductRecord)
		if lErr != nil {
			log.Println("InsertCSVIntoDB01 :", lErr.Error())
			return lErr

		}

	}
	if len(pOrderRecord) > 0 {

		lErr := HandleInsertRequest(pOrderRecord)
		if lErr != nil {
			log.Println("InsertCSVIntoDB01 :", lErr.Error())
			return lErr

		}

	}
	return nil

}

func HandleInsertRequest(data interface{}) error {
	log.Println("HandleInsertRequest started")

	rawDB, err := db.LocalDBConnect()
	if err != nil {
		log.Fatal("failed to get raw sql.DB:", err)
	}

	tx := rawDB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("Recovered from panic in HandleInsertRequest:", r)
		}
	}()

	insertResult := tx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(data, 1000)
	if insertResult.Error != nil {
		log.Println("MYSQL_INSERT_ERR:", insertResult.Error)
		return insertResult.Error
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("MYSQL_COMMIT_ERR:", err)
		return err
	}

	log.Println("HandleInsertRequest completed")
	return nil
}
