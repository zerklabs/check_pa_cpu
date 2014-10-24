check\_pa\_cpu
============

Nagios check for Palo Alto Management and Data CPU utilization

## Usage

```
Usage of check_pa_cpu:
  -H="127.0.0.1": Target host
  -community="public": SNMP community string
  -critical=90: Critical threshold
  -mode="": Specify which cpu mode to check. management-cpu or data-cpu
  -timeout=10: SNMP connection timeout
  -warning=80: Warning threshold
```


## Examples

```
$> check_pa_cpu -H 1.1.1.1 -community="public" -mode data-cpu -timeout 5 --warning 80 --critical 90
OK: Data plane cpu utilization - 4% | cpu=4%;;;80;90
```

## References
[Useful SNMP OIDs for Monitoring Palo Alto Networks Devices](https://live.paloaltonetworks.com/docs/DOC-1744)
