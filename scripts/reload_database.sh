#!/usr/bin/env sh
# This script will delete the exchange database and then rebuild it fresh
# from the schema.sql file.

psql <<EOF
\x
DROP DATABASE exchange;
CREATE DATABASE exchange;
EOF

psql -f ./schema.sql exchange
