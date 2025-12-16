# Changelog

All notable changes to this project will be documented in this file.

## [1.0.0] - 2025-12-16
### Added
- **Architecture**: Transformed from static file-based DNS to a Database-backed Local Registrar.
- **Database**: Introduced **PostgreSQL** to store Users, Domains, and DNS Records.
- **DNS Server**: Custom **CoreDNS** build with `pdsql` plugin to query records directly from PostgreSQL in real-time.
- **Backend API**: New **Go (Golang)** REST API for:
    -   User Authentication (JWT).
    -   Domain Registration.
    -   DNS Record Management.
- **Frontend UI**: Modern **React (Vite)** dashboard with **TailwindCSS** for:
    -   User Login/Registration.
    -   Domain Management (Add/List).
    -   DNS Record Editing.
-   **Admin Features**:
    -   Default Admin user (`admin` / `admin123`) seeded on startup.
    -   Role-Based Access Control (RBAC) middleware.
    -   Admin Dashboard with ability to view all domains and owners.
- **Infrastructure**: Updated `docker-compose.yml` to orchestrate all services (CoreDNS, Postgres, Backend, Frontend).

## [0.1.0] - Initial Version
### Features
- Basic CoreDNS configuration.
- Static zone files loaded from disk.
- Manual file editing required for updates.
