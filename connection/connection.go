package connection

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"log"
)

// TODO: Test connection.

// Connection has the basic configuration required to create the connection the the server.
type Connection struct {
	conn     *tls.Conn
	conf     *tls.Config
	endpoint string
}

// Dial will connect using with the server endpoint using the cert and key provided.
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

// Close the connection
func (c *Connection) Close() error {
	return c.conn.Close()
}

// Read bytes from the connection
func (c *Connection) Read(b []byte) (int, error) {
	return c.conn.Read(b)
}

// Write bytes to the connection
func (c *Connection) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

// WriteString writes a string to the connection and append a new line at the end.
func (c *Connection) WriteString(b string) (int, error) {
	log.Println(b)
	return c.conn.Write([]byte(b + "\n"))
}

// WriteSHA1String writes a two string seperated by a space where
// authdata, shaArg are used to generate the hash and stringArg is append after a space
func (c *Connection) WriteSHA1String(authdata, shaArg, stringArg string) (int, error) {
	hash := sha1.Sum([]byte(authdata + shaArg))
	return c.WriteString(hex.EncodeToString(hash[:]) + " " + stringArg)
}

// Reconnecte will reconnect to the server.
func (c *Connection) Reconnecte() error {
	conn, err := tls.Dial("tcp", c.endpoint, c.conf)
	if err != nil {
		fmt.Printf("failed to reconnect: %s", err.Error())
		return err
	}
	c.conn = conn
	return nil
}

// PrintConnState will print the TLS connection state.
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
