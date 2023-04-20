[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)

# check_prometheus
Monitoring Plugin to check the health of a Prometheus server and its data

## Usage
### Global Options
```
.\check_prometheus -h
NAME:
   check_prometheus - Checks different prometheus stats as well the data itself
   Copyright (c) 2017 Philip Griesbacher
   https://github.com/Griesbacher/check_prometheus

USAGE:
   check_prometheus [global options] command [command options] [arguments...]

VERSION:
   0.0.3

COMMANDS:
   mode, m  check mode
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -t value       Seconds till check returns unknown, 0 to disable (default: 10)
   -f value       If the checked data is older then this in seconds, unknown will be returned. Set to 0 to disable. (default: 300)
   --help, -h     show help
   --version, -v  print the version
```

### Command options

```
.\check_prometheus.exe mode -h
NAME:
   check_prometheus mode - check mode

USAGE:
   check_prometheus mode command [command options] [arguments...]

COMMANDS:
   ping, p         Returns the build informations
   query, q        Checks collected data
   targets_health  Returns the health of the targets

OPTIONS:
   --help, -h  show help
```

### Subcommand options example

```
NAME:
   check_prometheus mode query - Checks collected data

USAGE:
   check_prometheus mode query [command options] [arguments...]

DESCRIPTION:
   Your Promqlquery has to return a vector / scalar / matrix result. The warning and critical values are applied to every value.
   Examples:
       Vector:
           check_prometheus mode query -q 'up'
       --> OK - Query: 'up'|'up{instance="192.168.99.101:9245", job="iapetos"}'=1;;;; 'up{instance="0.0.0.0:9091", job="prometheus"}'=1;;;;

       Scalar:
           check_prometheus mode query -q 'scalar(up{job="prometheus"})'
       --> OK - OK - Query: 'scalar(up{job="prometheus"})' returned: '1'|'scalar'=1;;;;

       Matrix:
           check_prometheus mode query -q 'http_requests_total{job="prometheus"}[5m]'
       --> OK - Query: 'http_requests_total{job="prometheus"}[5m]'

       Search and Replace:
           check_prometheus m query -q 'up' --search '.*job=\"(.*?)\".*' --replace '$1'
       --> OK - Query: 'up'|'prometheus'=1;;;; 'iapetos'=0;;;;

           check_prometheus m q -q '{handler="prometheus",quantile="0.99",job="prometheus",__name__=~"http_.*bytes"}' --search '.*__name__=\"(.*?)\".*' --replace '$1' -a 'http_in_out'
       --> OK - Alias: 'http_in_out'|'http_request_size_bytes'=296;;;; 'http_response_size_bytes'=5554;;;;

     Use Alias to generate output with label values:
         Assumption that your query returns a label "hostname" and "details".
       IMPORTANT: To be able to use the value in more advanced output formatting, we just add a label "value" with the current value to the list of labels.
       If the specified Alias string cannot be processed by the text/template engine, the Alias string will be printed 1:1.
       check_prometheus m q -a 'Hostname: {{.hostname}} - Details: {{.details}}' --search '.*' --replace 'error_state'  -q 'http_requests_total{job="prometheus"}' -w 0 -c 0
     --> Critical - Hostname: Server01 - Details: Error404|'error_state'=1;0;0;;

     Use Alias with an if/else clause and the use of xvalue:
         If xvalue is 1, we output UP, else we output DOWN
         check_prometheus m q --search '.*' --replace 'up' -q 'up{instance="SUPERHOST"}' -a 'Hostname: {{.hostname}} Is {{if eq .xvalue "1"}}UP{{else}}DOWN{{end}}.\n' -w 1: -c 1:
       --> OK - Hostname: SUPERHOST Is UP.\n|'up'=1;1:;1:;;
       
     List all available labels to be used with Alias:
         Just use -a '{{.}}' and the whole map with all labels will be printed.
         check_prometheus m q -q 'up{instance="SUPERHOST"}' -a '{{.}}'
     --> OK - map[__name__:up hostname:SUPERHOST instance:SUPERHOST job:snmp mib:RittalCMC xvalue:1]|'{__name__="up", hostname="SUPERHOST", instance="SUPERHOST", job="snmp", mib="RittalCMC"}'=1;;;;
     
     Use Different Message and Status code for queries that return no data.
       If you have a query that only returns data in an error condition you can use this flags to return a custom message and status code.
           check_prometheus m q -eqm 'All OK' -eqs 'OK'  -q 'http_requests_total{job="prometheus"}' -w 0 -c 0
       --> OK - All OK
       Without -eqm, -eqs
           check_prometheus m q -q 'http_requests_total{job="prometheus"}' -w 0 -c 0
       --> UNKNOWN - The given States do not contain an State

     

OPTIONS:
   --address value  Prometheus address: Protocol + IP + Port. (default: "http://localhost:9100")
   -q value         Query to be executed
   -a value         Alias, will replace the query within the output, if set. You can use go text/template syntax to output label values (only for vector results).
   -w value         Warning value. Use nagios-plugin syntax here.
   -c value         Critical value. Use nagios-plugin syntax here.
   --search value   If this variable is set, the given Golang regex will be used to search and replace the result with the 'replace' flag content. This will be appied on the perflabels.
   --replace value  See search flag. If the 'search' flag is empty this flag will be ignored.
   --eqm value      Message if the query returns no data.
   --eqs value      Status if the query returns no data.

```
