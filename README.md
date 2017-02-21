# prometheus-kemp-exporter

[![Build Status](https://api.travis-ci.org/giantswarm/prometheus-kemp-exporter.svg)](https://travis-ci.org/giantswarm/prometheus-kemp-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/giantswarm/prometheus-kemp-exporter)](https://goreportcard.com/report/github.com/giantswarm/prometheus-kemp-exporter)
[![GoDoc](https://godoc.org/github.com/giantswarm/prometheus-kemp-exporter?status.svg)](http://godoc.org/github.com/giantswarm/prometheus-kemp-exporter)
[![Docker](https://img.shields.io/docker/pulls/giantswarm/prometheus-kemp-exporter.svg)](http://hub.docker.com/r/giantswarm/prometheus-kemp-exporter)
[![IRC Channel](https://img.shields.io/badge/irc-%23giantswarm-blue.svg)](https://kiwiirc.com/client/irc.freenode.net/#giantswarm)

`prometheus-kemp-exporter` exports Kemp statistics to Prometheus.

## Prerequisites

## Getting `prometheus-kemp-exporter`

Download the latest release: https://github.com/giantswarm/prometheus-kemp-exporter/releases/latest

Clone the git repository: https://github.com/giantswarm/prometheus-kemp-exporter.git

Download the latest docker image from here: https://hub.docker.com/r/giantswarm/prometheus-kemp-exporter/


### How to build

#### Dependencies

- [github.com/giantswarm/kemp-client](https://github.com/giantswarm/kemp-client)
- [github.com/prometheus/client_golang](https://github.com/prometheus/client_golang)
- [github.com/spf13/cobra](https://github.com/spf13/cobra)

#### Building the binary

```
make
```

#### Building the docker image

```
make docker-image
```


## Running `prometheus-kemp-exporter`

Running the binary directly:
```
$ prometheus-kemp-exporter server <ENDPOINT> <USERNAME> <PASSWORD>
2016/06/10 10:23:15 Listening on port 8000
```

Running in a Docker container:
```
$ docker run -p 8000:8000 giantswarm/prometheus-kemp-exporter:latest server <ENDPOINT> <USERNAME> <PASSWORD>
2016/06/10 09:24:03 Listening on port 8000
```

The `prometheus-kemp-exporter` have to authenticate with Basic-Auth at the XML Interface on the KEMP Load Master. See [KEMP - RESTful API Interface](https://support.kemptechnologies.com/hc/en-us/articles/201640799-RESTful-API-interface) for more information. You may add a read only in your LoadMaster user for this task. The <ENDPOINT> is the URL where you can get the XML stats. You may verify the URL by using a web browser or `curl`:

```
$ curl -u monitoring-user:monitoring-password https://loadbalancer.example.com/access/stats/
<?xml version="1.0" encoding="ISO-8859-1"?>
<Response stat="200" code="ok">
<Success><Data><CPU><total><User>1</User>
<System>1</System>
<Idle>98</Idle>
...
```

Help information can be found with the `--help` flag.

## Contact

- Mailing list: [giantswarm](https://groups.google.com/forum/!forum/giantswarm)
- IRC: #[giantswarm](irc://irc.freenode.org:6667/#giantswarm) on freenode.org
- Bugs: [issues](https://github.com/giantswarm/prometheus-kemp-exporter/issues)

## Contributing & Reporting Bugs

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches, the contribution workflow as well as reporting bugs.

## License

`prometheus-kemp-exporter` is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.
