package schedular

import (
	db "CSVPROJECT/Db"
	"fmt"
	"log"
	"net/http"
	"time"
)

func InitDailyCleanupTask() {
	log.Println("InitDailyCleanupTask (+)")
	go func() {
		for {
			currentTime := time.Now()
			if currentTime.Hour() == 0 && currentTime.Minute() == 0 {
				if err := executeDataPurge(); err != nil {
					log.Println(" data cleanup failed:", err)
				} else {
					log.Println(" data cleanup completed successfully")
				}
			}
			time.Sleep(1 * time.Minute)
		}
	}()
	log.Println("InitDailyCleanupTask (-)")

}

func HandleManualDataPurge(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleManualDataPurge (+)")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := executeDataPurge(); err != nil {
		log.Println(" data purge failed:", err)
		http.Error(w, " data purge failed", http.StatusInternalServerError)
		return
	}
	log.Println(" data purge completed successfully")
	w.Write([]byte("Manual purge completed"))

	log.Println("HandleManualDataPurge (-)")

}

func executeDataPurge() error {

	log.Println("executeDataPurge (+)")
	startTime := time.Now()
	log.Println(" Starting data purge at:", startTime)

	rawDB, lErr := db.LocalDBConnect()
	if lErr != nil {
		log.Println("executeDataPurge :", lErr.Error())

		return lErr

	}

	tx := rawDB.Begin()

	tables := []string{
		"all_data_struct",
		"product_data",
		"customer_data",
		"order_data",
	}

	for _, table := range tables {
		if err := tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table)).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to truncate table %s: %v", table, err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Println(" Data purge finished at:", time.Now())

	log.Println("executeDataPurge (-)")

	return nil
}
