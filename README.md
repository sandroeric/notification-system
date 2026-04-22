# Notification System Command Center

A full-stack, containerized Notification Dispatching System built to elegantly scale across multiple delivery channels (SMS, Email, Push). Designed with Domain-Driven/Clean Architecture principles, the platform dynamically resolves the delivery preferences of statically mocked users and tracks its fault-tolerant execution natively within SQLite.

---

## 🏗 Architecture & Tech Stack

This project was built focusing deeply on separation of concerns, scalability, and robust dependency structures:
- **Backend**: Go 1.22 (Standard `net/http` router for zero-bloat)
- **Frontend**: React.js / Vite (Functional components + Hooks + Custom Glassmorphism UI)
- **Database**: SQLite (Automated migrations for Delivery Logging)
- **Orchestration**: Docker & Docker Compose (Multi-stage optimized builds)
- **Core Patterns**: 
  - **Strategy Pattern**: Isolates channel dispatch behaviors (SMS, Email, Push) to allow isolated scaling.
  - **Factory / Registry Pattern**: Dynamic injection of channels preventing the Notification Service code from needing to know what channels exist.
  - **Clean Architecture / Domain-Driven**: Hard borders between Domain Models, Infrastructure (SQLite/Mock data), and Presentation (HTTP REST).

---

## 🚀 Quick Start Guide

### Prerequisites
You only need **Docker** and **Docker Compose** installed on your system. No local Go or Node environments are required.

### Running the Application

1. Open a terminal at the root of the project.
2. Build and launch the multi-container environment:
   ```bash
   docker-compose up --build
   ```
3. Wait momentarily while the Go binaries compile and the React static assets are bundled.
4. Access the Frontend Command Center in your browser:
   **[http://localhost:3000](http://localhost:3000)**

*(Note: The Backend API boots transparently on port `8080`, managed seamlessly via CORS).*

---

## 🎮 Using the Platform

### 1. The Dispatch Form
On the left side of the screen, you will find the **Submission Form**:
- **Category**: Select a topic you wish to broadcast (`Sports`, `Finance`, or `Movies`).
- **Message**: Type your exact notification payload.
- Hit **Dispatch**. 

*The UI prevents empty network submits. Once accepted, it fires the payload directly to the Go REST API.*

### 2. The Delivery Engine (Under the Hood)
When the backend receives the broadcast:
1. It queries predefined **Mock Users** stored exclusively in-memory.
2. If a user is actively subscribed to the dispatched Category, it triggers a broadcast pipeline.
3. For every channel the user opted into (e.g. they requested both *SMS* and *Email*), it looks up the specific notifier strategy via the **Registry**.
4. It attempts simulated delivery using an embedded **Fault-Tolerance loop** (performing standard retries upon simulated failure).

### 3. Log History Dashboard
On the right side of the screen, the **Delivery History** dashboard automatically flashes to reveal your live database records:
- View exactly which users were targeted based on the constraints.
- The dataset is queried and delivered **Newest to Oldest**.
- You can inspect the `Status`, tracking if a delivery executed smoothly (`Success`), or if it blew past the fault-tolerance limits returning `Failed` along with attached Error definitions and Retry loop metrics.

---

## ⚙️ Extending the System

Need to add a new Notification Channel (like *Slack* or *WhatsApp*)? 

Our architectural patterns make this a 3-step, highly scalable operation:
1. Implement the `NotifierStrategy` interface (found in `backend/internal/domain/interfaces.go`).
2. Add your new channel simulation logic inside `backend/internal/notification/your_channel.go`.
3. Register it dynamically inside `backend/cmd/server/main.go` on system Boot:
   ```go
   registry.Register(notification.NewWhatsAppStrategy())
   ```
Zero modifications are required to the core `notification_service.go` logic!
