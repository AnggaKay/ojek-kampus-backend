@echo off
REM Rollback Script for Windows: Drop all tables
REM This will undo all migrations

echo ================================================
echo ⏪ Database Rollback: Drop All Tables
echo ================================================
echo.
echo ⚠️  WARNING: This will DELETE all tables and data!
echo Press Ctrl+C to cancel, or any key to continue...
pause > nul

set DB_USER=AnggaKay
set DB_NAME=ojek_kampus_db
set CONTAINER=ojek_db

echo.
echo 📋 Running rollback migrations (reverse order)...

echo.
echo   → Rolling back: 006_create_refresh_tokens
docker exec -i %CONTAINER% psql -U %DB_USER% -d %DB_NAME% < db/migrations/006_create_refresh_tokens.down.sql
echo   ✅ 006 rolled back

echo.
echo   → Rolling back: 005_create_otp_codes
docker exec -i %CONTAINER% psql -U %DB_USER% -d %DB_NAME% < db/migrations/005_create_otp_codes.down.sql
echo   ✅ 005 rolled back

echo.
echo   → Rolling back: 004_create_admin_users
docker exec -i %CONTAINER% psql -U %DB_USER% -d %DB_NAME% < db/migrations/004_create_admin_users.down.sql
echo   ✅ 004 rolled back

echo.
echo   → Rolling back: 003_create_driver_profiles
docker exec -i %CONTAINER% psql -U %DB_USER% -d %DB_NAME% < db/migrations/003_create_driver_profiles.down.sql
echo   ✅ 003 rolled back

echo.
echo   → Rolling back: 002_create_passenger_profiles
docker exec -i %CONTAINER% psql -U %DB_USER% -d %DB_NAME% < db/migrations/002_create_passenger_profiles.down.sql
echo   ✅ 002 rolled back

echo.
echo   → Rolling back: 001_create_users_table
docker exec -i %CONTAINER% psql -U %DB_USER% -d %DB_NAME% < db/migrations/001_create_users_table.down.sql
echo   ✅ 001 rolled back

echo.
echo ================================================
echo ✅ Rollback completed successfully!
echo ================================================
echo.
echo 📊 Verifying database (should be empty)...
docker exec %CONTAINER% psql -U %DB_USER% -d %DB_NAME% -c "\dt"

echo.
echo 🎉 All tables have been dropped.
pause
