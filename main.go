package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"search/client"
	"search/server"
	"strconv"
	"strings"
)

// Sets up the flags
var numClients = flag.Int("num_clients", 5, "the number of clients for the server")
var lenMS = flag.Int("len_ms", 8, "the length of the master secret")
var lenSalt = flag.Int("len_salt", 8, "the length of the salts used to generate the PRFs")
var fpRate = flag.Float64("fp_rate", 0.000001, "the desired false positive rate for searchable encryption")
var numUniqWords = flag.Uint64("num_words", uint64(10000), "the expected number of unique words in all the documents")
var serverMountPoint = flag.String("server_mp", "server_fs", "the mount point for the server where all the server side data is stored")

var defaultClientNum = flag.Int("default_client_num", -1, "the dafault running client (set to -1 to initialize without a client)")
var clientMountPoint = flag.String("client_mp", "client_fs", "the mount point for the client where the client stores all the data")

func startServer() *server.Server {
	if _, err := os.Stat(path.Join(*serverMountPoint, "serverMD")); err == nil {
		fmt.Println("Server metadata found, loading server from mount point", *serverMountPoint)
		return server.LoadServer(*serverMountPoint)
	}
	if _, err := os.Stat(*serverMountPoint); os.IsNotExist(err) {
		if os.Mkdir(*serverMountPoint, 0777) != nil {
			fmt.Println("Failed to create the server mount point", *serverMountPoint)
		}
	}
	fmt.Println("No previous server metadata found, starting new server at mount point", *serverMountPoint)
	return server.CreateServer(*numClients, *lenMS, *lenSalt, *serverMountPoint, *fpRate, *numUniqWords)
}

func startClient(server *server.Server, clientNum int) *client.Client {
	if clientNum == -1 {
		return nil
	}
	fmt.Println("Now running client", clientNum)
	return client.CreateClient(server, clientNum, path.Join(*clientMountPoint, "client"+strconv.Itoa(clientNum)))
}

func main() {
	flag.Parse()

	// Initialize the server
	server := startServer()
	fmt.Printf("\nServer Started\n--------------\n")
	server.PrintServerInfo()
	fmt.Println()

	// Initialize the client
	if _, err := os.Stat(*clientMountPoint); os.IsNotExist(err) {
		if os.Mkdir(*clientMountPoint, 0777) != nil {
			fmt.Println("Failed to create the client mount point", *clientMountPoint)
		}
	}
	client := startClient(server, *defaultClientNum)
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		cmd, _ := reader.ReadString('\n')
		cmd = cmd[:len(cmd)-1]
		tokens := strings.Split(cmd, " ")
		switch tokens[0] {
		case "client":
			if len(tokens) < 2 {
				fmt.Printf("%s: client number missing\n", tokens[0])
				break
			}
			clientNum, err := strconv.Atoi(tokens[1])
			if err != nil || clientNum < 0 || clientNum >= server.GetNumClients() {
				fmt.Printf("%s: invalid client number \"%s\"\n", tokens[0], tokens[1])
				break
			}
			client = startClient(server, clientNum)
		case "ls":
			if client == nil {
				fmt.Printf("%s: client not running\n", tokens[0])
				break
			}

		case "info":
			server.PrintServerInfo()
		case "exit", "q":
			fmt.Println("Program exited.")
			return
		case "":
		default:
			fmt.Printf("%s: command not found\n", tokens[0])
		}
	}
}
