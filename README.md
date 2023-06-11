# HashGrab

![release](https://github.com/tauki/hashgrab/actions/workflows/test.yml/badge.svg)

HashGrab is a concurrent URL fetcher and MD5 hasher written in Go. Given a list of URLs, it fetches the contents of these URLs and computes their MD5 hash. The tool allows you to control the maximum number of concurrent requests, making it highly efficient in fetching and hashing multiple URLs.

## Requirements

- Go 1.16 or later

## Installation

### Build from Source

1. Clone this repository:

   ```
   git clone https://github.com/tauki/hashgrab.git
   ```

2. Change to the repository's directory:

   ```
   cd hashgrab
   ```

3. Change to the `cmd` directory where the `main.go` file is located:

   ```
   cd cmd
   ```

4. Build the project:

   ```
   go build -o hashgrab .
   ```

5. The `hashgrab` binary is now ready. You can move it to a directory in your `PATH` for easy access.

### Use `go install`

If you have Go installed, you can also directly install the package using `go install`.

```
go install github.com/tauki/hashgrab/cmd@latest
```

This command will install the `hashgrab` binary in your `GOBIN` or `GOPATH/bin` directory.

## Usage

HashGrab fetches URLs and calculates their MD5 hashes. The tool takes a list of URLs as command-line arguments.

```bash
hashgrab https://example.com https://another-example.com
```

This will fetch `https://example.com` and `https://another-example.com`, and print their MD5 hashes to the standard output.

### Flags

HashGrab supports the following flags:

- `-parallel=<number>`: Specifies the maximum number of parallel requests to fetch URLs. Defaults to 10 if not provided.

Example:

```bash
hashgrab -parallel=10 https://example.com https://another-example.com
```

This will fetch `https://example.com` and `https://another-example.com` in parallel (with a maximum of 10 parallel requests).
