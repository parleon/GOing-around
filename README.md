# GOing-around

A simple TCP messenger to communicate between processes. 

To initialize an instance of the messenger, run "./proccess <process number> <OPTIONAL: config path>"

The program starts by parsing the config file to define what constants should be used for simulated delays and create a map of `process_info` structs indexed by string `id`.

``
    type process_info struct {
        ip   string
        port string
    }
``


