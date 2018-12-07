#!/bin/sh

set -eux; \
    service rsyslog start; \
    logger -t api-service created by logger; \
    logger -t worker1 created by logger; \
    logger -t worker2 created by logger; \
    gosu postgres service postgresql start; \
    gosu postgres psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';"; \
    gosu postgres createdb -O docker app; \
    gosu docker psql -d app < /app/adventureworks-pg/install.sql > /dev/null; \
    gosu docker psql -d app < /app/adventureworks-pg/migration_0001.sql > /dev/null; \
    service redis-server start; \
    /app/build/_workspace/bin/api -lb syslog -ll trace & \
    /app/build/_workspace/bin/worker1 -lb syslog -ll trace & \
    /app/build/_workspace/bin/worker2 -lb syslog -ll trace & \
    /app/tmux.sh; \
    /bin/bash \
    ;
