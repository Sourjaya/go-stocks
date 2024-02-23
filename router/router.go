package router

import (
	"go-stocks/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/stock/{id}", middleware.GetStockByID).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/stock", middleware.GetAllStocks).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/createstock", middleware.CreateStock).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/stock/{id}", middleware.UpdateStock).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/delete/{id}", middleware.DeleteStock).Methods("DELETE", "OPTIONS")
	return router
}
