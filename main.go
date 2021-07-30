package main

import (
	"bufio"
	"crypto/tls"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/textproto"
	"strconv"
	"strings"

	"crypto/sha1"

	"github.com/google/uuid"
)

func createClientConfig(crt, key string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		// TODO: Investigate Error: 2021/07/27 18:24:27 failed to connect: x509: cannot validate certificate for 18.202.148.130 because it doesn't contain any IP SANs
		// Posible problem with using the IP.
		InsecureSkipVerify: true,
	}, nil
}

func printConnState(conn *tls.Conn) {
	log.Print(">>>>>>>>>>>>>>>> State <<<<<<<<<<<<<<<<")
	state := conn.ConnectionState()
	log.Printf("Version: %x", state.Version)
	log.Printf("HandshakeComplete: %t", state.HandshakeComplete)
	log.Printf("DidResume: %t", state.DidResume)
	log.Printf("CipherSuite: %x", state.CipherSuite)
	log.Printf("NegotiatedProtocol: %s", state.NegotiatedProtocol)

	log.Print("Certificate chain:")
	for i, cert := range state.PeerCertificates {
		subject := cert.Subject
		issuer := cert.Issuer
		log.Printf(" %d s:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", i, subject.Country, subject.Province, subject.Locality, subject.Organization, subject.OrganizationalUnit, subject.CommonName)
		log.Printf("   i:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", issuer.Country, issuer.Province, issuer.Locality, issuer.Organization, issuer.OrganizationalUnit, issuer.CommonName)
	}
	log.Print(">>>>>>>>>>>>>>>> State End <<<<<<<<<<<<<<<<")
}

var authdata = ""

func hexStartsWith(hash [20]byte, dificulty int) bool {
	// Improve method to use bit manipulation for more optimal comparion.
	sha1_hash := hex.EncodeToString(hash[:])
	fmt.Println(hash, sha1_hash)

	prefixDifficulty := strings.Repeat("0", dificulty)
	fmt.Println(prefixDifficulty)

	res := strings.HasPrefix(sha1_hash, prefixDifficulty)
	fmt.Println("Result:", res)

	return res
}

func main() {
	connect := flag.String("connect", "localhost:4433", "who to connect to")
	crt := flag.String("crt", "./configs/certs/public.crt", "certificate")
	key := flag.String("key", "./configs/certs/private.key", "key")
	flag.Parse()

	addr := *connect
	if !strings.Contains(addr, ":") {
		addr += ":443"
	}

	config, err := createClientConfig(*crt, *key)
	if err != nil {
		log.Fatalf("config failed: %s", err.Error())
	}

	conn, err := tls.Dial("tcp", addr, config)
	if err != nil {
		log.Fatalf("failed to connect: %s", err.Error())
	}
	defer conn.Close()

	log.Printf("connect to %s succeed", addr)
	printConnState(conn)

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	// defer conn.Close()

	for {
		// read one line (ended with \n or \r\n)
		line, _ := tp.ReadLine()
		fmt.Println(line)

		switch line {
		case "HELO":
			conn.Write([]byte("EHLO\n"))

		case "END":
			// if you get this command, then your data was submitted
			conn.Write([]byte("OK\n"))

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

				for {

					// generate short random string, server accepts all utf-8 characters,
					// except [\n\r\t ], it means that the suffix should not contain the
					// characters: newline, carriege return, tab and space
					suffix := uuid.New().String()
					fmt.Printf("String: %s\n", suffix)

					cksum_in_hex := sha1.Sum([]byte(authdata + suffix))
					fmt.Printf("  SHA1: %x\n", cksum_in_hex)

					// check if the checksum has enough leading zeros
					// (length of leading zeros should be equal to the difficulty)
					if hexStartsWith(cksum_in_hex, difficulty) {
						conn.Write([]byte(suffix + "\n"))
						break
					}

				}
			} else if strings.HasPrefix(line, "ERROR") {
				fmt.Println(line)
				break
			}
		}

		// do something with data here, concat, handle and etc...
	}

}
