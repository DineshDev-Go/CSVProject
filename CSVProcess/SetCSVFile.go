package csvprocess

import (
	modelstruct "CSVPROJECT/ModelStruct"
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func SetCSVFile(w http.ResponseWriter, r *http.Request) {

	log.Println("SetCSVFile (+)")

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", http.MethodPost)
	(w).Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSTF-Token,Authorization")

	var lOrderArray []modelstruct.OrderData
	var lProductArray []modelstruct.ProductData
	var lCustomerArray []modelstruct.CustomerData
	var lAllDataArray []modelstruct.AllDataStruct

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// var lInsertedIds string
	var lRespRec modelstruct.ResponseStruct
	lRespRec.Status = "S"

	lErr := r.ParseMultipartForm(10 << 20)
	if lErr != nil {
		log.Println("SetCSVFile SCF01 :", lErr.Error())
		lRespRec.Status = "E"
		lRespRec.ErrMsg = "SetCSVFile SCF01  :" + lErr.Error()
	}

	lFiles := r.MultipartForm.File["files"]

	lCustomerArray, lProductArray, lOrderArray, lAllDataArray, lErr = FileProcess(lFiles)
	if lErr != nil {
		log.Println("SetCSVFile SCF02 :", lErr.Error())
		lRespRec.Status = "E"
		lRespRec.ErrMsg = "SetCSVFile SCF02 :" + lErr.Error()
	}

	lErr = InsertCSVIntoDB(lCustomerArray, lProductArray, lOrderArray, lAllDataArray)
	if lErr != nil {
		log.Println("SetCSVFile SCF03 :", lErr.Error())
		lRespRec.Status = "E"
		lRespRec.ErrMsg = "SetCSVFile SCF03 :" + lErr.Error()
	}

	lData, lErr := json.Marshal(lRespRec)
	if lErr != nil {
		log.Println("SetCSVFile SCF04 :", lErr.Error())
	} else {
		fmt.Fprint(w, string(lData))
	}

	log.Println("SetCSVFile (-)")
}

func FileProcess(pFiles []*multipart.FileHeader) ([]modelstruct.CustomerData, []modelstruct.ProductData, []modelstruct.OrderData, []modelstruct.AllDataStruct, error) {

	var lCustomer []modelstruct.CustomerData
	var lProduct []modelstruct.ProductData
	var lOrder []modelstruct.OrderData
	var lAllData []modelstruct.AllDataStruct

	for _, lFileHeder := range pFiles {
		lFile, lErr := lFileHeder.Open()
		if lErr != nil {
			log.Println("FileProcess01 :", lErr.Error())
			return lCustomer, lProduct, lOrder, lAllData, lErr
		}
		defer lFile.Close()

		lReader := csv.NewReader(bufio.NewReader(lFile))

		lHeader, lErr := lReader.Read()
		if lErr != nil {
			log.Println("FileProcess02 :", lErr.Error())
			return lCustomer, lProduct, lOrder, lAllData, lErr
		}

		lIndex := make(map[string]int)
		for i, lVal := range lHeader {
			lIndex[lVal] = i
		}

		lRecords, lErr := lReader.ReadAll()
		if lErr != nil {
			log.Println("FileProcess03 :", lErr.Error())
			return lCustomer, lProduct, lOrder, lAllData, lErr
		}
		uniqueID := fmt.Sprintf("UID-%d", time.Now().Unix())

		lCustomerList, lProductList, lOrderList, lAllRecordList, lErr := ProcessSalesRecords(lRecords, lIndex, uniqueID)
		if lErr != nil {
			log.Println("FileProcess04 :", lErr.Error())
			return lCustomer, lProduct, lOrder, lAllData, lErr
		}

		if len(lAllRecordList) > 0 {
			lAllData = append(lAllData, lAllRecordList...)

		}

		if len(lOrderList) > 0 {
			lOrder = append(lOrder, lOrderList...)
		}
		if len(lProductList) > 0 {
			lProduct = append(lProduct, lProductList...)
		}
		if len(lCustomerList) > 0 {
			lCustomer = append(lCustomer, lCustomerList...)
		}

	}

	return lCustomer, lProduct, lOrder, lAllData, nil

}

func ProcessSalesRecords(pRecords [][]string, pIndex map[string]int, uniqueId string) (
	[]modelstruct.CustomerData,
	[]modelstruct.ProductData,
	[]modelstruct.OrderData,
	[]modelstruct.AllDataStruct, error) {
	var (
		customerArr []modelstruct.CustomerData
		productArr  []modelstruct.ProductData
		orderArr    []modelstruct.OrderData
		fileDataArr []modelstruct.AllDataStruct
	)

	createdDate := time.Now().Format("2006-01-02 15:04:05")

	for _, row := range pRecords {
		fileData := parseSalesRecord(row, pIndex, uniqueId, createdDate)
		customerArr = append(customerArr, toCustomerDetails(fileData))
		productArr = append(productArr, toProductDetails(fileData))
		orderArr = append(orderArr, toOrderDetails(fileData))
		fileDataArr = append(fileDataArr, fileData)
	}

	return customerArr, productArr, orderArr, fileDataArr, nil
}

func parseSalesRecord(row []string, cols map[string]int, uniqueId, createdDate string) modelstruct.AllDataStruct {
	return modelstruct.AllDataStruct{
		OrderId:         strings.TrimSpace(row[cols["Order ID"]]),
		ProductId:       strings.TrimSpace(row[cols["Product ID"]]),
		CustomerId:      strings.TrimSpace(row[cols["Customer ID"]]),
		ProductName:     strings.TrimSpace(row[cols["Product Name"]]),
		Category:        strings.TrimSpace(row[cols["Category"]]),
		Region:          strings.TrimSpace(row[cols["Region"]]),
		DateOfSale:      parseDate(row[cols["Date of Sale"]]),
		QuantitySold:    parseInt(row[cols["Quantity Sold"]]),
		UnitPrice:       parseFloat(row[cols["Unit Price"]]),
		Discount:        parseFloat(row[cols["Discount"]]),
		ShippingCost:    parseFloat(row[cols["Shipping Cost"]]),
		PaymentMethod:   strings.TrimSpace(row[cols["Payment Method"]]),
		CustomerName:    strings.TrimSpace(row[cols["Customer Name"]]),
		CustomerEmail:   strings.TrimSpace(row[cols["Customer Email"]]),
		CustomerAddress: strings.TrimSpace(row[cols["Customer Address"]]),
		UniqueId:        uniqueId,
		CreatedBy:       "AutoBot",
		CreatedDate:     createdDate,
	}
}

func toCustomerDetails(d modelstruct.AllDataStruct) modelstruct.CustomerData {
	return modelstruct.CustomerData{
		CustomerId:      d.CustomerId,
		CustomerName:    d.CustomerName,
		CustomerEmail:   d.CustomerEmail,
		CustomerAddress: d.CustomerAddress,
		CreatedDate:     d.CreatedDate,
		CreatedBy:       d.CreatedBy,
	}
}

func toProductDetails(d modelstruct.AllDataStruct) modelstruct.ProductData {
	return modelstruct.ProductData{
		ProductId:   d.ProductId,
		ProductName: d.ProductName,
		Category:    d.Category,
		UnitPrice:   d.UnitPrice,
		CreatedDate: d.CreatedDate,
		CreatedBy:   d.CreatedBy,
	}
}

func toOrderDetails(d modelstruct.AllDataStruct) modelstruct.OrderData {
	return modelstruct.OrderData{
		OrderId:       d.OrderId,
		ProductId:     d.ProductId,
		CustomerId:    d.CustomerId,
		QuantitySold:  d.QuantitySold,
		UnitPrice:     d.UnitPrice,
		Discount:      d.Discount,
		ShippingCost:  d.ShippingCost,
		DateOfSale:    d.DateOfSale,
		Region:        d.Region,
		PaymentMethod: d.PaymentMethod,
		CreatedDate:   d.CreatedDate,
		CreatedBy:     d.CreatedBy,
	}
}

func parseInt(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

func parseFloat(s string) float64 {
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

func parseDate(s string) string {
	t, _ := time.Parse("1/2/2006", s)
	return t.Format("2006-01-02")
}

// func InsertFiles() {

// }
