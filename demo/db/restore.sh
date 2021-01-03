#!/bin/bash

until pg_restore -h postgres -U postgres -d ${POSTGRES_DB} $1; do
    echo 'Wait 10 seconds...'
    sleep 10
done
