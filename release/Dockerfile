FROM ubuntu

WORKDIR /usr/src/app

COPY pgrokd .

COPY entrypoint.sh .

RUN mkdir -p certs

ENTRYPOINT [ "/usr/src/app/entrypoint.sh" ]
