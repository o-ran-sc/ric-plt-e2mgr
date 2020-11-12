##############################################################################
#
#   Copyright (c) 2019 AT&T Intellectual Property.
#   Copyright (c) 2020 Samsung Electronics Co., Ltd. All Rights Reserved.
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
${enb_ran_name}    enB_shortmacro_208__555540
${getNodeb}  /v1/nodeb
${nodeb_health_url}   /v1/nodeb/health
${empty_list_nodeb_body}   {}
${set_general_configuration}   /v1/nodeb/parameters
${set_general_configuration_body}   {"enableRic":false}
${update_gnb_url}   /v1/nodeb/gnb/${ranName}
${enb_url}    /v1/nodeb/enb
${nodeb_state_url}    /v1/nodeb/states/${enb_ran_name}
${update_gnb_body}  {"servedNrCells":[{"servedNrCellInformation":{"cellId":"abcd","choiceNrMode":{"fdd":{}},"nrMode":1,"nrPci":1,"servedPlmns":["whatever"]},"nrNeighbourInfos":[{"nrCgi":"one","choiceNrMode":{"fdd":{}},"nrMode":1,"nrPci":1}]}]}
${update_gnb_body_notvalid}  {"servedNrCells":[{"servedNrCellInformation":{"choiceNrMode":{"fdd":{}},"nrMode":1,"nrPci":1,"servedPlmns":["whatever"]},"nrNeighbourInfos":[{"nrCgi":"whatever","choiceNrMode":{"fdd":{}},"nrMode":1,"nrPci":1}]}]}
${add_enb_request_body}    {"ranName":"${enb_ran_name}","globalNbId":{"nbId":"abc","plmnId":"def"},"port":1234,"enb":{"enbType":3,"guGroupIds":["ghi"],"servedCells":[{"broadcastPlmns":["jkl"],"cellId":"mnop","choiceEutraMode":{"fdd":{"dlearFcn":1,"ulearFcn":1},"tdd":{"additionalSpecialSubframeExtensionInfo":{"additionalSpecialSubframePatternsExtension":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1},"additionalSpecialSubframeInfo":{"additionalSpecialSubframePatterns":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1},"earFcn":4,"specialSubframeInfo":{"specialSubframePatterns":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1}}},"eutraMode":1,"csgId":"string","mbmsServiceAreaIdentities":["sds"],"mbsfnSubframeInfos":[{"radioframeAllocationOffset":3,"subframeAllocation":"jhg"}],"multibandInfos":[4],"neighbourInfos":[{"earFcn":4,"ecgi":"klj","pci":5,"tac":"wew"}],"pci":2,"prachConfiguration":{"highSpeedFlag":true,"prachConfigurationIndex":5,"prachFrequencyOffset":6,"rootSequenceIndex":7,"zeroCorrelationZoneConfiguration":6},"tac":"asd","additionalCellInformation":{"cellLatitude":1,"cellLongitude":1,"antennaHeight":1,"antennaAzimuthDirection":2,"antennaTiltAngle":3,"antennaMaxTransmit":4,"antennaMaxGain":5,"sectorId":6}},{"broadcastPlmns":["jkl"],"cellId":"qrst","choiceEutraMode":{"fdd":{"dlearFcn":4,"ulearFcn":2},"tdd":{"additionalSpecialSubframeExtensionInfo":{"additionalSpecialSubframePatternsExtension":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1},"additionalSpecialSubframeInfo":{"additionalSpecialSubframePatterns":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1},"earFcn":4,"specialSubframeInfo":{"specialSubframePatterns":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1}}},"eutraMode":1,"csgId":"string","mbmsServiceAreaIdentities":["sds"],"mbsfnSubframeInfos":[{"radioframeAllocationOffset":5,"subframeAllocation":"jhg"}],"multibandInfos":[4],"neighbourInfos":[{"earFcn":2,"ecgi":"klj","pci":4,"tac":"wew"}],"pci":3,"prachConfiguration":{"highSpeedFlag":true,"prachConfigurationIndex":4,"prachFrequencyOffset":3,"rootSequenceIndex":3,"zeroCorrelationZoneConfiguration":2},"tac":"asd","additionalCellInformation":{"cellLatitude":3,"cellLongitude":3,"antennaHeight":3,"antennaAzimuthDirection":3,"antennaTiltAngle":4,"antennaMaxTransmit":4,"antennaMaxGain":5,"sectorId":5}}]}}
${add_enb_type_ng_request_body}    {"ranName":"${enb_ran_name}","globalNbId":{"nbId":"abc","plmnId":"def"},"port":1234,"enb":{"enbType":5,"guGroupIds":["ghi"],"servedCells":[{"broadcastPlmns":["jkl"],"cellId":"mnop","choiceEutraMode":{"fdd":{"dlearFcn":1,"ulearFcn":1},"tdd":{"additionalSpecialSubframeExtensionInfo":{"additionalSpecialSubframePatternsExtension":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1},"additionalSpecialSubframeInfo":{"additionalSpecialSubframePatterns":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1},"earFcn":4,"specialSubframeInfo":{"specialSubframePatterns":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1}}},"eutraMode":1,"csgId":"string","mbmsServiceAreaIdentities":["sds"],"mbsfnSubframeInfos":[{"radioframeAllocationOffset":3,"subframeAllocation":"jhg"}],"multibandInfos":[4],"neighbourInfos":[{"earFcn":4,"ecgi":"klj","pci":5,"tac":"wew"}],"pci":2,"prachConfiguration":{"highSpeedFlag":true,"prachConfigurationIndex":5,"prachFrequencyOffset":6,"rootSequenceIndex":7,"zeroCorrelationZoneConfiguration":6},"tac":"asd","additionalCellInformation":{"cellLatitude":1,"cellLongitude":1,"antennaHeight":1,"antennaAzimuthDirection":2,"antennaTiltAngle":3,"antennaMaxTransmit":4,"antennaMaxGain":5,"sectorId":6}},{"broadcastPlmns":["jkl"],"cellId":"qrst","choiceEutraMode":{"fdd":{"dlearFcn":4,"ulearFcn":2},"tdd":{"additionalSpecialSubframeExtensionInfo":{"additionalSpecialSubframePatternsExtension":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1},"additionalSpecialSubframeInfo":{"additionalSpecialSubframePatterns":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1},"earFcn":4,"specialSubframeInfo":{"specialSubframePatterns":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1}}},"eutraMode":1,"csgId":"string","mbmsServiceAreaIdentities":["sds"],"mbsfnSubframeInfos":[{"radioframeAllocationOffset":5,"subframeAllocation":"jhg"}],"multibandInfos":[4],"neighbourInfos":[{"earFcn":2,"ecgi":"klj","pci":4,"tac":"wew"}],"pci":3,"prachConfiguration":{"highSpeedFlag":true,"prachConfigurationIndex":4,"prachFrequencyOffset":3,"rootSequenceIndex":3,"zeroCorrelationZoneConfiguration":2},"tac":"asd","additionalCellInformation":{"cellLatitude":3,"cellLongitude":3,"antennaHeight":3,"antennaAzimuthDirection":3,"antennaTiltAngle":4,"antennaMaxTransmit":4,"antennaMaxGain":5,"sectorId":5}}]}}
${add_enb_response_body}    {"ranName":"${enb_ran_name}","port":1234,"connectionStatus":"DISCONNECTED","globalNbId":{"plmnId":"def","nbId":"abc"},"nodeType":"ENB","enb":{"enbType":"SHORT_MACRO_ENB","servedCells":[{"pci":2,"cellId":"mnop","tac":"asd","broadcastPlmns":["jkl"],"choiceEutraMode":{"fdd":{"ulearFcn":1,"dlearFcn":1},"tdd":{"earFcn":4,"specialSubframeInfo":{"specialSubframePatterns":"SSP0","cyclicPrefixDl":"NORMAL","cyclicPrefixUl":"NORMAL"},"additionalSpecialSubframeInfo":{"additionalSpecialSubframePatterns":"SSP0","cyclicPrefixDl":"NORMAL","cyclicPrefixUl":"NORMAL"},"additionalSpecialSubframeExtensionInfo":{"additionalSpecialSubframePatternsExtension":"SSP10","cyclicPrefixDl":"NORMAL","cyclicPrefixUl":"NORMAL"}}},"eutraMode":"FDD","prachConfiguration":{"rootSequenceIndex":7,"zeroCorrelationZoneConfiguration":6,"highSpeedFlag":true,"prachFrequencyOffset":6,"prachConfigurationIndex":5},"mbsfnSubframeInfos":[{"radioframeAllocationOffset":3,"subframeAllocation":"jhg"}],"csgId":"string","mbmsServiceAreaIdentities":["sds"],"multibandInfos":[4],"neighbourInfos":[{"ecgi":"klj","pci":5,"earFcn":4,"tac":"wew"}],"additionalCellInformation":{"cellLatitude":1,"cellLongitude":1,"antennaHeight":1,"antennaAzimuthDirection":2,"antennaTiltAngle":3,"antennaMaxTransmit":4,"antennaMaxGain":5,"sectorId":6}},{"pci":3,"cellId":"qrst","tac":"asd","broadcastPlmns":["jkl"],"choiceEutraMode":{"fdd":{"ulearFcn":2,"dlearFcn":4},"tdd":{"earFcn":4,"specialSubframeInfo":{"specialSubframePatterns":"SSP0","cyclicPrefixDl":"NORMAL","cyclicPrefixUl":"NORMAL"},"additionalSpecialSubframeInfo":{"additionalSpecialSubframePatterns":"SSP0","cyclicPrefixDl":"NORMAL","cyclicPrefixUl":"NORMAL"},"additionalSpecialSubframeExtensionInfo":{"additionalSpecialSubframePatternsExtension":"SSP10","cyclicPrefixDl":"NORMAL","cyclicPrefixUl":"NORMAL"}}},"eutraMode":"FDD","prachConfiguration":{"rootSequenceIndex":3,"zeroCorrelationZoneConfiguration":2,"highSpeedFlag":true,"prachFrequencyOffset":3,"prachConfigurationIndex":4},"mbsfnSubframeInfos":[{"radioframeAllocationOffset":5,"subframeAllocation":"jhg"}],"csgId":"string","mbmsServiceAreaIdentities":["sds"],"multibandInfos":[4],"neighbourInfos":[{"ecgi":"klj","pci":4,"earFcn":2,"tac":"wew"}],"additionalCellInformation":{"cellLatitude":3,"cellLongitude":3,"antennaHeight":3,"antennaAzimuthDirection":3,"antennaTiltAngle":4,"antennaMaxTransmit":4,"antennaMaxGain":5,"sectorId":5}}],"guGroupIds":["ghi"]}}
${update_enb_request_body}    {"enb":{"enbType":"HOME_ENB","guGroupIds":["ghi"],"servedCells":[{"broadcastPlmns":["jkl"],"cellId":"mnop","choiceEutraMode":{"fdd":{"dlearFcn":1,"ulearFcn":1},"tdd":{"additionalSpecialSubframeExtensionInfo":{"additionalSpecialSubframePatternsExtension":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1},"additionalSpecialSubframeInfo":{"additionalSpecialSubframePatterns":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1},"earFcn":4,"specialSubframeInfo":{"specialSubframePatterns":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1}}},"eutraMode":1,"csgId":"string","mbmsServiceAreaIdentities":["sds"],"mbsfnSubframeInfos":[{"radioframeAllocationOffset":3,"subframeAllocation":"jhg"}],"multibandInfos":[4],"neighbourInfos":[{"earFcn":4,"ecgi":"klj","pci":5,"tac":"wew"}],"pci":2,"prachConfiguration":{"highSpeedFlag":true,"prachConfigurationIndex":5,"prachFrequencyOffset":6,"rootSequenceIndex":7,"zeroCorrelationZoneConfiguration":6},"tac":"asd","additionalCellInformation":{"cellLatitude":1,"cellLongitude":1,"antennaHeight":1,"antennaAzimuthDirection":2,"antennaTiltAngle":3,"antennaMaxTransmit":4,"antennaMaxGain":5,"sectorId":6}},{"broadcastPlmns":["jkl"],"cellId":"qrst","choiceEutraMode":{"fdd":{"dlearFcn":4,"ulearFcn":2},"tdd":{"additionalSpecialSubframeExtensionInfo":{"additionalSpecialSubframePatternsExtension":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1},"additionalSpecialSubframeInfo":{"additionalSpecialSubframePatterns":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1},"earFcn":4,"specialSubframeInfo":{"specialSubframePatterns":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1}}},"eutraMode":1,"csgId":"string","mbmsServiceAreaIdentities":["sds"],"mbsfnSubframeInfos":[{"radioframeAllocationOffset":5,"subframeAllocation":"jhg"}],"multibandInfos":[4],"neighbourInfos":[{"earFcn":2,"ecgi":"klj","pci":4,"tac":"wew"}],"pci":3,"prachConfiguration":{"highSpeedFlag":true,"prachConfigurationIndex":4,"prachFrequencyOffset":3,"rootSequenceIndex":3,"zeroCorrelationZoneConfiguration":2},"tac":"asd","additionalCellInformation":{"cellLatitude":3,"cellLongitude":3,"antennaHeight":3,"antennaAzimuthDirection":3,"antennaTiltAngle":4,"antennaMaxTransmit":4,"antennaMaxGain":5,"sectorId":5}}]}}
${update_enb_type_ng_request_body}    {"enb":{"enbType":"SHORT_MACRO_NG_ENB","guGroupIds":["ghi"],"servedCells":[{"broadcastPlmns":["jkl"],"cellId":"mnop","choiceEutraMode":{"fdd":{"dlearFcn":1,"ulearFcn":1},"tdd":{"additionalSpecialSubframeExtensionInfo":{"additionalSpecialSubframePatternsExtension":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1},"additionalSpecialSubframeInfo":{"additionalSpecialSubframePatterns":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1},"earFcn":4,"specialSubframeInfo":{"specialSubframePatterns":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1}}},"eutraMode":1,"csgId":"string","mbmsServiceAreaIdentities":["sds"],"mbsfnSubframeInfos":[{"radioframeAllocationOffset":3,"subframeAllocation":"jhg"}],"multibandInfos":[4],"neighbourInfos":[{"earFcn":4,"ecgi":"klj","pci":5,"tac":"wew"}],"pci":2,"prachConfiguration":{"highSpeedFlag":true,"prachConfigurationIndex":5,"prachFrequencyOffset":6,"rootSequenceIndex":7,"zeroCorrelationZoneConfiguration":6},"tac":"asd","additionalCellInformation":{"cellLatitude":1,"cellLongitude":1,"antennaHeight":1,"antennaAzimuthDirection":2,"antennaTiltAngle":3,"antennaMaxTransmit":4,"antennaMaxGain":5,"sectorId":6}},{"broadcastPlmns":["jkl"],"cellId":"qrst","choiceEutraMode":{"fdd":{"dlearFcn":4,"ulearFcn":2},"tdd":{"additionalSpecialSubframeExtensionInfo":{"additionalSpecialSubframePatternsExtension":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1},"additionalSpecialSubframeInfo":{"additionalSpecialSubframePatterns":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1},"earFcn":4,"specialSubframeInfo":{"specialSubframePatterns":1,"cyclicPrefixDl":1,"cyclicPrefixUl":1}}},"eutraMode":1,"csgId":"string","mbmsServiceAreaIdentities":["sds"],"mbsfnSubframeInfos":[{"radioframeAllocationOffset":5,"subframeAllocation":"jhg"}],"multibandInfos":[4],"neighbourInfos":[{"earFcn":2,"ecgi":"klj","pci":4,"tac":"wew"}],"pci":3,"prachConfiguration":{"highSpeedFlag":true,"prachConfigurationIndex":4,"prachFrequencyOffset":3,"rootSequenceIndex":3,"zeroCorrelationZoneConfiguration":2},"tac":"asd","additionalCellInformation":{"cellLatitude":3,"cellLongitude":3,"antennaHeight":3,"antennaAzimuthDirection":3,"antennaTiltAngle":4,"antennaMaxTransmit":4,"antennaMaxGain":5,"sectorId":5}}]}}
${header}  {"Content-Type": "application/json"}

#K8S
${pods_number}    6
${pods_number-1}    5
${verify_all_pods_are_ready_command}  kubectl -n ricplt get pods | grep -E 'dbaas|e2mgr|rtmgr|gnbe2|e2term|e2adapter' | grep Running | grep 1/1 |wc --lines
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
${stop_e2adapter}  kubectl scale --replicas=0 deploy/e2adapter -n=ricplt
${start_e2adapter}  kubectl scale --replicas=1 deploy/e2adapter -n=ricplt
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
${delete_enb_log_message}    RAN name: ${enb_ran_name} - deleted successfully
${update_enb_log_message}    RAN name: ${enb_ran_name} - Successfully updated eNB
${update_gnb_log_message}    RAN name: ${ranName} - Successfully updated gNB cells
