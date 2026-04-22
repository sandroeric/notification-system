# Implementation Plan: Notification System

This document outlines the architecture, design patterns, and structure to fulfill the "Software Engineer Challenge" requirements. The solution involves a containerized full-stack application leveraging Go (backend), React.js (frontend), and SQLite (database).

## Proposed Architecture & Tech Stack

- **Backend**: Go (Golang)
- **Frontend**: React.js (Functional components, React Hooks) via Vite
- **Database**: SQLite (Local persistence for notification logs only)
- **Orchestration**: Docker & Docker Compose (Multi-container setup)
- **Patterns**: Clean Architecture / Domain-Driven Design, Dependency Injection, Strategy Pattern, Factory/Registry Pattern

---

## Execution Phases & Checklist

### Phase 1: Project Initialization & Infrastructure
This phase sets up the main folder structure, container orchestrations, and base configs.

- [x] Initialize Go modules and basic directory structure (`cmd/`, `internal/`).
- [x] Initialize React frontend using Vite (React + Javascript).
- [x] Create `docker-compose.yml` to orchestrate backend (port 8080) and frontend (port 3000).
- [x] Write `Dockerfile` for the Go backend (Multi-stage build).
- [x] Write `Dockerfile` for the React frontend (Node build + Nginx).

### Phase 2: Domain & Persistence Layers
Setting up the pure domain models, the mock users, and the SQLite repository.

- [x] **Domain Models**: Create `user.go` (mock domain entity) and `notification.go` (delivery log entity).
- [x] **Interfaces**: Define `NotifierStrategy`, `NotificationRegistry`, and `NotificationRepository`.
- [x] **Mock Users**: Implement `mock_users.go` with hardcoded users mapped to categories and channels (no DB Foreign Key constraint).
- [x] **SQLite Repository**: Implement `sqlite_repo.go` to interact with the database.
- [x] **DB Migrations**: Write schema for `notification_logs` including `delivery_status`, `error_message`, and `retry_count`.
- [x] **DB Indexes**: Apply an index on `timestamp` for optimal newest-to-oldest sorting.

### Phase 3: Core Logic & Channels (Strategy Pattern)
Implementing the registry and the specific sender logic.

- [x] **Registry**: Implement `registry.go` (Factory pattern) to register and resolve notification strategies.
- [x] **SMS Notifier**: Implement `sms_notifier.go` with simulated logic and fault tolerance (retries).
- [x] **Email Notifier**: Implement `email_notifier.go` with fault tolerance.
- [x] **Push Notifier**: Implement `push_notifier.go` with fault tolerance.
- [x] **Notification Service**: Implement `notification_service.go`. Logic: Resolve subscribed users, get needed channels, execute sending with retries, and save result logs to SQLite.
- [x] **Service Unit Tests**: Write table-driven tests for the service layer.

### Phase 4: API & Presentation Layer
Exposing the backend to the frontend.

- [x] **HTTP Handler**: Implement `http_handler.go` to handle REST routes matching the React application's needs.
- [x] **POST /api/notifications**: Endpoint to receive the category and message payload (includes empty payload validation).
- [x] **GET /api/notifications/log**: Endpoint to fetch logs sorted newest to oldest.
- [x] **Server Setup**: Hook up the router, inject dependencies, and run server in `main.go`.

### Phase 5: Frontend Layout & Formulation
Building the React views.

- [x] **Styling Foundation**: Setup base `index.css` for a clean, modern aesthetic.
- [x] **Main Layout**: Set up `App.jsx` structure.
- [x] **Submission Form**: Implement `SubmissionForm.jsx` (Dropdown for Movies/Finance/Sports, Textarea for message).
- [x] **Form Validation**: Client-side rejection of empty messages.
- [x] **API Tie-in**: Wire the form to accurately trigger the backend POST endpoint.

### Phase 6: Frontend Log History
Validating the delivery and showing records.

- [ ] **Log History Component**: Implement `LogHistory.jsx`.
- [ ] **Log Fetching**: Fetch data from `GET /api/notifications/log` on mount and after successful form submissions.
- [ ] **Sorting Requirements**: Ensure logs display properly ordered from newest to oldest.
- [ ] **Review**: End-to-end manual test confirming that mock users receive appropriate messages.

---

## Final Review Comments
This phased approach incorporates:
1. **Clean Architecture Terminology** naturally isolated by phase.
2. **Explicit User Mocking Strategy**, handled in Phase 2.
3. **Database Polish** via indexing constructed in Phase 2.
4. **Resiliency** implemented via mock-retries in Phase 3.
5. **Channel Scalability** heavily grounded in a Notification Registry (Factory) constructed in Phase 3.
