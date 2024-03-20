FROM postgres:15.2-alpine
COPY db/init.sql /docker-entrypoint-initdb.d/