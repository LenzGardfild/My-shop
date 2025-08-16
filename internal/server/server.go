package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"my-shop/internal/cache"
	"my-shop/internal/db"
	"my-shop/internal/kafka"

	"github.com/gorilla/mux"
)

type Server struct {
	cache *cache.OrderCache
	db    *db.DB
}

func New() *Server {
	connStr := "postgres://myshop_user:myshop_pass@localhost:5432/myshop_db?sslmode=disable"
	dbConn, err := db.New(connStr)
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	orderCache := cache.New()
	orders, err := dbConn.LoadAllOrders(context.Background())
	if err == nil {
		for _, o := range orders {
			orderCache.Set(o.OrderUID, o)
		}
	}
	// Запуск Kafka consumer
	go kafka.RunKafkaConsumer(
		[]string{"localhost:9092"},
		"myshop-group",
		"orders",
		dbConn,
		orderCache,
	)
	return &Server{cache: orderCache, db: dbConn}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
	router.HandleFunc("/order/{id}", s.handleGetOrder).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("internal/server/static")))
	router.ServeHTTP(w, r)
}

func (s *Server) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "order id required", http.StatusBadRequest)
		return
	}
	if order, ok := s.cache.Get(id); ok {
		json.NewEncoder(w).Encode(order)
		return
	}
	order, err := s.db.GetOrder(r.Context(), id)
	if err != nil || order == nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}
	s.cache.Set(id, order)
	json.NewEncoder(w).Encode(order)
}
