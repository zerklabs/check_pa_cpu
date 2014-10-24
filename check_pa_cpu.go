package main

import (
	"flag"
	"fmt"

	"github.com/alouca/gosnmp"
	"github.com/fractalcat/nagiosplugin"
)

func main() {
	var (
		host           string
		community      string
		timeout        int64
		management_cpu bool
		data_cpu       bool
	)

	flag.StringVar(&host, "H", "127.0.0.1", "Target host")
	flag.StringVar(&community, "community", "public", "SNMP community string")
	flag.Int64Var(&timeout, "timeout", 10, "SNMP connection timeout")

	flag.BoolVar(&management_cpu, "mode-management-cpu", false, "Check the management plane CPU utilization")
	flag.BoolVar(&data_cpu, "mode-data-cpu", false, "Check the data plane CPU utilization")

	flag.Parse()

	// Initialize the check - this will return an UNKNOWN result
	// until more results are added.
	check := nagiosplugin.NewCheck()
	// If we exit early or panic() we'll still output a result.
	defer check.Finish()

	// obtain data here
	c, err := gosnmp.NewGoSNMP(host, community, gosnmp.Version2c, timeout)
	if err != nil {
		check.AddResult(nagiosplugin.UNKNOWN, fmt.Sprintf("error: %v", err))
		return
	}

	if management_cpu {
		pkt, err := c.Get(oids["management-plane-cpu"])
		if err != nil {
			check.AddResult(nagiosplugin.UNKNOWN, fmt.Sprintf("error: %v", err))
			return
		}

		for _, v := range pkt.Variables {
			switch v.Type {
			case gosnmp.Integer:
				val := v.Value.(int)
				if val >= 80 && val < 90 {
					check.AddResult(nagiosplugin.WARNING, fmt.Sprintf("Management plane cpu utilization - %d%%", val))
				} else if val >= 90 {
					check.AddResult(nagiosplugin.CRITICAL, fmt.Sprintf("Management plane cpu utilization - %d%%", val))
				} else {
					check.AddResult(nagiosplugin.OK, fmt.Sprintf("Management plane cpu utilization - %d%%", val))
				}

				check.AddPerfDatum("cpu", "%", float64(val), 80.0, 90.0)
			}
		}
	}

	if data_cpu {
		pkt, err := c.Get(oids["data-plane-cpu"])
		if err != nil {
			check.AddResult(nagiosplugin.UNKNOWN, fmt.Sprintf("error: %v", err))
			return
		}

		for _, v := range pkt.Variables {
			switch v.Type {
			case gosnmp.Integer:
				val := v.Value.(int)
				if val >= 80 && val < 90 {
					check.AddResult(nagiosplugin.WARNING, fmt.Sprintf("Data plane cpu utilization - %d%%", val))
				} else if val >= 90 {
					check.AddResult(nagiosplugin.CRITICAL, fmt.Sprintf("Data plane cpu utilization - %d%%", val))
				} else {
					check.AddResult(nagiosplugin.OK, fmt.Sprintf("Data plane cpu utilization - %d%%", val))
				}

				check.AddPerfDatum("cpu", "%", float64(val), 80.0, 90.0)
			}
		}
	}
}
