# GOing-around

A simple TCP messenger. 

To initialize an instance of the messenger, run `./proccess <id> <OPTIONAL: config path>`
By default, the config path routes to "config" in the working directory

The program starts by parsing the config file into a slice of the simulated delay bounds and a proccess_info map indexed by string id

### process_info reference:
    type process_info struct {
        ip   string
        port string
    }

The instance starts to listen on the port mapped to the id provided when initializing execution.
A map is also created to 
