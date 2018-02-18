NAME="sws-confdb"

docker run -d --name $NAME -p 5432:5432 --rm \
  -e POSTGRES_USER=sws -e POSTGRES_DB=confdb -e POSTGRES_PASSWORD=$PG_PASSWORD \
  --mount source=sws-conf,target=/var/lib/postgresql/data postgres:10.2-alpine

sleep 2
docker exec -i $NAME psql -U sws confdb < ../resources/db-schema.sql
