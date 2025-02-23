# twlog

twlog is a small utility that analyzes Teeworlds and DDNet log files.
This is useful for determining a list of ip addresses of spam bots based on a specific phrase they use.

## installation

Download the executable from the releases page or install it using the Go toolchain:

```bash
go install github.com/jxsl13/twlog@latest
```

## building and installing from source

```bash
# building
go build .

# installing
go install .
```

## usage

```bash
$ ./twlog --help
Environment variables:
  OUTPUT    output format, one of 'json' or 'text' (default: "text")
Environment variables:
  SEARCH_DIR         directory to search for files recursively (default: ".")
  FILE_REGEX         regex to match files in the search dir (default: ".*\\.log$")
  ARCHIVE_REGEX      regex to match archive files in the search dir (default: "\\.(7z|bz2|gz|tar|xz|zip|xz|zst|lz)$")
  INCLUDE_ARCHIVE    search inside archive files (default: "false")
  CONCURRENCY        number of concurrent workers to use (default: "12")

Usage:
  twlog [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  who         who is the subcomand which allows to search for players, nicknames and their IPs

Flags:
  -a, --archive-regex string   regex to match archive files in the search dir (default "\\.(7z|bz2|gz|tar|xz|zip|xz|zst|lz)$")
  -t, --concurrency int        number of concurrent workers to use (default 12)
  -c, --config string          .env config file path (or via env variable CONFIG)
  -f, --file-regex string      regex to match files in the search dir (default ".*\\.log$")
  -h, --help                   help for twlog
  -A, --include-archive        search inside archive files
  -o, --output string          output format, one of 'json' or 'text' (default "text")
  -d, --search-dir string      directory to search for files recursively (default ".")

Use "twlog [command] --help" for more information about a command.
```

```shell
$ ./twlog who --help
who is the subcomand which allows to search for players, nicknames and their IPs

Usage:
  twlog who [flags]
  twlog who [command]

Available Commands:
  said        said searches for what players said in the chat

Flags:
  -h, --help   help for who

Global Flags:
  -a, --archive-regex string   regex to match archive files in the search dir (default "\\.(7z|bz2|gz|tar|xz|zip|xz|zst|lz)$")
  -t, --concurrency int        number of concurrent workers to use (default 12)
  -c, --config string          .env config file path (or via env variable CONFIG)
  -f, --file-regex string      regex to match files in the search dir (default ".*\\.log$")
  -A, --include-archive        search inside archive files
  -o, --output string          output format, one of 'json' or 'text' (default "text")
  -d, --search-dir string      directory to search for files recursively (default ".")

Use "twlog who [command] --help" for more information about a command.
```

```shell
$ ./twlog who said --help
Environment variables:
  DEDUPLICATE    deduplicate objects based on all fields (default: "false")
  EXTENDED       add two additional fields, file and id to the output (default: "false")
  IPS_ONLY       only print IP addresses (default: "false")

Usage:
  twlog who said [flags]

Flags:
  -D, --deduplicate   deduplicate objects based on all fields
  -e, --extended      add two additional fields, file and id to the output
  -h, --help          help for said
  -i, --ips-only      only print IP addresses

Global Flags:
  -a, --archive-regex string   regex to match archive files in the search dir (default "\\.(7z|bz2|gz|tar|xz|zip|xz|zst|lz)$")
  -t, --concurrency int        number of concurrent workers to use (default 12)
  -c, --config string          .env config file path (or via env variable CONFIG)
  -f, --file-regex string      regex to match files in the search dir (default ".*\\.log$")
  -A, --include-archive        search inside archive files
  -o, --output string          output format, one of 'json' or 'text' (default "text")
  -d, --search-dir string      directory to search for files recursively (default ".")
```

example:

```bash

# get minimal information about the players that said the phrase 'https?://bot.xyz\..+'
./twlog who said 'https?://bot.xyz'

# get all information about the players that said the phrase 'https?://bot.xyz\..+'
./twlog who said -e 'https?://bot.xyz'

# get all information of all players that said the phrase 'https?://bot.xyz\..+' but deduplicate all entries
./twlog who said -e -D 'https?://bot.xyz'

# get all deduplicated ip addresses of all players that said the phrase 'https?://bot.xyz\..+'
./twlog who said -D -i -o json 'https?://bot.xyz'
````
