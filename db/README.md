# Database Migrations

This directory contains database migration files for the Ojek Kampus project.

## 📁 Structure

```
db/
├── migrations/
│   ├── 001_create_users_table.up.sql           # Create users table
│   ├── 001_create_users_table.down.sql         # Rollback users table
│   ├── 002_create_passenger_profiles.up.sql    # Create passenger_profiles
│   ├── 002_create_passenger_profiles.down.sql
│   ├── 003_create_driver_profiles.up.sql       # Create driver_profiles
│   ├── 003_create_driver_profiles.down.sql
│   ├── 004_create_admin_users.up.sql           # Create admin_users
│   ├── 004_create_admin_users.down.sql
│   ├── 005_create_otp_codes.up.sql             # Create otp_codes
│   ├── 005_create_otp_codes.down.sql
│   ├── 006_create_refresh_tokens.up.sql        # Create refresh_tokens
│   └── 006_create_refresh_tokens.down.sql
├── migrate.bat          # Windows script to run all migrations
├── rollback.bat         # Windows script to rollback all migrations
├── migrate.sh           # Linux/Mac script (for production)
└── README.md           # This file
```

## 🚀 How to Run Migrations

### Prerequisites
- Docker is running
- Container `ojek_db` is up and running
- Database credentials are correct

### Windows (Development)

**Option 1: Run All Migrations (Fresh Start)**
```bash
# Double-click or run in terminal:
db\migrate.bat
```

This will:
1. ⚠️ **DROP** the existing database
2. Create a fresh database
3. Run all migrations in order
4. Show database structure

**Option 2: Rollback All Migrations**
```bash
# Double-click or run in terminal:
db\rollback.bat
```

This will:
1. Drop all tables in reverse order
2. Remove all data
3. Leave database empty

### Linux/Mac (Production)

```bash
# Make script executable
chmod +x db/migrate.sh

# Run migrations
./db/migrate.sh
```

### Manual Migration (Advanced)

If you want to run migrations manually:

```bash
# Run a specific migration
docker exec -i ojek_db psql -U AnggaKay -d ojek_kampus_db < db/migrations/001_create_users_table.up.sql

# Rollback a specific migration
docker exec -i ojek_db psql -U AnggaKay -d ojek_kampus_db < db/migrations/001_create_users_table.down.sql
```

## 📊 Database Schema

After running all migrations, you will have these tables:

### 1. `users`
- Core user authentication table
- Supports both PASSENGER and DRIVER roles
- Password-based authentication (bcrypt)
- Phone number as primary identifier

### 2. `passenger_profiles`
- Passenger-specific data
- FCM token for push notifications
- Order metrics and ratings

### 3. `driver_profiles`
- Driver-specific data
- Vehicle information
- Document verification status
- Location tracking (lat/long)
- Online/offline status

### 4. `admin_users`
- Separate admin authentication
- Role-based access (SUPER_ADMIN, ADMIN, VIEWER)
- Dashboard access only

### 5. `otp_codes`
- OTP for phone verification
- OTP for password reset
- Expiration and attempt tracking

### 6. `refresh_tokens`
- JWT refresh token storage
- Multi-device support
- Token revocation capability

## 🔧 Troubleshooting

### Error: "role postgres does not exist"
**Solution:** Use the correct database user (`AnggaKay` in our case). Check `.env` file.

### Error: "database does not exist"
**Solution:** The migrate script will create it. Make sure Docker container is running.

### Error: "relation already exists"
**Solution:** Run rollback first, then migrate again:
```bash
db\rollback.bat
db\migrate.bat
```

### How to check current database structure
```bash
docker exec ojek_db psql -U AnggaKay -d ojek_kampus_db -c "\dt"
docker exec ojek_db psql -U AnggaKay -d ojek_kampus_db -c "\d users"
```

## ⚠️ Important Notes

1. **Development vs Production**
   - Use `migrate.bat` for development (Windows)
   - Use `migrate.sh` for production (Linux/Mac)

2. **Data Loss Warning**
   - `migrate.bat` will **DELETE ALL DATA**
   - Only use in development or with database backup

3. **Migration Order**
   - Always run migrations in order (001 → 002 → 003...)
   - Dependencies between tables require correct order

4. **Rollback Order**
   - Rollback must be in reverse order (006 → 005 → 004...)
   - This prevents foreign key constraint errors

## 📝 Creating New Migrations

When adding new tables or modifying schema:

1. Create two files:
   - `00X_description.up.sql` (apply changes)
   - `00X_description.down.sql` (undo changes)

2. Follow naming convention:
   - Use sequential numbering (001, 002, 003...)
   - Use descriptive names (create_table_name, add_column_name)

3. Test both UP and DOWN:
   ```bash
   # Test UP
   docker exec -i ojek_db psql -U AnggaKay -d ojek_kampus_db < db/migrations/00X_new_migration.up.sql
   
   # Test DOWN
   docker exec -i ojek_db psql -U AnggaKay -d ojek_kampus_db < db/migrations/00X_new_migration.down.sql
   ```

4. Update `migrate.bat` and `rollback.bat` with new migration

## 🎓 Learning Resources

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Database Migrations Best Practices](https://www.prisma.io/dataguide/types/relational/what-are-database-migrations)
- [Docker Exec Commands](https://docs.docker.com/engine/reference/commandline/exec/)
