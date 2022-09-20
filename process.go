package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type process_info struct {
	ip   string
	port string
}

func unicast_send(destination net.Conn, message string) {
	_, err := destination.Write([]byte(message + "\n"))
	if err != nil {
		log.Fatal()
	}
}

// assigns connections to individual reader goroutines that route messages into the propper channel
func unicast_recieve(source net.Listener, message_channels map[string]chan []byte) {
	for {
		conn, err := source.Accept()
		if err != nil {
			log.Fatal(err)
		}
		
		go func(conn net.Conn, message_channels map[string]chan []byte) {
			conn.Read(<-message_channels["PLACEHOLDER"])
		}(conn, message_channels)

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
	conn, err := net.Dial("tcp", "golang.org:80")
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

// parse_config takes the path pointing to a config file and translates it into a map indexed by process id containing process_info structs
func parse_config(path string) (map[string]process_info, []int) {
	// initialize empty proccess map
	var processes map[string]process_info

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
	if len(args) > 2 {
		process_infomap, delay_bounds = parse_config("config")
	} else {
		process_infomap, delay_bounds = parse_config(args[2])
	}

	// assign current process to listen to the port defined by id provided as command line argument
	source_server := initialize_source(process_infomap[self])
	defer source_server.Close()

	// initialize a map of outgoing connections for each non-source process id
	var outgoing map[string]net.Conn
	for key, value := range process_infomap {
		if key != self {
			outgoing[key] = initialize_outgoing(value)
		}
	}
	for {

		select {}
	}

}
