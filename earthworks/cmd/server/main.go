package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/jmoiron/sqlx"
	"github.com/namsral/flag"
	"github.com/stephenhillier/geoprojects/earthworks/db"
	"github.com/stephenhillier/geoprojects/earthworks/server"
)

func main() {
	// set config parameters
	// the flag library grabs values either from command line args, env variables, or the default specified here
	// see github.com/namsral/flag
	conf := &server.Config{}
	var dbuser, dbpass, dbname, dbhost, dbport, dbsslmode string

	flag.StringVar(&conf.DBDriver, "dbdriver", "postgres", "database driver")
	flag.StringVar(&dbuser, "dbuser", "geo", "database username")
	flag.StringVar(&dbpass, "dbpass", "", "database password")
	flag.StringVar(&dbname, "dbname", "geo", "database name")
	flag.StringVar(&dbhost, "dbhost", "127.0.0.1", "database service host")
	flag.StringVar(&dbport, "dbport", "5432", "database service port")
	flag.StringVar(&dbsslmode, "dbsslmode", "disable", "database ssl mode")
	flag.StringVar(&conf.AuthAudience, "auth_audience", "https://earthworks.islandcivil.com", "authentication service audience claim")
	flag.StringVar(&conf.AuthIssuer, "auth_issuer", "https://earthworks.auth0.com/", "authentication service issuer claim")
	flag.StringVar(&conf.AuthJWKSEndpoint, "jwks_endpoint", "https://earthworks.auth0.com/.well-known/jwks.json", "authentication JWKS endpoint")
	flag.Parse()

	conf.DBConn = fmt.Sprintf("%s://%s:%s@%s:%v/%s?sslmode=%s", conf.DBDriver, dbuser, dbpass, dbhost, dbport, dbname, dbsslmode)

	// get new certificate when server initially starts
	// see auth.go
	cert, err := conf.GetCert(nil)
	if err != nil {
		log.Panic(err)
	}

	conf.AuthCert = cert

	conf.AuthGroupClaim = conf.AuthAudience + "/claims/authorization/groups"
	conf.AuthRoleClaim = conf.AuthAudience + "/claims/authorization/roles"

	dbConf := db.Config{
		Conn:   conf.DBConn,
		Driver: conf.DBDriver,
	}

	database, err := sqlx.Open(dbConf.Driver, dbConf.Conn)
	if err != nil {
		log.Panic(err)
	}

	store, err := db.NewDatastore(database)
	if err != nil {
		log.Panic(err)
	}

	svc, err := server.NewEarthworksService(store, conf)
	if err != nil {
		log.Panic(err)
	}

	h := http.Server{Addr: ":8000", Handler: svc.Router}

	log.Printf("Starting HTTP server on port 8000.\n")
	log.Printf("Press CTRL+C to stop.")
	go func() {
		if err := h.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Server is listening; Wait here for interrupt signal...
	<-stop
	log.Println("Shutting down...")
	h.Shutdown(context.Background())
	log.Println("Server stopped.")

}
