package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time" // ADD THIS IMPORT

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-uuid"
)

type Order struct {
	ID      string   `json:"id"`
	ItemIDs []string `json:"item_ids"`
	Status  string   `json:"status"`
}

var orders = make(map[string]Order)

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
	registration.ID = "order-service"
	registration.Name = "order-service"
	registration.Port = 8081
	
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

// Discover other services using Consul
func findService(serviceName string) (string, error) {
	config := consulapi.DefaultConfig()
	config.Address = "consul-server:8500"

	consul, err := consulapi.NewClient(config)
	if err != nil {
		return "", err
	}

	services, _, err := consul.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", err
	}

	if len(services) == 0 {
		return "", fmt.Errorf("could not find any healthy instance of %s", serviceName)
	}

	// In a real app, you'd implement load balancing here.
	// For now, we just take the first healthy instance.
	addr := services[0].Service.Address
	port := services[0].Service.Port
	return fmt.Sprintf("http://%s:%d", addr, port), nil
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
			"service": "order-service",
		})
	})

	r.Post("/orders", func(w http.ResponseWriter, r *http.Request) {
		var newOrder Order
		if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Example of inter-service communication
		// Here you would call the food-catalog-service to validate ItemIDs
		catalogAddr, err := findService("food-catalog-service")
		if err != nil {
			http.Error(w, "Food catalog service not available", http.StatusInternalServerError)
			log.Printf("Error finding catalog service: %v", err)
			return
		}
		log.Printf("Found food-catalog-service at: %s. Would validate items here.", catalogAddr)

		orderID, _ := uuid.GenerateUUID()
		newOrder.ID = orderID
		newOrder.Status = "received"
		orders[orderID] = newOrder

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newOrder)
	})

	log.Println("Order Service starting on port 8081...")
	http.ListenAndServe(":8081", r)
}