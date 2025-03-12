canlink
======================

`canlink` contains utilities for managing CAN traffic for the duration of a HIL test. 

BusManager
---------------
`BusManager` is a centralized node responsible for orchestrating all interactions with a CAN bus.
It acts as a message broker supporting the transmission of bus traffic to registered handlers and writing frames onto the bus.

Tracer
---------------
`Tracer` writes traffic on a CAN bus into trace files.
Tracer takes in a struct that must implement the `Converter` interface.
A `Converter` must have a method `GetFileExtension` to return the proper file extension and `FrameToString` to convert a frame to a string. Currently implemented converters support `Jsonl` (https://jsonlines.org/) and `Text` for basic text logging.

__NOTE: If seeking to trace traffic into an unsupported format, implement a `Converter` for that specific format.__

Handler
---------------
`Handler` is an interface implemented by structs if they wish to receive data from the bus manager. The `Name` method simply returns a string used for logging purposes and error messages. The `Handle` method is used to pass in a broadcast channel for communicating frames between the handler and the bus manager. The first parameter is the frame channel which receives frames broadcast by the bus manager.

The bus manager is designed to be the __single point of contact__ to a CAN interface. All interaction to a bus should be done through the bus manager.


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
A CAN interface, logger, and converter must be provided as arguments.
A timeout and a file name are optional parameters. If no file name is provided, a name will be generated using the current time and date. The tracer will timeout after 30 minutes by default. The trace file will be created in the same directory.
Functional options are available of type `TracerOption` if required.

    ```go
    func main() {
        tracer := canlink.NewTracer(
            "vcan0",
            logger,
            &canlink.Text{},
            canlink.WithTimeout(1*time.Second),
    
	    )
    }
    ```
3) Register the `Tracer` instance as a handler for the bus manager by calling `Register`, passing in the `Tracer`. The `Register` function will create a broadcast channel for the handler. Once started, the manager will send frames from the CAN bus through the broadcast channel to every registered handler.

    ```go
    func main() {
        manager.Register(tracer)
    }
    ```

5) Call `manager.Start` to start the traffic broadcast and incoming frame listener for each of the registered handlers. This will also run the `Handle` method on all registered `Handler`'s. The `manager` will now listen for CAN frames and automatically transmit them to registered `Handler`'s.

    ```go
    func main() {
        manager.Start(ctx)
        defer manager.Close()
    }
    ```
5) The `Send` method on the bus manager can be used to transmit frames onto the CAN bus.

    ```go
    func main() {
        manager.Send(frame)
    }
    ```
6) Use the `Stop` method to close all registered handlers. If seeking to stop and start the manager, handlers will have to be instantiated and registered with the bus manager before calling `manager.Start` again.

    ```go
    func main() {
        manager.Stop()
    }
    ```

__NOTE: The Tracer object is reusable. If a new CAN Trace is desired, repeat steps 2-3, creating and registering a new trace object. A new file will be written to with this traffic.__
