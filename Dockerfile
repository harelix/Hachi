FROM debian:bullseye-slim

RUN apt-get update

RUN apt-get install ca-certificates -y

ADD hachi /usr/local/bin

RUN chmod 755 /usr/local/bin/hachi

RUN mkdir -p /etc/hachi

ENTRYPOINT ["/bin/bash", "-c"]

CMD ["/usr/local/bin/hachi"]
