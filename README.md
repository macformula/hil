# MAC Formula HIL
Hardware-in-the-loop tester to accelerate testing of the electrical-software test-bench 

## Starting the SIL

Navigate to `hil/cmd/hilapp` then in the terminal run 
```bash
go build
``` 

This will build the binary that instantiates all of the structs needed for a sil test, and start the cli application. 
Next run the binary, with the `--config` flag as the path to the config file in `hil/macformula/config` directory.

```bash
./hilapp --config=path/to/config.yaml
```

The cli application should appear in the terminal.


## Demo/BasicIO testing

There is a go script in the `hil/cmd/basicio` directory that mimics the logic of the `Demo/BasicIO` project in firmware. The go script first initializes a TCP socket connection to the SIL. The script then enters a while loop, which reads the value of the button input, and sets the value of the indicator LED to match the button value. It then waits 50 milliseconds to ensure SIL has enough time to handle the request. 

On the SIL side, there is a test state in `hil/macformula/state/basicio` that is used to validate the firmware's logic. There are two main states to test. When the input button is set to high, the indicator led should also be set to high. When the input button is set to low, the indicator led should also be set to low. The test state checks these two conditions then updates the value of the result tags accordingly.

### Running the BasicIO script

1. First follow the steps above to run the SIL. Keep the SIL cli running persistently. When the SIL cli is running. It automatically starts the TCP server to listen for connections. 
2. Open a separate terminal and navigate to the `hil/cmd/basicio` directory. Then run:

```bash
go build
./basicio
```

__NOTE__: It's important you start the SIL first before running the basicio binary. Otherwise the basicio binary won't be able to establish a connection to the TCP socket.

3. Now that the SIL server is running, and the basicio script is running at the same time in different terminals. Navigate to the terminal with the SIL cli, and run the BasicIO test state. A log, and an HTML output file will be produced in the `hil/macformula/results` directory. All test cases should be passed!