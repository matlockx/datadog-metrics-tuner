datadog-metrics-tuner
=====================

Docker image which can be used to report metrics to datadog which have been defined in a config file. Why? It might be useful when you are reporting business metrics from your services to datadog. After some time the service might not be used anymore but you still need the business metrics on your dashboards. It's easier to have an app running which reports them than having all your old services running to just report metrics. The container reads all files from a given directory and reports the configured metrics to datadog.

Currently only gauges are provided.

### Config Format

Config files are yaml files. Just add a list of your metrics with the following format:

```yaml
- name: METRIC_NAME
  tags: ["TAG_NAME:TAG_VALUE", "FOO:BAR"]
  value: 123

- name: METRIC_NAME2
  tags: ["TAG_NAME2:TAG_VALUE2", "FOO2:BAR2"]
  value: 111
```

### Usage

Put all your yaml config files in a directory. Then mount it into the container and provide datadog app and api key:

```
docker run -d -v [PATH_TO_METRICS_DIR]:/metrics.d matlockx/datadog-metrics-tuner -api-key [API_KEY] -app-key [APP_KEY]
```
