package main

import (
        "context"
        "fmt"
        "log"
        "net/http"
        "os"
        "os/signal"
        "syscall"
        "time"
	"flag"

        "github.com/gorilla/mux"
        "gopkg.in/natefinch/lumberjack.v2"
	"github.com/nats-io/nats.go"

var natsServer = flag.String("natsServer", "nats://operator-nats:4222", "Comma seperated nats servers")

func handler(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query()
        name := query.Get("name")
        if name == "" {
                name = "Guest"
        }
        log.Printf("Received request for %s\n", name)
        
        log.Printf("Connect to nats streaming server")
        connectToNats()
        
        w.Write([]byte(fmt.Sprintf("Hello, %s\n", name)))
}

func main() {
        // Create Server and Route Handlers
        r := mux.NewRouter()

        r.HandleFunc("/", handler)

        srv := &http.Server{
                Handler:      r,
                Addr:         ":8080",
                ReadTimeout:  10 * time.Second,
                WriteTimeout: 10 * time.Second,
        }

        // Configure Logging
        LOG_FILE_LOCATION := os.Getenv("LOG_FILE_LOCATION")
        if LOG_FILE_LOCATION != "" {
                log.SetOutput(&lumberjack.Logger{
                        Filename:   LOG_FILE_LOCATION,
                        MaxSize:    500, // megabytes
                        MaxBackups: 3,
                        MaxAge:     28,   //days
                        Compress:   true, // disabled by default
                })
        }

        // Start Server
        go func() {
                log.Println("Starting Server")
                if err := srv.ListenAndServe(); err != nil {
                        log.Fatal(err)
                }
        }()

        // Graceful Shutdown
        waitForShutdown(srv)
}

func connectToNats() {
        //connect to nats server
        fmt.Printf("Establishing connection to the Nats server....")
        nc, err := nats.Connect(*natsServer)
        
        checkErr(err, "At nats connect")
        
        //publish to nats
        nc.Publish("test", "some test data to nats")
        
        nc.Close()
        
}

func waitForShutdown(srv *http.Server) {
        interruptChan := make(chan os.Signal, 1)
        signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

        // Block until we receive our signal.
        <-interruptChan

        // Create a deadline to wait for.
        ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
        defer cancel()
        srv.Shutdown(ctx)

        log.Println("Shutting down")
        os.Exit(0)
}

func checkErr(err error, point string) {
	if err != nil {
		log.Println(point)
		fmt.Println(err)
		log.Println(err)
		log.Fatal(err)

	}
}
