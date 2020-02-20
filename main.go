package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/anothrnick/machinable/config"
	"github.com/anothrnick/machinable/dsi/postgres"
	"github.com/anothrnick/machinable/events"
	"github.com/anothrnick/machinable/management"
	"github.com/anothrnick/machinable/projects"
	"github.com/go-redis/redis"
)

// HostSwitch is used to switch routers based on sub domain
type HostSwitch map[string]http.Handler

// Implement the ServeHTTP method on our new type
func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if a http.Handler is registered for the given host.
	// If yes, use it to handle the request.

	hostParts := strings.Split(r.Host, ".")

	// {sub}.{domain}.{tld}
	if len(hostParts) < 3 {
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
	// load config
	configPath := os.Getenv("MACHINABLE_CONFIG_PATH")
	file, _ := ioutil.ReadFile(configPath)

	config := &config.AppConfig{}
	json.Unmarshal([]byte(file), &config)

	// secrets come from environment
	config.LoadSecrets()

	// use postgres client
	datastore, err := postgres.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PW"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"),
	)

	if err != nil {
		log.Fatal(err)
	}

	// create a new redis client
	cache := redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// ping the redis server
	pong, err := cache.Ping().Result()

	// fail quickly if ping fails
	if err != nil {
		log.Fatal(pong, err)
	}

	// create event processor
	processor := events.NewProcessor(cache, datastore)

	// process web hook results
	go func() {
		err := processor.ProcessResults()
		// fail out
		log.Fatal(err)
	}()

	// switch routers based on subdomain
	hostSwitch := make(HostSwitch)

	// manage is for the management application api, i.e. project/team management
	hostSwitch["manage"] = management.CreateRoutes(datastore, cache, config)
	// all other subdomains will be treated as project names, and use the project routes
	hostSwitch["*"] = projects.CreateRoutes(datastore, cache, processor, config)

	log.Fatal(http.ListenAndServe(":5001", hostSwitch))
}
