FROM devopsfaith/krakend

USER root

RUN apt-get update && apt-get install -y ca-certificates

COPY ./conf/certs/pem/* /usr/share/ca-certificates/nistagram/
RUN realpath --relative-to=/usr/share/ca-certificates/ /usr/share/ca-certificates/nistagram/* >> /etc/ca-certificates.conf && update-ca-certificates && rm -rf /var/lib/apt/lists/* 

VOLUME [ "/etc/krakend" ]

WORKDIR /etc/krakend

ENTRYPOINT [ "/usr/bin/krakend" ]
CMD [ "run", "-c", "/etc/krakend/krakend.json" ]

EXPOSE 8000 8090 80