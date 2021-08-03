[![Build Status](https://github.com/MihaiLupoiu/interview-exasol/workflows/Test/badge.svg)](https://github.com/MihaiLupoiu/interview-exasol/actions)
[![codecov](https://codecov.io/gh/MihaiLupoiu/interview-exasol/branch/main/graph/badge.svg?token=N15FQSWTAW)](https://codecov.io/gh/MihaiLupoiu/interview-exasol)

As part of our application process, we would like you to write a test
program that we could later discuss in your interview. The data you
send using the program will be used for further communication with
you. You can write your program in any programming language you
prefer, but you should be able to show and explain your solution later
in the interview. To connect to the server, you need the keys included
in this README.

The following pseudocode represents the program you need to write. It
is a full implementation with all required elements. It is written
with Python language syntax and semantics in mind, but it is not a
correct implementation and needs to be extended to actually run in a
Python interpreter. The main purpose of this pseudocode is to give you
an idea of what you need to develop.

### Pseudocode:

```
conn = tls_connect("IP:PORT", cert, key)
authdata = ""
while true:
    args = conn.read().strip().split(' ')
    if args[0] == "HELO":
        conn.write("EHLO\n")
    elif args[0] == "ERROR":
        print("ERROR: " + " ".join(args[1:]))
        break
    elif args[0] == "POW":
        authdata, difficulty = args[1], args[2]
        while true:
            # generate short random string, server accepts all utf-8 characters,
            # except [\n\r\t ], it means that the suffix should not contain the
            # characters: newline, carriege return, tab and space
            suffix = random_string()
            cksum_in_hex = SHA1(authdata + suffix)
            # check if the checksum has enough leading zeros
            # (length of leading zeros should be equal to the difficulty)
            if cksum_in_hex.startswith("0"*difficulty):
                conn.write(suffix + "\n")
                break
    elif args[0] == "END":
        # if you get this command, then your data was submitted
        conn.write("OK\n")
        break
    # the rest of the data server requests are required to identify you
    # and get basic contact information
    elif args[0] == "NAME":
       # as the response to the NAME request you should send your full name
       # including first and last name separated by single space
       conn.write(SHA1(authdata + args[1]) + " " + "My name\n")
    elif args[0] == "MAILNUM":
       # here you specify, how many email addresses you want to send
       # each email is asked separately up to the number specified in MAILNUM
       conn.write(SHA1(authdata + args[1]) + " " + "2\n")
    elif args[0] == "MAIL1":
       conn.write(SHA1(authdata + args[1]) + " " + "my.name@example.com\n")
    elif args[0] == "MAIL2":
       conn.write(SHA1(authdata + args[1]) + " " + "my.name2@example.com\n")
    elif args[0] == "SKYPE":
       # here please specify your Skype account for the interview, or N/A
       # in case you have no Skype account
       conn.write(SHA1(authdata + args[1]) + " " + "my.name@example.com\n")
    elif args[0] == "BIRTHDATE":
       # here please specify your birthdate in the format %d.%m.%Y
       conn.write(SHA1(authdata + args[1]) + " " + "01.02.2017\n")
    elif args[0] == "COUNTRY":
       # country where you currently live and where the specified address is
       # please use only the names from this web site:
       #   https://www.countries-ofthe-world.com/all-countries.html
       conn.write(SHA1(authdata + args[1]) + " " + "Germany\n")
    elif args[0] == "ADDRNUM":
       # specifies how many lines your address has, this address should
       # be in the specified country
       conn.write(SHA1(authdata + args[1]) + " " + "2\n")
    elif args[0] == "ADDRLINE1":
       conn.write(SHA1(authdata + args[1]) + " " + "Long street 3\n")
    elif args[0] == "ADDRLINE2":
       conn.write(SHA1(authdata + args[1]) + " " + "32345 Big city\n")
conn.close()
```

Notes:
- This pseudocode is written with Python semantics in mind
- Only TLS connections with valid keys are allowed
- All communication needs to be in valid UTF-8
- Protocol is line oriented and each line should end with \n
- Only if you see the END request from the server is the data fully sent
- On problems, the server sends the ERROR command and closes the connection
- If the data is not fully sent, no data will be recorded on server
- There are no logs about connections on the server side
- HELO and POW commands always come first (handshake)
- END command is always the last command and confirms successful application
- Other commands come from the server in random order
- List of acceptable country names:
    https://www.countries-ofthe-world.com/all-countries.html
- The timeout of the POW command is 2 hours
- All other commands have a timeout of 6 seconds
- It is possible to reach this service on the following ports:
  3336, 8083, 8446, 49155, 3481, 65532


### RUN code

```
go run main.go -connect 18.202.148.130:3336
```

## TODOS
- [x] Create package.
- [x] Connect to server.
- [x] Generate randon string.
   - ~~Using random uuid from google.~~
   - Generating random string using math/rand for speed performance. Also available the SecureRandomString to generate a cryptography secure random. 
- [x] Calculate sha1.
- [x] Calculate compare hash with dificulty.
- [ ] Improve code structure to split responsibilities.
- [ ] Test functions. 
   - Especially the solver for each dificulty.
- [ ] Improve loggin.
- [ ] Improve flags to specify log level.
- [x] Add contact information in a configuration file. 
- [ ] Benchmark functions like string generator and solver and profile to check what could be improved in the secuencial model.
- [ ] Improve speed by implementing concurrency using a worker pool to calculate hash in multiple corutines.
- [ ] Check performance increase and ajust the number of working coroutines in the worker pool.
- [ ] Execute against the server in a multicore CPU with more than 2 cores than my 2013 i5 Macbook pro. 
- [ ] Improve state machine processing. Low priority for now.
- [x] Show number of Hash/second.
- [ ] Documentation in code and architecture diagram.
- [ ] Implement a SHA1 calculation using GPU.
   - If time available, if not create a small explination on how it could be done.
