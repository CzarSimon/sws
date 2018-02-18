package main

import (
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/CzarSimon/sws/pkg/swsutil"
)

var (
	upsertAgentMetadataQuery = getUpsertAgentMetadataQuery()
	NoSuchAgentErr           = errors.New("No such agent exists")
)

// AgentMetadata threadsafe metadata.
type AgentMetadata struct {
	AgentID     string
	LastUpdated time.Time
	Mutex       sync.Mutex
	Locked      bool
}

// NewAgentMetadata creates a new agent metadata.
func NewAgentMetadata(name string) *AgentMetadata {
	return &AgentMetadata{
		AgentID: name,
		Mutex:   sync.Mutex{},
		Locked:  false,
	}
}

// Lock locks the metatdata to prevent concurrent updates.
func (am *AgentMetadata) Lock() {
	am.Mutex.Lock()
	am.Locked = true
}

// Unlock unlocka metadata to allow updates.
func (am *AgentMetadata) Unlock() {
	am.Mutex.Unlock()
	am.Locked = false
}

// GetAgentMetadata gets agent metadata from the database.
func GetAgentMetadata(db *sql.DB) (*AgentMetadata, error) {
	agent := &AgentMetadata{
		AgentID: AgentName,
		Mutex:   sync.Mutex{},
		Locked:  false,
	}
	query := "SELECT LAST_UPDATED FROM AGENT_METADATA WHERE AGENT_ID = $1"
	err := db.QueryRow(query).Scan(agent.LastUpdated)
	if err == sql.ErrNoRows {
		return nil, NoSuchAgentErr
	}
	return agent, nil
}

// upsertAgentMetadata updates the last updated timestamp for the agent
// and stores in the database.
func (env *Env) upsertAgentMetadata(updateTime bool) error {
	if updateTime {
		env.Agent.LastUpdated = swsutil.GetNow()
	}
	stmt, err := env.DB.Prepare(upsertAgentMetadataQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(env.Agent.AgentID, env.Agent.LastUpdated)
	if err != nil {
		return err
	}
	if updateTime && env.Agent.Locked {
		env.Agent.Unlock()
	}
	return nil
}

// getUpsertAgentMetadataQuery gets query to update agent metatdata.
func getUpsertAgentMetadataQuery() string {
	return `
    INSERT INTO AGENT_METADATA(AGENT_ID, LAST_UPDATED) VALUES($1, $2)
      ON CONFLICT(AGENT_ID) DO UPDATE SET LAST_UPDATED = $2`
}
