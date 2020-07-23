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
Variables  ../Scripts/variables.py
Suite Setup   Prepare Enviorment  ${True}
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Library     OperatingSystem
Library     REST        ${url}


*** Variables ***
${url}  ${e2mgr_address}

*** Test Cases ***

Prepare Redis Monitor Log
    Start Redis Monitor

Update gNB
    Sleep  2s
    Update Gnb request
    Integer  response status  200
    String   response body ranName    ${ranname}
    String   response body connectionStatus    CONNECTED
    String   response body nodeType     GNB
    String   response body gnb servedNrCells 0 servedNrCellInformation cellId   abcd
    String   response body gnb servedNrCells 0 nrNeighbourInfos 0 nrCgi  one
    String   response body gnb servedNrCells 0 servedNrCellInformation servedPlmns 0  whatever

prepare logs for tests
    Remove log files
    Save logs

E2M Logs - Verify Update
    ${result}    log_scripts.verify_log_message   ${EXECDIR}/${e2mgr_log_filename}  ${update_gnb_log_message}
    Should Be Equal As Strings    ${result}      True

Redis Monitor Logs - Verify Publish
    Redis Monitor Logs - Verify Publish To Manipulation Channel    ${ranName}    UPDATED

[Teardown]
    Stop Redis Monitor