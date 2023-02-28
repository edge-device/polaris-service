FROM ubuntu:22.04

RUN apt-get update 
RUN apt-get upgrade -y
RUN apt-get install -y mysql-client
COPY ./polaris /opt/polaris/polaris
RUN chmod +x /opt/polaris/polaris
ENTRYPOINT ["/opt/polaris/polaris"]
