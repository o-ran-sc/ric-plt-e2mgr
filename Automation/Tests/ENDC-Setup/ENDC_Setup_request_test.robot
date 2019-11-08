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
Suite Setup   Prepare Enviorment
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Resource    ../Resource/scripts_variables.robot
Library     OperatingSystem
Library     ../Scripts/find_rmr_message.py
Library     ../Scripts/rsmscripts.py
Library     REST        ${url}




*** Test Cases ***
Run Endc setup request
    Post Request setup node b endc-setup
    Integer     response status       204

Get request gnb
    Sleep    1s
    Get Request node b enb test2
    Integer  response status  200
    String   response body ranName    test2
    String   response body ip    ${ip_e2adapter} 
    String   response body connectionStatus    CONNECTED
    Integer  response body port     49999
    String   response body nodeType     GNB
    String   response body globalNbId plmnId    42f490
    String   response body globalNbId nbId    000004
    String   response body gnb servedNrCells 0 servedNrCellInformation cellId  42f490:000007fff0
    String   response body gnb servedNrCells 0 servedNrCellInformation configuredStac  0000
    String   response body gnb servedNrCells 0 servedNrCellInformation servedPlmns 0   "42f490"
    String   response body gnb servedNrCells 0 servedNrCellInformation nrMode  TDD
    String   response body gnb servedNrCells 0 servedNrCellInformation choiceNrMode tdd nrFreqInfo nrArFcn   650056
    Integer  response body gnb servedNrCells 0 servedNrCellInformation choiceNrMode tdd nrFreqInfo frequencyBands 0 nrFrequencyBand   78
    String   response body gnb servedNrCells 0 servedNrCellInformation choiceNrMode tdd transmissionBandwidth nrscs   SCS30
    String  response body gnb servedNrCells 0 servedNrCellInformation choiceNrMode tdd transmissionBandwidth ncnrb   NRB162

prepare logs for tests
    Remove log files
    Save logs


ENDC - RAN Connected message going to be sent
    ${result}    find_rmr_message.verify_logs   ${EXECDIR}   ${e2mgr_log_filename}  ${RAN_CONNECTED_message_type}    ${Meid_test2}
    Should Be Equal As Strings    ${result}      True

RSM RESOURCE STATUS REQUEST message not sent
    ${result}    find_rmr_message.verify_logs     ${EXECDIR}    ${rsm_log_filename}  ${RIC_RES_STATUS_REQ_message_type_successfully_sent}    ${RAN_NAME_test2}
    Should Be Equal As Strings    ${result}      False

Verify RSM RAN info doesn't exist in redis
   ${result}=   rsmscripts.verify_rsm_ran_info_start_false
   Should Be Equal As Strings  ${result}    False

