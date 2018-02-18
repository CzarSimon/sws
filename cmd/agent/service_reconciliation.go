package main

import "log"

func (env *Env) ReconcileServices() {
	if env.Agent.Locked {
		return
	}
	log.Println("Reconciling services")
	env.Agent.Lock()
	env.upsertAgentMetadata(true)
}
