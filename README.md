# twlog-who-said

twlog-who-said is a small utility that analyzes Teeworlds log files in order to determine who said a specific phrase.
This is useful for determining a list of ip addresses of spam bots based on a specific phrase they use.

## installation

Download the executable from the releases page or install it using the Go toolchain:

```bash
go install github.com/jxsl13/twlog-who-said@latest
```

## usage

```bash
$ twlog-who-said --help
Environment variables:
  PHRASE_REGEX    regex to search for that a player said
  SEARCH_DIR      directory to search for files recursively (default: ".")
  FILE_REGEX      regex to match files in the search dir (default: ".*\\.log")
  DEDUPLICATE     deduplicate objects based on all fields (default: "false")
  EXTENDED        add two additional fields, file and id to the output (default: "false")
  IPS_ONLY        only print IP addresses (default: "false")
  OUTPUT          output format, one of 'json' or 'text' (default: "text")

Usage:
  twlog-who-said [flags]

Flags:
  -c, --config string         .env config file path (or via env variable CONFIG)
  -D, --deduplicate           deduplicate objects based on all fields
  -e, --extended              add two additional fields, file and id to the output
  -f, --file-regex string     regex to match files in the search dir (default ".*\\.log")
  -h, --help                  help for twlog-who-said
  -i, --ips-only              only print IP addresses
  -o, --output string         output format, one of 'json' or 'text' (default "text")
  -p, --phrase-regex string   regex to search for that a player said
  -d, --search-dir string     directory to search for files recursively (default ".")
```

example:
```bash

# get all information about the players that said the phrase 'https?://bot.xyz\..+'
./twlog-who-said -e -p 'https?://bot.xyz'

# get all information of all players that said the phrase 'https?://bot.xyz\..+' but deduplicate all entries
./twlog-who-said -e -D -p 'https?://bot.xyz'

# get all deduplicated ip addresses of all players that said the phrase 'https?://bot.xyz\..+'
./twlog-who-said -D -p 'https?://bot.xyz' -i -o json
````

## building and installing from source

```bash
# building
go build .

# installing
go install .
```
