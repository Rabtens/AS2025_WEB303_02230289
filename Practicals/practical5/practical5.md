# Practical 5 Report — Refactoring a Monolithic Web Server to Microservices

**Module:** WEB303 – Microservices & Serverless Applications  
**Practical:** 5 – Refactoring a Monolithic Student Café App into Microservices  

---

## 1. Introduction

In this practical, the goal was to convert a working monolithic "Student Café" application into a microservices architecture. The task required building the monolith first, understanding its structure and dependencies, and then slowly breaking it into independent, domain-based services.

The entire practical helped me understand why real companies break monoliths and how to do it safely without breaking the system.

---

## 2. Objectives

The main objectives were:

- Understand the differences between monolithic and microservices architecture.
- Practice Domain-Driven Design (DDD) to identify service boundaries.
- Extract individual features from a monolith into independent services.
- Use service discovery via Consul.
- Orchestrate the services using Docker Compose.
- Prepare the system for future migration toward gRPC and Kubernetes.

---

## 3. Why Refactoring Matters

In the real world, most systems start as monoliths because they are easier to build and deploy. But as applications grow, monoliths become:

- Difficult to scale independently
- Too risky to change (one bug can break everything)
- Hard to maintain and deploy
- A blocker for parallel team development

This practical teaches the safe step-by-step migration approach used by most companies today.

---

## 4. Understanding the Monolithic Architecture

Before refactoring, we built a full monolithic application for the Student Café. This single Go application handled:

- Users
- Menu items
- Orders

Everything was stored in one PostgreSQL database.

### 4.1 Monolithic Architecture Diagram

```
                 ┌──────────────────────────┐
                 │    Student Café App      │
                 │       (Monolith)         │
                 ├──────────┬───────────────┤
                 │ Users    │ Menu          │
                 │ Orders   │ Handlers      │
                 │ Models   │ DB Access     │
                 └──────────┴───────────────┘
                           │
                           ▼
                  ┌─────────────────┐
                  │  PostgreSQL DB  │
                  └─────────────────┘
```

### 4.2 Monolithic Problems

- Tight coupling
- One failure affects entire system
- Code becomes huge and difficult to understand
- Scaling requires scaling the whole system
- Slow deployments

This helped me understand why microservices solve these problems.

---

## 5. Applying Domain-Driven Design (DDD)

To break the monolith, we applied DDD to identify bounded contexts and business capabilities.

### 5.1 Identified Bounded Contexts

| Domain / Context | Description | Becomes a Microservice |
|-----------------|-------------|------------------------|
| User Context | Authentication, profiles | user-service |
| Menu Context | Food items, pricing | menu-service |
| Order Context | Order creation, line items | order-service |

### 5.2 Reasoning for Boundaries

**User Service**
- Changes relate only to authentication
- Independent scaling
- Used by order-service

**Menu Service**
- Frequently updated
- High read traffic
- Needed by order-service for validation

**Order Service**
- Depends on both user and menu
- Best to split last because it has the most dependencies

---

## 6. Building the Monolith (Baseline Work)

Before splitting the system, I created:

- Database models (User, MenuItem, Order, OrderItem)
- Handlers for all routes
- A router using Chi
- A PostgreSQL connection
- Dockerfile + docker-compose.yml

### 6.1 Monolithic Workflow Diagram

```
Client → API → Monolith → DB → Response
```

Everything goes through one binary and one database.

---

## 7. Refactoring Approach: Strangler Fig Pattern

We followed a gradual extraction pattern:

- Build monolith
- Identify boundaries
- Extract one service at a time
- Replace internal database access with API calls
- Update docker-compose to run multiple services
- Introduce service discovery via Consul

This approach avoids breaking the system during refactoring.

---

## 8. Extracting the Menu Service

The first service to extract was menu-service, because:

- It has simple CRUD operations
- It has no dependencies on user-service
- Order-service depends on it (but only for reading)
- It is mostly read-heavy (better for scaling)

### 8.1 Microservice Architecture After Extraction

```
                ┌──────────────────┐
                │  User Service    │
                └──────────────────┘
                         │
                         │
  Client → API Gateway → │
                         │
                ┌──────────────────┐
                │  Menu Service    │
                └──────────────────┘
                         │
                         │
                ┌──────────────────┐
                │ Order Service    │
                └──────────────────┘
```

(Each service has its own DB)

### 8.2 Menu Service: Workflow Diagram

```
Client → Menu-Service API → Menu DB → Response
```

---

## 9. Service Discovery with Consul (Later Stage)

Consul helps each service find other services dynamically.

- Services register → Consul
- Order-service queries → "Where is menu-service?"
- Consul returns → "http://menu-service:8082"

This removes the need for hardcoded URLs.

---

## 10. Updated System Architecture after Extraction

Here is the multi-service architecture after refactoring:

```
                     ┌──────────────────────┐
                     │      API Gateway     │
 Client ───────────▶ │ (Optional in future) │
                     └─────────┬────────────┘
                               │
         ┌─────────────────────┼─────────────────────┐
         ▼                     ▼                     ▼
 ┌─────────────────┐   ┌─────────────────┐   ┌─────────────────┐
 │   User Service   │   │   Menu Service  │   │  Order Service   │
 └─────────┬────────┘   └────────┬────────┘   └────────┬────────┘
           │                     │                     │
           ▼                     ▼                     ▼
    ┌───────────┐        ┌───────────┐         ┌──────────────┐
    │ user_db   │        │ menu_db   │         │ order_db      │
    └───────────┘        └───────────┘         └──────────────┘
```

Each service:

- Has its own database
- Runs as a separate container
- Is discovered via Consul
- Communicates via REST APIs

---

## 11. Workflows After Microservices Migration

### 11.1 User Creation Workflow

```
Client → user-service → user_db → response
```

### 11.2 Menu Browsing Workflow

```
Client → menu-service → menu_db → response
```

### 11.3 Order Creation Workflow

**Steps:**

1. Client sends order request
2. order-service receives request
3. order-service queries user-service to validate user
4. order-service queries menu-service to validate items
5. order-service stores order in its DB
6. Returns final response

**Workflow:**

```
Client
   │
   ▼
Order Service
   │
   ├──→ User Service
   │       │
   │       └──→ user_db
   │
   ├──→ Menu Service
   │       │
   │       └──→ menu_db
   │
   └──→ order_db
```

---

## 12. Docker Compose Setup

Each service runs as:

- Independent container
- Independent database container
- All connected by Docker network
- Consul added for discovery

**Example:**

```yaml
services:
  user-service
  menu-service
  order-service
  consul
  user_db
  menu_db
  order_db
```

---

## 13. Challenges Faced

### 1. Breaking database relationships

The monolith stored everything in one database. After splitting, each service needed its own database, so I had to:

- Remove foreign keys
- Replace them with simple IDs
- Use API calls instead of joins

### 2. Inter-service communication

Order-service required data from:

- user-service
- menu-service

This required careful API calling and error handling.

### 3. Deployment coordination

Running multiple services required:

- A clean Docker Compose file
- Each service exposing the correct ports
- Health checks

### 4. Maintaining backwards compatibility

I had to ensure the app still worked while extracting services.

---

## 14. Learning Outcomes Achieved

After completing the practical, I learned:

### Understanding service boundaries

DDD helped me identify correct microservices.

### Building and refactoring a monolith

I learned how to restructure code safely.

### Microservice orchestration with Docker Compose

I deployed multiple services working together.

### Inter-service communication patterns

I now understand how services validate each other's data.

### The importance of independent databases

Each microservice owns its data.

---

## 15. Output Screenshot

API Testing

![alt text](<Screenshot from 2025-11-26 11-48-16.png>)

Create user 

![alt text](<Screenshot from 2025-11-26 11-48-03.png>)

Create menu

![alt text](<Screenshot from 2025-11-26 11-47-50.png>)

![alt text](<Screenshot from 2025-11-26 11-47-36.png>)

Create order

![alt text](<Screenshot from 2025-11-26 11-47-21.png>)

![alt text](<Screenshot from 2025-11-26 11-46-59.png>)

Inter-Service Communication Logs

![alt text](<Screenshot from 2025-11-26 11-46-30.png>)

# Quick Test Commands

```
# Start services
docker-compose up --build -d

# Test complete flow
curl -X POST http://localhost:8080/api/users \
  -d '{"name": "Alice", "email": "alice@test.com", "is_cafe_owner": false}'

curl -X POST http://localhost:8080/api/menu \
  -d '{"name": "Coffee", "description": "Hot coffee", "price": 3.00}'

curl -X POST http://localhost:8080/api/orders \
  -d '{"user_id": 1, "items": [{"menu_item_id": 1, "quantity": 2}]}'

# Verify
curl http://localhost:8080/api/orders

```

## 16. Conclusion

This practical gave me a strong understanding of how to break a monolithic application into microservices step-by-step. By applying Domain-Driven Design, I selected the correct service boundaries. By gradually extracting services, I minimized risks. I learned how to use Docker Compose to run multiple services and how inter-service communication works in real systems.

The final architecture is modular, scalable, and ready for future enhancements such as API gateway, Consul service discovery, gRPC, and Kubernetes.

This practical greatly improved my confidence in building real microservices architecture.