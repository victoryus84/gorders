GOrders Backend Module (Go)

Professional backend servier for the GOrders ecosystem, providing a robust API for order management, contract tracking, and legacy ERP integration.
🚀 Overview

The backend is built using Go 1.26, leveraging the Gin web framework for high-performance routing and GORM for seamless PostgreSQL integration. It is designed following Clean Architecture principles to ensure scalability, maintainability, and ease of testing.
📂 Project Structure

Following the Clean Architecture pattern, the code is organized as follows:

    cmd/server/main.go — Entry Point: Initializes configuration, database connections, auto-migrations, and starts the HTTP server.

    internal/models/ — Entities: Core business objects (User, Client, Contract, Product, Order) and GORM tags.

    internal/repository/ — Data Layer: Direct database interactions and CRUD operations using GORM.

    internal/service/ — Domain Logic: Business rules, validations, and orchestration between repositories.

    internal/api/ — Delivery Layer: HTTP handlers, Gin routes registration, and request/response DTOs.

    internal/config/ — Infrastructure: Application configuration management (DSN, JWT secrets, Environment variables).

🏗 Architectural Principles

    Dependency Separation: Clear boundaries between Handlers → Services → Repositories → Models.

    Security: JWT-based authentication. The /login endpoint provides a token required in the Authorization: Bearer <token> header for all protected routes.

    Data Integrity: Automatic schema synchronization via GORM AutoMigrate during startup.

    Legacy Sync: Specialized logic for handling XML data and encoding conversions (Windows-1251 to UTF-8) for 1C:Enterprise integration.

🔐 Core API Endpoints

Base URL: http://localhost:8080/api/v1
Authentication

    POST /login — Authenticates user and returns a JWT.

Client Management

    GET /clients — List all clients (paginated).

    POST /clients — Register a new client (Owner ID automatically linked via JWT).

    GET /clients/:id — Detailed client view including associated contracts.

Contract & Order Management

    POST /clients/:id/contracts — Create a contract for a specific client.

    GET /contracts/:id — Retrieve contract details.

    POST /orders — Place a new order linked to a contract and client.

🛠 Tech Stack
Technology	Purpose
Go 1.26	Core Programming Language
Gin Gonic	HTTP Web Framework
GORM	Object-Relational Mapper (ORM)
PostgreSQL	Primary Relational Database
Docker	Containerization & Deployment
JWT	Secure Authentication
💻 Quick Start (cURL)

1. Login to get Token:
Bash

curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@gorders.com","password":"securepassword"}'

2. Create a Client (using Token):
Bash

curl -X POST http://localhost:8080/api/v1/clients \
  -H "Authorization: Bearer <your_token>" \
  -d '{"name":"Business Partner Ltd","email":"contact@partner.com","address":"Main St 101"}'

📈 Future Roadmap

    [ ] Implementation of golang-migrate for version-controlled DB migrations.

    [ ] Expanded integration for automated synchronization with 1C:Enterprise 7.7/8.x.

    [ ] Comprehensive unit testing for the Service layer.

    [ ] Structured logging with Zap or Zerolog for production monitoring.

📄 License

This project is licensed under the MIT License - see the LICENSE file for details.