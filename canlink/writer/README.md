writer
======================

`Writer` is a struct which manages the creation, writing to, and closing of a trace file.
Currently supported formats are ascii and jsonl. Each writer instance manages its own trace file. Multiple instances can be created if multiple trace files are desired. 

### Usage

1) Create a Writer instance with `NewWriter()`. A logger and a file extension are passed in as arguments.

__NOTE: The only currently supported file extensions are `.asc` and `.jsonl`__

    ```go
    func main() {
        writer := NewWriter(
            logger,
            ".jsonl",
        )
    }
    ```

2) Initialize the trace file with `CreateTraceFile()`.
A trace directory and bus name are passed as arguments.

    ```go
    func main() {
        writer.CreateTraceFile(traceDir, busName)
    }
    ```
3) Write to the trace file with `WriteFrameToFile()`, passing in a frame to write.


    ```go
    func main() {
        writer.WriteFrameToFile(timestampedFrame)
    }
    ```
4) Once the trace has ended, close the file with `CloseTraceFile()`.

    ```go
    func main() {
        writer.CloseTraceFile()
    }
    ```

