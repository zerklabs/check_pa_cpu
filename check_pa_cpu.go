package main

import (
	"flag"
	"fmt"

	"github.com/alouca/gosnmp"
	"github.com/fractalcat/nagiosplugin"
)

func main() {
	var (
		host      string
		community string
		timeout   int64
		mode      string
		critical  int64
		warning   int64
	)

	flag.StringVar(&host, "H", "127.0.0.1", "Target host")
	flag.StringVar(&community, "community", "public", "SNMP community string")
	flag.Int64Var(&timeout, "timeout", 10, "SNMP connection timeout")

	flag.StringVar(&mode, "mode", "", "Specify which cpu mode to check. management-cpu or data-cpu")
	flag.Int64Var(&warning, "warning", 80, "Warning threshold")
	flag.Int64Var(&critical, "critical", 90, "Critical threshold")

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

	utilization, err := getData(c, mode)
	if err != nil {
		check.AddResult(nagiosplugin.UNKNOWN, fmt.Sprintf("error: %v", err))
		return
	}

	check.AddPerfDatum("utilization", "%", float64(utilization), float64(warning), float64(critical))

	crit, _, _, err := parseRange(critical, utilization)
	if err != nil {
		check.AddResult(nagiosplugin.UNKNOWN, fmt.Sprintf("error: %v", err))
		return
	}

	if crit {
		check.AddResult(nagiosplugin.CRITICAL, fmt.Sprintf("%s utilization - %d%%", mode, utilization))
		return
	}

	warn, _, _, err := parseRange(warning, utilization)
	if err != nil {
		check.AddResult(nagiosplugin.UNKNOWN, fmt.Sprintf("error: %v", err))
		return
	}

	if warn {
		check.AddResult(nagiosplugin.WARNING, fmt.Sprintf("%s utilization - %d%%", mode, utilization))
		return
	}

	check.AddResult(nagiosplugin.OK, fmt.Sprintf("%s utilization - %d%%", mode, utilization))
}

func getData(s *gosnmp.GoSNMP, oidType string) (int, error) {
	val := -1

	pkt, err := s.Get(oids[oidType])
	if err != nil {
		return val, err
	}

	for _, v := range pkt.Variables {
		switch v.Type {
		case gosnmp.Integer:
			val = v.Value.(int)
		}
	}

	return val, nil
}

func parseRange(r int64, val int) (bool, float64, float64, error) {
	nr, err := nagiosplugin.ParseRange(fmt.Sprintf("%d", r))
	if err != nil {
		return false, 0, 0, err
	}

	return nr.CheckInt(val), nr.Start, nr.End, nil
}
