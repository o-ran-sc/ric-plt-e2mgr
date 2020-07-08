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
#
#   This source code is part of the near-RT RIC (RAN Intelligent Controller)
#   platform project (RICP).
#

*** Settings ***
Documentation    Resource file


*** Variables ***
${docker_number}    5
${docker_number-1}    4
${url}   http://localhost:3800
${ranName}  gnb_208_092_303030
${getNodeb}  /v1/nodeb/${ranName}
${set_general_configuration}   /v1/nodeb/parameters
${set_general_configuration_body}   {"enableRic":false}
${update_gnb_url}   /v1/nodeb/${ranName}/update
${update_gnb_body}  {"servedNrCells":[{"servedNrCellInformation":{"cellId":"abcd","choiceNrMode":{"fdd":{}},"nrMode":1,"nrPci":1,"servedPlmns":["whatever"]},"nrNeighbourInfos":[{"nrCgi":"one","choiceNrMode":{"fdd":{}},"nrMode":1,"nrPci":1}]}]}
${update_gnb_body_notvalid}  {"servedNrCells":[{"servedNrCellInformation":{"choiceNrMode":{"fdd":{}},"nrMode":1,"nrPci":1,"servedPlmns":["whatever"]},"nrNeighbourInfos":[{"nrCgi":"whatever","choiceNrMode":{"fdd":{}},"nrMode":1,"nrPci":1}]}]}
${E2tInstanceAddress}   10.0.2.15:38000
${header}  {"Content-Type": "application/json"}
${docker_command}  kubectl -n ricplt get pods | grep -E 'dbaas|e2mgr|rtmgr|gnbe2|e2term' | grep Running |wc --lines
${docker_Remove}    docker rm gnbe2_oran_simu
${stop_simu}  kubectl scale --replicas=0 deploy/oran-simulator-gnbe2-oran-simu -n=ricplt
${run_simu_regular}  kubectl scale --replicas=1 deploy/oran-simulator-gnbe2-oran-simu -n=ricplt
${restart_simu}  kubectl -n ricplt delete po $(${gnbe2_sim_pod})
${docker_restart}  kubectl -n ricplt delete po ${e2mgr_pod}
${start_e2}  kubectl scale --replicas=1 deploy/deployment-ricplt-e2term-alpha -n=ricplt
${stop_e2}      kubectl scale --replicas=0 deploy/deployment-ricplt-e2term-alpha -n=ricplt
${dbass_start}   docker run -d --name dbass -p 6379:6379 --env DBAAS_SERVICE_HOST=10.0.2.15  snapshot.docker.ranco-dev-tools.eastus.cloudapp.azure.com:10001/dbass:1.0.0
${dbass_remove}    docker rm dbass
${dbass_stop}      docker stop dbass
${stop_routing_manager}  kubectl scale --replicas=0 deploy/deployment-ricplt-rtmgr -n=ricplt
${start_routing_manager}  kubectl scale --replicas=1 deploy/deployment-ricplt-rtmgr -n=ricplt

${gnbe2_sim_pod}  kubectl -n ricplt get pods |grep gnbe2 | awk '{print $1}'
${dbaas_pod}  kubectl -n ricplt get pods |grep dbaas | awk '{print $1}'
${e2mgr_pod}  kubectl -n ricplt get pods |grep e2mgr | awk '{print $1}'
${e2term_pod}  kubectl -n ricplt get pods |grep e2term | awk '{print $1}'

#${gnbe2_sim_pod_logs}  kubectl -n ricplt logs -f ${gnbe2_sim_pod}
#${e2mgr_pod_logs}  kubectl -n ricplt logs -f ${e2mgr_pod}
#${e2term_pod_logs}  kubectl -n ricplt logs -f ${e2term_pod}


