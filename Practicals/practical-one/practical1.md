# WEB303 Practical 1: Microservices Development Environment Setup and Inter-Service Communication

## Table of Contents

- [Introduction](#introduction)
- [Learning Objectives](#learning-objectives)
- [What I Built](#what-i-built)
- [Part 1: Setting Up My Environment](#part-1-setting-up-my-environment)
- [Part 2: Building My Microservices](#part-2-building-my-microservices)
- [Testing My Application](#testing-my-application)
- [Problems I Faced and How I Solved Them](#problems-i-faced-and-how-i-solved-them)
- [What I Learned](#what-i-learned)
- [Conclusion](#conclusion)

## Introduction

This report documents my work on WEB303 Practical 1, where I set up a complete development environment for microservices and built my first multi-container application. I created two services that talk to each other using gRPC, packaged them in Docker containers, and managed them with Docker Compose.

## Learning Objectives

Through this practical, I achieved the following learning outcomes:

- **Learning Outcome 1:** I now understand the basic concepts of microservices architecture
- **Learning Outcome 2:** I can design and build microservices using gRPC and Protocol Buffers
- **Learning Outcome 6:** I learned how to deploy microservices using Docker Compose

## What I Built

I created a simple but complete microservices system with two components:

- **Time Service** - A service that tells you the current time
- **Greeter Service** - A service that greets people and asks the Time Service for the current time

When someone asks the Greeter Service to say hello, it contacts the Time Service, gets the current time, and returns a personalized greeting with the time included.

## Part 1: Setting Up My Environment

### Installing Go Programming Language

First, I had to install Go on my computer:

- I went to the official Go website (https://go.dev/dl/) and downloaded the installer
- I ran the installer which put Go in the right place on my system
- I tested it worked by opening a terminal and typing:

```bash
go version
go env
```

Both commands worked perfectly, showing me that Go was properly installed

### Installing Protocol Buffers Tools

Next, I needed tools to work with Protocol Buffers (the way services define how to talk to each other):

- I downloaded the Protocol Buffers compiler (protoc) from GitHub
- I installed two Go plugins that help generate Go code from proto files:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

- I had to update my system PATH so these tools could be found

### Installing Docker

Finally, I installed Docker to run my services in containers:

- I downloaded Docker Desktop from their website
- I installed it and started the Docker service
- I tested it by running the hello-world container:

```bash
docker run hello-world
```

It worked perfectly, downloading and running the test container

## Part 2: Building My Microservices

### Creating My Project Structure

I organized my project like this:

```
practical-one/
├── proto/                 # Service definitions
│   ├── gen/              # Generated Go code goes here
│   ├── time.proto        # Time service definition
│   └── greeter.proto     # Greeter service definition
├── time-service/         # Time service code
│   ├── main.go
│   ├── Dockerfile
│   └── go.mod
├── greeter-service/      # Greeter service code
│   ├── main.go
│   ├── Dockerfile
│   └── go.mod
└── docker-compose.yml    # Orchestration file
```

### Defining My Service Contracts

I created two .proto files to define how my services work:

#### Time Service (time.proto):
- Has one function called GetTime
- Takes no input
- Returns the current time as a string

#### Greeter Service (greeter.proto):
- Has one function called SayHello
- Takes a person's name as input
- Returns a greeting message

After writing these files, I generated the Go code by running:

```bash
protoc --go_out=./proto/gen --go_opt=paths=source_relative \
    --go-grpc_out=./proto/gen --go-grpc_opt=paths=source_relative \
    proto/*.proto
```

This created all the Go code I needed to build my services.

### How I Implemented the Services

#### My Time Service

I built this service to be simple and focused:

- It listens on port 50052
- When someone asks for the time, it gets the current time and sends it back
- I included logging so I can see when requests come in

The main parts of my code:
- A server struct that implements the time service
- A GetTime function that returns the current time in RFC3339 format
- A main function that starts the server and listens for requests

#### My Greeter Service

This service is more interesting because it talks to another service:

- It listens on port 50051
- When someone asks it to say hello, it:
  - Connects to the Time Service
  - Asks for the current time
  - Creates a greeting message that includes the time
  - Sends the greeting back
- I included error handling in case the Time Service isn't available

### Containerizing My Services

I created Dockerfile for each service:

- I used multi-stage builds to make smaller container images
- I used Alpine Linux as the base image for security and size
- Each Dockerfile builds the Go application and creates a clean runtime container

### Orchestrating with Docker Compose

I created a docker-compose.yml file to run both services together:

- The Time Service runs internally (no external access needed)
- The Greeter Service exposes port 50051 for external access
- I set up dependencies so the Time Service starts before the Greeter Service
- I used service names for network communication

## Testing My Application

### How I Tested It

I used a tool called grpcurl to test my complete system:

```bash
grpcurl -plaintext \
    -import-path ./proto -proto greeter.proto \
    -d '{"name": "WEB303 Student"}' \
    0.0.0.0:50051 greeter.GreeterService/SayHello
```

### My Results

When I ran the test, I got exactly what I expected:

```json
{
  "message": "Hello WEB303 Student! The current time is 2025-07-24T09:45:00Z"
}
```

I could also see in the Docker logs that:

- The Greeter Service received my request
- It successfully called the Time Service
- The Time Service responded with the current time
- Everything worked together perfectly

### Output Screenshot

![alt text](<Screenshot from 2025-08-25 23-41-00.png>)

![alt text](<Screenshot from 2025-08-25 23-52-58.png>)

## Problems I Faced and How I Solved Them

### Environment Setup Problems

**Problem:** At first, the protoc-gen-go tools weren't found when I tried to generate code.

**How I Fixed It:** I had to add the Go binary directory to my system PATH. I edited my shell profile file to include: 
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Understanding Protocol Buffers

**Problem:** I was confused about the protoc command and all its options.

**How I Fixed It:** I read the documentation carefully and broke down each part of the command. The --go_out and --go-grpc_out options specify where to put the generated files.

### Docker Networking Issues

**Problem:** My Greeter Service couldn't find the Time Service when running in containers.

**How I Fixed It:** I learned that in Docker Compose, services can find each other using their service names. So I used `time-service:50052` as the address instead of `localhost:50052`.

### Go Module Management

**Problem:** I had trouble with import paths when my services tried to use the generated proto code.

**How I Fixed It:** I set up separate Go modules for each service and made sure the import paths matched the go_package option in my proto files.

## What I Learned

### Technical Skills I Gained

- **Microservices Architecture:** I now understand how to break applications into small, independent services
- **gRPC Communication:** I learned how services can talk to each other efficiently using gRPC
- **Protocol Buffers:** I can now define service contracts and generate code from them
- **Docker Containerization:** I know how to package applications into containers
- **Service Orchestration:** I can use Docker Compose to manage multiple containers

### Development Practices I Learned

- **Planning First:** Defining the service contracts before writing code made everything clearer
- **Single Responsibility:** Each service does one thing well
- **Error Handling:** Important to handle cases where services can't communicate
- **Logging:** Essential for understanding what's happening in distributed systems
- **Testing:** Using tools like grpcurl to verify everything works

### Tools I Now Know How to Use

- **Go Programming Language:** For building efficient microservices
- **Protocol Buffers:** For defining how services communicate
- **Docker:** For packaging and running applications
- **Docker Compose:** For managing multiple containers
- **gRPC:** For high-performance service communication

## Conclusion

I successfully completed this practical and learned a lot about microservices development. I built a working system where two services communicate with each other, packaged them in Docker containers, and orchestrated them with Docker Compose.

### What I Accomplished

- Set up a complete development environment for microservices
- Designed service contracts using Protocol Buffers
- Implemented two communicating services in Go
- Containerized both services with Docker
- Orchestrated the system with Docker Compose
- Successfully tested the complete system

### My Key Takeaways

- **Microservices are powerful but complex** - They give you flexibility but require careful design
- **Good tooling is essential** - Tools like Protocol Buffers and gRPC make communication much easier
- **Containers solve deployment problems** - Docker makes it easy to run services consistently
- **Testing is crucial** - You need good tools to verify that distributed systems work correctly

### What I Want to Learn Next

- How to add monitoring and logging to see what's happening in production
- How to handle errors and failures gracefully
- How to deploy these services to Kubernetes
- How to add authentication and security
- How to manage configuration across multiple services

This practical gave me a solid foundation in microservices development, and I'm excited to build more complex systems using these concepts.

---

**Author:** [Kuenzang Rabten]  
**Module:** WEB303 Microservices & Serverless Applications  
**Completed:** [22/08/202]  