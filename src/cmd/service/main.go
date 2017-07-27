// The command line interface to the microservice
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/saschagrunert/ccli"
	"github.com/urfave/cli"

	"github.com/quaponatech/golang-extensions/grpcservice"
	"github.com/quaponatech/golang-extensions/server"
	"github.com/quaponatech/microservice-template/src"
)

var (
	authorName     = "Prename Surname"               // TODO
	authorMail     = "p.surname@quapona.com"         // TODO
	appUsage       = "A microservice template"       // TODO
	appName        = "microservice"                  // TODO
	appDescription = "This microservice provides..." // TODO
)

func main() {
	// Create a new command line application
	app := ccli.NewApp()

	// Set the basic CLI options
	app.Authors = []cli.Author{{Name: authorName, Email: authorMail}}
	app.Copyright = fmt.Sprintf("© %d quapona technologies GmbH", time.Now().Year())
	app.Description = appDescription
	app.Name = appName
	app.Usage = appUsage
	app.Version = microservice.Version

	// Create a new microservice instance
	service := microservice.MicroService{}

	// Add and parse the flags
	var (
		port         int
		useTLS       bool
		certFile     string
		privKeyFile  string
		caFile       string
		logDirectory string
		logLevel     int
		dryRun       bool
	)

	// TODO: This is just for demo purposes
	clientCredentials := &grpcservice.ConnectionInfo{}

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "port,p",
			Usage:       "server port of the service",
			Value:       42302,
			Destination: &port,
		},
		cli.BoolFlag{
			Name:        "usetls,t",
			Usage:       "use TLS if true, else plain TCP communication",
			Destination: &useTLS,
		},
		cli.StringFlag{
			Name:        "certfile,c",
			Usage:       "TLS certificate file",
			Destination: &certFile,
		},
		cli.StringFlag{
			Name:        "privkeyfile,k",
			Usage:       "TLS private key file",
			Destination: &privKeyFile,
		},
		cli.StringFlag{
			Name:        "cafile,k",
			Usage:       "TLS CA file",
			Destination: &caFile,
		},
		cli.StringFlag{
			Name:        "logdir,d",
			Usage:       "log directory",
			Destination: &logDirectory,
		},
		cli.IntFlag{
			Name:        "loglevel,l",
			Usage:       "defines the log output level from Debug (0) to Quiet (5)",
			Destination: &logLevel,
			Value:       0,
		},
		cli.BoolFlag{
			Name:        "dry-run, n",
			Usage:       "do a dry run without starting the server",
			Destination: &dryRun,
		},
		// TODO: Adapt all client connections here
		cli.StringFlag{
			Name:        "client-ip",
			Usage:       "client IP to connect",
			Value:       "localhost",
			Destination: &clientCredentials.IP,
		},
		cli.StringFlag{
			Name:        "client-port",
			Usage:       "client port to connect",
			Value:       "42303",
			Destination: &clientCredentials.Port,
		},
		cli.BoolFlag{
			Name:        "client-usetls",
			Usage:       "use TLS to connect with the client if true, else plain TCP communication",
			Destination: &clientCredentials.UseTLS,
		},
		cli.StringFlag{
			Name:        "client-certfile",
			Usage:       "TLS certificate file for the client connection",
			Destination: &clientCredentials.CertFile,
		},
	}

	// TODO: Get external service credentials from CLI arguments
	// This could be any other microservices
	client := &microservice.GrpcClientMicroService{AccessInfo: clientCredentials}

	// Run the library when CLI parsing is done
	app.Action = func(_ *cli.Context) error {
		// Setup the service
		serverName := "Microservice"
		serverInstance := grpcservice.NewMutualGRPCServer(useTLS, certFile, privKeyFile, caFile, port)
		serverLogger := server.NewLogger(serverName,
			logDirectory,
			appName+".log",
			make(chan server.Status),
			make(chan error, 100),
			make(chan string, 1000),
			make(chan string, 10000),
			make(chan string, 10000),
			logLevel)

		if err := service.Setup(serverName, serverInstance, serverLogger, client); err != nil {
			log.Fatal(err)
		}

		// Start the microservice
		if !dryRun {
			if err := service.Serve(); err != nil {
				log.Fatal("could not start microservice")
			}
		} else {
			// TODO: Extend the module testing if needed
			log.Println("Everything seems okay.")
		}
		return nil
	}

	// Run the command line application
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
