FROM golang:1.13-alpine AS build

COPY main.go .

ARG LDFLAGS

RUN GOOS=linux GOARCH=386 go build -ldflags "${LDFLAGS}" -o owlet

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/owlet /go/owlet
ENTRYPOINT ["/go/owlet"]

ARG NAME
ARG VERSION
ARG COMMIT
ARG BUILD_DATE

LABEL maintainer="Kamil Sindi" repository="https://github.com/ksindi/owlet" homepage="https://github.com/ksindi/owlet"

LABEL org.label-schema.name="${NAME}" org.label-schema.build-date="${BUILD_DATE}" org.label-schema.vcs-ref="${COMMIT}" org.label-schema.version="${VERSION}" org.label-schema.schema-version="1.0"

LABEL com.github.alerts.name="${NAME}" com.github.alerts.description="enable vulnerability alerts" com.github.alerts.icon="github" com.github.alerts.color="black"
