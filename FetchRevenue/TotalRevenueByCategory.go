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

func FetchTotalRevenueByCategory(w http.ResponseWriter, r *http.Request) {
	log.Println("FetchTotalRevenueByCategory (+)")

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", http.MethodGet)
	(w).Header().Set("Access-Control-Allow-Headers", "startDate,endDate,CategoryName,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSTF-Token,Authorization")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	startStr := r.Header.Get("startDate")
	endStr := r.Header.Get("endDate")
	Category := r.Header.Get("CategoryName")

	var lResponse modelstruct.RevenueResponse
	var lErr error
	lResponse.Status = "S"

	rawDB, lErr := db.LocalDBConnect()
	if lErr != nil {
		log.Println("FetchTotalRevenueByCategory01 :", lErr.Error())
		lResponse.Status = "E"
		lResponse.ErrMsg = "FetchTotalRevenueByCategory :" + lErr.Error()
		return

	}

	lResponse.TotalRevenue, lErr = GetTotalRevenueByCategory(startStr, endStr, rawDB, Category)
	if lErr != nil {
		log.Println("FetchTotalRevenueByCategory02 :", lErr.Error())
		lResponse.Status = "E"
		lResponse.ErrMsg = "FetchTotalRevenueByCategory02 :" + lErr.Error()
		return

	}

	lData, lErr := json.Marshal(lResponse)
	if lErr != nil {
		log.Println("FetchTotalRevenueByCategory03 :", lErr.Error())
		return

	} else {
		fmt.Fprint(w, string(lData))
	}

	log.Println("FetchTotalRevenueByCategory (-)")
}

func GetTotalRevenueByCategory(startDate, endDate string, pDB *gorm.DB, pCategory string) (float64, error) {
	log.Println("GetTotalRevenueByCategory (+)")

	var lProductId string
	var revenue float64 = 0.0

	// Step 1: Get the product ID for the given product name
	err := pDB.
		Table("product_data").
		Select("product_id").
		Where("product_name = ?", pCategory).
		Scan(&lProductId).Error

	if err != nil {
		log.Println("GetTotalRevenueByCategory01:", err)
		return 0.0, err
	}

	// Step 2: Use product_id to get revenue from order_details
	err = pDB.Table("order_data").Where("product_id = ? and date_of_sale >= ? and date_of_sale <= ?", lProductId, startDate, endDate).
		Select("coalesce(SUM((unit_price * quantity_sold) * (1 - discount)),0) AS revenue_Category_total").
		Scan(&revenue).Error

	if err != nil {
		log.Println("GetTotalRevenueByCategory02:", err)
		return 0.0, err
	}

	log.Println("GetTotalRevenueByCategory (-)")
	return revenue, nil
}
