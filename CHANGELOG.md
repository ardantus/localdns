# Changelog

All notable changes to this project will be documented in this file.

## [1.1.0] - 2025-12-18
### Added
- **Contact Info Auto-Copy**: Domain registrant contact information is automatically copied from user profile when creating a new domain.
- **Expiry Date Management**: Domains now have automatic expiry dates (default: 365 days) calculated from creation date.
- **Enhanced WHOIS Display**: Frontend WHOIS modal now properly displays all contact information with fallback to user profile data.

### Changed
- **Database Volume**: Changed from bind mount (`./pgdata`) to Docker named volume (`postgres_data`) for better data management.
  - Data is now automatically removed with `docker-compose down -v`.
  - Easier backup and restore operations.
- **User Data Loading**: Backend now preloads user/owner information when fetching domains for better frontend display.
- **Frontend Improvements**: 
  - Better handling of expired date display (fixes "1/1/1" issue).
  - Contact info fallback logic in WHOIS modal.
  - Admin sees "All Domains" while users see "Your Domains".

### Fixed
- **Domain Creation**: Fixed missing `updated_at` column error during domain registration.
- **Contact Data**: Fixed missing contact information in domain WHOIS data.
- **Expired Date**: Fixed invalid expired date display (showing "1/1/1") in frontend.
- **Database Schema**: Updated `init.sql` to match current database schema exactly.

### Improved
- **Database Schema**: `init.sql` now uses correct data types (TEXT for user fields, BIGINT for registrar config).
- **Migration Support**: Added manual migration fallbacks in backend for better compatibility.

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
