# Hubble Scan Server

*Hubble Scan Server that checks blocks of [Burrow Blockchain](https://github.com/hyperledger/burrow) and saves them in Postgre database.

## Compiling the code

You need to make sure you have install [Go](https://golang.org/) (version 1.10.1 or higher) and [postgre](https://www.postgresql.org). After installing them, import HubbleScan.sql from script folder into postgre to create database and then you can follow these steps to compile and build the project:

```bash
mkdir -p $GOPATH/src/github.com/hubbleServer
cd $GOPATH/src/github.com/hubbleServer
git clone https://github.com/BurrowBlocks.git .
make
```
