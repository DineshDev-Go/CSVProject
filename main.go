package main

import (
	csvprocess "CSVPROJECT/CSVProcess"
	fetchrevenue "CSVPROJECT/FetchRevenue"
	schedular "CSVPROJECT/FetchRevenue/Schedular"
	tablecreation "CSVPROJECT/TableCreation"
	"fmt"
	"net/http"
)

func main() {

	tablecreation.ExecuteSchemaUpdates()

	go schedular.InitDailyCleanupTask()

	http.HandleFunc("/insertCSV", csvprocess.SetCSVFile)
	http.HandleFunc("/totalRevenue", fetchrevenue.FetchTotalRevenue)
	http.HandleFunc("/productRevenue", fetchrevenue.FetchTotalRevenueByProduct)
	http.HandleFunc("/categoryRevenue", fetchrevenue.FetchTotalRevenueByCategory)
	http.HandleFunc("/regionRevenue", fetchrevenue.FetchTotalRevenueByRegion)
	http.HandleFunc("/restart", schedular.HandleManualDataPurge)

	http.ListenAndServe(":8080", nil)
	fmt.Println("server end")

}
