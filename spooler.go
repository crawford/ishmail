package ishmail

import (
	"log"
	"net/smtp"
	"sync"
)

const (
	DefaultSpoolLength = 16
)

// Start the default spooler.
// This will create a default spooler if one does not already exist.
func StartSpooler() {
	getSpooler().Start()
}

// Stop the default spooler and wait for all remaining messages to be sent.
// This will create a default spooler if one does not already exist.
func StopSpooler() {
	getSpooler().Stop()
}

// Stop the default spooler, clear all queued messages, and wait for the
// currently-pending message to finish.
// This will create a default spooler if one does not already exist.
func TerminateSpooler() {
	getSpooler().Terminate()
}

// Queue a message into the default spooler. This will block until the spool
// isn't full.
// This will create a default spooler if one does not already exist.
func Spool(email Emailer) {
	getSpooler().Spool(email)
}

// Configure the default spooler to use the specified authentication and
// remote host.
func ConfigSpooler(auth smtp.Auth, addr string) {
	s := getSpooler()
	s.auth = auth
	s.addr = addr
}

var spool *Spooler

func getSpooler() *Spooler {
	if spool == nil {
		spool = CreateSpooler(nil, "", DefaultSpoolLength)
	}
	return spool
}

type Spooler struct {
	spool chan Emailer
	wg    sync.WaitGroup
	auth  smtp.Auth
	addr  string
}

// Create a new spooler with the specified queue length.
func CreateSpooler(auth smtp.Auth, addr string, spoolLength int) *Spooler {
	return &Spooler{
		spool: make(chan Emailer, spoolLength),
		auth:  auth,
		addr:  addr,
	}
}

// Start the spooler process in a goroutine. This can be called multiple times
// to create multiple, concurrent spooler processes. All of these processes
// will be owned by this spooler.
func (s *Spooler) Start() {
	s.wg.Add(1)
	go s.run()
}

// Prevent new messages from being queued and wait for all of the queued
// messages to be sent.
func (s *Spooler) Stop() {
	close(s.spool)
	s.wg.Wait()
}

// Prevent new messages from being queued, clear the queue of pending messages,
// and wait for all of the current operations to complete.
func (s *Spooler) Terminate() {
	close(s.spool)
	for _ = range s.spool {
	}
	s.wg.Wait()
}

// Queue a message.
func (s *Spooler) Spool(email Emailer) {
	s.spool <- email
}

func (s *Spooler) run() {
SendLoop:
	for {
		select {
		case email, open := <-s.spool:
			if !open {
				break SendLoop
			}
			if err := Send(email, s.auth, s.addr); err != nil {
				log.Printf("Failed to send email (%s)\n", err)
			}
		}
	}
	s.wg.Done()
}
