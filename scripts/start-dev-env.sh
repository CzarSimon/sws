NETWORK_NAME="sws-dev-net"
docker network create $NETWORK_NAME --driver bridge

# DB config
DB_NAME="sws-confdb"
DB_PORT="5432"
DB_USER="sws"
DATABASE="confdb"
docker stop $DB_NAME

docker run -d --name $DB_NAME -p $DB_PORT:$DB_PORT --rm --network $NETWORK_NAME \
  -e POSTGRES_USER=$DB_USER -e POSTGRES_DB=$DATABASE -e POSTGRES_PASSWORD=$PG_PASSWORD \
  postgres:10.2-alpine

echo "Waitng 5 seconds for $DB_NAME to be ready"
sleep 5

echo "Installing schema"
docker exec -i $DB_NAME psql -U sws confdb < ../resources/db-schema.sql

echo "Inserting seed data"
docker exec -i $DB_NAME psql -U sws confdb < ../resources/test/test-data.sql

APISERVER_PORT="10430"
APISERVER_VERSION="v0.1"
docker stop sws-apiserver

docker run -d --name sws-apiserver --rm --network $NETWORK_NAME \
  -p $APISERVER_PORT:$APISERVER_PORT -e SWS_API_SERVER_PORT=$APISERVER_PORT \
  -e SWS_CONFDB_NAME=$DATABASE -e SWS_CONFDB_USER=$DB_USER \
  -e SWS_CONFDB_HOST=$DB_NAME -e SWS_CONFDB_PASSWORD=$PG_PASSWORD \
  -e SWS_CONFDB_PORT=$DB_PORT czarsimon/sws-apiserver:$APISERVER_VERSION

echo "Waitng 5 seconds for apiserver to be ready"
sleep 5

curl http://localhost:$APISERVER_PORT/health
