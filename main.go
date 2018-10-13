package main

import (
	"log"
	"net/http"
	"strings"

	"bitbucket.org/nsjostrom/machinable/management"
	"bitbucket.org/nsjostrom/machinable/projects"
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
	appRoutes := management.CreateManagementRoutes()
	projectRoutes := projects.CreateProjectRoutes()

	// switch routers based on subdomain
	hostSwitch := make(HostSwitch)

	// manage is for the management application api, i.e. project/team management
	hostSwitch["manage"] = appRoutes
	// all other subdomains will be treated as project names, and use the project routes
	hostSwitch["*"] = projectRoutes

	log.Fatal(http.ListenAndServe(":5001", hostSwitch))
}
