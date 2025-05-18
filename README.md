# CSVProject
For CSV Set and Get


Go Project for insert csv


Tables

1.all_data_struct
2.customer_data
3.product_data
4.order_data

How to Run:

1.setup DataBase :
      "CREATE DATABASE dbName"

2.Mention database credentials in config.env
    DB_USER=root  
    DB_PASSWORD=Dinesh@21
    DB_HOST=localhost
    DB_PORT=3306
    DB_NAME=dinesh_db

3.Run with below command
     go run  main.go

4.Endpoint: /insertCSV
Method: POST
Description: Uploads a CSV file into the database table.
Headers:
Content-Type: multipart/form-data
Request: Form-data with a file input named files
Response:
{
  "status":"S",
  "errMsg":"",
}

5. Get Total Revenue by Date
Endpoint: /totalRevenue
Method: GET
Description: Returns the total revenue between a given date range.
Headers:
start: YYYY-MM-DD
end: YYYY-MM-DD
Response:
{
  "status":"S",
  "errMsg":"",
  "totRevenue": 12345.67
}


6. Get Revenue by Product
Endpoint: /productRevenue
Method: GET
Description: Returns total revenue for a specific product between two dates.
header:
  "productName": "Product1",
  "startDate": "2024-01-01",
  "endDate": "2024-01-31"


Response:
{
  "status":"S",
  "errMsg":"",
  "totRevenue": 4567.89
}


7. Get Revenue by Category
Endpoint: /categoryRevenue
Method: GET
Description: Returns total revenue for a specific category between two dates.
Header:

  "CategoryName": "tools",
  "startDate": "2024-01-01",
  "endDate": "2024-01-31"

Response:
{
  "status":"S",
  "errMsg":"",
  "totRevenue": 7890.12
}


8. Get Revenue by Region
Endpoint: /regionRevenue
Method: GET
Description: Returns total revenue for a specific region between two dates.
Header:
 
  "region": "North",
  "startDate": "2024-01-01",
  "endDate": "2024-01-31"

Response:
{
  "status":"S",
  "errMsg":"",
  "totRevenue": 3456.78
}


9.manual refresh the table:
Endpoint: /restart
Method: PUT
Description:reset the table 
  

  *note: for db config use config.env for reference
    

   