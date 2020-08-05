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
    Prepare Enviorment  ${True}


Restart Simulator
    Restart simulator
    wait until keyword succeeds  2 min    10 sec    Validate Required Dockers

Prepare Logs For Tests
    Remove log files
    Save logs

Redis Monitor Logs - Verify Publish
    Redis Monitor Logs - Verify Publish To Connection Status Channel   ${ranName}    CONNECTED
    Redis Monitor Logs - Verify NOT Published To Manipulation Channel    ${ranName}    UPDATED


[Teardown]
    Stop Redis Monitor





