services=("sws-loadbalancer" "sws-proxy-1" "sws-proxy-2" "sws-apiserver" "sws-confdb")
for service in "${services[@]}"
do
  stop_result=$(docker stop $service)
  echo "Stopping: $stop_result"
  rm_result=$(docker rm $service)
  echo "Removing: $rm_result"
done

NETWORK_NAME="sws-net"
if [ "$1" = "start-network" ]
then
  docker network create $NETWORK_NAME --driver bridge
fi

health_check() {
  name=$1
  uri=$2
  echo "Health checking: $name"
  curl $uri
  echo ""
}

# sws-confdb commands
DB_NAME="sws-confdb"
DB_PORT="5432"
DB_USER="sws"
DATABASE="confdb"

docker run -d --name $DB_NAME -p $DB_PORT:$DB_PORT --network $NETWORK_NAME \
  -e POSTGRES_USER=$DB_USER -e POSTGRES_DB=$DATABASE -e POSTGRES_PASSWORD=$PG_PASSWORD \
  postgres:10.2-alpine

echo "Waitng 5 seconds for $DB_NAME to be ready"
sleep 5

echo "Installing schema"
docker exec -i $DB_NAME psql -U sws confdb < ../resources/db-schema.sql

echo "Inserting seed data"
docker exec -i $DB_NAME psql -U sws confdb < ../resources/test/test-data.sql

# sws-apiserver commands
APISERVER_PORT="10430"
APISERVER_VERSION="v0.3"

docker run -d --name sws-apiserver --network $NETWORK_NAME \
  -p $APISERVER_PORT:$APISERVER_PORT -e SWS_API_SERVER_PORT=$APISERVER_PORT \
  -e SWS_CONFDB_NAME=$DATABASE -e SWS_CONFDB_USER=$DB_USER \
  -e SWS_CONFDB_HOST=$DB_NAME -e SWS_CONFDB_PASSWORD=$PG_PASSWORD \
  -e SWS_CONFDB_PORT=$DB_PORT czarsimon/sws-apiserver:$APISERVER_VERSION

echo "Waitng 2 seconds for apiserver to be ready"
sleep 2
health_check "sws-apiserver" "http://localhost:$APISERVER_PORT/health"

# sws-proxy commands
proxies=("sws-proxy-1" "sws-proxy-2")
for proxy in "${proxies[@]}"
do
  docker run -d --name $proxy --network $NETWORK_NAME czarsimon/sws-proxy:init
done

echo "Waitng 2 seconds for proxies to be ready"
sleep 2

# sws-loadbalancer commands
LB_PORT="81"
docker run -d --name sws-loadbalancer --network $NETWORK_NAME \
  -p $LB_PORT:$LB_PORT czarsimon/sws-lb:dev

echo "Waitng 2 seconds for sws-loadbalancer to be ready"
sleep 2
health_check "sws-loadbalancer" "http://localhost:$LB_PORT/sws-lb/health"

# List running sws services
docker ps | grep sws

create_dir_if_missing() {
  DIR=$1
  if [ ! -d "$DIR" ]
  then
    mkdir $DIR
  fi
}

SWS_DIR="$HOME/.sws"
create_dir_if_missing $SWS_DIR

SWS_PROXY_DIR="$SWS_DIR/proxy"
create_dir_if_missing $SWS_PROXY_DIR

rm $SWS_PROXY_DIR/Dockerfile
echo "FROM nginx:1.13.8-alpine\nWORKDIR /etc/nginx\n\nCOPY nginx.conf nginx.conf" > $SWS_PROXY_DIR/Dockerfile
