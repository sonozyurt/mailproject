package main

import (
	"github.com/alexedwards/scs/v2"
	"log"
	"mailprojesi/datamanagement"
	"sync"
)

type Config struct {
	Users        datamanagement.User
	errorLog     *log.Logger
	infoLog      *log.Logger
	session      *scs.SessionManager
	mailSendChan chan message
	mailDoneChan chan bool
	MailServer   mailServer
	wg           *sync.WaitGroup
}
