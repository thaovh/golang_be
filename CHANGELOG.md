# Changelog

## [1.2.0] - 2024-01-15

### Added
- **BMSF Table Prefix**: Implemented BMSF_ prefix for all application tables
- **Table Naming Convention**: All tables now use BMSF_ prefix (BMSF_USER, BMSF_ORDER, etc.)
- **Constraint Naming**: Updated constraint naming to include BMSF_ prefix (PK_BMSF_USER, UK_BMSF_USER_USERNAME)
- **Index Naming**: Updated index naming to include BMSF_ prefix (IDX_BMSF_USER_USERNAME)
- **Trigger Naming**: Updated trigger naming to include BMSF_ prefix (TR_BMSF_USER_UPDATED_AT)

### Changed
- **Database Schema**: Updated all table references from USER to BMSF_USER
- **Repository Queries**: Updated all SQL queries to use BMSF_USER table
- **Migration Script**: Updated to create BMSF_USER table with proper naming
- **Entity Comments**: Added comments explaining BMSF_ table mapping
- **Documentation**: Updated README with BMSF_ prefix convention explanation

### Database Changes
- Table name: `USER` → `BMSF_USER`
- Constraints: `PK_USER` → `PK_BMSF_USER`, `UK_USER_USERNAME` → `UK_BMSF_USER_USERNAME`
- Indexes: `IDX_USER_*` → `IDX_BMSF_USER_*`
- Trigger: `TR_USER_UPDATED_AT` → `TR_BMSF_USER_UPDATED_AT`

### Breaking Changes
- All database queries now reference BMSF_USER table
- Migration script creates BMSF_USER instead of USER table
- Existing USER table needs to be renamed or recreated

### Migration Required
```sql
-- Drop old table if exists
DROP TABLE USER;

-- Run new migration with BMSF_ prefix
@migrations/001_create_users.sql
```

## [1.1.0] - 2024-01-15

### Added
- **BaseEntity**: Implemented base entity with common fields (ID, CreatedAt, UpdatedAt, CreatedBy, UpdatedBy, DeletedAt, Version, TenantID)
- **Oracle Naming Conventions**: Updated all database queries to use UPPERCASE table and column names
- **Soft Delete**: Implemented soft delete pattern with DELETED_AT field
- **Audit Fields**: Added CreatedBy, UpdatedBy fields for audit trail
- **Versioning**: Added Version field for optimistic locking
- **Multi-tenancy Support**: Added TenantID field for tenant isolation

### Changed
- **Database Schema**: Updated USER table with new fields and Oracle naming conventions
- **User Entity**: Refactored to use BaseEntity composition
- **Repository Queries**: Updated all SQL queries to use UPPERCASE naming
- **Migration Script**: Updated to create USER table with proper Oracle conventions
- **User Methods**: Updated to use BaseEntity methods for versioning and audit

### Database Changes
- Table name: `users` → `USER`
- Column names: `id` → `ID`, `username` → `USERNAME`, etc.
- Added fields: `CREATED_BY`, `UPDATED_BY`, `DELETED_AT`, `VERSION`, `TENANT_ID`
- Updated constraints: `PK_USERS` → `PK_USER`, `UK_USERS_USERNAME` → `UK_USER_USERNAME`
- Updated indexes: `IDX_USER_*` naming convention
- Updated trigger: `TR_USER_UPDATED_AT` with version increment

### Breaking Changes
- All database queries now use UPPERCASE naming
- User entity now embeds BaseEntity
- Soft delete is now the default behavior
- Version field is required for optimistic locking

### Migration Required
```sql
-- Drop old table
DROP TABLE users;

-- Run new migration
@migrations/001_create_users.sql
```

## [1.0.0] - 2024-01-15

### Added
- Initial project setup with Clean Architecture
- User entity and CRUD operations
- Oracle database integration
- RESTful API with Gin framework
- Error handling with standardized codes
- Configuration management
- Logging with Zap
- Dependency injection container
