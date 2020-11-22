FROM golang:1.13-alpine AS build

COPY main.go .

ARG LDFLAGS

RUN GOOS=linux GOARCH=386 go build -ldflags "${LDFLAGS}" -o snoo

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/snoo /go/snoo
ENTRYPOINT ["/go/snoo"]

ARG NAME
ARG VERSION
ARG COMMIT
ARG BUILD_DATE

LABEL maintainer="Kamil Sindi" repository="https://github.com/ksindi/snoo" homepage="https://github.com/ksindi/snoo"

LABEL org.label-schema.name="${NAME}" org.label-schema.build-date="${BUILD_DATE}" org.label-schema.vcs-ref="${COMMIT}" org.label-schema.version="${VERSION}" org.label-schema.schema-version="1.0"

LABEL com.github.alerts.name="${NAME}" com.github.alerts.description="API client for the SNOO" com.github.alerts.icon="github" com.github.alerts.color="black"
