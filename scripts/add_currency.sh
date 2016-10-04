#!/usr/bin/env sh

var=$1
if [ -z ${var+x} ]; then 
	echo "Command format: add_currency <3 letter identifier of currency>"; 
	exit 1
fi

psql -d exchange <<EOF
\x
ALTER TYPE currency ADD VALUE '$1';
ALTER TABLE users ADD COLUMN $1 bigint;
EOF
