package miner

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MihaiLupoiu/interview-exasol/connection"
	"github.com/MihaiLupoiu/interview-exasol/utils"
	"github.com/MihaiLupoiu/interview-exasol/worker"
	"github.com/paulbellamy/ratecounter"

	crypto_rand "crypto/rand"
	math_rand "math/rand"
)

// Miner has all the basic information to start the search for the SHA1
type Miner struct {
	Authdata   string
	Conn       *connection.Connection
	Counter    *ratecounter.RateCounter
	UserConfig UserConfig
	WPool      worker.Pool
}

var (
	minRandomStringLength = 5
	maxRandomStringLength = 64
)

func initSeed() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

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
	initSeed()
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
	}, err
}

// Run miner will connect to the server, initialize the workes when request POW is received and
// start to search for the SHA1 with the given dificulty.
func (ctx *Miner) Run() error {
	defer ctx.Conn.Close()

	stop := make(chan bool, 1)
	defer close(stop)
	go utils.HashRate(ctx.Counter, stop)

	connTextReader := textproto.NewReader(bufio.NewReader(ctx.Conn))
	for {
		// read one line (ended with \n or \r\n)
		line, _ := connTextReader.ReadLine()
		if len(line) > 0 {
			fmt.Println(line)
		}

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
				args := strings.Fields(line)

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
					break
				}

				if err == nil && suff != "" {
					fmt.Println("Suff: ", suff)
					ctx.Conn.WriteString(suff)
					break
				}
			} else if strings.HasPrefix(line, "ERROR") {
				fmt.Println(line)
				return nil
			}
		}
		// Stop goroutine hashRate
		stop <- true
	}
}
