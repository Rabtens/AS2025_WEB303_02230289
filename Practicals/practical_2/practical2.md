# WEB303: Microservices & Serverless Applications

## Practical 2: API Gateway with Service Discovery

---

## Table of Contents

- [Overview](#overview)
- [Learning Outcomes](#learning-outcomes)
- [System Architecture](#system-architecture)
- [Tools and Technologies](#tools-and-technologies)
- [Project Structure](#project-structure)
- [Setup and Installation](#setup-and-installation)
- [Implementation](#implementation)
  - [Part 1: Users Service](#part-1-users-service)
  - [Part 2: Products Service](#part-2-products-service)
  - [Part 3: API Gateway](#part-3-api-gateway)
- [Testing the System](#testing-the-system)
- [Troubleshooting](#troubleshooting)
- [Demonstrating Resilience](#demonstrating-resilience)
- [Required Screenshots](#required-screenshots)
- [Key Takeaways](#key-takeaways)
- [Author](#author)

---

## Overview

This practical demonstrates how to build a small **microservices ecosystem** using **Go (Golang)**, **Consul** for service discovery, and an **API Gateway** for dynamic routing.

### Goal

Create two independent microservices:
- `users-service`
- `products-service`

Both services will register themselves with Consul, which acts as a **service registry**. The **API Gateway** will then route incoming client requests to the correct microservice dynamically by checking Consul for available, healthy services.

This setup shows how modern systems can be made **scalable, fault-tolerant, and easily maintainable** without manual configuration changes.

---

## Learning Outcomes

This practical supports the following module outcomes:

- **Outcome 2:** Design and implement microservices using gRPC and Protocol Buffers for efficient inter-service communication.
- **Outcome 8:** Implement observability solutions for microservices and serverless applications, including distributed tracing, metrics, and logging.

---

## System Architecture

### Components

#### 1. Users Service
- Provides user data (e.g., `/users/{id}`)
- Registers itself to Consul and responds to health checks

#### 2. Products Service
- Provides product data (e.g., `/products/{id}`)
- Also registers to Consul and responds to health checks

#### 3. Consul (Service Registry)
- Keeps track of all running services and their health status
- Provides service discovery for the gateway

#### 4. API Gateway
- Acts as a single entry point for all client requests
- Uses Consul to find and forward requests to the correct microservice

### Architecture Diagram (Conceptual)

```
Client → API Gateway → Consul → {users-service / products-service}
```

---

## Tools and Technologies

| Tool / Technology | Purpose |
|-------------------|---------|
| **Go (Golang)** | For building the microservices and API gateway |
| **Docker / Consul** | Service registry for service discovery |
| **Chi Router** | Lightweight HTTP router for building REST endpoints |
| **Postman / cURL** | For testing API endpoints |
| **VS Code / GoLand** | Recommended IDEs for development |

---

## Project Structure

```
go-microservices-demo/
├── api-gateway/
│   └── main.go
└── services/
    ├── products-service/
    │   └── main.go
    └── users-service/
        └── main.go
```

---

## Setup and Installation

### Prerequisites

- Go installed on your system
- Docker installed (for running Consul)
- Basic understanding of REST APIs
- Terminal/Command Line access

### Step 1: Install Go

1. Download Go from the official site: [https://golang.org/dl/](https://golang.org/dl/)
2. Follow installation instructions for your operating system
3. Verify installation:

```bash
go version
```

### Step 2: Install Docker

1. Install Docker Desktop or Docker CLI based on your OS
2. Verify installation:

```bash
docker --version
```

### Step 3: Run Consul

Run Consul in development mode using Docker:

```bash
docker run -d -p 8500:8500 --name=consul hashicorp/consul agent -dev -ui -client=0.0.0.0
```

Access Consul UI at: [http://localhost:8500](http://localhost:8500)

---

## Implementation

### Part 1: Users Service

**File Location:** `services/users-service/main.go`

#### Features:
- Runs on port **8081**
- Registers with Consul
- Exposes endpoints:
  - `/health` - For Consul health checks
  - `/users/{id}` - Returns mock user details

#### Running the Service:

```bash
cd services/users-service
go run .
```

#### Expected Output:

```
Successfully registered 'users-service' with Consul
'users-service' starting on port 8081...
```

---

### Part 2: Products Service

**File Location:** `services/products-service/main.go`

#### Features:
- Runs on port **8082**
- Registers with Consul
- Exposes endpoints:
  - `/health` - For Consul health checks
  - `/products/{id}` - Returns mock product details

#### Running the Service:

```bash
cd services/products-service
go run .
```

#### Expected Output:

```
Successfully registered 'products-service' with Consul
'products-service' starting on port 8082...
```

---

### Part 3: API Gateway

**File Location:** `api-gateway/main.go`

#### Features:
- Runs on port **8080**
- Forwards requests to correct microservices using Consul service discovery
- Routes:
  - `/api/users/*` → forwards to users-service
  - `/api/products/*` → forwards to products-service

#### Running the Gateway:

```bash
cd api-gateway
go run .
```

#### Expected Output:

```
API Gateway starting on port 8080...
```

---

## Testing the System

### Setup

Open **four terminals** for the following:

1. **Terminal 1:** Consul (via Docker or local installation)
2. **Terminal 2:** Users Service (`go run .`)
3. **Terminal 3:** Products Service (`go run .`)
4. **Terminal 4:** API Gateway (`go run .`)

### Test Cases

#### Test 1: Users Service (through Gateway)

```bash
curl http://localhost:8080/api/users/123
```

**Expected Response:**

```
Response from 'users-service': Details for user 123
```

#### Test 2: Products Service (through Gateway)

```bash
curl http://localhost:8080/api/products/abc
```

**Expected Response:**

```
Response from 'products-service': Details for product abc
```

---

## Troubleshooting

### Common Issue: Address Mismatch

**Problem:**
Services register to Consul running inside Docker, but the gateway cannot reach them (address mismatch).

#### Solution Option 1: Run Consul Locally

1. Stop Docker Consul container:
   ```bash
   docker stop consul
   ```

2. Install Consul locally: [https://developer.hashicorp.com/consul/install](https://developer.hashicorp.com/consul/install)

3. Start Consul in dev mode:
   ```bash
   consul agent -dev
   ```

4. Visit [http://localhost:8500/ui](http://localhost:8500/ui)

#### Solution Option 2: Containerize All Services

1. Create Dockerfiles for each service
2. Use Docker Compose to define all services
3. Run all together so they share the same network

---

## Demonstrating Resilience

### Testing Service Recovery

1. **Stop the Users Service:**
   - Press `Ctrl + C` in the users-service terminal

2. **Check Consul UI:**
   - The service will turn red (unhealthy)

3. **Test the Endpoint:**
   ```bash
   curl http://localhost:8080/api/users/123
   ```

   **Expected Output:**
   ```
   no healthy instances of service 'users-service' found in Consul
   ```

4. **Restart the Service:**
   ```bash
   cd services/users-service
   go run .
   ```

5. **Verify Recovery:**
   - Health status turns green in Consul UI
   - Service works again without restarting the gateway

**This demonstrates automatic recovery and dynamic routing powered by Consul.**

---

## Required Screenshots

Include the following screenshots in your submission:

- [ ] Consul UI showing both `users-service` and `products-service` as healthy
- [ ] Postman or cURL output for `/api/users/{id}` request
- [ ] Postman or cURL output for `/api/products/{id}` request
- [ ] API Gateway terminal showing request logs (received and forwarded)

---

## Key Takeaways

- Consul enables dynamic service discovery and health monitoring
- API Gateway acts as a central routing point for microservices
- Loose coupling allows services to be updated or restarted without downtime
- Go and Docker provide a lightweight environment for building scalable microservices systems

---

## GitHub Repository

**Repository URL:** [Insert your public GitHub repository link here]

---

## Author

| Field | Details |
|-------|---------|
| **Name** | Kuenzang Rabten |
| **Module** | WEB303 - Microservices & Serverless Applications |
| **Practical** | 2 - API Gateway with Service Discovery |

---
