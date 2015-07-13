FROM ubuntu:14.04


COPY requirements.txt /tmp/
COPY ./src /tmp/src

RUN apt-get update
RUN apt-get install -yq python python-pip
RUN pip install -r /tmp/requirements.txt
RUN dlibBullshitInstallScript

CMD bash
