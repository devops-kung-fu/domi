FROM debian:bullseye-slim

LABEL maintainer="DevOps Kung Fu Masters"

ADD bin/domi domi
EXPOSE 8080

RUN groupadd -r domi && useradd -r -g domi domi
USER domi

CMD ["./domi"]
