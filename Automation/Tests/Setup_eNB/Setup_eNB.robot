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
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Library     OperatingSystem
Library     ../Scripts/find_rmr_message.py
Library     ../Scripts/log_scripts.py
Library     REST        ${url}

*** Variables ***
${url}  ${e2mgr_address}

*** Test Cases ***
[Setup]
    Start Redis Monitor
    AND Prepare Enviorment

Send eNB setup request via e2adapter
    Send eNB Setup Request

Redis Monitor Logs - Verify Publish
    Redis Monitor Logs - Verify Publish To Connection Status Channel    ${enb_ran_name}    CONNECTED
    Redis Monitor Logs - Verify NOT Published To Manipulation Channel    ${enb_ran_name}    UPDATED

Get request eNB
    Sleep    2s
    Get Request nodeb    ${enb_ran_name}
    Integer  response status  200
    String   response body ranName    ${enb_ran_name}
    String   response body connectionStatus    CONNECTED
    String   response body nodeType     ENB
    String   response body enb enbType    SHORT_MACRO_ENB
    Boolean  response body setupFromNetwork    true

Prepare Logs For Tests
    Remove log files
    Save logs

[Teardown]
    Stop Redis Monitor





