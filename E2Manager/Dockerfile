##############################################################################
#
#   Copyright (c) 2019 AT&T Intellectual Property.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#
##############################################################################

FROM nexus3.o-ran-sc.org:10004/bldr-ubuntu16-c-go:1-u16.04-nng1.1.1 as ubuntu

WORKDIR /opt/E2Manager
COPY . . 
ENV PATH=$PATH:/usr/local/go/bin:/usr/lib/go-1.12/bin
#RUN git clone https://gerrit.o-ran-sc.org/r/ric-plt/lib/rmr && cd rmr/; mkdir build; cd build; /opt/bin/cmake -DDEV_PKG=1 ..; make install
# Install RMr library and dev files
RUN wget --content-disposition  https://packagecloud.io/o-ran-sc/staging/packages/debian/stretch/rmr_1.10.0_amd64.deb/download.deb
RUN dpkg -i rmr_1.10.0_amd64.deb
RUN wget --content-disposition https://packagecloud.io/o-ran-sc/staging/packages/debian/stretch/rmr-dev_1.10.0_amd64.deb/download.deb
RUN dpkg -i rmr-dev_1.10.0_amd64.deb

RUN cd asn1codec && make  
RUN go build app/main.go

# Execute UT
ENV LD_LIBRARY_PATH=/usr/local/lib

# cgocheck=2 enables expensive checks that should not miss any errors, but will cause your program to run slower.
# clobberfree=1 causes the garbage collector to clobber the memory content of an object with bad content when it frees the object.
# gcstoptheworld=1 disables concurrent garbage collection, making every garbage collection a stop-the-world event.
# Setting gcstoptheworld=2 also disables concurrent sweeping after the garbage collection finishes.
# Setting allocfreetrace=1 causes every allocation to be profiled and a stack trace printed on each object's allocation and free.
ENV GODEBUG=cgocheck=2,clobberfree=1,gcstoptheworld=2,allocfreetrace=0
ENV RIC_ID="bbbccc-abcd0e/20"
RUN go test ./...

FROM ubuntu:16.04

RUN apt-get update && apt-get install -y \
  net-tools \
  iputils-ping \
  curl \
  tcpdump
  
COPY --from=ubuntu /opt/E2Manager/router.txt /opt/E2Manager/router.txt
COPY --from=ubuntu /opt/E2Manager/main /opt/E2Manager/main
COPY --from=ubuntu /opt/E2Manager/resources /opt/E2Manager/resources
COPY --from=ubuntu /usr/local/lib/librmr_nng.so.1 /usr/local/lib/librmr_nng.so.1
COPY --from=ubuntu /usr/local/lib/libnng.so.1 /usr/local/lib/libnng.so.1
WORKDIR /opt/E2Manager
ENV LD_LIBRARY_PATH=/usr/local/lib \
    port=3800 
ENV RMR_SEED_RT=router.txt
EXPOSE 3800
CMD ["sh", "-c", "./main  -port=$port"]
