FROM scratch
MAINTAINER Joseph Salisbury <joseph@giantswarm.io>

COPY ./prometheus-kemp-exporter /prometheus-kemp-exporter

EXPOSE 8000

ENTRYPOINT ["/prometheus-kemp-exporter"]
