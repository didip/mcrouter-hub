[![GoDoc](https://godoc.org/github.com/didip/mcrouter-hub?status.svg)](http://godoc.org/github.com/didip/mcrouter-hub)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/didip/mcrouter-hub/master/LICENSE)


# mcrouter-hub

An HTTP companion to Facebook's McRouter. It performs CRUD operation on McRouter config JSON.

The daemon has 2 mode: agent and central.

While in agent mode, and if central URL is provided, mcrouter-hub gathers configuration data and reports it to central. Here's how to run the agent mode:
```
MCROUTER_CONFIG_FILE=./tests/mcrouter.json \
MCRHUB_CENTRAL_URL=http://localhost:5002 \
/path/to/mcrouter-hub
```

While in central mode, mcrouter-hub receives configuration data from individual McRouter host and performs CRUD operations on them. Here's how to run the central mode:
```
/path/to/mcrouter-hub
```

## Environment Variables

mcrouter-hub uses environment variables as configuration:

### Required

* **MCROUTER_CONFIG_FILE:** Path to McRouter config file (Required).


### Optional

* **MCRHUB_READ_ONLY:** Read only flag. Default: true

* **MCRHUB_CENTRAL_URL:** URL to mcrouter-hub central. Default: ""

* **MCRHUB_LOG_LEVEL:** Log level. Default: "info"

* **MCRHUB_ADDR:** The HTTP server host and port. Default: Agent: ":5001", Central: ":5002"

* **MCRHUB_CERT_FILE:** Path to cert file. Default: ""

* **MCRHUB_KEY_FILE:** Path to key file. Default: ""


### Newrelic Insights

* **NR_INSIGHTS_URL:** Newrelic Insights endpoint. Default: "https://insights-collector.newrelic.com/v1/accounts/1/events"

* **NR_INSIGHTS_INSERT_KEY:** Newrelic Insights insert key. Default: ""
