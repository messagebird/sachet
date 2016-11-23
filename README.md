# Sachet
Sachet (or सचेत) is Hindi for conscious. Sachet is an SMS alerting tool for the [Prometheus Alertmanager](https://github.com/prometheus/alertmanager).

## The problem
There are many SMS providers and Alertmanager supporting all of them would make the code noisy. To get around this issue a new service needed to be created dedicated only for SMS.

## The solution
An HTTP API that accepts Alertmanager webhook calls and allows an end-user to configure it for the SMS provider of their dreams.

## Play with it

```
curl -H "Content-type: application/json" -X POST -d '{"receiver": "team-sms", "status": "error", "commonLabels": {"key": "value"}}' http://127.0.0.1:9876/alert

```


### Alert manager configuration

```
...
receivers:
- name: 'team-sms'
  webhook_configs:
  - url: 'http://localhost:9876/alert'
...
```

## License

Pushprom is licensed under [The BSD 2-Clause License](http://opensource.org/licenses/BSD-2-Clause). Copyright (c) 2016, MessageBird
