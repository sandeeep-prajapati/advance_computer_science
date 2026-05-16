# Chapter 2: Setup & Installation

## System Requirements

Before installing Pimcore, make sure your environment meets these requirements:

| Requirement | Minimum Version |
|---|---|
| PHP | ≥ 8.1 (8.3 recommended for latest Pimcore) |
| MySQL | ≥ 8.0 |
| MariaDB | ≥ 10.4 |
| Composer | Latest |
| Node.js / npm | Optional (for front-end assets) |
| Web Server | Apache or Nginx |

---

## Option A: Using Composer (Recommended for production)

```bash
# Create a new project (skeleton = empty starter)
composer create-project pimcore/skeleton my-project

# OR use the demo project (has sample data)
composer create-project pimcore/demo my-project

# Move into the project folder
cd my-project

# Install dependencies
composer install
```

---

## Option B: Using Docker (Easiest for beginners)

### Step 1: Install Docker + Docker Compose

On Ubuntu/Debian:

```bash
# Install dependencies
sudo apt update
sudo apt install -y ca-certificates curl gnupg lsb-release

# Add Docker GPG key and repo
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | \
  sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] \
  https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Allow your user to run docker without sudo
sudo usermod -aG docker $USER
newgrp docker
```

Verify:
```bash
docker run --rm hello-world
docker compose version
```

### Step 2: Create a Pimcore Project Using the Official Image

```bash
mkdir ~/pimcore && cd ~/pimcore

# Use the Pimcore Docker image to run Composer (no PHP needed locally)
docker run -u $(id -u):$(id -g) --rm -v $(pwd):/var/www/html \
  pimcore/pimcore:php8.3-latest composer create-project pimcore/demo my-project

cd my-project
```

> **Why `-u $(id -u):$(id -g)`?** This ensures files created inside Docker are owned by your user, preventing permission headaches.

### Step 3: Configure docker-compose.yaml

```bash
# Check your UID:GID
echo $(id -u):$(id -g)

# Open docker-compose.yaml and uncomment the `user:` lines
# Replace with your uid:gid (e.g., user: "1000:1000")
```

### Step 4: Start the Stack

```bash
docker compose up -d

# Check running containers
docker compose ps
```

### Step 5: Run the Pimcore Installer

```bash
docker compose exec php vendor/bin/pimcore-install \
  --admin-username=admin --admin-password=admin \
  --mysql-username=pimcore --mysql-password=pimcore --mysql-database=pimcore \
  --mysql-host-socket=db --no-interaction
```

> This creates the database schema, admin user, and all base configuration.

### Step 6: Access Pimcore

- **Frontend:** `http://localhost`
- **Admin UI:** `http://localhost/admin`
- **Login:** Use the credentials you set in the installer (e.g., `admin` / `admin`)

---

## Useful Development Commands

```bash
# Follow logs
docker compose logs -f php
docker compose logs -f db

# Open a shell inside the PHP container
docker compose exec php bash

# Clear Symfony/Pimcore cache
docker compose exec php bin/console cache:clear

# Check system requirements
docker compose exec php bin/console pimcore:system:requirements:check

# Stop and remove containers + volumes (clean slate)
docker compose down -v --remove-orphans
```

---

## Common Pitfalls & Fixes

### 1. Permission Errors (files owned by root)
```bash
# On host
sudo chown -R $(id -u):$(id -g) .

# Or inside container
docker compose exec php bash -lc "chown -R www-data:www-data var public vendor"
```

### 2. Database Connection Error (SQLSTATE[HY000] [2002])
- Check that the `db` container is running: `docker compose ps`
- Ensure your DSN in `.env` uses `db` as the host (that's the Docker network hostname)
- Check DB logs: `docker compose logs db`

### 3. Admin Page 404 After Install
- Ensure your web server document root points to `my-project/public`
- Run `composer install` if you set up without Docker
- Check file permissions on `var/` and `public/var/`

### 4. `docker compose` Not Found
- Install the Docker Compose plugin: `sudo apt install docker-compose-plugin`

### 5. Elasticsearch Issues (Demo Project)
- The demo project may require Elasticsearch for search indexing
- You can disable the search backend locally until needed
- Add an `elasticsearch` service to `docker-compose.yaml` if required

---

## Production Architecture (Services Needed)

| Service | Purpose |
|---|---|
| PHP-FPM + Nginx/Apache | Web server + runtime |
| MySQL / MariaDB | Primary data store |
| Redis | Cache + sessions (highly recommended) |
| Elasticsearch | Search + product indexing (for large catalogs) |
| S3 / NFS | Scalable asset storage |
| RabbitMQ / Symfony Messenger | Async task queue (optional) |
| Varnish / CDN | Full-page caching + global delivery |
| Grafana + Prometheus | Monitoring |

---

## Quick Troubleshooting Cheatsheet

```bash
docker compose ps                    # Check containers
docker compose logs -f php           # Check PHP/Pimcore errors
docker compose exec php vendor/bin/pimcore-install --help  # Installer options
docker compose exec php bin/console cache:clear             # Clear cache
docker compose down -v               # Destroy volumes (fresh start)
```
