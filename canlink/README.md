canlink
======================

`canlink` contains utilities for managing CAN traffic for the duration of a HIL test.

BusManager
---------------
`BusManager` is a centralized node responsible for orchestrating all interactions with a CAN bus.
It acts as a message broker supporting the transmission of bus traffic to registered handlers and receiving incoming messages from these handlers to write to the bus.

Tracer
---------------
`Tracer` writes traffic on a CAN bus into trace files.
Tracer takes in a struct that must implement the `Converter` interface. 
A `Converter` must have a method `GetFileExtension`to return the proper file extension and `FrameToString`to convert a frame to a string. Currently implemented converters support `Jsonl` (https://jsonlines.org/) and `Text` for basic text logging. 

__NOTE: If seeking to trace traffic into an unsupported format, implement a converter for that specific format.__


### Usage

1) Create a Bus Manager using `NewBusManager()` function.
A logger, and pointer to a socketcan connection are passed as arguments.

    ```go
    func main() {
        ctx := context.Background()

        loggerConfig := zap.NewDevelopmentConfig()
        logger, err := loggerConfig.Build()

        // Create a network connection for vcan0
        conn, err := socketcan.DialContext(context.Background(), "can", "vcan0")
        if err != nil {
            return
        }

        manager := canlink.NewBusManager(logger, &conn)
    }
    ```

2) Create an instance of Tracer with the `NewTracer()` function. 
A can interface, logger, and converter must be provided as arguments. 
A timeout and a file name are optional parameters. If no file name is provided, a name will be generated using the current time and date. The tracer will timeout after 30 minutes by default. The trace file will be created in the same directory.
Functional options are available of type `TracerOption` if required. 

    ```go
    func main() {
        tracer := canlink.NewTracer(
            "vcan0",
            logger,
            canlink.Text{}
            canlink.WithTimeout(1*time.Second)
            canlink.WithFileName("trace_sample")
	    )
    }
    ```
3) Register the `Tracer` instance as a handler for the bus manager by calling `Register`, passing in the `Tracer`. Once started, the manager will send frames from the CAN bus through the broadcast channel to every registered handler. The `transmit` channel can be used by handlers to transmit frames onto the CAN bus.

    ```go
    func main() {
        broadcast, transmit := manager.Register(tracer)
        go tracer.Handle(broadcast, transmit)
	    defer tracer.Close()
    }
    ```
4) Then call `tracer.Handle` on a separate go routine with the broadcast and transmit channel as arguments. Once `tracer.Handle` is called any frames transmitted on the broadcast channel will be traced. 

    ```go
    func main() {
        go tracer.Handle(broadcast, transmit)
	    defer tracer.Close()
    }
    ```

5) Call `manager.Start` to start the traffic broadcast and incoming frame listener for each of the registered handlers.

    ```go
    func main() {
        manager.Start(ctx)
        defer manager.Close()
    }
    ```

    __NOTE: The Tracer object is reusable. If a new CAN Trace is desired, repeat steps 2-3, creating and registering a new trace object. A new file will be written to with this traffic.__
