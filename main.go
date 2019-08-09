package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

// subscribeStreaming subscribes to specified NATS streaming subjects.
func subscribeStreaming(connection *nats.Conn, clusterID string, clientID string, subjects []string) {
	log.Printf("NATS streaming using clusterID=%q and clientID=%q", clusterID, clientID)
	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(connection))
	if err != nil {
		log.Fatal(err)
	}

	for _, subject := range subjects {
		log.Printf("Subscribe to stream %q", subject)
		_, err = sc.Subscribe(subject, func(msg *stan.Msg) {
			log.Println(msg)
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

// subscribe subscribes to specified NATS subjects.
func subscribe(connection *nats.Conn, subjects []string) {
	for _, subject := range subjects {
		log.Printf("Subscribe to %q", subject)
		_, err := connection.Subscribe(subject, func(msg *nats.Msg) {
			log.Println(msg)
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	clusterID := flag.String("cluster_id", "test-cluster", "NATS streaming cluster ID")
	clientID := flag.String("client_id", "nats-tester-id", "NATS streaming client ID")
	needStreaming := flag.Bool("streaming", false, "Use NATS streaming")
	natsURL := nats.DefaultURL
	if os.Getenv("NATS_URL") != "" {
		natsURL = os.Getenv("NATS_URL")
	}
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of nats-subscribe <subject>...\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nFor NATS connection use "+
			"environment variable NATS_URL (default %q)\n", nats.DefaultURL)
	}
	flag.Parse()
	subjects := flag.Args()
	if len(subjects) == 0 {
		log.Fatal("No subjects provided (see --help)")
	}

	log.Printf("Connecting to %q...", natsURL)
	connection, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	if *needStreaming {
		subscribeStreaming(connection, *clusterID, *clientID, subjects)
	} else {
		subscribe(connection, subjects)
	}

	exit := make(chan bool, 1)
	connection.SetDisconnectErrHandler(func(conn *nats.Conn, err error) {
		log.Fatal("disconnected", err)
		exit <- true
	})
	<-exit
}
