# Practical 5B – Complete gRPC Migration (HTTP Gateway → Pure gRPC Backend)

**WEB303:** Microservices & Serverless Applications  
**Student Report (README Format)**

---

## 1. Introduction

This practical focuses on completing the migration of our Student Café application from:

**A hybrid architecture (HTTP + gRPC)**

to

**A fully gRPC-based backend with an HTTP → gRPC translation layer.**

In Practical 5A, we added gRPC for inter-service communication but still kept HTTP servers in each service. In Practical 5B, we remove all HTTP servers from backend services and let the API Gateway become the only HTTP entry point while internally communicating using pure gRPC.

This reflects real-world architectures used by Google, Uber, Netflix, etc.

---

## 2. Objective

The goals of Practical 5B are:

- Convert the API Gateway from an HTTP reverse proxy into a proper gRPC client.
- Remove all HTTP endpoints from backend services (User, Menu, Order).
- Keep external clients using HTTP/REST but internally use only gRPC.
- Learn protocol translation from HTTP → gRPC and gRPC → HTTP.
- Build a production-like microservices setup.

---

## 3. Architecture Overview

### 3.1 Architecture Before (Practical 5A – Hybrid)

```
Clients (HTTP)
    |
API Gateway (HTTP Reverse Proxy)
    |
Backend Services = HTTP + gRPC
```

**Problems:**

- Services had to run two servers (HTTP & gRPC).
- Code duplication.
- Harder to maintain.

### 3.2 Architecture After (Practical 5B – Pure gRPC Backend)

```
Clients (HTTP/REST)
         |
API Gateway (HTTP → gRPC translator)
         |
Backend Microservices (gRPC only)
         |
gRPC Inter-Service Communication
```

**Benefits:**

- Backend services become simpler (only gRPC).
- Gateway handles all HTTP logic.
- Strong type-safety (Protocol Buffers).
- Higher performance (binary protocol).
- Maintains backward compatibility with existing HTTP clients.

---

## 4. System Workflow (End-to-End)

### Example: Creating an Order

1. A client sends an HTTP POST request to the gateway:

```bash
POST /api/orders
```

2. **API Gateway:**
   - Parses JSON
   - Converts it into a gRPC request
   - Calls OrderService via gRPC

3. **OrderService:**
   - Calls UserService (gRPC) to check if user exists
   - Calls MenuService (gRPC) to fetch item prices
   - Calculates total price
   - Saves order in the database

4. OrderService returns a gRPC response to the gateway

5. Gateway converts gRPC → JSON and returns HTTP 201 Created

This full pipeline is now gRPC inside, HTTP outside.

---

## 5. What I Implemented in This Practical

### 1. Converted API Gateway to use gRPC clients

The gateway now imports generated proto files and creates gRPC client connections to:

- User Service
- Menu Service
- Order Service

### 2. Implemented HTTP → gRPC translation

Each HTTP endpoint now:

- Reads HTTP JSON request
- Converts to gRPC protobuf type
- Calls backend services over gRPC
- Converts protobuf response into JSON
- Returns HTTP response

### 3. Added gRPC → HTTP error mapping

gRPC status codes are mapped to correct HTTP codes (e.g., NotFound → 404).

### 4. Removed HTTP servers from User, Menu, Order services

Now each backend service only:

- Starts a gRPC server
- Does not run any HTTP endpoints

### 5. Updated Docker configuration

- Only gRPC ports (9091–9093) are exposed for services.
- Gateway remains on HTTP port 8080.

---

## 6. API Gateway Architecture

The Gateway Now Performs:

- Request parsing
- JSON ↔ Protobuf conversion
- gRPC client calls
- Error translation
- Response formatting

This promotes proper separation of concerns.

---

## 7. Backend Service Architecture (Simplified)

Every service now contains:

- gRPC server
- Service implementation (business logic)
- Database layer

All HTTP code, routers, handlers have been removed.

---

## 8. Testing the Architecture

### Test 1: Create a User

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","is_cafe_owner":false}'
```

**Expected:**

- The gateway receives JSON
- Calls user-service via gRPC
- Returns a JSON response

### Test 2: Verify internal communication

Try HTTP on gRPC port (should fail):

```bash
curl http://localhost:9091
```

Try gateway (should work):

```bash
curl http://localhost:8080/api/users
```

### Test 3: Create an order (Full workflow)

```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{ "user_id": 1, "items": [{"menu_item_id": 1, "quantity": 2}] }'
```

**Behind the scenes:**

- Gateway → OrderService (gRPC)
- OrderService → UserService (gRPC)
- OrderService → MenuService (gRPC)
- All internal calls are gRPC only.

---

## 9. Troubleshooting Summary

### Gateway cannot connect

Check:

```bash
docker-compose logs api-gateway
```

### Proto import errors

Fix: Update go.mod replace directive.

### Accidentally still running HTTP servers

Search for:

```go
startHTTPServer
```

Remove them.

---

## 10. Key Learnings (Student Reflection)

### 1. API Gateway acts as a protocol translator

This is common practice at Google and Uber.

### 2. Removing HTTP from backend simplifies the entire system

Services become smaller, cleaner, and easier to maintain.

### 3. gRPC is much more efficient

Fast binary communication, strong contracts via Protobuf.

### 4. Backward compatibility is preserved

Clients still use HTTP — no breaking changes.

### 5. This practical teaches real industry migration patterns

Hybrid → Full gRPC is a common evolutionary path.

---

## 11. Conclusion

Practical 5B successfully completes the migration of the Student Café system into a modern, production-style microservices architecture with:

- A clean gRPC-only backend
- A powerful HTTP → gRPC API Gateway
- Reduced complexity in services
- Improved performance and structure
- Real-world architectural patterns

This practical gave hands-on experience with protocol translation, microservice simplification, gRPC contracts, and communication reliability — all essential skills for professional backend development.