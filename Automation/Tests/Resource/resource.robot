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

#REST
${ranName}  gnb_208_092_303030
${getNodeb}  /v1/nodeb/${ranName}
${set_general_configuration}   /v1/nodeb/parameters
${set_general_configuration_body}   {"enableRic":false}
${update_gnb_url}   /v1/nodeb/${ranName}/update
${update_gnb_body}  {"servedNrCells":[{"servedNrCellInformation":{"cellId":"abcd","choiceNrMode":{"fdd":{}},"nrMode":1,"nrPci":1,"servedPlmns":["whatever"]},"nrNeighbourInfos":[{"nrCgi":"one","choiceNrMode":{"fdd":{}},"nrMode":1,"nrPci":1}]}]}
${update_gnb_body_notvalid}  {"servedNrCells":[{"servedNrCellInformation":{"choiceNrMode":{"fdd":{}},"nrMode":1,"nrPci":1,"servedPlmns":["whatever"]},"nrNeighbourInfos":[{"nrCgi":"whatever","choiceNrMode":{"fdd":{}},"nrMode":1,"nrPci":1}]}]}
${header}  {"Content-Type": "application/json"}

#K8S
${pods_number}    5
${pods_number-1}    4
${verify_all_pods_are_ready_command}  kubectl -n ricplt get pods | grep -E 'dbaas|e2mgr|rtmgr|gnbe2|e2term' | grep Running | grep 1/1 |wc --lines
${stop_simu}  kubectl scale --replicas=0 deploy/oran-simulator-gnbe2-oran-simu -n=ricplt
${start_simu}  kubectl scale --replicas=1 deploy/oran-simulator-gnbe2-oran-simu -n=ricplt
${start_e2mgr}  kubectl scale --replicas=1 deploy/deployment-ricplt-e2mgr -n=ricplt
${stop_e2mgr}      kubectl scale --replicas=0 deploy/deployment-ricplt-e2mgr -n=ricplt
${start_e2}  kubectl scale --replicas=1 deploy/deployment-ricplt-e2term-alpha -n=ricplt
${stop_e2}      kubectl scale --replicas=0 deploy/deployment-ricplt-e2term-alpha -n=ricplt
${dbass_start}   kubectl -n ricplt scale statefulsets statefulset-ricplt-dbaas-server --replicas=1
${dbass_stop}      kubectl -n ricplt scale statefulsets statefulset-ricplt-dbaas-server --replicas=0
${stop_routing_manager}  kubectl scale --replicas=0 deploy/deployment-ricplt-rtmgr -n=ricplt
${start_routing_manager}  kubectl scale --replicas=1 deploy/deployment-ricplt-rtmgr -n=ricplt
${gnbe2_sim_pod}  kubectl -n ricplt get pods |grep gnbe2 | awk '{print $1}'
${e2mgr_pod}  kubectl -n ricplt get pods |grep e2mgr | awk '{print $1}'
${e2term_pod}  kubectl -n ricplt get pods |grep e2term | awk '{print $1}'
${rtmgr_pod}  kubectl -n ricplt get pods |grep rtmgr | awk '{print $1}'


#Logs
${E2_INIT_message_type}    MType: 1100
${Setup_failure_message_type}    MType: 12003
${first_retry_to_retrieve_from_db}      RnibDataService.retry - retrying 1 GetNodeb
${third_retry_to_retrieve_from_db}      RnibDataService.retry - after 3 attempts of GetNodeb
${RIC_RES_STATUS_REQ_message_type_successfully_sent}     Message type: 10090 - Successfully sent RMR message
${E2_TERM_KEEP_ALIVE_REQ_message_type_successfully_sent}     Message type: 1101 - Successfully sent RMR message
${save_general_configuration}      SetGeneralConfigurationHandler.Handle - save general configuration to rnib: {EnableRic:false}
${set_and_publish_disconnect}      RnibDataService.UpdateNodebInfoOnConnectionStatusInversion - stateChangeMessageChannel: RAN_CONNECTION_STATUS_CHANGE, event: gnb_208_092_303030_DISCONNECTED

