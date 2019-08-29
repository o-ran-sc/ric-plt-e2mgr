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

*** Settings ***
Documentation    Resource file


*** Variables ***
${url}   http://localhost:3800
${json}    {"ranIp": "10.0.2.15","ranPort": 5577,"ranName":"test1"}
${endcjson}    {"ranIp": "10.0.2.15","ranPort": 49999,"ranName":"test2"}
${resetcausejson}   {"cause": "misc:not-enough-user-plane-processing-resources"}
${header}  {"Content-Type": "application/json"}
${docker_command}  docker ps | grep 1.0 | wc --lines
${run_simu}   docker run -d --name gnbe2_simu --env gNBipv4=localhost  --env gNBport=36422  --env duration=600000000000 --env indicationReportRate=1000000000 --env indicationInsertRate=0 -p 5577:36422/sctp snapshot.docker.ranco-dev-tools.eastus.cloudapp.azure.com:10001/gnbe2_sim:1.0.5
${stop_simu}  docker stop gnbe2_simu
${docker_Remove}    docker rm gnbe2_simu
