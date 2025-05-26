# aggreGATOR

aggreGATOR is a command-line RSS feed aggregator built in Go, developed as part of a [Boot.dev](https://www.boot.dev) guided project. It enables users to register, log in, add and follow RSS feeds, and aggregate posts for offline reading.

## Prerequisites

Ensure you have installed:

* The [Go](https://go.dev) programming language to build and run the application.
* A [PostgreSQL](https://www.postgresql.org) database to store all the data.

The PostgreSQL server has to be running before starting the application and a database named `gator` also must exist.

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/marekmchl/aggreGATOR.git
cd aggreGATOR
```

### 2. Install the `aggreGATOR` CLI

```bash
go install .
```

This will compile the application and place the executable in your Go `bin` directory. Ensure that this directory is included in your system's `PATH`.

## Configuration

The application requires a configuration file named `gatorconfig.json` located in the root directory. This file should contain your database connection string and the current user's username. Example:

```json
{
    "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
    "current_user_name": "your_username"
}
```

Replace `username` and `password` with your PostgreSQL credentials.

## 🧪 Usage

* **Register a new user:**

  ```bash
  aggreGATOR register your_username
  ```

* **Log in as a user:**

  ```bash
  aggreGATOR login username
  ```

* **Add a new RSS feed:**

  ```bash
  aggreGATOR addfeed https://example.com/rss
  ```

* **Aggregate feeds:**

  This command fetches and stores posts from active user's followed feeds. The duration parameter specifies how often to refresh the feeds. Beware too short refresh durations, they may overwhelm the servers.

  ```bash
  aggreGATOR agg 1h
  ```

* **Browse aggregated posts:**

  ```bash
  aggreGATOR browse
  ```
