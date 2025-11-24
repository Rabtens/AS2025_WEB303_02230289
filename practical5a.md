# Practical 5A: gRPC Migration with Centralized Proto Repository

---

## 1. Introduction

This practical (Practical 5A) builds on Practical 5. In Practical 5, the Student Café application was refactored into independent microservices. In this practical, the communication between these services is upgraded from REST/HTTP to gRPC and a centralized protocol buffer repository is created.

By completing this practical, the Student Café system now supports:

- Faster inter-service communication
- Type-safe request/response messages
- Centralized proto version management
- Dual-protocol support (REST + gRPC)

---

## 2. Objective

The objectives of this practical are:

1. Convert internal service-to-service communication from REST to gRPC
2. Create a centralized proto repository for all .proto files
3. Generate Go gRPC code using protoc
4. Update each microservice to run both HTTP and gRPC servers
5. Implement gRPC clients inside services (especially order-service)
6. Fix build issues caused by proto dependencies
7. Run the whole system using Docker and Docker Compose

---

## 3. Why gRPC?

gRPC gives several benefits for microservices:

| Feature | REST | gRPC |
|---------|------|------|
| Message Type | JSON (text) | Protobuf (binary) |
| Speed | Normal | Faster |
| Type Safety | Low | Very High |
| Data Size | Larger | Smaller |
| Internal Services | OK | Excellent |

gRPC is mainly used for internal communication, while REST is still used for public API calls.

---

## 4. System Architecture

### 4.1 Overview

The final architecture contains:

- A central proto repository (`student-cafe-protos`)
- Three services:
  - user-service
  - menu-service
  - order-service
- Each service:
  - Exposes HTTP API for external clients
  - Exposes gRPC API for internal communication

### 4.2 Centralized Proto Repository Structure

```
student-cafe-protos/
├── proto/
│   ├── user/v1/
│   ├── menu/v1/
│   └── order/v1/
├── gen/go/
├── go.mod
└── Makefile
```

This repository acts as the single source of truth for all proto files.

---

## 5. Implementation Summary

### 5.1 Creating the Central Proto Repository

I created a new folder `student-cafe-protos` and inside it:

- Wrote all `.proto` files for user/menu/order services
- Added versioning using folders `/v1/`
- Created a `Makefile` to generate Go code
- Added a `go.mod` file so microservices can import it as a Go module

Running:

```bash
make generate
```

generates all gRPC code into `gen/go`.

### 5.2 Updating Microservices

Each microservice was updated to:

#### A. Import proto module

In each `go.mod`:

```go
require github.com/douglasswm/student-cafe-protos v0.0.0
replace github.com/douglasswm/student-cafe-protos => ../student-cafe-protos
```

The replace line is very important for local development.

#### B. Add gRPC server

Example:

```go
grpcServer := grpc.NewServer()
userv1.RegisterUserServiceServer(grpcServer, &UserServer{})
```

Each microservice now listens on:

- HTTP Port → `808X`
- gRPC Port → `909X`

#### C. Run dual servers

Each service runs both servers:

```go
go startGRPCServer()
startHTTPServer()
```

This allows a gradual migration from REST to gRPC.

### 5.3 Implementing gRPC Clients in Order-Service

Order-service now calls:

- user-service via gRPC
- menu-service via gRPC

Example:

```go
resp, err := GrpcClients.UserClient.GetUser(ctx, &userv1.GetUserRequest{Id: 1})
```

This replaced the old HTTP JSON calls.

---

## 6. Docker Setup

The biggest challenge was Docker failing to copy proto files and modules.

To solve this:

- I changed the build context in docker-compose
- I copied the `student-cafe-protos` folder inside each Docker build
- The `replace` directive allowed Go to use the local proto directory

Result: All services build successfully and communicate via gRPC inside Docker.

---

## 7. Workflow of an Order Creation (After gRPC Migration)

1. Client sends HTTP request to API Gateway
2. API Gateway forwards request to order-service
3. Order-service internally calls:
   - user-service (gRPC → validate user)
   - menu-service (gRPC → fetch price)
4. Order-service stores the order in the database
5. Response sent back to client via REST

This improves performance and reduces errors.

---

## 8. Testing Performed

### 8.1 REST Testing

```bash
POST /api/users
POST /api/menu
POST /api/orders
```

All worked correctly.

### 8.2 gRPC Testing

Using grpcurl:

```bash
grpcurl -plaintext -d '{"id":1}' localhost:9091 user.v1.UserService/GetUser
```

Result: Returned correct user data.

---

## 9. Challenges Faced

| Problem | Cause | Solution |
|---------|-------|----------|
| Proto mismatch | Duplicate files across services | Centralized proto repo |
| Docker build error | Go module not found | Copy proto repo + replace directive |
| Connection refused | Wrong port mapping | Exposed correct ports in compose |
| Proto not updating | Old generated code | `make clean && make generate` |

---

## 10. Key Learnings

- Centralized proto repository prevents sync issues
- gRPC is faster and more structured than REST
- Dual-protocol design is useful during migration
- Docker requires careful context and module setup
- Service-to-service calls become type-safe and reliable

---

## 11. Conclusion

Practical 5A successfully upgraded the Student Café microservices architecture by introducing gRPC with a centralized proto repository. The new system is faster, more efficient, and more maintainable.

This practical gave hands-on experience with:

- Modern microservices communication
- Protocol buffers
- gRPC server & client implementation
- Docker multi-service builds
- Realistic production patterns

Overall, this practical significantly improved the quality and performance of the Student Café application.