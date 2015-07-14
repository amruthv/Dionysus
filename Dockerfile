FROM ubuntu:14.04


COPY requirements.txt /tmp/
COPY . /project

RUN apt-get update
RUN apt-get install -yq python \
                        python-pip \
                        xorg \
                        openbox \
                        cmake \
                        libboost-python-dev
RUN pip install -r /tmp/requirements.txt
RUN bash /project/dlib_source/dlibBullshitInstallScript

CMD bash
