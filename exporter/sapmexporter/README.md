# SAPM Exporter

<!-- status autogenerated section -->
| Status        |           |
| ------------- |-----------|
| Stability     | [beta]: traces   |
| Distributions | [contrib] |

[beta]: https://github.com/open-telemetry/opentelemetry-collector#beta
[contrib]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
<!-- end autogenerated section -->

The SAPM exporter builds on the Jaeger proto and adds additional batching on top. This allows
the collector to export traces from multiples nodes/services in a single batch. The SAPM proto
and some useful related utilities can be found [here](https://github.com/signalfx/sapm-proto/).

> Please review the Collector's [security
> documentation](https://github.com/open-telemetry/opentelemetry-collector/blob/main/docs/security-best-practices.md),
> which contains recommendations on securing sensitive information such as the
> API key required by this exporter.

## Configuration

The following configuration options are required:

- `access_token` (no default): AccessToken is the authentication token provided by SignalFx or
another backend that supports the SAPM proto. The SignalFx access token can be obtained from the
web app. For details on how to do so please refer the documentation [here](https://docs.signalfx.com/en/latest/admin-guide/tokens.html#access-tokens).
- `endpoint` (no default): This is the destination to where traces will be sent to in SAPM
format. It must be a full URL and include the scheme, port and path e.g,
<!-- markdown-link-check-disable-line -->https://ingest.us0.signalfx.com/v2/trace. This can be pointed to the SignalFx 
backend or to another Otel collector that has the SAPM receiver enabled.

The following configuration options can also be configured:

- `max_connections` (default = 100): MaxConnections is used to set a limit to the maximum
idle HTTP connection the exporter can keep open.
- `num_workers` (default = 8): NumWorkers is the number of workers that should be used to
export traces. Exporter can make as many requests in parallel as the number of workers. Note
that this will likely be removed in future in favour of processors handling parallel exporting.
- `access_token_passthrough`: (default = `true`) Whether to use `"com.splunk.signalfx.access_token"`
trace resource attribute, if any, as SFx access token.  In either case this attribute will be deleted
during final translation.  Intended to be used in tandem with identical configuration option for
[SAPM receiver](../../receiver/sapmreceiver/README.md) to preserve trace origin.
- `timeout` (default = 5s): Is the timeout for every attempt to send data to the backend.
- `log_detailed_response` (default = `false`): Option to log detailed response from Splunk APM.
In addition to setting this option to `true`, debug logging at the Collector level needs to be enabled.

In addition, this exporter offers queued retry which is enabled by default.
Information about queued retry configuration parameters can be found
[here](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md).

Example:

```yaml
exporters:
  sapm:
    access_token: YOUR_ACCESS_TOKEN
    access_token_passthrough: true
    endpoint: https://ingest.YOUR_SIGNALFX_REALM.signalfx.com/v2/trace
    max_connections: 100
    num_workers: 8
    log_detailed_response: true
```

The full list of settings exposed for this exporter are documented [here](config.go)
with detailed sample configurations [here](testdata/config.yaml).

This exporter also offers proxy support as documented
[here](https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter#proxy-support).
