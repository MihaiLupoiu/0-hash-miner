package connection

import (
	"crypto/tls"
	"fmt"
	"log"
)

// TODO: Test connection.

type Connection struct {
	conn     *tls.Conn
	conf     *tls.Config
	endpoint string
}

func Dial(certFile, keyFile, endpoint string) (*Connection, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		fmt.Printf("load cert failed: %s", err.Error())
		return nil, err
	}

	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		// TODO: Investigate Error: 2021/07/27 18:24:27 failed to connect: x509: cannot validate certificate for 18.202.148.130 because it doesn't contain any IP SANs
		// Posible problem with using the IP.
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", endpoint, tlsConf)
	if err != nil {
		fmt.Printf("failed to connect: %s", err.Error())
		return nil, err
	}
	return &Connection{
		conn:     conn,
		conf:     tlsConf,
		endpoint: endpoint,
	}, nil
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

func (c *Connection) Read(b []byte) (int, error) {
	return c.conn.Read(b)
}

func (c *Connection) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *Connection) WriteString(b string) (int, error) {
	return c.conn.Write([]byte(b + "\n"))
}

func (c *Connection) Reconnecte() error {
	conn, err := tls.Dial("tcp", c.endpoint, c.conf)
	if err != nil {
		fmt.Printf("failed to reconnect: %s", err.Error())
		return err
	}
	c.conn = conn
	return nil
}

func (c *Connection) PrintConnState() {
	log.Print(">>>>>>>>>>>>>>>> State <<<<<<<<<<<<<<<<")
	state := c.conn.ConnectionState()
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
