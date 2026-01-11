# 📰 RSS Aggregator

A lightweight Go service that aggregates **RSS/Atom feeds** into a **PostgreSQL** database and exposes a clean **REST API** to manage users, feeds, subscriptions, and posts.

---

## ✨ Key Features

* 👤 Create users (server generates an API key)
* 📡 Register RSS/Atom feeds and follow/unfollow them
* 🔄 Background scraper that periodically fetches feeds and stores posts
* 🚀 Simple RESTful HTTP API with lightweight JSON responses

---

## 🧰 Prerequisites

* **Go** 1.20+ (module-aware)
* **PostgreSQL** database
* *(Optional)* **sqlc** — required only if you want to regenerate typed queries

---

## ⚙️ Configuration

The service reads configuration from environment variables. A local `.env` file is also loaded automatically if present.

| Variable | Description                                        |
| -------- | -------------------------------------------------- |
| `PORT`   | TCP port the HTTP server listens on (**required**) |
| `DB_URL` | PostgreSQL connection string (**required**)        |

**Example:**

```
postgres://user:password@localhost:5432/rss_aggregator?sslmode=disable
```

---

## 🏗️ Build & Run

### 1️⃣ Build

```bash
go build -o rss_aggregator
```

### 2️⃣ Run

```bash
PORT=8080 \
DB_URL="postgres://user:pass@localhost:5432/rss_aggregator?sslmode=disable" \
./rss_aggregator
```

* The HTTP server starts immediately
* The background scraper runs **every 1 minute** using **10 goroutines** (configurable in code)

---

## 🗄️ Database

This project uses **sqlc-generated typed queries**, located under:

```
internal/database
```

If you modify SQL files, regenerate queries using **sqlc** (see sqlc documentation).

### Minimal SQL Schema

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  name TEXT NOT NULL,
  api_key TEXT NOT NULL UNIQUE
);

CREATE TABLE feeds (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  name TEXT NOT NULL,
  url TEXT NOT NULL,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  last_fetched_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE feed_follows (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  feed_id UUID REFERENCES feeds(id) ON DELETE CASCADE
);

CREATE TABLE posts (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  title TEXT NOT NULL,
  description TEXT,
  published_at TIMESTAMP WITH TIME ZONE NOT NULL,
  url TEXT NOT NULL UNIQUE,
  feed_id UUID REFERENCES feeds(id) ON DELETE CASCADE
);
```

---

## 🔐 API Authentication

Some endpoints require an API key passed via the `Authorization` header:

```
Authorization: ApiKey <your_api_key>
```

---

## 🌐 API Reference

**Base Path:** `/v1`

### Health & Debug

| Method | Endpoint   | Description               |
| ------ | ---------- | ------------------------- |
| GET    | `/healthz` | Health/readiness check    |
| GET    | `/error`   | Sample 400 error response |

### Users

| Method | Endpoint | Auth | Description            |
| ------ | -------- | ---- | ---------------------- |
| POST   | `/users` | ❌    | Create a new user      |
| GET    | `/users` | ✅    | Get authenticated user |

### Feeds

| Method | Endpoint | Auth | Description    |
| ------ | -------- | ---- | -------------- |
| POST   | `/feeds` | ✅    | Create a feed  |
| GET    | `/feeds` | ❌    | List all feeds |

### Feed Follows

| Method | Endpoint             | Auth | Description         |
| ------ | -------------------- | ---- | ------------------- |
| POST   | `/feed_follows`      | ✅    | Follow a feed       |
| GET    | `/feed_follows`      | ✅    | List followed feeds |
| DELETE | `/feed_follows/{id}` | ✅    | Unfollow a feed     |

### Posts

| Method | Endpoint | Auth | Description                 |
| ------ | -------- | ---- | --------------------------- |
| GET    | `/posts` | ✅    | Get latest posts (limit 10) |

---

## 📌 API Usage Examples

### Create a User

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice"}' \
  http://localhost:8080/v1/users
```

### Create a Feed (Authenticated)

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey <your_api_key>" \
  -d '{"name":"Golang Blog","url":"https://blog.golang.org/feed.atom"}' \
  http://localhost:8080/v1/feeds
```

### Follow a Feed

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey <your_api_key>" \
  -d '{"feed_id":"<feed-uuid>"}' \
  http://localhost:8080/v1/feed_follows
```

### Get User Posts

```bash
curl -X GET \
  -H "Authorization: ApiKey <your_api_key>" \
  http://localhost:8080/v1/posts
```

---

## 🛠️ Development Notes

* Background scraper is started via `startScraping` and fetches feeds concurrently
* Scraper settings (interval & concurrency) are currently **hardcoded** in `main.go`
* CORS is enabled for **all origins** (adjust in `main.go` for production)
* JSON helpers and error utilities live in `json.go`
* Authorization logic is implemented in `internal/auth/auth.go`

---

## 🤝 Contributing

Contributions are welcome!

1. Open an issue describing the bug or enhancement
2. Create a feature/topic branch
3. Submit a pull request

---

⭐ If you find this project useful, consider giving it a star!
