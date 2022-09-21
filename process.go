package main

import (
	"bufio"
	"fmt"
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

func process_message(message string) string {
	mparsed := strings.SplitN(message, " ", 3)
	return `Recieved "` + mparsed[2] + `" from process ` + mparsed[0] + `, system time is` + mparsed[1]  
}

func process_send(message string) (string, string) {
	message = strings.TrimSuffix(message, "\n")
	mparsed := strings.SplitN(message, " ", 3)
	
	fmt.Println(mparsed[2] + mparsed[1] + mparsed[0])
	return mparsed[1], mparsed[2]
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
		conn, err := source.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(conn net.Conn) {
			for {
				message, _ := bufio.NewReader(conn).ReadString('\n')
				fmt.Println(process_message(message))
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
	conn, err := net.Dial("tcp", "golang.org:80")
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
	if len(args) < 3 {
		process_infomap, delay_bounds = parse_config("config")
	} else {
		process_infomap, delay_bounds = parse_config(args[2])
	}

	// assign current process to listen to the port defined by id provided as command line argument
	source_server := initialize_source(process_infomap[self])
	defer source_server.Close()

	// initialize a map of outgoing connections for each non-source process id
	outgoing := make(map[string]net.Conn)
	/*for key, value := range process_infomap {
		if key != self {
			outgoing[key] = initialize_outgoing(value)
		}
	}*/

	// activate reciever
	go unicast_recieve(source_server)

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		going_to, text := process_send(text)

		go func () {
			if _, ok := outgoing[going_to]; !ok {
				outgoing[going_to] = initialize_outgoing(process_infomap[going_to])
			}
			time.Sleep(time.Duration(rand.Intn(delay_bounds[1]-delay_bounds[0]) + delay_bounds[0])*time.Millisecond)
			unicast_send(outgoing[going_to], text)
		}()
		
	}

}
