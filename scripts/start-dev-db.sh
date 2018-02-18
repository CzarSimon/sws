NAME="sws-confdb"
docker stop $NAME

docker run -d --name $NAME -p 5432:5432 --rm \
  -e POSTGRES_USER=sws -e POSTGRES_DB=confdb -e POSTGRES_PASSWORD=$PG_PASSWORD \
  postgres:10.2-alpine

echo "Waitng 5 seconds for $NAME to be ready"
sleep 5

echo "Installing schema"
docker exec -i $NAME psql -U sws confdb < ../resources/db-schema.sql

echo "Inserting seed data"
docker exec -i $NAME psql -U sws confdb < ../resources/test/test-data.sql
