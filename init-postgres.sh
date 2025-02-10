#!/bin/bash

# Wait for PostgreSQL to start up
echo "Waiting for PostgreSQL to start..."
su postgres
until pg_isready -U postgres -d postgres; do
  sleep 1
done

# Connect to the default 'postgres' database to create the new role and database
echo "Creating role and database 'greenlight'..."

psql -U postgres -d postgres <<EOF
-- Create a new role (user) for the greenlight database
CREATE ROLE greenlight WITH LOGIN PASSWORD 'greenlight_password';

-- Create the greenlight database owned by the greenlight
CREATE DATABASE greenlight OWNER greenlight;

-- Enable the citext extension in the greenlight database
\c greenlight
CREATE EXTENSION IF NOT EXISTS citext;
EOF

# Apply PostgreSQL configuration settings to the greenlight database
echo "Applying custom PostgreSQL configurations to the 'greenlight' database..."

psql -U postgres -d greenlight <<EOF
ALTER SYSTEM SET max_connections = '100';
ALTER SYSTEM SET shared_buffers = '1GB';
ALTER SYSTEM SET effective_cache_size = '3GB';
ALTER SYSTEM SET maintenance_work_mem = '256MB';
ALTER SYSTEM SET checkpoint_completion_target = '0.9';
ALTER SYSTEM SET wal_buffers = '16MB';
ALTER SYSTEM SET default_statistics_target = '100';
ALTER SYSTEM SET random_page_cost = '1.1';
ALTER SYSTEM SET effective_io_concurrency = '200';
ALTER SYSTEM SET work_mem = '5242kB';
ALTER SYSTEM SET huge_pages = 'off';
ALTER SYSTEM SET min_wal_size = '1GB';
ALTER SYSTEM SET max_wal_size = '4GB';
EOF

# Reload PostgreSQL to apply the changes
pg_ctl reload -D /var/lib/postgresql/data

echo "Custom configurations applied successfully."