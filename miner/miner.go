package miner

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MihaiLupoiu/interview-exasol/connection"
	"github.com/MihaiLupoiu/interview-exasol/worker"
	"github.com/paulbellamy/ratecounter"
)

// Miner has all the basic information to start the search for the SHA1
type Miner struct {
	Authdata   string
	Conn       *connection.Connection
	Counter    *ratecounter.RateCounter
	UserConfig UserConfig
	WPool      worker.Pool
	incoming   chan string
	outcoming  chan string
}

var (
	minRandomStringLength = 5
	maxRandomStringLength = 64
)

// connect creates the TLS connection required to the server in order to process the work.
func connect(configuration Data) (*connection.Connection, error) {
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
func Init(configuration Data) (*Miner, error) {
	conn, err := connect(configuration)
	if err != nil {
		return nil, err
	}
	return &Miner{
		Authdata:   "",
		Conn:       conn,
		Counter:    ratecounter.NewRateCounter(1 * time.Second),
		UserConfig: configuration.UserConfig,
		WPool:      worker.New(configuration.Workers),
		incoming:   make(chan string, 1),
		outcoming:  make(chan string, 1),
	}, err
}

// Run miner will connect to the server, initialize the workes when request POW is received and
// start to search for the SHA1 with the given dificulty.
func (ctx *Miner) Run() error {
	defer ctx.Conn.Close()

	precessingInterval := time.Hour * time.Duration(2)
	respondInterval := time.Second * time.Duration(6)

	t := time.NewTimer(time.Hour)

	go ctx.readConnData()
	for {
		select {
		case line := <-ctx.incoming:
			log.Println(line)
			args := strings.Fields(line)
			t.Reset(respondInterval)

			switch args[0] {
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
				ctx.Conn.WriteSHA1String(ctx.Authdata, args[1], ctx.UserConfig.Name)

			case "MAILNUM":
				// here you specify, how many email addresses you want to send
				// each email is asked separately up to the number specified in MAILNUM
				ctx.Conn.WriteSHA1String(ctx.Authdata, args[1], strconv.Itoa(len(ctx.UserConfig.Mails)))

			case "MAIL1":
				ctx.Conn.WriteSHA1String(ctx.Authdata, args[1], ctx.UserConfig.Mails[0])

			case "MAIL2":
				ctx.Conn.WriteSHA1String(ctx.Authdata, args[1], ctx.UserConfig.Mails[1])

			case "SKYPE":
				// here please specify your Skype account for the interview, or N/A
				// in case you have no Skype account
				ctx.Conn.WriteSHA1String(ctx.Authdata, args[1], ctx.UserConfig.Skype)

			case "BIRTHDATE":
				// here please specify your birthdate in the format %d.%m.%Y
				ctx.Conn.WriteSHA1String(ctx.Authdata, args[1], ctx.UserConfig.BirthDate)

			case "COUNTRY":
				// country where you currently live and where the specified address is
				// please use only the names from this web site:
				//   https://www.countries-ofthe-world.com/all-countries.html
				ctx.Conn.WriteSHA1String(ctx.Authdata, args[1], ctx.UserConfig.Country)

			case "ADDRNUM":
				// specifies how many lines your address has, this address should
				// be in the specified country
				ctx.Conn.WriteSHA1String(ctx.Authdata, args[1], strconv.Itoa(len(ctx.UserConfig.Address)))

			case "ADDRLINE1":
				ctx.Conn.WriteSHA1String(ctx.Authdata, args[1], ctx.UserConfig.Address[0])

			case "ADDRLINE2":
				ctx.Conn.WriteSHA1String(ctx.Authdata, args[1], ctx.UserConfig.Address[1])

			case "POW":
				log.Println("Searching for HASH:")
				t.Reset(precessingInterval)
				go ctx.pow(line)
			case "ERROR":
				return errors.New(line)
			default:
				log.Println("Unkown command")
				return errors.New("unkown command")
			}
		case suff := <-ctx.outcoming:
			ctx.Conn.WriteString(suff)
		case <-t.C:
			return errors.New("time expired")
		}
	}
}

func (ctx *Miner) readConnData() {
	connReader := textproto.NewReader(bufio.NewReader(ctx.Conn))

	for {
		// read one line (ended with \n or \r\n)
		line, err := connReader.ReadLine()
		if err != nil {
			fmt.Printf("incoming error: %v\n", err)
			return
		}

		if len(line) > 0 {
			fmt.Println(line)
			ctx.incoming <- line
		}
	}
}

func (ctx *Miner) pow(line string) {
	args := strings.Fields(line)

	stop := make(chan bool, 1)
	defer close(stop)
	// go utils.HashRate(ctx.Counter, stop)

	ctx.Authdata = args[1]
	difficulty, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatalf("Difficulty of POW not integer: %s", err.Error())
	}
	fmt.Println("Authdata: ", ctx.Authdata, "Dificulty: ", difficulty)

	// create context fro workerPool
	minerCtx, cancelWorkerPool := context.WithTimeout(context.Background(), time.Hour*2)
	defer cancelWorkerPool()

	// Start workers
	go ctx.WPool.Run(minerCtx)

	jobs := GenerateWorkerJobs(ctx.WPool.GetWorkerCount(), difficulty, minRandomStringLength, maxRandomStringLength, ctx.Authdata, ctx.Counter)
	go ctx.WPool.SendBulkJobs(jobs)

	suff, err := GetResults(ctx.WPool)
	if err == context.DeadlineExceeded {
		fmt.Println("Dedline reached: ", err.Error())
	}

	if err == nil && suff != "" {
		fmt.Println("Suff: ", suff)
		// ctx.Conn.WriteString(suff)
		ctx.outcoming <- suff
	}

	// Stop goroutine hashRate
	stop <- true
}
