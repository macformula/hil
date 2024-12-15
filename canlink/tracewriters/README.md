tracewriters
======================

`tracewriters` contains structs to create, write to, and close a trace file of a specific format.
Currently supported writers are `JsonWriter` and `AsciiWriter`.

### Usage

__NOTE: Usage is the same for both JsonWriter and AsciiWriter but shown for only JsonWriter.__

1) Create a JsonWriter instance with _NewJsonWriter()_ function.
A logger must be provided as an argument.

    ```go
    func main() {
        jsonWriter := NewJsonWriter(logger)
    }
    ```
2) Initialize the trace file with _CreateTraceFile()_.
A trace directory and bus name are passed as arguments.

    ```go
    func main() {
        jsonWriter.CreateTraceFile(traceDir, busName)
    }
    ```
3) Write to the trace file with _WriteFrameToFile()_, passing in a frame to write.


    ```go
    func main() {
        jsonWriter.WriteFrameToFile(timestampedFrame)
    }
    ```
4) Once the trace has ended, close the file with _CloseTraceFile()_.

    ```go
    func main() {
        jsonWriter.CloseTraceFile()
    }
    ```