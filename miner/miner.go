package miner

import (
	"bufio"
	"fmt"
	"log"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MihaiLupoiu/interview-exasol/config"
	"github.com/MihaiLupoiu/interview-exasol/connection"
	"github.com/MihaiLupoiu/interview-exasol/solver"
	"github.com/MihaiLupoiu/interview-exasol/utils"
	"github.com/paulbellamy/ratecounter"
)

type Miner struct {
	Authdata   string
	Conn       *connection.Connection
	Counter    *ratecounter.RateCounter
	UserConfig config.UserConfig
}

func connect(configuration config.Data) (*connection.Connection, error) {
	conn, err := connection.Dial(configuration.Crt, configuration.Key, configuration.Endpoint)
	if err != nil {
		log.Fatalf("failed to connect: %s", err.Error())
		return nil, err
	}

	log.Printf("connect to %s succeed", configuration.Endpoint)
	conn.PrintConnState()

	return conn, err
}

// Init miner with configuration with connection data and user information.
func Init(configuration config.Data) (*Miner, error) {
	conn, err := connect(configuration)
	if err != nil {
		return nil, err
	}
	return &Miner{
		Authdata:   "",
		Conn:       conn,
		Counter:    ratecounter.NewRateCounter(1 * time.Second),
		UserConfig: configuration.UserConfig,
	}, err
}

func (ctx *Miner) Run() error {
	defer ctx.Conn.Close()

	stopHashRate := make(chan bool)
	go utils.HashRate(ctx.Counter, stopHashRate)

	connTextReader := textproto.NewReader(bufio.NewReader(ctx.Conn))
	for {
		// read one line (ended with \n or \r\n)
		line, _ := connTextReader.ReadLine()
		fmt.Println(line)

		switch line {
		case "HELO":
			ctx.Conn.WriteString("EHLO")

		case "END":
			// if you get this command, then your data was submitted
			ctx.Conn.WriteString("OK")
			os.Exit(0)

		// the rest of the data server requests are required to identify you
		// and get basic contact information
		case "NAME":
			// as the response to the NAME request you should send your full name
			// including first and last name separated by single space
			ctx.Conn.WriteSHA1String(ctx.Authdata, line, ctx.UserConfig.Name)

		case "MAILNUM":
			// here you specify, how many email addresses you want to send
			// each email is asked separately up to the number specified in MAILNUM
			ctx.Conn.WriteSHA1String(ctx.Authdata, line, strconv.Itoa(len(ctx.UserConfig.Mails)))

		case "MAIL1":
			ctx.Conn.WriteSHA1String(ctx.Authdata, line, ctx.UserConfig.Mails[1])

		case "MAIL2":
			ctx.Conn.WriteSHA1String(ctx.Authdata, line, ctx.UserConfig.Mails[2])

		case "SKYPE":
			// here please specify your Skype account for the interview, or N/A
			// in case you have no Skype account
			ctx.Conn.WriteSHA1String(ctx.Authdata, line, ctx.UserConfig.Skype)

		case "BIRTHDATE":
			// here please specify your birthdate in the format %d.%m.%Y
			ctx.Conn.WriteSHA1String(ctx.Authdata, line, ctx.UserConfig.BirthDate)

		case "COUNTRY":
			// country where you currently live and where the specified address is
			// please use only the names from this web site:
			//   https://www.countries-ofthe-world.com/all-countries.html
			ctx.Conn.WriteSHA1String(ctx.Authdata, line, ctx.UserConfig.Country)

		case "ADDRNUM":
			// specifies how many lines your address has, this address should
			// be in the specified country
			ctx.Conn.WriteSHA1String(ctx.Authdata, line, strconv.Itoa(len(ctx.UserConfig.Addess)))

		case "ADDRLINE1":
			ctx.Conn.WriteSHA1String(ctx.Authdata, line, ctx.UserConfig.Addess[1])

		case "ADDRLINE2":
			ctx.Conn.WriteSHA1String(ctx.Authdata, line, ctx.UserConfig.Addess[2])

		default:
			if strings.HasPrefix(line, "POW") {
				fmt.Println(line)
				fmt.Println("Time to work")
				args := strings.Fields(line)

				ctx.Authdata = args[1]
				difficulty, err := strconv.Atoi(args[2])
				if err != nil {
					log.Fatalf("Difficulty of POW not integer: %s", err.Error())
				}

				fmt.Println(ctx.Authdata, difficulty)
				for {
					// generate short random string, server accepts all utf-8 characters,
					// except [\n\r\t ], it means that the suffix should not contain the
					// characters: newline, carriege return, tab and space
					suffix, _ := utils.RandStringRunes(30)
					ctx.Counter.Incr(1)

					if solver.Check(ctx.Authdata, suffix, difficulty) != "" {
						fmt.Printf("Authdata: %s\n Suffix: %s\n", ctx.Authdata, suffix)
						ctx.Conn.WriteString(suffix)
						// Stop goroutine hashRate
						stopHashRate <- true
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
