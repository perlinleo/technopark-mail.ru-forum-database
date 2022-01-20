FROM golang:latest AS builder

ADD . /app
WORKDIR /app
RUN go build ./cmd/main.go

FROM ubuntu:20.04

RUN apt-get -y update && apt-get install -y tzdata
ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER
USER postgres

RUN /etc/init.d/postgresql start &&\
    createdb -O postgres forum &&\
    /etc/init.d/postgresql stop


EXPOSE 5432
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
USER root

WORKDIR /usr/src/app

COPY . .
COPY --from=builder /app/main/ .

EXPOSE 5000
ENV PGPASSWORD password
CMD service postgresql start && psql -h localhost -d forum -U marvin -p 5432 -a -q -f ./db/dump.sql && ./main