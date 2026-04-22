# Notification System Command Center

A full-stack, containerized Notification Dispatching System built to elegantly scale across multiple delivery channels (SMS, Email, Push). Designed with Clean Architecture principles, the platform utilizes a robust relational database to manage users, categories, and subscription preferences with localized fault-tolerant execution natively within SQLite.

---

## 🏗 Architecture & Tech Stack

This project was built focusing deeply on separation of concerns, scalability, and robust dependency structures:
- **Backend**: Go 1.23+ (Refined Clean Architecture with strictly decoupled layers)
- **Frontend**: React.js / Vite (Vitest + React Testing Library + Glassmorphism UI)
- **Database**: SQLite (Relational schema with Foreign Keys, Migrations, and Automated Seeding)
- **Orchestration**: Docker & Docker Compose (Multi-stage optimized builds)
- **Core Patterns**: 
  - **Strategy Pattern**: Isolates channel dispatch behaviors (SMS, Email, Push).
  - **Registry Pattern**: Dynamic discovery of delivery strategies.
  - **Repository Pattern**: Abstracted data access for both Notification Logs and User Metadata.
  - **Clean Architecture**: Domain-driven borders between Application logic, Infrastructure, and Presentation.

---

## 🚀 Quick Start Guide

### Prerequisites
You only need **Docker** and **Docker Compose** installed on your system.

### Running the Application

1. Open a terminal at the root of the project.
2. Build and launch the multi-container environment:
   ```bash
   docker-compose up --build
   ```
3. The system will automatically:
   - Run relational SQLite migrations ensuring Foreign Key integrity.
   - Seed the database with initial Categories, Channels, and Users.
   - Start the Go API on port `8080` and the React UI on port `3000`.

4. Access the Command Center: **[http://localhost:3000](http://localhost:3000)**

---

## 🎮 Using the Platform

### 1. The Dispatch Form
On the left side of the screen, you will find the **Submission Form**:
- **Category**: Select a topic (`Sports`, `Finance`, or `Movies`).
- **Message**: Type your notification payload.
- Hit **Dispatch**. 

### 2. The Delivery Engine (Fault Tolerance)
When the backend receives the broadcast:
1. It queries the **SQLite Repository** to find users matching the subscription.
2. For every user, it identifies their preferred channels.
3. It resolves the specific notifier strategy via the **Registry**.
4. It performs delivery using an **execute-with-retry** logic (3 attempts) to handle transient failures.
5. Errors are aggregated using `errors.Join`. If some notifications fail but others succeed, the API returns a **207 Multi-Status** to preserve visibility into partial successes.
6. Every attempt and final result is persisted to the `notification_logs` table.

### 3. Log History Dashboard
On the right side, the **Delivery History** reveals live records:
- Inspect who received what, through which channel, and when.
- Track fine-grained metrics including `RetryCount` and raw `ErrorMessages` for failed attempts.

---

## ⚙️ Extending the System

Our architectural patterns make adding new features a simple, isolated operation:

### Adding a New Channel (e.g. Slack)
1. Implement the `NotifierStrategy` interface in `internal/domain/interfaces.go`.
2. Register the strategy in `cmd/server/main.go`.
3. The system will automatically begin routing messages to Slack for any user configured with that channel id in the DB.

## 🧪 Testing

We maintain a high standard of quality with comprehensive testing across both tiers:

### Backend (Go)
```bash
cd backend && go test ./... -v -cover
```
- **Internal/Notification**: 100% Coverage
- **Application/Infrastructure**: Mock-driven isolation

### Frontend (Vitest)
```bash
cd frontend && npm run test
```
- Component testing for `SubmissionForm` and `LogHistory`.
- Mocking of API fetch responses and error boundary verification.
