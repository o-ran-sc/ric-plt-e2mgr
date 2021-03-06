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
Library     ../Scripts/log_scripts.py
Library     REST        ${url}


*** Variables ***
${url}  ${e2mgr_address}

*** Test Cases ***
[Setup]
    Start Redis Monitor
    Prepare Enviorment  ${False}

Add eNB
    Sleep  2s
    Add eNb Request
    Integer  response status  201
    String   response body ranName    ${enb_ran_name}
    String   response body connectionStatus    DISCONNECTED
    String   response body nodeType     ENB
    String   response body enb enbType    SHORT_MACRO_ENB
    Missing  response body setupFromNetwork

prepare logs for tests
    Remove log files
    Save logs

Redis Monitor Logs - Verify Publish
    Redis Monitor Logs - Verify Publish To Manipulation Channel    ${enb_ran_name}    ADDED

[Teardown]
    Stop Redis Monitor









