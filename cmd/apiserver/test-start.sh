export SWS_CONFDB_NAME=confdb
export SWS_CONFDB_USER=sws
export SWS_CONFDB_HOST=localhost
export SWS_CONFDB_PASSWORD=$PG_PASSWORD
export SWS_CONFDB_PORT=5432

export SWS_API_SERVER_PORT=10430

echo "Building SWS_API_SERVER"
go build
./apiserver
