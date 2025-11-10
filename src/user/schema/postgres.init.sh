#!/bin/sh

# This script prepares the PostgreSQL database initialization
# It's typically not needed as PostgreSQL databases are created externally,
# but this script can be used for initial setup if needed.

cat <<EOF
-- Create database if it doesn't exist (run this as postgres superuser)
-- CREATE DATABASE ${POSTGRES_DB};

-- Create user if it doesn't exist (run this as postgres superuser)
-- CREATE USER ${POSTGRES_USER} WITH PASSWORD '${POSTGRES_PASSWORD}';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE ${POSTGRES_DB} TO ${POSTGRES_USER};

-- Grant schema privileges (run this while connected to the database)
GRANT ALL ON SCHEMA public TO ${POSTGRES_USER};
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ${POSTGRES_USER};
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ${POSTGRES_USER};

-- Ensure future tables and sequences are also granted
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO ${POSTGRES_USER};
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO ${POSTGRES_USER};
EOF
