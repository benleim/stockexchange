FROM mongo:latest

ENV MONGO_INITDB_ROOT_USERNAME admin-user
ENV MONGO_INITDB_ROOT_PASSWORD admin-password
ENV MONGO_INITDB_DATABASE exchange

ADD mongo-init.js /docker-entrypoint-initdb.d/