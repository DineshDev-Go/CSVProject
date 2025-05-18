package fetchrevenue

import (
	db "CSVPROJECT/Db"
	modelstruct "CSVPROJECT/ModelStruct"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gorm.io/gorm"
)

func FetchTotalRevenue(w http.ResponseWriter, r *http.Request) {
	log.Println("FetchTotalRevenue (+)")

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", http.MethodGet)
	(w).Header().Set("Access-Control-Allow-Headers", "startDate,endDate,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSTF-Token,Authorization")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	startStr := r.Header.Get("startDate")
	endStr := r.Header.Get("endDate")

	var lResponse modelstruct.RevenueResponse
	var lErr error
	lResponse.Status = "S"

	rawDB, lErr := db.LocalDBConnect()
	if lErr != nil {
		log.Println("FetchTotalRevenue :", lErr.Error())
		lResponse.Status = "E"
		lResponse.ErrMsg = "FetchTotalRevenue :" + lErr.Error()
		return

	}

	lResponse.TotalRevenue, lErr = GetTotalRevenue(startStr, endStr, rawDB)
	if lErr != nil {
		log.Println("FetchTotalRevenue :", lErr.Error())
		lResponse.Status = "E"
		lResponse.ErrMsg = "FetchTotalRevenue :" + lErr.Error()
		return

	}

	lData, lErr := json.Marshal(lResponse)
	if lErr != nil {
		log.Println("FetchTotalRevenue :", lErr.Error())
		return

	} else {
		fmt.Fprint(w, string(lData))
	}

	log.Println("FetchTotalRevenue (-)")
}

func GetTotalRevenue(startDate, endDate string, pDB *gorm.DB) (float64, error) {
	log.Println("GetTotalRevenue (+)")

	var revenue float64 = 0.0

	err := pDB.
		Table("order_Data").
		Select("IFNULL(SUM(unit_price * quantity_sold * (1 - discount)), 0) AS revenue_total").
		Where("date_of_sale BETWEEN ? AND ?", startDate, endDate).
		Scan(&revenue).Error

	if err != nil {
		log.Println("GetTotalRevenue01:", err)
		return 0.0, err
	}

	log.Println("GetTotalRevenue (-)")
	return revenue, nil
}
