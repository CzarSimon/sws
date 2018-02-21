package main

import (
	"os"
	"path/filepath"
)

// Common infrastructure constants.
const (
	NetworkName = "sws-net"
	AgentName   = "sws-agent"
)

// Common infrastructure variables
var (
	ContainerFolder = filepath.Join(os.Getenv("HOME"), ".sws", "containers")
)

// BootupAgent starts the sws agent and ensures the host is ready to set up and serve services.
func (env *Env) BootupAgent() error {
	err := env.setupAgentMetadata()
	if err != nil {
		return err
	}
	err = env.reconcileState()
	if err != nil {
		return nil
	}
	return env.upsertAgentMetadata(true)
}

// setupAgentMetadata initalizes agent metadata.
func (env *Env) setupAgentMetadata() error {
	am, err := GetAgentMetadata(env.DB)
	if err != nil && err != NoSuchAgentErr {
		return err
	}
	if err == NoSuchAgentErr {
		am = NewAgentMetadata(AgentName)
	}
	env.Agent = am
	return nil
}
