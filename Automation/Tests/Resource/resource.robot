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
${run_simu_load}   docker run -d --name gnbe2_simu --env gNBipv4=localhost  --env gNBport=36422  --env duration=600000000000 --env indicationReportRate=1000000000 --env indicationInsertRate=0 -p 5577:36422/sctp snapshot.docker.ranco-dev-tools.eastus.cloudapp.azure.com:10001/gnbe2_simu:1.0.6
${stop_simu}  docker stop gnbe2_simu
${run_simu_regular}  docker run -d --name gnbe2_simu --env gNBipv4=localhost  --env gNBport=36422 --env duration=600000000000 --env indicationReportRate=0 --env indicationInsertRate=0 -p 5577:36422/sctp snapshot.docker.ranco-dev-tools.eastus.cloudapp.azure.com:10001/gnbe2_simu:1.0.6
${docker_Remove}    docker rm gnbe2_simu
${docker_restart}   docker restart e2mgr
${restart_simu}  docker restart gnbe2_simu
${restart_e2adapter}  docker restart e2adapter
${start_e2}  docker start e2
${stop_docker_e2}      docker stop e2
${start_redis}   docker run -d --name redis -p 6379:6379 --env DBAAS_SERVICE_HOST=10.0.2.15  snapshot.docker.ranco-dev-tools.eastus.cloudapp.azure.com:10001/dbass:1.0.0
${redis_remove}    docker rm redis
${restart_simu}  docker restart gnbe2_simu
${start_e2}  docker start e2
${stop_docker_e2}      docker stop e2
${Run_Config}       docker exec gnbe2_simu pkill gnbe2_simu -INT
${Save_e2_log}      docker logs gnbe2_simu > gnb.log
${Save_e2mgr_log}   docker logs e2mgr > e2mgr.log




