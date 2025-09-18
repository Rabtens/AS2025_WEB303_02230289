# WEB303 Practical 3: Full-Stack Microservices Implementation

**Student Name:** [Kuenzang Rabten]  
**Course:** WEB303 Microservices & Serverless Applications  
**Practical:** 3 - Full-Stack Microservices with gRPC, Databases, and Service Discovery  
**Date:** [26/08/2025]

## Executive Summary

This project implements a complete microservices ecosystem featuring:
- Two independent services (users and products) communicating through gRPC
- Individual PostgreSQL databases for each service
- Consul for service discovery and health checking
- API Gateway serving as the single entry point with HTTP-to-gRPC conversion
- Data aggregation capabilities across multiple services

**Main Achievement:** Successfully fixed the API Gateway to use Consul service discovery instead of hardcoded addresses and implemented proper service-to-service communication.

## Table of Contents

- [Project Overview](#project-overview)
- [Architecture Design](#architecture-design)
- [Technologies Used](#technologies-used)
- [Implementation Details](#implementation-details)
- [Problem Solving](#problem-solving-and-fixes)
- [Testing & Validation](#testing-and-validation)
- [Deployment](#deployment-and-operations)
- [Learning Outcomes](#learning-outcomes)
- [Future Improvements](#future-improvements)

## Project Overview

### What I Built

- **API Gateway** - HTTP server acting as the main entry point
- **Users Service** - Manages user data with dedicated PostgreSQL database
- **Products Service** - Manages product data with dedicated PostgreSQL database
- **Consul** - Service discovery and registration system
- **PostgreSQL Databases** - Separate databases ensuring service independence

### Key Features

- gRPC communication between services
- Service discovery with Consul
- Database per service pattern
- HTTP-to-gRPC translation
- Concurrent data aggregation
- Containerized deployment with Docker

## Architecture Design

### System Architecture

```
┌─────────────────┐    HTTP    ┌──────────────────┐
│   Client/User   │ ────────→ │   API Gateway    │
└─────────────────┘           │   (Port 8080)    │
                               └──────────────────┘
                                        │
                                        │ gRPC
                                        ▼
                               ┌──────────────────┐
                               │     Consul       │
                               │ Service Discovery│
                               │   (Port 8500)    │
                               └──────────────────┘
                                        │
                        ┌───────────────┼───────────────┐
                        │               │               │
                        ▼               ▼               ▼
               ┌─────────────────┐ ┌─────────────────┐
               │  Users Service  │ │Products Service │
               │  (Port 50051)   │ │  (Port 50052)   │
               └─────────────────┘ └─────────────────┘
                        │                       │
                        ▼                       ▼
               ┌─────────────────┐ ┌─────────────────┐
               │   Users DB      │ │  Products DB    │
               │  (Port 5432)    │ │  (Port 5433)    │
               └─────────────────┘ └─────────────────┘
```

### Communication Flow

1. **Client Request** → API Gateway receives HTTP requests
2. **Service Discovery** → API Gateway queries Consul for service locations
3. **gRPC Communication** → API Gateway makes gRPC calls to appropriate services
4. **Database Operations** → Services perform CRUD operations on their databases
5. **Response Aggregation** → API Gateway combines responses and returns JSON

## Technologies Used

- **Go (Golang)** - Primary programming language
- **gRPC & Protocol Buffers** - Efficient service-to-service communication
- **PostgreSQL** - Database management system
- **GORM** - Go ORM for database operations
- **Consul** - Service discovery and health checking
- **Docker & Docker Compose** - Containerization and orchestration
- **Gorilla Mux** - HTTP router for API Gateway

## Implementation Details

### Project Structure

```
practical-three/
├── docker-compose.yml
├── proto/
│   ├── users.proto
│   ├── products.proto
│   └── gen/
│       ├── users.pb.go
│       ├── users_grpc.pb.go
│       ├── products.pb.go
│       └── products_grpc.pb.go
├── api-gateway/
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
├── services/
│   ├── users-service/
│   │   ├── main.go
│   │   ├── go.mod
│   │   ├── go.sum
│   │   └── Dockerfile
│   └── products-service/
│       ├── main.go
│       ├── go.mod
│       ├── go.sum
│       └── Dockerfile
```

### Core Features Implementation

#### 1. Protocol Buffers Definition
- **Users Service**: CreateUser, GetUser operations
- **Products Service**: CreateProduct, GetProduct operations
- Generated Go code for type-safe communication

#### 2. Service Discovery with Consul
```go
func registerServiceWithConsul() error {
    config := consulapi.DefaultConfig()
    consul, err := consulapi.NewClient(config)
    // Registration logic implemented here
}
```

#### 3. Database Integration
- **Users DB**: Stores user information (ID, Name, Email)
- **Products DB**: Stores product information (ID, Name, Price)
- GORM ORM for clean database operations
- Separate databases ensuring service independence

#### 4. API Gateway Endpoints
- `POST /api/users` - Create new user
- `GET /api/users/{id}` - Retrieve user by ID
- `POST /api/products` - Create new product  
- `GET /api/products/{id}` - Retrieve product by ID
- `GET /api/purchases/user/{userId}/product/{productId}` - Aggregate data from both services

#### 5. Data Aggregation
Concurrent processing for efficient multi-service data retrieval:
```go
func getPurchaseDataHandler(w http.ResponseWriter, r *http.Request) {
    var wg sync.WaitGroup
    // Concurrent calls to both services
    wg.Add(2)
    go func() { /* Get user data */ }()
    go func() { /* Get product data */ }()
    wg.Wait()
    // Combine and return aggregated data
}
```
![alt text](<Screenshot from 2025-08-28 14-16-59.png>)

![alt text](<Screenshot from 2025-08-28 14-19-54.png>)

## Problem Solving and Fixes

### Issues Encountered and Resolved

#### 1. Hardcoded Service Addresses
- **Problem**: API Gateway directly connecting to services using hardcoded ports
- **Solution**: Implemented Consul service discovery integration
- **Impact**: Improved scalability and service location flexibility

#### 2. Missing Service Communication
- **Problem**: Composite endpoint not properly aggregating data from multiple services
- **Solution**: Implemented concurrent gRPC calls with proper synchronization
- **Impact**: Enabled efficient multi-service data retrieval

#### 3. Service Registration Issues
- **Problem**: Services not properly registering with Consul
- **Solution**: Added comprehensive service registration and health check logic
- **Impact**: Enabled dynamic service discovery and health monitoring

### Solutions Implemented

- **Consul Integration**: Dynamic service location instead of hardcoded addresses
- **Concurrent Processing**: Goroutines and WaitGroups for parallel service calls
- **Comprehensive Error Handling**: Robust error handling for all service interactions
- **Health Checks**: Proper service health monitoring through Consul

## Testing and Validation

### Test Scenarios

#### 1. User Management
```bash
# Create User
curl -X POST -H "Content-Type: application/json" \
     -d '{"name": "John Doe", "email": "john@example.com"}' \
     http://localhost:8080/api/users

# Get User
curl http://localhost:8080/api/users/1
```

#### 2. Product Management
```bash
# Create Product
curl -X POST -H "Content-Type: application/json" \
     -d '{"name": "Laptop", "price": 1299.99}' \
     http://localhost:8080/api/products

# Get Product
curl http://localhost:8080/api/products/1
```

#### 3. Data Aggregation
```bash
# Get Aggregated Purchase Data
curl http://localhost:8080/api/purchases/user/1/product/1
```

### Test Results

- User and product creation return complete object details
- Retrieval endpoints successfully fetch data from databases
- Aggregated endpoint returns combined user and product information
- Services properly register and are discoverable through Consul
- All HTTP status codes and JSON responses are properly formatted

## Deployment and Operations

### Quick Start

#### Prerequisites
- Docker installed and running
- Docker Compose available
- Ports 8080, 8500, 50051, 50052, 5432, 5433 available

#### Running the System
```bash
# Clone the repository
git clone [your-repository-url]
cd practical-three

# Start all services
docker-compose up --build

# Verify deployment
curl http://localhost:8080/api/users  # API Gateway
open http://localhost:8500           # Consul UI
```

### Container Status Verification

| Service | Port | Status Check |
|---------|------|-------------|
| **Consul** | 8500 | `http://localhost:8500` |
| **API Gateway** | 8080 | `http://localhost:8080/api/users` |
| **Users Service** | 50051 | Internal gRPC |
| **Products Service** | 50052 | Internal gRPC |
| **Users Database** | 5432 | Internal PostgreSQL |
| **Products Database** | 5433 | Internal PostgreSQL |

## Learning Outcomes

### Technical Skills Developed

- **Microservices Architecture**
  - Service decomposition and independence
  - Database per service pattern
  - Inter-service communication strategies

- **gRPC Implementation**
  - Protocol buffer design and code generation
  - Efficient binary communication
  - Type-safe service contracts

- **Service Discovery**
  - Consul integration and configuration
  - Dynamic service location
  - Health check implementation

- **Database Management**
  - Multi-database architecture
  - GORM ORM usage
  - Database isolation strategies

- **Containerization**
  - Docker multi-service setup
  - Docker Compose orchestration
  - Container networking

### Problem-Solving Skills Enhanced

- **Service Integration**: Resolved communication issues between distributed services
- **Concurrent Programming**: Implemented efficient parallel data processing
- **Error Handling**: Added comprehensive error management across service boundaries
- **Debugging**: Diagnosed and fixed service discovery and connectivity issues

## Future Improvements

### Short-term Enhancements
- **Health Checks**: Comprehensive health check endpoints for all services
- **Logging**: Structured logging with correlation IDs for request tracing
- **Configuration**: External configuration management for different environments

### Medium-term Additions
- **Caching**: Redis caching layer for frequently accessed data
- **Authentication**: JWT-based authentication and authorization
- **Rate Limiting**: API rate limiting and throttling mechanisms

### Long-term Scaling
- **Monitoring**: Prometheus and Grafana integration for metrics and alerting
- **Load Balancing**: Multiple service instances with intelligent load distribution
- **Message Queues**: Asynchronous communication with RabbitMQ or Apache Kafka
- **API Versioning**: Backward-compatible API version management

## Conclusion

This practical successfully demonstrates a production-ready microservices implementation following industry best practices:

### Key Achievements
- **gRPC Communication**: Efficient inter-service communication
- **Service Discovery**: Dynamic service location with Consul
- **Database Independence**: Separate databases per service
- **API Gateway**: HTTP-to-gRPC translation and aggregation
- **Containerization**: Docker-based deployment
- **Concurrent Processing**: Parallel data retrieval optimization

### Technical Impact
- **Scalability**: Services can be scaled independently
- **Maintainability**: Clear service boundaries and responsibilities
- **Reliability**: Fault isolation and robust error handling
- **Performance**: Efficient gRPC communication and concurrent processing

### Learning Value
This project provides hands-on experience with:
- Modern microservices architecture patterns
- Service discovery and registry concepts
- Database per service implementation
- Container orchestration with Docker Compose
- gRPC protocol buffer communication

---

**Repository**: [Link to your submission repository]  
**Documentation**: Complete codebase with inline comments  
**Testing**: Postman collection and cURL examples included