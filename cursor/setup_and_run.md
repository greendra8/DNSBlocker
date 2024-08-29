# Setup and Run Guide for DNS Manager Project

This guide provides step-by-step instructions for setting up and running the DNS Manager project, including the CoreDNS server, PostgreSQL database, and SvelteKit frontend.

## Prerequisites

- Go (1.16 or later)
- PostgreSQL
- Node.js and npm
- CoreDNS source code

## 1. PostgreSQL Setup

1. Start PostgreSQL service:
   ```
   sudo service postgresql start
   ```

2. Connect to PostgreSQL:
   ```
   sudo -u postgres psql
   ```

3. Create database and user:
   ```sql
   CREATE DATABASE dns_manager;
   CREATE USER dns_user WITH ENCRYPTED PASSWORD 'your_password_here';
   GRANT ALL PRIVILEGES ON DATABASE dns_manager TO dns_user;
   ```

4. Connect to the new database:
   ```
   \c dns_manager
   ```

5. Create necessary tables:
   ```sql
   CREATE TABLE users (
       id SERIAL PRIMARY KEY,
       username TEXT UNIQUE NOT NULL,
       password_hash TEXT NOT NULL
   );

   CREATE TABLE blocked_sites (
       id SERIAL PRIMARY KEY,
       username TEXT NOT NULL,
       domain TEXT NOT NULL,
       FOREIGN KEY (username) REFERENCES users(username)
   );
   ```

6. Exit PostgreSQL:
   ```
   \q
   ```

## 2. CoreDNS Setup

1. Navigate to CoreDNS directory:
   ```
   cd /path/to/coredns
   ```

2. Update Corefile:
   ```
   .:53 {
       custom_user_plugin dns.domain.com :443 "postgres://dns_user:your_password_here@localhost/dns_manager?sslmode=disable"
       forward . 8.8.8.8 9.9.9.9
       log
       errors
   }
   ```

3. Generate SSL certificates (for development):
   ```
   openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
   ```

4. Build CoreDNS:
   ```
   go generate
   go build
   ```

5. Run CoreDNS (with sudo for port 53):
   ```
   sudo ./coredns
   ```

## 3. SvelteKit Frontend Setup

1. Navigate to the frontend directory:
   ```
   cd /path/to/dns-manager-frontend
   ```

2. Install dependencies:
   ```
   npm install
   ```

3. Start the development server:
   ```
   npm run dev
   ```

   The server typically starts on `http://localhost:5173`.

## Useful Commands

- Restart PostgreSQL:
  ```
  sudo service postgresql restart
  ```

- View CoreDNS logs:
  ```
  tail -f /var/log/coredns.log
  ```

- Test DNS resolution:
  ```
  dig @localhost -p 53 example.com
  ```

- Test DoH (DNS over HTTPS) with curl:
  ```
  echo -n 'q80BAAABAAAAAAAAA3d3dwdleGFtcGxlA2NvbQAAAQAB' | base64 -d | curl -H 'content-type: application/dns-message' -k --data-binary @- https://localhost/user1
  ```

## Troubleshooting

- If CoreDNS fails to start, check if port 53 is already in use:
  ```
  sudo lsof -i :53
  ```

- To view PostgreSQL logs:
  ```
  sudo tail -f /var/log/postgresql/postgresql-[version]-main.log
  ```

- If the frontend can't connect to the backend, ensure CORS is properly configured in the CoreDNS plugin.

Remember to replace placeholders like `your_password_here` and `/path/to/...` with actual values.