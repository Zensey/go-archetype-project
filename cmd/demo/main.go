package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"time"

	"github.com/Zensey/go-archetype-project/pkg/domain"
	"github.com/Zensey/slog"

	"net/smtp"
	//"github.com/emersion/go-smtp"
	dkim "github.com/toorop/go-dkim"
	gomail "gopkg.in/mail.v2"
)

var version string

type A struct {
	t     int64
	state int
}
type M struct {
	t [1024 * 1024 * 8]int32
	m [1024 * 1024 * 8]byte
}

func statemachine() {
	m := M{}
	for i := 0; i < 1024*1024*8; i++ {
		m.m[int64(i)] = ^m.m[int64(i)]
	}

	//m := make(map[int64]A, 0)
	//for i := 0; i < 10000*1000; i++ {
	//	m[int64(i)] = A{}
	//}

	time.Sleep(time.Minute)
}

func idgen(l slog.Logger) {
	b := bytes.Buffer{}
	p := domain.Producer{}
	for i := 1; i <= 30; i++ {
		b.Reset()
		p.GetNewMsgID(&b)

		l.Infof("%s", b.String())
		time.Sleep(10 * time.Millisecond)
	}
	l.Info(b.String())
}

//func Example() {
//	// Connect to the remote SMTP server.
//	c, err := smtp.Dial("gmail-smtp-in.l.google.com:25")
//	if err != nil {
//		log.Fatal("1>", err)
//	}
//
//	// Set the sender and recipient first
//	if err := c.Mail("jnashicq@zensey.cloudns.cl"); err != nil {
//		log.Fatal("2>", err)
//	}
//	if err := c.Rcpt("jnashicq@gmail.com"); err != nil {
//		log.Fatal("3>", err)
//	}
//
//	// Send the email body.
//	wc, err := c.Data()
//	if err != nil {
//		log.Fatal("4>", err)
//	}
//	_, err = fmt.Fprintf(wc, "This is the email body")
//	if err != nil {
//		log.Fatal("5>", err)
//	}
//	err = wc.Close()
//	if err != nil {
//		log.Fatal("6>", err)
//	}
//
//	// Send the QUIT command and close the connection.
//	err = c.Quit()
//	if err != nil {
//		log.Fatal(err)
//	}
//}

var (
	from       = "jnashicq@zensey.cloudns.cl"
	msg        = []byte("dummy message")
	recipients = []string{"jnashicq@gmail.com"}
	hostname   = "gmail-smtp-in.l.google.com:25" //"smtp.gmail.com:25" //
)

func ExamplePlainAuth() {

	//err := smtp.SendMail(hostname+":25", nil, from, recipients, msg)
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Set up authentication information.
	//auth := sasl.NewPlainClient("", "user@example.com", "password")
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.

	//msg := strings.NewReader("To: jnashicq@gmail.com\r\n" +
	//	"Subject: discount Gophers!\r\n" +
	//	"\r\n" +
	//	"This is the email body.\r\n")
	//err := smtp.SendMail(hostname+":25", nil, from, recipients, msg)
	//if err != nil {
	//	log.Fatal(err)
	//}

	c, err := smtp.Dial(hostname)
	if err != nil {
		log.Fatal("1>", err)
	}

	// Set the sender and recipient first
	if err := c.Mail(from); err != nil {
		log.Fatal("2>", err)
	}
	if err := c.Rcpt("jnashicq@gmail.com"); err != nil {
		log.Fatal("3>", err)
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal("4>", err)
	}
	_, err = fmt.Fprintf(wc, "This is the email body")
	if err != nil {
		log.Fatal("5>", err)
	}
	err = wc.Close()
	if err != nil {
		log.Fatal("6>", err)
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		log.Fatal(err)
	}
}

func goMail() {

	// generate key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("Cannot generate RSA key\n")
		return
	}
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privatekey),
	}
	privateKeyBytes := pem.EncodeToMemory(privateKeyBlock)

	const CRLF = "\r\n"
	var emailBase = "Received: (qmail 28277 invoked from network); 1 May 2015 09:43:37 -0000" + CRLF +
		"Received: (qmail 21323 invoked from network); 1 May 2015 09:48:39 -0000" + CRLF +
		"Received: from mail483.ha.ovh.net (b6.ovh.net [213.186.33.56])" + CRLF +
		" by mo51.mail-out.ovh.net (Postfix) with SMTP id A6E22FF8934" + CRLF +
		" for <toorop@toorop.fr>; Mon,  4 May 2015 14:00:47 +0200 (CEST)" + CRLF +
		"MIME-Version: 1.0" + CRLF +
		"Date: Fri, 1 May 2015 11:48:37 +0200" + CRLF +
		"Message-ID: <CADu37kTXBeNkJdXc4bSF8DbJnXmNjkLbnswK6GzG_2yn7U7P6w@tmail.io>" + CRLF +
		"Subject: Test DKIM" + CRLF +
		"From: =?UTF-8?Q?St=C3=A9phane_Depierrepont?= <toorop@tmail.io>" + CRLF +
		"To: =?UTF-8?Q?St=C3=A9phane_Depierrepont?= <toorop@toorop.fr>" + CRLF +
		"Content-Type: text/plain; charset=UTF-8" + CRLF + CRLF +
		"Hello world" + CRLF +
		"line with trailing space         " + CRLF +
		"line with           space         " + CRLF +
		"-- " + CRLF +
		"Toorop" + CRLF + CRLF + CRLF + CRLF + CRLF + CRLF

	email := []byte(emailBase)
	emailToTest := append([]byte(nil), email...)

	// email is the email to sign (byte slice)
	// privateKey the private key (pem encoded, byte slice )
	options := dkim.NewSigOptions()
	options.PrivateKey = privateKeyBytes
	options.Domain = "mydomain.tld"
	options.Selector = "myselector"
	//options.SignatureExpireIn = 3600
	//options.BodyLength = 7
	options.Headers = []string{"from"}
	options.Algo = "rsa-sha256"
	options.AddSignatureTimestamp = true
	options.Canonicalization = "relaxed/relaxed"

	err = dkim.Sign(&emailToTest, options)
	fmt.Println(string(emailToTest))
	if err != nil {
		fmt.Println("err>", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", recipients[0])
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/plain", "Hello!")
	d := gomail.Dialer{Host: "gmail-smtp-in.l.google.com", Port: 25}
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func main() {
	l := slog.ConsoleLogger()
	l.SetLevel(slog.LevelTrace)
	l.Infof("Hello, World ! Version: %s", version)

	//statemachine()
	//idgen(l)

	//Example()
	//ExamplePlainAuth()
	goMail()

	return
}
