# Sachet
Sachet (or सचेत) is Hindi for conscious. Sachet is an SMS alerting tool for the [Prometheus Alertmanager](https://github.com/prometheus/alertmanager).

## The problem
There are many SMS providers and Alertmanager supporting all of them would make the code noisy. To get around this issue a new service needed to be created dedicated only for SMS.

## The solution
An HTTP API that accepts Alertmanager webhook calls and allows an end-user to configure it for the SMS provider of their dreams.

## Authors
* Marcel Corso
* Sam Wierema