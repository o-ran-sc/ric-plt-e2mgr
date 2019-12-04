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
${docker_number}    7
${docker_number-1}    6
${ip_gnb_simu}	10.0.2.15
${ip_e2adapter}	10.0.2.15
${url}   http://localhost:3800
${url_rsm}   http://localhost:4800
${json_setup_rsm_tests}    {"ranIp": "10.0.2.15","ranPort": 36422,"ranName":"test1"}
${json}    {"ranIp": "10.0.2.15","ranPort": 5577,"ranName":"test1"}
${endcbadjson}    {"ranIp": "a","ranPort": 49999,"ranName":"test2"}
${endcjson}    {"ranIp": "10.0.2.15","ranPort": 49999,"ranName":"test2"}
${resetcausejson}   {"cause": "misc:not-enough-user-plane-processing-resources"}
${resetbadcausejson}   {"cause": "bla" }
${resetbad1causejson}   {"cause":  }
${resource_status_start_json}   {"enableResourceStatus":true}
${resource_status_stop_json}    {"enableResourceStatus":false}
${header}  {"Content-Type": "application/json"}
${docker_command}  docker ps | grep Up | wc --lines
${run_simu_load}   docker run -d --name gnbe2_simu --env gNBipv4=localhost  --env gNBport=36422  --env duration=600000000000 --env indicationReportRate=1000000000 --env indicationInsertRate=0 -p 5577:36422/sctp snapshot.docker.ranco-dev-tools.eastus.cloudapp.azure.com:10001/gnbe2_simu:1.0.6
${stop_e2e_simu}  docker stop e2e_simu
${stop_simu}  docker stop gnbe2_simu
${run_simu_regular}  docker run -d --name gnbe2_simu --env gNBipv4=localhost  --env gNBport=36422 --env duration=600000000000 --env indicationReportRate=0 --env indicationInsertRate=0 -p 5577:36422/sctp snapshot.docker.ranco-dev-tools.eastus.cloudapp.azure.com:10001/gnbe2_simu:1.0.6
${run_e2e_simu_regular}  docker run -d --name e2e_simu -p 36422:36422 --net host -it snapshot.docker.ranco-dev-tools.eastus.cloudapp.azure.com:10001/e2sim:1.4.0 sh -c "./build/e2sim 10.0.2.15 36422"
${docker_e2e_simu_remove}    docker rm e2e_simu
${docker_Remove}    docker rm gnbe2_simu
${docker_remove_e2e_simu}    docker rm e2e_simu
${docker_restart}   docker restart e2mgr
${restart_simu}  docker restart gnbe2_simu
${restart_e2e_simu}  docker restart e2e_simu
${restart_e2adapter}  docker restart e2adapter
${restart_rsm}  docker restart rsm
${start_e2}  docker start e2
${stop_docker_e2}      docker stop e2
${dbass_start}   docker run -d --name dbass -p 6379:6379 --env DBAAS_SERVICE_HOST=10.0.2.15  snapshot.docker.ranco-dev-tools.eastus.cloudapp.azure.com:10001/dbass:1.0.0
${dbass_remove}    docker rm dbass
${dbass_stop}      docker stop dbass
${start_e2}  docker start e2
${stop_docker_e2}      docker stop e2
${Run_Config}       docker exec gnbe2_simu pkill gnbe2_simu -INT
${403_reset_message}    "Activity X2_RESET rejected. RAN current state DISCONNECTED does not allow its execution "




