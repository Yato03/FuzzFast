# FastFuzz

Tool for fuzzing directories of a web application in a **fast** way thanks to the advantages offered by a language like **Golang**. It is possible to assign a number of **Workers** to try to discover directories and files in public paths in parallel.

## Installation

```bash
git clone https://github.com/Yato03/FuzzFast
cd FuzzFast
go build app/fuzzer.go
```

## Use

### Help

```bash
fuzzer --help
Usage of fuzzer:
  -output string
        Output file
  -t int
        Workers to use (default 100)
  -url string
        URL to fuzz
  -wordlist string
        Wordlist to use
```

As the help command indicates, there are the following parameters:

- **url** [required]: Specifies the url to fuzz. It should be of the form: `http(s):///domain/`, either `http` or `https`.
- **wordlist** [required]: File where the directories and files to be fuzzed on the web page are determined.
- **output** [optional]: Output file with the results. It details the paths found by the fuzzer. If no path is specified, no output file will be generated.
- **t** [optional]: Determines the number of *Workers* to be used to parallelize the job. The recommended number will depend on the performance of the device where it is executed. By default it will be 100.

### Example

```bash
fuzzer --url http://localhost --wordlist ./wordlist.txt --t 50 --output ./output.txt
```

Spanish README [here](README-es.md)
