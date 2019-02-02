package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/anothrnick/machinable/dsi/mongo"
	"github.com/anothrnick/machinable/dsi/mongo/database"
	"github.com/anothrnick/machinable/management"
	"github.com/anothrnick/machinable/projects"
)

// HostSwitch is used to switch routers based on sub domain
type HostSwitch map[string]http.Handler

// Implement the ServeHTTP method on our new type
func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if a http.Handler is registered for the given host.
	// If yes, use it to handle the request.

	hostParts := strings.Split(r.Host, ".")

	// {sub}.{domain}.{tld}
	if len(hostParts) != 3 {
		http.Error(w, "Project Not Found", 404)
		return
	}

	subDomain := hostParts[0]

	handler, ok := hs[subDomain]

	if ok {
		handler.ServeHTTP(w, r)
	} else {
		hs["*"].ServeHTTP(w, r)
	}
}

func main() {
	// use mongoDB connector
	// if another connector is needed, the Datastore interface can be implemented and these 2 lines changes to instantiate the new connector
	// potential connectors: InfluxDB, Postgres JSON, Redis, CouchDB
	mongoDB := database.Connect()
	datastore := mongo.New(mongoDB)

	// switch routers based on subdomain
	hostSwitch := make(HostSwitch)

	// manage is for the management application api, i.e. project/team management
	hostSwitch["manage"] = management.CreateRoutes(datastore)
	// all other subdomains will be treated as project names, and use the project routes
	hostSwitch["*"] = projects.CreateRoutes(datastore)

	log.Fatal(http.ListenAndServe(":5001", hostSwitch))
}
