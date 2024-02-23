package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-stocks/models"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading env variables")
	}
	url := os.Getenv("POSTGRES_URL")
	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to Postgres DB")
	return db
}
func GetStockByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string id to integer type: %v", err)
	}
	stock, err := getStock(int64(id))
	if err != nil {
		log.Fatalf("Unable to get stock: %v", err)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "	")
	enc.Encode(stock)
}
func GetAllStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := allstocks()
	if err != nil {
		log.Fatalf("Unable to get stocks: %v", err)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "	")
	enc.Encode(stocks)
}
func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock
	if err := json.NewDecoder(r.Body).Decode(&stock); err != nil {
		log.Fatalf("Decoding json failed,%v", err)
	}
	insertID := insertStock(stock)
	res := response{
		ID:      insertID,
		Message: "stock created successfully",
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "	")
	enc.Encode(res)

}
func DeleteStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string id to integer type: %v", err)
	}
	deletedRows := delete(int64(id))
	msg := fmt.Sprintf("Stock deleted successfully. Total rows/records affected %v", deletedRows)
	res := response{
		ID:      int64(id),
		Message: msg,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "	")
	enc.Encode(res)
}
func UpdateStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string id to integer type: %v", err)
	}
	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	updatedRows := update(int64(id), stock)
	msg := fmt.Sprintf("Stock updated successfully. Total rows/records affected %v", updatedRows)
	res := response{
		ID:      int64(id),
		Message: msg,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "	")
	enc.Encode(res)
}

func insertStock(stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()
	var id int64
	statement := `INSERT INTO stocks(name,price,company) VALUES ($1,$2,$3) RETURNING stockid`
	if err := db.QueryRow(statement, stock.Name, stock.Price, stock.Company).Scan(&id); err != nil {
		log.Fatalf("Unable to execute the query: %v", err)
	}
	fmt.Printf("Inserted a single record %v", id)
	return id
}
func getStock(id int64) (models.Stock, error) {
	db := createConnection()
	defer db.Close()
	var stock models.Stock
	statement := `SELECT * FROM stocks WHERE stockid=$1`
	err := db.QueryRow(statement, id).Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("Now rows were returned!")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("Unable to scan the row: %v", err)
	}

	return stock, err

}
func allstocks() ([]models.Stock, error) {
	db := createConnection()
	defer db.Close()
	var stocks []models.Stock
	statement := `SELECT * FROM stocks`
	rows, err := db.Query(statement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var stock models.Stock
		err = rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		stocks = append(stocks, stock)
	}
	return stocks, err
}
func delete(id int64) int64 {
	db := createConnection()
	defer db.Close()
	statement := `DELETE FROM stocks WHERE stockid=$1`
	res, err := db.Exec(statement, id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows/record affected %v", rowsAffected)
	return rowsAffected
}
func update(id int64, stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1`
	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows/record affected %v", rowsAffected)
	return rowsAffected
}
