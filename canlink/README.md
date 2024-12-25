canlink
======================

`canlink` contains utilities for managing CAN traffic for the duration of a HIL test.

BusManager and Tracer
---------------
`BusManager` is a centralized node responsible for orchestrating all interactions with a CAN bus.
It acts as a message broker supporting the transmission of bus traffic to registered handlers and receiving incoming messages from these handlers to write to the bus.

`Tracer` captures writes traffic on a CAN bus into trace files.
Currently supported formats are Json and Ascii.


### Usage

1) Create a Bus Manager using `NewBusManager()` function.
A logger, and pointer to a socketcan connection are passed in.

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
A can interface, trace directory and logger must be provided as arguments. 
The trace directory must be a folder within the current directory with the same name as the argument provided.

Functional options are available of type `TracerOption` if required. 

    ```go
    func main() {
        tracer := canlink.NewTracer(
            "vcan0",
            "traces",
            logger,
            canlink.WithBusName("veh"),
            canlink.WithConverter(canlink.NewAscii()),
	    )
    }
    ```
3) Register `Tracer` instance as a handler for the bus manager by calling `Register`, passing in `tracer`. The `broadcast` channel receives frames from the socket connection. The `transmit` channel can be used by handlers to transmit frames onto the CAN bus.

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
