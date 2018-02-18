export SWS_CONFDB_NAME=confdb
export SWS_CONFDB_USER=sws
export SWS_CONFDB_HOST=localhost
export SWS_CONFDB_PASSWORD=$PG_PASSWORD
export SWS_CONFDB_PORT=5432

export SWS_AGENT_UPDATE_FREQ=30

echo "Building SWS_AGENT"
go build
./agent
