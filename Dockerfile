FROM debian:bullseye-slim

LABEL maintainer="DevOps Kung Fu Masters"

RUN apt-get update && apt-get upgrade -y && apt-get install jq curl wget -y 
RUN curl -s https://api.github.com/repos/open-policy-agent/conftest/releases/latest | jq -r ".assets[] | select(.name | contains(\"_linux_amd64.deb\")) | .browser_download_url" | wget -i -
RUN dpkg -i *.deb
ADD bin/domi domi

EXPOSE 8080

RUN apt-get autoremove -y && apt-get remove jq curl wget -y
RUN groupadd -r domi && useradd -r -g domi domi
USER domi

CMD ["./domi"]