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

func FetchTotalRevenueByProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("FetchTotalRevenueByProduct (+)")

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
	ProductName := r.Header.Get("productName")

	var lResponse modelstruct.RevenueResponse
	var lErr error
	lResponse.Status = "S"

	rawDB, lErr := db.LocalDBConnect()
	if lErr != nil {
		log.Println("FetchTotalRevenueByProduct01 :", lErr.Error())
		lResponse.Status = "E"
		lResponse.ErrMsg = "FetchTotalRevenueByProduct :" + lErr.Error()
		return

	}

	lResponse.TotalRevenue, lErr = GetTotalRevenueByProduct(startStr, endStr, rawDB, ProductName)
	if lErr != nil {
		log.Println("FetchTotalRevenueByProduct02 :", lErr.Error())
		lResponse.Status = "E"
		lResponse.ErrMsg = "FetchTotalRevenueByProduct02 :" + lErr.Error()
		return

	}

	lData, lErr := json.Marshal(lResponse)
	if lErr != nil {
		log.Println("FetchTotalRevenueByProduct03 :", lErr.Error())
		return

	} else {
		fmt.Fprint(w, string(lData))
	}

	log.Println("FetchTotalRevenueByProduct (-)")
}

func GetTotalRevenueByProduct(startDate, endDate string, pDB *gorm.DB, pProductName string) (float64, error) {
	log.Println("GetTotalRevenueByProduct (+)")

	var productId string
	var revenue float64 = 0.0

	// Step 1: Get the product ID for the given product name
	err := pDB.
		Table("product_data").
		Select("product_id").
		Where("product_name = ?", pProductName).
		Scan(&productId).Error

	if err != nil {
		log.Println("GetRevenueByProduct01:", err)
		return 0.0, err
	}

	// Step 2: Use product_id to get revenue from order_details
	err = pDB.
		Table("order_data").
		Select("IFNULL(SUM(unit_price * quantity_sold * (1 - discount)), 0) AS revenue_Product_total").
		Where("product_id = ? AND date_of_sale BETWEEN ? AND ?", productId, startDate, endDate).
		Scan(&revenue).Error

	if err != nil {
		log.Println("GetRevenueByProduct02:", err)
		return 0.0, err
	}

	log.Println("GetTotalRevenueByProduct (-)")
	return revenue, nil
}
