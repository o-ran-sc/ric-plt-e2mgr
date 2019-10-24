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

WORKDIR /opt/xapp_mock
COPY . . 
ENV PATH=$PATH:/usr/local/go/bin:/usr/lib/go-1.12/bin
# Install RMr library and dev files
RUN wget --content-disposition https://packagecloud.io/o-ran-sc/master/packages/debian/stretch/rmr_1.3.0_amd64.deb/download.deb
RUN dpkg -i rmr_1.3.0_amd64.deb
RUN wget --content-disposition https://packagecloud.io/o-ran-sc/master/packages/debian/stretch/rmr-dev_1.3.0_amd64.deb/download.deb
RUN dpkg -i rmr-dev_1.3.0_amd64.deb

RUN go get github.com/pkg/errors 
RUN go build main/xapp_mock.go 
    


FROM ubuntu:16.04
COPY --from=ubuntu /opt/xapp_mock/xapp_mock /opt/xapp_mock/xapp_mock
COPY --from=ubuntu /opt/xapp_mock/resources /opt/xapp_mock/resources
COPY --from=ubuntu /usr/local/lib/librmr_nng.so.1 /usr/local/lib/librmr_nng.so.1
COPY --from=ubuntu /usr/local/lib/libnng.so.1 /usr/local/lib/libnng.so.1
WORKDIR /opt/xapp_mock
ENV LD_LIBRARY_PATH=/usr/local/lib 
ENV RMR_SEED_RT=resources/router.txt
ENV RMR_PORT=5001
#CMD mkdir -p resources/conf  exec ./xapp_mock
CMD mkdir -p resources/conf &&  exec /bin/bash 
