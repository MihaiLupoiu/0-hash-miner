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
			conn.WriteSHA1String(authdata, line, configuration.UserConfig.Name)

		case "MAILNUM":
			// here you specify, how many email addresses you want to send
			// each email is asked separately up to the number specified in MAILNUM
			conn.WriteSHA1String(authdata, line, strconv.Itoa(len(configuration.UserConfig.Mails)))

		case "MAIL1":
			conn.WriteSHA1String(authdata, line, configuration.UserConfig.Mails[1])

		case "MAIL2":
			conn.WriteSHA1String(authdata, line, configuration.UserConfig.Mails[2])

		case "SKYPE":
			// here please specify your Skype account for the interview, or N/A
			// in case you have no Skype account
			conn.WriteSHA1String(authdata, line, configuration.UserConfig.Skype)

		case "BIRTHDATE":
			// here please specify your birthdate in the format %d.%m.%Y
			conn.WriteSHA1String(authdata, line, configuration.UserConfig.BirthDate)

		case "COUNTRY":
			// country where you currently live and where the specified address is
			// please use only the names from this web site:
			//   https://www.countries-ofthe-world.com/all-countries.html
			conn.WriteSHA1String(authdata, line, configuration.UserConfig.Country)

		case "ADDRNUM":
			// specifies how many lines your address has, this address should
			// be in the specified country
			conn.WriteSHA1String(authdata, line, strconv.Itoa(len(configuration.UserConfig.Addess)))

		case "ADDRLINE1":
			conn.WriteSHA1String(authdata, line, configuration.UserConfig.Addess[1])

		case "ADDRLINE2":
			conn.WriteSHA1String(authdata, line, configuration.UserConfig.Addess[2])

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
