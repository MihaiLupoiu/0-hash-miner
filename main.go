package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MihaiLupoiu/interview-exasol/config"
	"github.com/MihaiLupoiu/interview-exasol/connection"
	"github.com/MihaiLupoiu/interview-exasol/solver"
	"github.com/google/uuid"
	"github.com/paulbellamy/ratecounter"
)

var authdata = ""
var counter = ratecounter.NewRateCounter(1 * time.Second)

func hashRate() {
	for range time.Tick(time.Second * 1) {
		fmt.Println("Hash rate:", counter.Rate())
	}
}

const (
	// exitFail is the exit code if the program
	// fails.
	exitFail = 1
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}

func run(args []string, stdout io.Writer) error {
	configuration := config.Get()

	// -----

	conn, err := connection.Dial(configuration.Crt, configuration.Key, configuration.Endpoint)
	if err != nil {
		log.Fatalf("failed to connect: %s", err.Error())
	}
	defer conn.Close()
	log.Printf("connect to %s succeed", configuration.Endpoint)
	conn.PrintConnState()

	// -----

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	// defer conn.Close()

	for {
		// read one line (ended with \n or \r\n)
		line, _ := tp.ReadLine()
		fmt.Println(line)

		switch line {
		case "HELO":
			conn.WriteString("EHLO")

		case "END":
			// if you get this command, then your data was submitted
			conn.WriteString("OK")
			os.Exit(0)

		// the rest of the data server requests are required to identify you
		// and get basic contact information
		case "NAME":
			// as the response to the NAME request you should send your full name
			// including first and last name separated by single space
			// conn.write(SHA1(authdata+args[1]) + " " + "My name\n")

		case "MAILNUM":
			// here you specify, how many email addresses you want to send
			// each email is asked separately up to the number specified in MAILNUM
			// conn.write(SHA1(authdata+args[1]) + " " + "2\n")

		case "MAIL1":
			// conn.write(SHA1(authdata+args[1]) + " " + "my.name@example.com\n")

		case "MAIL2":
			// conn.write(SHA1(authdata+args[1]) + " " + "my.name2@example.com\n")

		case "SKYPE":
			// here please specify your Skype account for the interview, or N/A
			// in case you have no Skype account
			// conn.write(SHA1(authdata+args[1]) + " " + "my.name@example.com\n")

		case "BIRTHDATE":
			// here please specify your birthdate in the format %d.%m.%Y
			// conn.write(SHA1(authdata+args[1]) + " " + "01.02.2017\n")

		case "COUNTRY":
			// country where you currently live and where the specified address is
			// please use only the names from this web site:
			//   https://www.countries-ofthe-world.com/all-countries.html
			// conn.write(SHA1(authdata+args[1]) + " " + "Germany\n")

		case "ADDRNUM":
			// specifies how many lines your address has, this address should
			// be in the specified country
			// conn.write(SHA1(authdata+args[1]) + " " + "2\n")

		case "ADDRLINE1":
			// conn.write(SHA1(authdata+args[1]) + " " + "Long street 3\n")

		case "ADDRLINE2":
			// conn.write(SHA1(authdata+args[1]) + " " + "32345 Big city\n")

		default:
			if strings.HasPrefix(line, "POW") {
				fmt.Println(line)
				fmt.Println("Time to work")
				args := strings.Fields(line)

				authdata = args[1]
				difficulty, err := strconv.Atoi(args[2])
				if err != nil {
					log.Fatalf("Difficulty of POW not integer: %s", err.Error())
				}

				fmt.Println(authdata, difficulty)
				go hashRate()
				for {
					// generate short random string, server accepts all utf-8 characters,
					// except [\n\r\t ], it means that the suffix should not contain the
					// characters: newline, carriege return, tab and space
					suffix := uuid.New().String()
					counter.Incr(1)

					if solver.Check(authdata, suffix, difficulty) != "" {
						fmt.Printf("Authdata: %s\n Suffix: %s\n", authdata, suffix)
						conn.WriteString(suffix)
						return nil
					}
				}
			} else if strings.HasPrefix(line, "ERROR") {
				fmt.Println(line)
				return nil
			}
		}
	}
}
