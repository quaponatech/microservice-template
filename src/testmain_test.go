package microservice

import (
	"flag"
	"log"
	"os"
	"testing"
)

/* Test Data */
var clusterIP string
var clusterPort int

/* Arguments */
func init() {
	flag.StringVar(&clusterIP, "ip", "127.0.0.1", "The service ip.")
	flag.IntVar(&clusterPort, "port", 0, "The service port.")
}

/* Test Funtions */
func TestMain(m *testing.M) {
	flag.Parse()

	log.Println("* Using parameters:")
	log.Println("* -ip = ", clusterIP)
	log.Println("* -port = ", clusterPort)

	// Run tests
	testreturn := m.Run()

	os.Exit(testreturn)
}
