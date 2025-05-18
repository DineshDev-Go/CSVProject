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

func FetchTotalRevenueByRegion(w http.ResponseWriter, r *http.Request) {
	log.Println("FetchTotalRevenueByRegion (+)")

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", http.MethodGet)
	(w).Header().Set("Access-Control-Allow-Headers", "startDate,endDate,region,Accept,Content-Type,Content-Length,Accept-Encoding,X-CSTF-Token,Authorization")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	startStr := r.Header.Get("startDate")
	endStr := r.Header.Get("endDate")
	Region := r.Header.Get("region")

	var lResponse modelstruct.RevenueResponse
	var lErr error
	lResponse.Status = "S"

	rawDB, lErr := db.LocalDBConnect()
	if lErr != nil {
		log.Println("FetchTotalRevenueByRegion01 :", lErr.Error())
		lResponse.Status = "E"
		lResponse.ErrMsg = "FetchTotalRevenueByRegion01 :" + lErr.Error()
		return

	}

	lResponse.TotalRevenue, lErr = GetTotalRevenueByRegion(startStr, endStr, rawDB, Region)
	if lErr != nil {
		log.Println("FetchTotalRevenueByRegion02 :", lErr.Error())
		lResponse.Status = "E"
		lResponse.ErrMsg = "FetchTotalRevenueByRegion02 :" + lErr.Error()
		return

	}

	lData, lErr := json.Marshal(lResponse)
	if lErr != nil {
		log.Println("FetchTotalRevenueByRegion03 :", lErr.Error())
		return

	} else {
		fmt.Fprint(w, string(lData))
	}

	log.Println("FetchTotalRevenueByRegion (-)")
}

func GetTotalRevenueByRegion(startDate, endDate string, pDB *gorm.DB, pProductName string) (float64, error) {
	log.Println("GetTotalRevenueByRegion (+)")

	var productId string
	var revenue float64 = 0.0

	// Step 1: Get the product ID for the given product name
	err := pDB.
		Table("product_data").
		Select("product_id").
		Where("product_name = ?", pProductName).
		Scan(&productId).Error

	if err != nil {
		log.Println("GetTotalRevenueByRegion01:", err)
		return 0.0, err
	}

	// Step 2: Use product_id to get revenue from order_details
	err = pDB.Table("order_data").Where("region = ? and date_of_sale >= ? and date_of_sale <= ?", productId, startDate, endDate).Select("coalesce(SUM((unit_price * quantity_sold) * (1 - discount)),0) AS total_sales").Scan(&revenue).Error

	if err != nil {
		log.Println("GetTotalRevenueByRegion02:", err)
		return 0.0, err
	}

	log.Println("GetTotalRevenueByRegion (-)")
	return revenue, nil
}
