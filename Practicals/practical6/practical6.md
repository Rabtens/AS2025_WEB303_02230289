# Practical 6: Comprehensive Testing for Microservices 

**Module:** WEB303 – Microservices & Serverless Applications  
**Practical:** 6  
**Title:** Comprehensive Testing for Distributed Microservices

---

## 1. Introduction

This practical focuses on implementing a complete and production-ready testing strategy for a microservices-based Student Café system. Building on Practical 5A, the objective was to verify system correctness at all levels—individual services, combined services, and full end-to-end workflows.

### The work completed includes:

- Unit Tests for isolated gRPC service methods
- Integration Tests validating interactions between multiple microservices
- End-to-End (E2E) Tests executed through the API Gateway
- Mocking, test isolation, and in-memory databases
- Automated test execution with a Makefile

This testing architecture reflects real-world DevOps and microservice deployment environments, where reliability and resilience are essential.

---

## 2. System Architecture Overview

The microservices system is composed of four core services:

```
practical6/
├── user-service
├── menu-service
├── order-service
├── api-gateway
└── tests/
```

### 2.1 Service Responsibilities

| Service | Responsibility |
|---------|----------------|
| User Service | Manages student and café owner accounts |
| Menu Service | Manages menu items, descriptions, and pricing |
| Order Service | Creates orders, validates user & menu, snapshots prices |
| API Gateway | Exposes REST APIs and communicates with all gRPC services |

Each service uses:

- Go (Golang)
- gRPC for communication
- GORM ORM
- SQLite/PostgreSQL depending on test environment

---

## 3. Testing Strategy (Testing Pyramid)

The testing model follows the industry-standard Testing Pyramid:

```
                GUI Tests (Few / Manual)
                    /\
                   /  \
             E2E Tests (10%)
                 /      \
          Integration Tests (20%)
               /            \
      Unit Tests (70%) — Largest, Fastest Layer
```

### Why this structure is important

- Unit tests are fast, isolated, and form the majority
- Integration tests confirm services communicate correctly
- E2E tests validate the entire real-world workflow

This structured pyramid ensures high confidence during development and supports continuous integration pipelines.

---

## 4. Unit Testing

Unit tests focus on individual gRPC service methods, ensuring correctness without external dependencies.

### Tools and Approaches Used

- Go testing framework
- Testify (assert, require, mocking)
- In-memory SQLite databases
- Table-driven test patterns

### 4.1 User Service Unit Tests

**Located in:** `user-service/grpc/server_test.go`

**Key Features:**

In-memory SQLite for clean, isolated test runs:

```go
db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"))
db.AutoMigrate(&models.User{})
```

**Tests for:**

- Successful user creation
- Handling café owner roles
- Retrieving users
- Error cases for non-existent users

gRPC error validation using:

```go
st, _ := status.FromError(err)
assert.Equal(t, codes.NotFound, st.Code())
```

### 4.2 Menu Service Unit Tests

**Focus areas:**

- Creating menu items
- Validating price boundaries
- Testing floating-point precision

Floating-point checks use:

```go
assert.InDelta(t, expectedPrice, actualPrice, 0.001)
```

### 4.3 Order Service Unit Tests (with Mocks)

The order service depends on the user and menu services. To isolate tests, mock clients were used:

```go
mockUserClient.On("GetUser", mock.Anything, &userv1.GetUserRequest{Id: 1}).
Return(&userv1.GetUserResponse{User: &userv1.User{Id: 1}}, nil)
```

**Covered Scenarios:**

- Valid order creation
- Rejecting orders from invalid users
- Rejecting orders with invalid menu items
- Price snapshot verification

Mocking ensures no external service calls are made.

---

## 5. Integration Testing

Integration tests validate service interactions using **bufconn**, an in-memory gRPC transport.

### Why bufconn?

- No network setup required
- Extremely fast
- Ideal for CI
- Zero port conflicts

### Test Setup

Each microservice is launched in-memory:

```go
listener := bufconn.Listen(1024 * 1024)
server := grpc.NewServer()
```

### 5.1 Complete Order Flow Integration Test

Simulates a real multi-service workflow:

1. Create a user
2. Create menu items
3. Create an order
4. Retrieve order
5. Validate all computed fields

### 5.2 Validation and Error Integration Tests

Covers invalid inputs such as:

| Case | Expected Result |
|------|-----------------|
| Non-existent user | gRPC error |
| Invalid menu item | gRPC error |
| Empty order | Validation failure |

These ensure services reject incorrect data.

---

## 6. End-to-End (E2E) Testing

E2E tests run through the API Gateway, covering:

```
HTTP → API Gateway → gRPC Services → DB → Response
```

### This validates:

- HTTP routing
- JSON serialization
- gRPC translation
- Multi-service interaction
- Full business logic behavior

### Example E2E Scenario

1. Create user (HTTP)
2. Add menu items (HTTP)
3. Place order (HTTP)
4. Get final order (HTTP)

These tests simulate real user behavior.

---

## 7. Makefile Automation

To support CI/CD pipelines, the following commands were created:

- `make test-unit`
- `make test-unit-user`
- `make test-unit-menu`
- `make test-unit-order`
- `make test-integration`
- `make test-e2e`
- `make test-coverage`

These allow consistent, repeatable test runs.

---

## 8. Workflow Diagram

```
          ┌──────────────────────┐
          │        API           │
          │       GATEWAY        │
          └──────────┬───────────┘
                     │
             HTTP Request Layer
                     │
      ┌──────────────┼────────────────┐
      │              │                │
 gRPC │         gRPC │           gRPC │
      ▼              ▼                ▼
┌──────────┐   ┌─────────┐    ┌────────────┐
│ UserSvc  │   │ MenuSvc │    │ OrderSvc   │
└─────┬────┘   └────┬────┘    └──────┬─────┘
      │ DB          │ DB             │ DB
      ▼             ▼                ▼
 SQLite/Postgres  SQLite/Postgres  SQLite/Postgres
```

### Testing layers:

- **Unit tests:** Each box individually
- **Integration tests:** Arrows between services
- **E2E tests:** Entering from the API Gateway

---

# 9. Output Screenshot

User Service Unit Tests

![alt text](<Screenshot from 2025-11-26 20-48-17.png>)

Menu Service Unit Tests

![alt text](<Screenshot from 2025-11-26 20-48-43.png>)

Order Service Unit Tests

![alt text](<Screenshot from 2025-11-26 20-49-01.png>)

Terminal Coverage Reports

![alt text](<Screenshot from 2025-11-26 20-51-03.png>)

![alt text](<Screenshot from 2025-11-26 20-51-48.png>)

![alt text](<Screenshot from 2025-11-26 20-52-03.png>)

Integration Testing

![alt text](<Screenshot from 2025-11-26 21-09-31.png>)

HTML Coverage Reports

![alt text](<Screenshot from 2025-11-26 21-15-45.png>)

![alt text](<Screenshot from 2025-11-26 21-15-58.png>)

## 10. Conclusion

This practical provided hands-on experience with designing and executing a complete microservices testing strategy. It demonstrated:

- Isolated and realistic unit tests
- Mocking for dependent services
- Accurate floating-point validation
- Efficient integration tests using bufconn
- Robust E2E workflows through the API Gateway
- Test automation suitable for CI/CD pipelines