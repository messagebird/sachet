# Sachet
Sachet (or सचेत) is Hindi for conscious. Sachet is an SMS alerting tool for the [Prometheus Alertmanager](https://github.com/prometheus/alertmanager).

## The problem
There are many SMS providers and Alertmanager supporting all of them would make the code noisy. To get around this issue a new service needed to be created dedicated only for SMS.

## The solution
An HTTP API that accepts Alertmanager webhook calls and allows an end-user to configure it for the SMS provider of their dreams.

## Usage

Running Sachet is as easy as executing `sachet` on the command line.

```
$ sachet
```

Use the `-h` flag to get help information.

```
$ sachet -h
Usage of sachet:
  -config string
        The configuration file (default "config.yaml")
  -listen-address string
        The address to listen on for HTTP requests. (default ":9876")
```

## Testing

Sachet expects a JSON object from Alertmanager. The format of this JSON is described in [the Alertmanager documentation](https://prometheus.io/docs/alerting/configuration/#webhook-receiver-<webhook_config>), or, alternatively, in [the Alertmanager GoDoc](https://godoc.org/github.com/prometheus/alertmanager/template#Data).

To quickly test Sachet is working you can also run:
```bash
$ curl -H "Content-type: application/json" -X POST \
  -d '{"receiver": "team-sms", "status": "firing", "alerts": [{"status": "firing", "labels": {"alertname": "test-123"} }], "commonLabels": {"key": "value"}}' \
  http://localhost:9876/alert
```

## Alertmanager configuration

To enable Sachet you need to configure a webhook in Alertmanager. You can do that by adding a webhook receiver to your Alertmanager configuration. 

```yaml
receivers:
- name: 'team-sms'
  webhook_configs:
  - url: 'http://localhost:9876/alert'
```

## License

Sachet is licensed under [The BSD 2-Clause License](http://opensource.org/licenses/BSD-2-Clause). Copyright (c) 2016, MessageBird
