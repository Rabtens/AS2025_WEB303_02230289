package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	consulapi "github.com/hashicorp/consul/api"
)

type FoodItem struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var foodItems = []FoodItem{
	{ID: "1", Name: "Coffee", Price: 2.50},
	{ID: "2", Name: "Sandwich", Price: 5.00},
	{ID: "3", Name: "Muffin", Price: 3.25},
}

// REPLACE THE ENTIRE registerServiceWithConsul FUNCTION WITH THIS:
func registerServiceWithConsul() {
	config := consulapi.DefaultConfig()
	config.Address = "consul-server:8500"

	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Printf("Warning: Could not create Consul client: %v", err)
		return
	}

	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = "food-catalog-service"
	registration.Name = "food-catalog-service"
	registration.Port = 8080
	
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Could not get hostname: %v", err)
	}
	registration.Address = hostname

	// Use TTL health check instead of HTTP
	registration.Check = &consulapi.AgentServiceCheck{
		TTL: "15s",
	}

	if err := consul.Agent().ServiceRegister(registration); err != nil {
		log.Printf("Warning: Failed to register service with Consul: %v", err)
		return
	}
	log.Println("Successfully registered service with Consul")

	// Start a goroutine to periodically update the health check
	go func() {
		for {
			time.Sleep(10 * time.Second)
			checkID := "service:" + registration.ID
			if err := consul.Agent().UpdateTTL(checkID, "Service is healthy", "pass"); err != nil {
				log.Printf("Warning: Failed to update TTL health check: %v", err)
			} else {
				log.Printf("TTL health check updated successfully")
			}
		}
	}()
}

func main() {
	// Try to register with Consul, but don't fail if it's not available
	go registerServiceWithConsul()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
			"service": "food-catalog-service",
		})
	})

	r.Get("/items", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(foodItems)
	})

	log.Println("Food Catalog Service starting on port 8080...")
	http.ListenAndServe(":8080", r)
}