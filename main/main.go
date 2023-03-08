package main

import (
	"database/sql"
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"mailprojesi/datamanagement"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"net/http"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"os"
)

const Port = ":8050"

func main() {
	db := initDB()

	config := Config{
		Users:        datamanagement.SetUpUsers(db),
		errorLog:     log.New(os.Stdout, "errors\t", log.Ldate|log.Ltime|log.Llongfile),
		infoLog:      log.New(os.Stdout, "infoLog:\t", log.Ldate|log.Ltime|log.Llongfile),
		session:      createSession(),
		mailSendChan: make(chan message, 10),
		mailDoneChan: make(chan bool),
		wg:           &sync.WaitGroup{},
	}

	config.MailServer = createMailServer()

	go config.concurrencyMailing()

	config.infoLog.Println("Mail server created.")

	go config.shuttingDown()

	config.serve()

}

func initDB() *sql.DB {
	dsn := os.Getenv("DSN")
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func (config *Config) serve() {
	server := http.Server{
		Addr:    Port,
		Handler: config.routing(),
	}
	config.infoLog.Println("Starting web server...")
	err := server.ListenAndServe()
	config.err(err)
}

func createSession() *scs.SessionManager {
	gob.Register(datamanagement.User{})
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = true
	sessionManager.Store = memstore.New()
	return sessionManager
}

func createMailServer() mailServer {
	mailServer := mailServer{
		host:        "localhost",
		port:        1025,
		ConnTimeOut: 10 * time.Second,
		SendTimeOut: 10 * time.Second,
	}
	return mailServer
}

func (config *Config) err(err error) {
	if err != nil {
		config.errorLog.Println(err)
	}
}

func (config *Config) shuttingDown() {
	q := make(chan os.Signal, 1)
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
	_ = <-q
	config.shutdown()
}
func (config *Config) shutdown() {
	config.wg.Wait()
	config.infoLog.Println("closing channels")
	config.mailDoneChan <- true
	close(config.mailDoneChan)
	close(config.mailSendChan)
}
