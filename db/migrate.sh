#!/bin/bash
# Migration Script: Drop & Recreate Database
# WARNING: This will DELETE all data in the database!

set -e  # Exit on error

echo "================================================"
echo "🗑️  Database Migration: Fresh Start"
echo "================================================"
echo ""
echo "⚠️  WARNING: This will DELETE all existing data!"
echo "Press Ctrl+C to cancel, or Enter to continue..."
read

DB_USER="AnggaKay"
DB_NAME="ojek_kampus_db"
CONTAINER="ojek_db"

echo ""
echo "📋 Step 1: Dropping existing database..."
docker exec $CONTAINER psql -U $DB_USER -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;"
echo "✅ Database dropped"

echo ""
echo "📋 Step 2: Creating fresh database..."
docker exec $CONTAINER psql -U $DB_USER -d postgres -c "CREATE DATABASE $DB_NAME;"
echo "✅ Database created"

echo ""
echo "📋 Step 3: Running migrations..."

migrations=(
    "001_create_users_table"
    "002_create_passenger_profiles"
    "003_create_driver_profiles"
    "004_create_admin_users"
    "005_create_otp_codes"
    "006_create_refresh_tokens"
)

for migration in "${migrations[@]}"
do
    echo ""
    echo "  → Running migration: $migration"
    docker exec -i $CONTAINER psql -U $DB_USER -d $DB_NAME < db/migrations/${migration}.up.sql
    echo "  ✅ $migration applied"
done

echo ""
echo "================================================"
echo "✅ Migration completed successfully!"
echo "================================================"
echo ""
echo "📊 Verifying database structure..."
docker exec $CONTAINER psql -U $DB_USER -d $DB_NAME -c "\dt"

echo ""
echo "🎉 All done! Your database is ready."
