# Code Review — Notification System vs. ADR Requirements

> Reviewed every source file against every section of [docs/adr.md](file:///home/sandro/test/docs/adr.md).  
> Severity legend: 🔴 Critical · 🟠 Major · 🟡 Minor · ✅ Pass

---

## Scorecard

| ADR Criteria | Verdict | Issues |
|---|---|---|
| Best practices (validations, exceptions, naming, OOP) | 🟠 Partial | 3 issues |
| SOLID principles | 🟠 Partial | 2 issues |
| Design patterns (Strategy + Registry) | ✅ Pass | — |
| Architecture (folder structure, separation) | ✅ Pass | 1 minor |
| Unit testing | 🟠 Partial | 2 issues |
| Database (migrations, indexing, types) | ✅ Pass | 1 minor |
| Challenge fulfillment (fault tolerance, scalability) | 🟠 Partial | 2 issues |
| User Interface (form + log history) | 🟠 Partial | 1 issue |

---

## Detailed Findings

### 1. 🔴 Retry count is always hardcoded to 3 — never reflects actual attempts

[notification_service.go:48](file:///home/sandro/test/backend/internal/application/notification_service.go#L46-L55)

```go
retryCount := 3 // Simplified for metric visibility
```

The `retryCount` in the log is always `3`, regardless of whether the send succeeded on the 1st attempt or failed after 3. The `executeWithRetry` method does not return the actual attempt count. The inline comment even acknowledges the issue (*"Wait, actually we can track actual tries"*), meaning this was left as known tech debt.

> [!IMPORTANT]
> The ADR explicitly requires logging **"any other pertinent information"** for delivery verification. A misleading retry count undermines log integrity and the fault-tolerance narrative the evaluators look for.

**Fix:** Return `(attempts int, err error)` from `executeWithRetry` and use the real value.

---

### 2. 🔴 No input validation for `Category` — accepts arbitrary strings

[http_handler.go:44](file:///home/sandro/test/backend/internal/presentation/http_handler.go#L44)

```go
err := h.Service.Send(r.Context(), domain.Category(req.Category), req.Message)
```

The handler checks that `Category` is not empty, but **never validates it** against the known set (`Sports`, `Finance`, `Movies`). A POST with `{"category":"Hacking", "message":"hi"}` will silently persist invalid logs and send notifications to zero users — with no error returned.

> [!WARNING]
> The ADR specifies exactly three categories. Best practices criteria explicitly call out **validations**. This is a clear gap.

**Fix:** Add a `ValidCategories` set and reject unknown categories with `400 Bad Request`.

---

### 3. 🟠 `Send()` always returns `nil` — errors are swallowed

[notification_service.go:75](file:///home/sandro/test/backend/internal/application/notification_service.go#L75)

```go
return nil
```

Even when every single delivery fails, `Send()` returns `nil`. The handler therefore always responds `200 OK` to the client. The frontend has no way to know that all notifications failed.

**Fix:** Collect errors and return a summary error (or at minimum return an error if *all* deliveries failed).

---

### 4. 🟠 `HTTPHandler` depends on concrete `*application.NotificationService`, not an interface

[http_handler.go:12-13](file:///home/sandro/test/backend/internal/presentation/http_handler.go#L12-13)

```go
type HTTPHandler struct {
    Service *application.NotificationService
}
```

This breaks **Dependency Inversion** (the `D` in SOLID). The presentation layer directly imports and couples to the application layer concrete type. It also makes the handler untestable without spinning up a real service.

> [!NOTE]
> The ADR evaluates: *"use of interfaces, inversion of dependencies"*. The domain layer defines interfaces beautifully — but this last mile in the handler undoes the pattern.

**Fix:** Define a `NotificationServiceInterface` (in domain or application) with `Send` and `GetLogs`, and depend on it in the handler.

---

### 5. 🟠 Test coverage is narrow — handler, repository, and registry are untested

The ADR states: *"Tests for each service and each function, multiple test scenarios per function."*

| Layer | Test file | Verdict |
|---|---|---|
| `application/` | [notification_service_test.go](file:///home/sandro/test/backend/internal/application/notification_service_test.go) | ✅ 3 tests (happy, retry, failure) |
| `presentation/` | — | ❌ No tests |
| `infrastructure/` | — | ❌ No tests |
| `notification/` | — | ❌ No tests (registry, strategies) |

Missing scenarios include:
- Handler: invalid JSON, empty category, empty message, method-not-allowed
- Repository: `Save` error handling, `FindAll` with populated data
- Registry: duplicate registration, `Get` for missing channel
- Service: unsubscribed user filtering, no channels configured, empty user list

---

### 6. 🟠 Tests import concrete `notification.NewRegistry()` instead of using a mock

[notification_service_test.go:58](file:///home/sandro/test/backend/internal/application/notification_service_test.go#L58)

```go
reg := notification.NewRegistry()
```

The service tests create mock strategies and mock repos (good), but use the **real** `Registry` from the `notification` package. This cross-package dependency means the application tests are no longer true unit tests — they also exercise the registry implementation.

**Fix:** Create a `MockRegistry` that implements `domain.NotificationRegistry`.

---

### 7. 🟡 `User.Name` is missing from `NotificationLog` — can't identify users in the log

The ADR logging section requires **"User data"** to be stored. The log captures `user_id`, `user_email`, and `user_phone` — but **not the user's name**. The log table and the `NotificationLog` struct both omit `UserName`.

The frontend log table column header says "User Config" and only shows email + phone, making it harder for evaluators to verify *who* received the notification.

---

### 8. 🟡 `App.css` is leftover Vite boilerplate — dead code

[App.css](file:///home/sandro/test/frontend/src/App.css) contains the default Vite scaffold styles (`.logo`, `.read-the-docs`, `logo-spin` animation) that are never referenced anywhere in the application. It's imported nowhere in the current code but remains in the repo.

**Fix:** Delete `App.css` or clean it up.

---

### 9. 🟡 Frontend hardcodes `localhost:8080` — breaks in Docker

[SubmissionForm.jsx:3](file:///home/sandro/test/frontend/src/components/SubmissionForm.jsx#L3) and [LogHistory.jsx:3](file:///home/sandro/test/frontend/src/components/LogHistory.jsx#L3)

```js
const API_URL = 'http://localhost:8080/api/notifications';
```

When running via `docker-compose`, the frontend is served from nginx on port 3000. The browser will call `localhost:8080` directly — which works only if the backend port is also exposed to the host. This is fragile; if the backend container is only accessible to the Docker network, the frontend breaks.

**Fix:** Use a relative URL (`/api/notifications`) and add an nginx reverse-proxy config so `/api/*` routes to the backend container.

---

### 10. 🟡 `docker-compose.yml` uses deprecated `version` key

```yaml
version: "3.8"
```

The `version` key is [deprecated since Docker Compose v2](https://docs.docker.com/compose/releases/migrate/) and generates a warning. Not a functional issue, but signals outdated knowledge to evaluators.

---

### 11. 🟡 No `UserName` displayed in Log History UI

The ADR requires logging **"User data"**. The frontend table shows email and phone but not the user's name. Even if the backend stored the name, the log table has no column for it.

---

### 12. 🟡 CORS middleware uses wildcard `*` origin

[http_handler.go:76](file:///home/sandro/test/backend/internal/presentation/http_handler.go#L76)

```go
w.Header().Set("Access-Control-Allow-Origin", "*")
```

Acceptable for local development, but evaluators looking at **best practices** may flag this. Consider making it configurable via an environment variable.

---

## What's Done Well

| Aspect | Details |
|---|---|
| **Strategy Pattern** | Clean `NotifierStrategy` interface with per-channel implementations — exactly what the ADR asks for |
| **Registry Pattern** | Thread-safe registry with `sync.RWMutex`; trivial to add new channels |
| **Clean Architecture layers** | `domain` → `application` → `infrastructure` / `notification` / `presentation` — well separated |
| **Domain types** | `Category`, `Channel`, `DeliveryStatus` as typed constants — idiomatic Go |
| **Database** | Auto-migration, timestamp index, parameterized queries, proper `rows.Err()` check |
| **Fault tolerance** | Retry mechanism with exponential… well, fixed 50ms backoff, but the structure is there |
| **Frontend** | Functional form with client-side validation, loading/error states, status badges, responsive grid |
| **Docker** | Multi-stage builds for both services, SQLite volume persistence |

---

## Priority Action Items

1. **Fix retry count tracking** — return actual attempts from `executeWithRetry`
2. **Add category validation** in the handler — reject invalid categories
3. **Define a service interface** for the handler to depend on — complete DI
4. **Add handler and registry tests** — the ADR explicitly evaluates test breadth
5. **Add `UserName` to log** — both DB column and struct field
6. **Fix Docker API URL** — use nginx proxy or environment-based config
