package main

import (
	"github.com/consol-monitoring/check_prometheus/helper"
	"github.com/consol-monitoring/check_prometheus/mode"
	"github.com/consol-monitoring/check_x"
	"github.com/urfave/cli"
	"os"
	"time"
)

var (
	address           string
	timeout           int
	warning           string
	critical          string
	query             string
	alias             string
	search            string
	replace           string
	label             string
	emptyQueryMessage string
	emptyQueryStatus  string
)

func startTimeout() {
	if timeout != 0 {
		check_x.StartTimeout(time.Duration(timeout) * time.Second)
	}
}

func getStatus(state string) check_x.State {
	switch state {
	case "OK":
		return check_x.OK
	case "WARNING":
		return check_x.Warning
	case "CRITICAL":
		return check_x.Critical
	default:
		return check_x.Unknown
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "check_prometheus"
	app.Usage = "Checks different prometheus stats as well the data itself\n   Copyright (c) 2017 Philip Griesbacher\n   https://github.com/Griesbacher/check_prometheus"
	app.Version = "0.0.3"
	flagAddress := cli.StringFlag{
		Name:        "address",
		Usage:       "Prometheus address: Protocol + IP + Port.",
		Destination: &address,
		Value:       "http://localhost:9100",
	}
	flagWarning := cli.StringFlag{
		Name:        "w",
		Usage:       "Warning value. Use nagios-plugin syntax here.",
		Destination: &warning,
	}
	flagCritical := cli.StringFlag{
		Name:        "c",
		Usage:       "Critical value. Use nagios-plugin syntax here.",
		Destination: &critical,
	}
	flagQuery := cli.StringFlag{
		Name:        "q",
		Usage:       "Query to be executed",
		Destination: &query,
	}
	flagAlias := cli.StringFlag{
		Name:        "a",
		Usage:       "Alias, will replace the query within the output, if set. You can use go text/template syntax to output label values (only for vector results).",
		Destination: &alias,
	}
	flagEmptyQueryMessage := cli.StringFlag{
		Name:        "eqm",
		Usage:       "Message if the query returns no data.",
		Destination: &emptyQueryMessage,
	}
	flagEmptyQueryStatus := cli.StringFlag{
		Name:        "eqs",
		Usage:       "Status if the query returns no data.",
		Destination: &emptyQueryStatus,
	}
	flagLabel := cli.StringFlag{
		Name:        "l",
		Usage:       "Prometheus-Label, which will be used for the performance data label. By default job and instance should be available.",
		Destination: &label,
		Value:       mode.DefaultLabel,
	}
	app.Commands = []cli.Command{
		{
			Name:    "mode",
			Aliases: []string{"m"},
			Usage:   "check mode",
			Subcommands: []cli.Command{
				{
					Name:        "ping",
					Aliases:     []string{"p"},
					Usage:       "Returns the build informations",
					Description: `This check requires that the prometheus server itself is listetd as target. Following query will be used: 'prometheus_build_info{job="prometheus"}'`,
					Action: func(c *cli.Context) error {
						startTimeout()
						return mode.Ping(address)
					},
					Flags: []cli.Flag{
						flagAddress,
					},
				}, {
					Name:    "query",
					Aliases: []string{"q"},
					Usage:   "Checks collected data",
					Description: `Your Promqlquery has to return a vector / scalar / matrix result. The warning and critical values are applied to every value.
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

	   `,

					Action: func(c *cli.Context) error {
						startTimeout()
						return mode.Query(address, query, warning, critical, alias, search, replace, emptyQueryMessage, getStatus(emptyQueryStatus))
					},
					Flags: []cli.Flag{
						flagAddress,
						flagQuery,
						flagAlias,
						flagWarning,
						flagCritical,
						cli.StringFlag{
							Name:        "search",
							Usage:       "If this variable is set, the given Golang regex will be used to search and replace the result with the 'replace' flag content. This will be appied on the perflabels.",
							Destination: &search,
						},
						cli.StringFlag{
							Name:        "replace",
							Usage:       "See search flag. If the 'search' flag is empty this flag will be ignored.",
							Destination: &replace,
						},
						flagEmptyQueryMessage,
						flagEmptyQueryStatus,
					},
				}, {
					Name:        "targets_health",
					Usage:       "Returns the health of the targets",
					Description: `The warning and critical thresholds are appied on the health_rate. The health_rate is calculted: sum(healthy) / sum(targets).`,
					Action: func(c *cli.Context) error {
						startTimeout()
						return mode.TargetsHealth(address, label, warning, critical)
					},
					Flags: []cli.Flag{
						flagAddress,
						flagWarning,
						flagCritical,
						flagLabel,
					},
				},
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "t",
			Usage:       "Seconds till check returns unknown, 0 to disable",
			Value:       10,
			Destination: &timeout,
		},
		cli.IntFlag{
			Name:        "f",
			Usage:       "If the checked data is older then this in seconds, unknown will be returned. Set to 0 to disable.",
			Value:       300,
			Destination: &helper.TimestampFreshness,
		},
	}

	if err := app.Run(os.Args); err != nil {
		check_x.ErrorExit(err)
	}
}
