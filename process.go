package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type process_info struct {
	ip   string
	port string
}

func process_received_m(message string) string {
	mparsed := strings.SplitN(message, " ", 2)
	return `Recieved "` + strings.TrimSuffix(mparsed[1], "\n") + `" from process ` + mparsed[0] + `, system time is ` + time.Now().Format("15:04:05.000000")
}

func process_send_c(command string) (string, string) {
	command = strings.TrimSuffix(command, "\n")
	cparsed := strings.SplitN(command, " ", 3)
	return cparsed[1], cparsed[2]
}

func unicast_send(destination net.Conn, message string) {
	_, err := destination.Write([]byte(message + "\n"))
	if err != nil {
		log.Fatal()
	}
}

// assigns connections to individual reader goroutines that route messages into the proper channel
func unicast_recieve(source net.Listener) {
	for {

		// accept incoming connections
		conn, err := source.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// pass connection into subproccess to handle incoming messages
		go func(conn net.Conn) {
			for {
				message, err := bufio.NewReader(conn).ReadString('\n')

				if err == io.EOF {
					fmt.Println("connection to " + conn.LocalAddr().String() + " has shut down, please reboot program to re-establish connection")
					conn.Close()
				}

				if err == nil {
					fmt.Println(process_received_m(message))
				}
			}
		}(conn)

	}
}

func initialize_source(info process_info) net.Listener {
	ln, err := net.Listen("tcp", info.ip+":"+info.port)
	if err != nil {
		log.Fatal(err)
	}
	return ln
}

func initialize_outgoing(info process_info) net.Conn {
	conn, err := net.Dial("tcp", info.ip+":"+info.port)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

// parse_config takes the path pointing to a config file and translates it into a map indexed by process id containing process_info structs
func parse_config(path string) (map[string]process_info, []int) {

	// initialize empty proccess map
	processes := make(map[string]process_info)

	// open config file and initialize a scanner
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// read delays into a delay slice
	scanner.Scan()
	delays := func() []int {
		raw_delays := strings.Split(scanner.Text(), " ")
		low, err := strconv.Atoi(raw_delays[0])
		if err != nil {
			log.Fatal()
		}
		high, err := strconv.Atoi(raw_delays[1])
		if err != nil {
			log.Fatal()
		}
		return []int{low, high}
	}()

	// read remaining values into process map
	for scanner.Scan() {
		splitline := strings.Split(scanner.Text(), " ")
		processes[splitline[0]] = process_info{splitline[1], splitline[2]}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	return processes, delays
}

func main() {

	args := os.Args
	self := args[1]

	// get procress_info map and delay bound slice
	var process_infomap map[string]process_info
	var delay_bounds []int
	switch len(args) {
	case 2:
		process_infomap, delay_bounds = parse_config("config")
	case 3:
		process_infomap, delay_bounds = parse_config(args[2])
	default:
		log.Fatal("\nusage: ./process <id> <optional: config path>\n")
	}

	// assign current process to listen to the port defined by id provided as command line argument
	source_server := initialize_source(process_infomap[self])
	defer source_server.Close()

	// initialize an empty map to store active outgoing connections
	outgoing := make(map[string]net.Conn)

	// activate reciever
	go unicast_recieve(source_server)

	// parse stdin for send commands
	for {

		// parse input in stdin into necessary strings
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		going_to, raw_text := process_send_c(input)
		text_with_header := self + " " + raw_text

		// handle simulated delay and actual send in a separate goroutine to prevent blocking
		go func() {

			// establish outgoing connection if not established yet
			if _, ok := outgoing[going_to]; !ok {
				outgoing[going_to] = initialize_outgoing(process_infomap[going_to]) // this may not be safe
			}
			fmt.Println(`sending "` + raw_text + `" to ` + going_to + ". System time is " + time.Now().Format("15:04:05.000000"))
			time.Sleep(time.Duration(rand.Intn(delay_bounds[1]-delay_bounds[0])+delay_bounds[0]) * time.Millisecond)
			unicast_send(outgoing[going_to], text_with_header)

		}()

	}

}
