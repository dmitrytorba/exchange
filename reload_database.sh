#!/usr/bin/env sh
# This script will delete the exchange database and then rebuild it fresh
# from the schema.sql file.

sudo -u postgres psql <<EOF
\x
DROP DATABASE exchange;
CREATE DATABASE exchange;
EOF

sudo -u postgres psql -f schema.sql exchange
