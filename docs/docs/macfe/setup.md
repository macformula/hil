# Developer Setup

Follow these steps to begin developing in `racecar/hil`.

## Dependencies

### Go
In order to compile and run Go code, you must download Go.

Go to [Go documentation](https://go.dev/doc/install) page and follow the installation guide.

To verify Go is installed correctly, run:

```bash
go version
```

### gRPC
You will need a Unix development environment (Unix machine, WSL, or remote into the Raspberry Pi).

Go through the [gRPC Go Quickstart Guide](https://grpc.io/docs/languages/go/quickstart/). Build the example project.

As long as you can successfully build the gRPC example and have Go installed you are ready to get started with the SIL!

## Precommit setup

We have a precommit that will ensure `go fmt` is ran before every commit to maintain style consistency. To 

=== "Windows"
    1. If you don't already have Python installed, get it from [https://www.python.org/downloads](https://www.python.org/downloads).
    2. Open a terminal and run:
    ```bash
    pip install pre-commit
    ```
    3. Open a terminal in the `hil` directory and run:
    ```bash
    pre-commit install
    ```


=== "Mac/Linux"
    1. If you don't already have Python installed, get it from [https://www.python.org/downloads](https://www.python.org/downloads), or by using a package manager like [Brew](https://docs.brew.sh).
    2. Open a terminal and run:
    ```bash
    pip3 install pre-commit
    ```
    3. Open a terminal in the `hil` directory and run:
    ```bash
    pre-commit install
    ```

Now if you try committing to your local repository with formatting errors in a `.go` file, an error will be thrown and your code will automatically be reformated!
