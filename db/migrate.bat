@echo off
echo ================================================
echo Database Migration: Fresh Start
echo ================================================
echo.
echo WARNING: This will DELETE all existing data!
echo Press Ctrl+C to cancel, or any key to continue...
pause > nul

echo.
echo Terminating active connections...
docker exec ojek_db psql -U AnggaKay -d postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'ojek_kampus_db' AND pid <> pg_backend_pid();"

echo.
echo Dropping database...
docker exec ojek_db psql -U AnggaKay -d postgres -c "DROP DATABASE IF EXISTS ojek_kampus_db;"

echo.
echo Creating database...
docker exec ojek_db psql -U AnggaKay -d postgres -c "CREATE DATABASE ojek_kampus_db;"

echo.
echo Running migrations...
docker exec -i ojek_db psql -U AnggaKay -d ojek_kampus_db < db/migrations/001_create_users_table.up.sql
docker exec -i ojek_db psql -U AnggaKay -d ojek_kampus_db < db/migrations/002_create_passenger_profiles.up.sql
docker exec -i ojek_db psql -U AnggaKay -d ojek_kampus_db < db/migrations/003_create_driver_profiles.up.sql
docker exec -i ojek_db psql -U AnggaKay -d ojek_kampus_db < db/migrations/004_create_admin_users.up.sql
docker exec -i ojek_db psql -U AnggaKay -d ojek_kampus_db < db/migrations/005_create_otp_codes.up.sql
docker exec -i ojek_db psql -U AnggaKay -d ojek_kampus_db < db/migrations/006_create_refresh_tokens.up.sql
docker exec -i ojek_db psql -U AnggaKay -d ojek_kampus_db < db/migrations/007_add_ktm_to_driver_profiles.up.sql
docker exec -i ojek_db psql -U AnggaKay -d ojek_kampus_db < db/migrations/008_enforce_motor_only.up.sql

echo.
echo ================================================
echo SUCCESS! Migration completed
echo ================================================
docker exec ojek_db psql -U AnggaKay -d ojek_kampus_db -c "\dt"
pause
