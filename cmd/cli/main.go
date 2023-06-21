package main

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/emersion/go-smtp"
	"golang.org/x/net/proxy"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Config struct for SMTP credentials
type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Server   string `yaml:"server"`
	Port     string `yaml:"port"`
}

// LoadConfig function reads SMTP credentials from a YAML file
func LoadConfig() Config {
	file, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config file: %v", err)
	}

	return config
}

// The Backend implements SMTP server methods.
type Backend struct{}

func (bkd *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return &Session{}, nil
}

// A Session is returned after AUTH.
type Session struct {
	from string
	to   []string
}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	s.to = append(s.to, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	config := LoadConfig()

	// Establish a connection to the SOCKS5 proxy.
	dialer, err := proxy.SOCKS5("tcp", "localhost:1080", nil, proxy.Direct)
	if err != nil {
		log.Fatal(err)
	}

	// Set up an SMTP client.
	smtpClient, err := smtp.Dial(config.Server+":"+config.Port, smtp.DialWithDialer(dialer.Dial))
	if err != nil {
		log.Fatal(err)
	}
	defer smtpClient.Close()

	// Authenticate with the server.
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Server)
	if err := smtpClient.Auth(auth); err != nil {
		log.Fatal(err)
	}

	// Set the sender and recipients.
	if err := smtpClient.Mail(s.from); err != nil {
		log.Fatal(err)
	}
	for _, recipient := range s.to {
		if err := smtpClient.Rcpt(recipient); err != nil {
			log.Fatal(err)
		}
	}

	// Send the email body.
	wc, err := smtpClient.Data()
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(wc, r)
	if err != nil {
		log.Fatal(err)
	}
	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func main() {
	be := &Backend{}

	s := smtp.NewServer(be)

	s.Addr = ":1025"
	s.Domain = "localhost"
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

