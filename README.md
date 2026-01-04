# Gator 🐊

Gator is a command-line RSS feed aggregator written in Go. It allows users to register, follow RSS feeds, fetch posts, and browse aggregated content using a PostgreSQL backend.

This project is intended to be used as a compiled CLI binary, not via `go run`.

---

## Requirements

Before using `gator`, make sure you have the following installed:

* **Go** (1.25 or newer recommended)
  [https://go.dev/dl/](https://go.dev/dl/)

* **PostgreSQL** (17 or newer recommended)
  [https://www.postgresql.org/download/](https://www.postgresql.org/download/)

You will also need a running PostgreSQL database that `gator` can connect to.

---

## Installation

Since `gator` is a Go CLI tool, you install it using `go install`.

From anywhere on your system, run:

```bash
go install github.com/Rachit-Gandhi/gator@latest
```

Make sure your `$GOPATH/bin` (or `$HOME/go/bin`) is in your `PATH`. After installation, you should be able to run:

```bash
gator
```

without using `go run` or being inside the project directory.

---

## Configuration

`gator` uses a JSON configuration file located in your **home directory**:

```text
~/.gatorconfig.json
```

### Example config file

```json
{
  "db_url": "postgres://postgres:@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

### Fields

* `db_url`
  PostgreSQL connection string used by the application.

* `current_user_name`
  The currently active user for CLI commands, not needed to be set.

The config file path is resolved using `os.UserHomeDir()` in the code, so it must exist in your home directory.

---

## Database Setup

Create a PostgreSQL database (example):

```sql
CREATE DATABASE gator;
```

Ensure your database credentials match the `db_url` in `.gatorconfig.json`.

Run database migrations if required by the project before using the CLI.

---

## Usage

Once installed and configured, you can use `gator` directly from your terminal.

### Common Commands

```bash
gator register <username>
```

Registers a new user.

```bash
gator login <username>
```

Sets the active user in the config file.

```bash
gator addfeed <feed_name> <feed_url>
```

Adds a new RSS feed.

```bash
gator follow <feed_url>
```

Follow an existing feed.

```bash
gator unfollow <feed_url>
```

Unfollow a feed.

```bash
gator fetch
```

Fetches the latest posts from followed feeds.

```bash
gator browse [limit]
```

Browse aggregated posts (default limit applies if not specified).

---

## Development vs Production

* `go run .`
  Used only during development.

* `go build` / `go install`
  Produces a statically compiled binary suitable for production use.

Once built or installed, the `gator` binary does **not** require the Go toolchain to run.

---
