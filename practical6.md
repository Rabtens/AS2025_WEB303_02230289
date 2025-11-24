# Practical 6 Report – Comprehensive Testing for Microservices

**Module:** WEB303 – Microservices & Serverless Applications  
**Practical:** 6  
**Title:** Comprehensive Testing for Distributed Microservices  
**Student Perspective Report**

---

## 1. Introduction

This practical focuses on building a complete testing strategy for a microservices-based Student Café system, continuing from Practical 5A. The work completed in this practical helps ensure each microservice behaves correctly on its own, works properly when combined, and functions as a full end-to-end system.

The practical implements:

- Unit Tests for gRPC service methods
- Integration Tests for multiple microservices working together
- End-to-End (E2E) Tests from API Gateway → all services
- Mocking, test isolation, in-memory databases, and coverage reporting
- Makefile automation to support CI/CD pipelines

The goal is to understand real-world testing practices required in microservice deployments, where systems must be reliable and resilient.

---

## 2. Architecture Overview

The microservices system contains four main components:

```
practical6/
├── user-service
├── menu-service
├── order-service
├── api-gateway
└── tests/
```

### 2.1 Microservices

| Service | Responsibilities |
|---------|------------------|
| User Service | Manages student/cafe owner accounts |
| Menu Service | Manages menu items and prices |
| Order Service | Creates orders, validates users and menu items |
| API Gateway | Exposes REST APIs and proxies to gRPC services |

Each service uses:

- Golang
- gRPC for communication
- GORM ORM
- SQLite/PostgreSQL (depending on environment)

---

## 3. Testing Strategy (Testing Pyramid)

The practical applies the standard testing pyramid:

```
GUI Tests (Manual)
        (Few)
      /\
     /  \
    /----\    End-to-End Tests  (10%)
   /      \
  /--------\ Integration Tests  (20%)
 /          \
/------------\ Unit Tests (70%)
```

### Why this matters

- Unit tests are the fastest → most of the test coverage
- Integration tests confirm services work together
- E2E tests ensure the full system behaves as the user expects

This testing structure gives confidence during refactoring, reduces regressions, and prepares the project for CI/CD pipelines.

---

## 4. Unit Testing

Unit tests validate individual functions and gRPC methods inside each service without external dependencies.

Each service has:

```
/service-name/grpc/server_test.go
```

### Tools Used

- Go Testing Framework
- Testify (assert, require, mocking)
- In-memory SQLite (for isolated databases)

### 4.1 User Service Unit Tests

The user service tests cover:

- Creating new users
- Retrieving users
- gRPC error handling
- Validations

#### Test Setup

An in-memory SQLite database is created for each test:

```go
db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"))
db.AutoMigrate(&models.User{})
```

This ensures:

- tests run fast
- no external DB needed
- each test is isolated

#### Highlights of Test Coverage

- Successful user creation
- Creating cafe-owner users
- Getting non-existent users (return gRPC NotFound)
- Table-driven test format

gRPC error code validation example:

```go
st, _ := status.FromError(err)
assert.Equal(t, codes.NotFound, st.Code())
```

### 4.2 Menu Service Unit Tests

The menu service tests check:

- creating menu items
- handling floats (prices)
- validating descriptions and price boundaries

Key detail: using floating-point InDelta checks:

```go
assert.InDelta(t, tc.price, resp.MenuItem.Price, 0.001)
```

This avoids float precision errors.

### 4.3 Order Service Unit Tests (with Mocks)

The order service depends on:

- User Service
- Menu Service

To avoid calling real services, the practical uses testify mocks.

#### Mocking Example

```go
mockUserClient.On("GetUser", mock.Anything, &userv1.GetUserRequest{Id: 1}).
Return(&userv1.GetUserResponse{User: &userv1.User{Id: 1}}, nil)
```

#### Tests covered

| Scenario | Status |
|----------|--------|
| Order created successfully with valid user & menu item | ✓ |
| Reject order if user doesn't exist | ✓ |
| Reject order if menu item is invalid | ✓ |
| Price snapshotting | ✓ |

Mocking ensures unit tests run instantly and are completely isolated.

---

## 5. Integration Testing

Integration tests validate multiple services working together using in-memory gRPC connections (bufconn).

### Why bufconn?

- No network ports required
- No latency
- No Docker needed
- Useful for CI

### Test Architecture

Each service is started in-memory:

```go
listener := bufconn.Listen(1024 * 1024)
server := grpc.NewServer()
```

Tests then connect to the services like real gRPC clients.

### 5.1 Integration Test: Complete Order Flow

This test simulates a realistic workflow:

1. Create user
2. Create menu items
3. Create order
4. Retrieve order
5. Validate prices and quantities

This is a full multi-service workflow (user → menu → order).

Key validations:

- user existence
- menu item existence
- price snapshot
- order item quantities

### 5.2 Validation Integration Test

Covers negative cases:

| Invalid Input | Expected Result |
|---------------|-----------------|
| Create order with non-existent user | gRPC error |
| Invalid menu item in order | gRPC error |
| Empty order items | Validation error |

These tests ensure that the system rejects invalid requests.

---

## 6. End-to-End Testing (E2E)

E2E tests execute the entire system through the API Gateway.

### Workflow:

```
HTTP Request → API Gateway → gRPC Services → Database → Response
```

E2E covers:

- HTTP → gRPC translation
- API Gateway routing
- Authentication / validation
- Multi-service communication
- Final JSON output

These tests are slowest, so few are included.

### Example:

1. Create user (HTTP)
2. Create menu items (HTTP)
3. Place an order (HTTP)
4. Check final order details (HTTP)

---

## 7. Makefile Automation

The Makefile provides commands to simplify testing:

```bash
make test-unit
make test-unit-user
make test-unit-menu
make test-unit-order
make test-integration
make test-e2e
make test-coverage
```

This prepares the system for CI/CD pipelines such as:

- GitHub Actions
- GitLab CI
- Jenkins

---

## 8. Workflow Diagram

```
          ┌──────────────────────┐
          │        API           │
          │       GATEWAY        │
          └──────────┬───────────┘
                     │ HTTP
     ┌───────────────┼────────────────┐
     │               │                │
 gRPC↓             gRPC↓            gRPC↓
┌──────────┐   ┌─────────┐    ┌────────────┐
│UserSvc   │   │MenuSvc  │    │OrderSvc    │
└─────┬────┘   └────┬────┘    └──────┬─────┘
      │ DB          │ DB             │ DB
      ▼             ▼                ▼
 SQLite/Postgres  SQLite/Postgres  SQLite/Postgres
```

Testing layers interact with this architecture as:

- Unit tests → each box individually
- Integration tests → arrows between boxes
- E2E tests → through entire diagram

---

## 9. Conclusion

This practical provided hands-on experience in implementing a complete professional-grade testing system for microservices. Key takeaways:

- How to write isolated unit tests
- Using mocks for service dependencies
- Testing floating-point values safely
- Running integration tests using bufconn
- Designing realistic E2E workflows
- Automating test execution via Makefile
- Importance of the Testing Pyramid in real microservice systems

This practical significantly improved understanding of how reliable distributed systems are tested and prepared for production-level DevOps pipelines.