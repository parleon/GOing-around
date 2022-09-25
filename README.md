# GOing-around

A simple TCP messenger. 

### Definitions: 
 - instance: an instance refers to a program and all of its running subprocesses that is instantiated after executing the initializing command `./process ...`
 - Subprocess: refers to a goroutine that is spawned and runs throughout the entire runtime of the program. This includes the unicast_reciever, and all outgoing connections, as well as the main process.
 - goroutine: when a goroutine is referred to as a goroutine, it means the purpose of the thread is short lived. (e.g., a send command is processed in its own goroutine to ensure its non-blocking)

To initialize an instance of the messenger, run `./proccess <id> <OPTIONAL: config path>`
By default, the config path routes to "config" in the working directory

The program starts by parsing the config file into a slice of the simulated delay bounds and a proccess_info map indexed by string id

### process_info reference:
    type process_info struct {
        ip   string
        port string
    }

The main subprocess initializes a server listener for the source port associated with the id provided, as well as an empty map to track outgoing connections.
The main subprocess then creates another subprocess "unicase-reciever" to handle incoming messages. 

### unicast_reciever
The unicast reciever works by accepting incoming connections and passing the newly established connection into a new subproccess that handles printing anything sent and closing the connection when the other party disconnects


A message is sent by entering the input `send <id> <message>` into the standard input of a running instance



