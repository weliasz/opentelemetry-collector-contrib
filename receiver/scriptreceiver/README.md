# Script Receiver

<!-- status autogenerated section -->
| Status        |                        |
| ------------- |------------------------|
| Stability     | [In Development]: logs |
| Distributions | [contrib]              |

[beta]: https://github.com/open-telemetry/opentelemetry-collector#beta
[contrib]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
<!-- end autogenerated section -->

## Configuration

The following settings are required:

- `exec_file` : Name of the script to be executed. Available script: bandwidth.sh, cpu.sh, df.sh, hardware.sh, interfaces.sh, iostat.sh, lastlog.sh, lsof.sh, netstat.sh, nfsiostat.sh, openPorts.sh, openPortsEnhanced.sh, package.sh, passwd.sh, protocol.sh, ps.sh, rlog.sh, selinuxChecker.sh, service.sh, setup.sh, sshdChecker.sh, time.sh, top.sh, update.sh, uptime.sh, usersWithLoginPrivs.sh, version.sh, vmstat.sh, vsftpdChecker.sh, who.sh
- `interval` : (default = `60s`) how often the script should be executed


The following settings are optional:

- `source` : source of the event
- `sourcetype` : sourcetype of the event
- `multiline` : how the standard output of the script is split
Example:

```yaml
receivers:
  script/df:
    exec_file: df.sh
    interval: 10s
    source: df
    sourcetype: df
    multiline:
      line_end_pattern: '\n'
```


```yaml
service:
  pipelines:
    logs:
      receivers: [script/df]
      processors: [memory_limiter, batch]
      exporters: [splunk_hec]
```