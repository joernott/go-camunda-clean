# camunda-clean - Clean leftover camunda processes

## Purpose
This little tool tries to delete all leftover processes from the camunda
process-instance list.

Feel free to contribute and improve.

## License
BSD 3-Clause

## Usage
  camunda-clean [flags]

### Flags:
|Short | Long          | Type  | Purpose                                       |
|------|---------------|-------|-----------------------------------------------|
| -h   | --help        |       | help for camunda-clean                        |
| -c   | --config      |string | Configuration file (default camunda-clean.yml in the current working directory)|
| -u   | --user        |string | Username for Camunda                          |
| -p   | --password    |string | Password for the Camunda user                 |
| -s   | --ssl         |bool   | Use SSL                                       |
| -v   | --validatessl |bool   | Validate SSL certificate (default true)       |
| -H   | --host        |string | Hostname of the server (default "localhost")  |
| -P   | --port        |int    | Network port (default 8080)                   |
| -B   | --baseendpoint|string | Base endpoint (default /engine-rest)          |
| -y   | --proxy       |string | Proxy (defaults to none)                      |
| -Y   | --socks       |bool   | This is a SOCKS proxy                         |
| -l   | --loglevel    |int    | Log level (default 5)                         |
| -L   | --logfile     |string | Log file (defaults to stdout)                 |

