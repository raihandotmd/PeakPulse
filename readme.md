## Overview

A simple RESTful API built with Gin, GORM, and PostgreSQL that lets authenticated users create and manage their own workout plans. Each plan consists of a title, description, and a list of exercises (with sets and reps). JWT authentication secures the endpoints.

## Features

* **User Authentication**
  * Sign up & log in using Supabase GoTrue
  * JWT issuance and middleware to protect routes

* **Workout Plans**
  * Create, read, update, and delete plans
  * Each plan contains multiple exercises with sets and reps

* **Exercises Catalog**
  * Master list of exercises seeded at startup

## Tech Stack

* **Language & Framework:** Go 1.20+, [Gin](https://github.com/gin-gonic/gin)
* **ORM:** [GORM](https://gorm.io) with PostgreSQL driver
* **Authentication:** [Supabase GoTrue](https://github.com/supabase-community/gotrue-go) for signup/login, custom JWT for API access
* **Migrations & Seeding:** GORM AutoMigrate + CSV import
* **API Docs:** [Swagger / Swaggo](https://github.com/swaggo/swag)

## Getting Started

### Prerequisites

* Go 1.20+
* PostgreSQL instance
* Supabase project (for Auth)

### Installation

1. **Clone & enter directory**

   ```bash
   git clone https://github.com/yourorg/peakPulse.git
   cd peakPulse
   ```

2. **Environment variables**
   Copy `.env.example` to `.env` and fill in:

   ```ini
   DB_DSN=host=<HOST> user=<USER> password=<PWD> dbname=<DB> sslmode=require
   SUPABASE_URL=https://<YOUR>.supabase.co
   SUPABASE_KEY=<anon_or_service_key>
   JWT_SECRET=<your_jwt_secret>
   ```

3. **Install dependencies & build**

   ```bash
   go mod tidy
   go build ./cmd/server
   ```

4. **Run the server**

   ```bash
   ./server
   # default: listens on http://localhost:8080
   ```

## API Reference

### Authentication

| Method | Path      | Description          | Body                  | Response                |
| ------ | --------- | -------------------- | --------------------- | ----------------------- |
| POST   | `/signup` | Register new user    | `{ email, password }` | `201 { message, user }` |
| POST   | `/login`  | Log in & receive JWT | `{ email, password }` | `200 { token }`         |

### Workout Plans

All `/plans` routes require `Authorization: Bearer <token>`.

| Method | Path          | Description                      | Body                                         | Response                  |
| ------ | ------------- | -------------------------------- | -------------------------------------------- | ------------------------- |
| POST   | `/plans`      | Create a new plan                | `{ title, description, exercises: [{...) }]` | `201 { plan, exercises }` |
| GET    | `/plans`      | List all plans                   | —                                            | `200 { plans: […] }`      |
| GET    | `/plans/{id}` | Get a single plan by ID          | —                                            | `200 { plan: […] }`       |
| PUT    | `/plans/{id}` | Update plan title/desc/exercises | `{ title?, description?, exercises? }`       | `200 { plan, exercises }` |
| DELETE | `/plans/{id}` | Delete a plan and its exercises  | —                                            | `200 { message }`         |

### Exercises Catalog

| Method | Path         | Description        | Body | Response                      |
| ------ | ------------ | ------------------ | ---- | ----------------------------- |
| GET    | `/exercises` | List all exercises | —    | `200 { data: […], count: N }` |

## Seeding Exercise Data

Import the provided `exercises.csv`:

```sql
COPY public.exercises(name,description,category,muscle_group)
FROM '/path/to/exercises_data-seed.csv' WITH (FORMAT csv, HEADER true);
```

## API Documentation

After running:

```bash
swag init
```

Serve Swagger UI at:

```go
r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```

Visit: `http://localhost:8080/docs/index.html`

## Future Improvements

* Pagination & filtering on list endpoints
* Role-based access or sharing of plans
* Export plans to PDF or CSV
* Webhooks for reminders or progress tracking