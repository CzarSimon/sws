check_val() {
  var_name=$1
  val=$2
  if [[ ! -n "$val" ]]; then
    echo "$var_name empty"
    exit 1
  fi
}

health_check() {
  name=$1
  uri=$2
  echo "Health checking: $name"
  curl $uri
  echo ""
}

create_dir_if_missing() {
  DIR=$1
  if [ ! -d "$DIR" ]
  then
    mkdir $DIR
  fi
}

check_val "SWS_CONFDB_PASSWORD" $SWS_CONFDB_PASSWORD

echo "Removing previously runinng sws components"
services=("sws-loadbalancer" "sws-proxy-1" "sws-proxy-2" "sws-apiserver" "sws-confdb")
for service in "${services[@]}"
do
  stop_result=$(docker stop $service)
  echo "Stopping: $stop_result"
  rm_result=$(docker rm $service)
  echo "Removing: $rm_result"
done

NETWORK_NAME="sws-net"
docker network create $NETWORK_NAME --driver bridge

# sws-confdb commands
DB_NAME="sws-confdb"
DB_PORT="5432"
DB_USER="sws"
DATABASE="confdb"

DB_VOLUME="sws-confdb-volume"
docker volume create $DB_VOLUME

echo "Starting $DB_NAME"
docker run -d --name $DB_NAME -p $DB_PORT:$DB_PORT --network $NETWORK_NAME \
  -e POSTGRES_USER=$DB_USER -e POSTGRES_DB=$DATABASE \
  -e POSTGRES_PASSWORD=$SWS_CONFDB_PASSWORD \
  --mount source=$DB_VOLUME,target=/var/lib/postgresql/data \
  postgres:10.2-alpine

echo "Waitng 5 seconds for $DB_NAME to be ready"
sleep 5

echo "Installing schema"
docker exec -i $DB_NAME psql -U $DB_USER $DATABASE < db-schema.sql

# sws-apiserver commands
APISERVER_PORT="10430"
SWS_VERSION="v0.4"

docker run -d --name sws-apiserver --network $NETWORK_NAME \
  -p $APISERVER_PORT:$APISERVER_PORT -e SWS_API_SERVER_PORT=$APISERVER_PORT \
  -e SWS_CONFDB_NAME=$DATABASE -e SWS_CONFDB_USER=$DB_USER \
  -e SWS_CONFDB_HOST=$DB_NAME -e SWS_CONFDB_PASSWORD=$SWS_CONFDB_PASSWORD \
  -e SWS_CONFDB_PORT=$DB_PORT czarsimon/sws-apiserver:$SWS_VERSION

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
LB_PORT="80"
docker run -d --name sws-loadbalancer --network $NETWORK_NAME \
  -p $LB_PORT:$LB_PORT czarsimon/sws-lb:$SWS_VERSION

echo "Waitng 2 seconds for sws-loadbalancer to be ready"
sleep 2
health_check "sws-loadbalancer" "http://localhost:$LB_PORT/sws-lb/health"

# List running sws services
docker ps | grep sws

SWS_DIR="$HOME/.sws"
create_dir_if_missing $SWS_DIR

SWS_PROXY_DIR="$SWS_DIR/proxy"
create_dir_if_missing $SWS_PROXY_DIR

rm $SWS_PROXY_DIR/Dockerfile
cp Dockerfile $SWS_PROXY_DIR/Dockerfile

SWS_AGENT_DIR="/usr/local/sbin/sws-agent"
create_dir_if_missing $SWS_AGENT_DIR
cp sws-agent $SWS_AGENT_DIR/sws-agent

python subsitute_service_values.py
cp sws-agent.service /etc/systemd/system/sws-agent.service
systemctl enable sws-agent.service
systemctl start sws-agent.service
