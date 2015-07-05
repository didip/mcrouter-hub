[![GoDoc](https://godoc.org/github.com/didip/mcrouter-hub?status.svg)](http://godoc.org/github.com/didip/mcrouter-hub)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/didip/mcrouter-hub/master/LICENSE)


# mcrouter-hub

It's an HTTP companion to Facebook's McRouter.

It allows CRUD operation on McRouter config JSON as well as providing live stats.


## Agent Mode

### Getting Started
```
MCROUTER_ADDR=localhost:5000 \
MCROUTER_CONFIG_FILE=./tests/mcrouter.json \
MCRHUB_CENTRAL_URLS=http://localhost:5002 \
/path/to/mcrouter-hub
```

The agent mode allows mcrouter-hub to collect live information about McRouter. It can also modify the configuration of local McRouter.

While in agent mode, mcrouter-hub has the following required configurations:

* **MCROUTER_ADDR:** McRouter host and port address.

* **MCROUTER_CONFIG_FILE:** Path to McRouter config file.


## Central Mode

### Getting Started
```
MCRHUB_MODE=central \
/path/to/mcrouter-hub
```

While in central mode, mcrouter-hub can receives configuration or stats data from all mcrouter-hub agents.


## Complete Configuration List

* **MCROUTER_ADDR:** McRouter host and port address (Required for agent mode).

* **MCROUTER_CONFIG_FILE:** Path to McRouter config file (Required for agent mode).

* **MCRHUB_MODE:** 2 modes: agent or central. Default: `"agent"`

* **MCRHUB_READ_ONLY:** Read only flag. When true, `POST` endpoints are disabled. Default: Agent: `true`, Central: `false`

* **MCRHUB_LOG_LEVEL:** Log level. Default: `"info"`

* **MCRHUB_ADDR:** The HTTP server host and port. Default: Agent: `":5001"`, Central: `":5002"`

* **MCRHUB_CERT_FILE:** Path to cert file. Default: `""`

* **MCRHUB_KEY_FILE:** Path to key file. Default: `""`

* **MCRHUB_TOKENS_DIR:** Directory of token files. If it is not empty, each line on the file will be treated as one token. These tokens will be used for HTTP POST actions from agent to central. Tokens that exist on both sides are considered valid. Default: `""`
    ```
    # Token is passed as HTTP user. Example on updating config in agent mode:
    curl -u 0b79bab50daca910b000d4f1a2b675d604257e42: https://localhost:5001/config

    # Example on updating config in central mode:
    curl -u 0b79bab50daca910b000d4f1a2b675d604257e42: https://localhost:5002/configs

    # Example on updating stats in central mode:
    curl -u 0b79bab50daca910b000d4f1a2b675d604257e42: https://localhost:5002/stats
    ```

* **MCRHUB_CENTRAL_URLS:** URLs to mcrouter-hub central. mcrouter-hub does not use persistent storage, therefore agents report to one or more central daemons (Useful on agent mode only). Default: `""`

* **MCRHUB_REPORT_INTERVAL:** The frequency of agent reporting (Useful on agent mode only). Default: `"3s"`

* **NR_INSIGHTS_URL:** Newrelic Insights URL endpoint (Useful on agent mode only). Example: https://insights-collector.newrelic.com/v1/accounts/{AccountID}/events. Default: `""`

* **NR_INSIGHTS_INSERT_KEY:** Newrelic Insights insert key (Useful on agent mode only). Default: `""`
