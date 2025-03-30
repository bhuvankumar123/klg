# Jammy - ubuntu-22.04 LTS
FROM ubuntu:jammy 

WORKDIR /service

RUN mkdir -p /service/bin \
    && mkdir -p /service/scripts

ADD bin/* /service/bin/
ADD scripts/* /service/scripts/

CMD ["/service/scripts/run.sh"]