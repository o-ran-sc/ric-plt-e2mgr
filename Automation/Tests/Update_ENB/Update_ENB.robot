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
Library     ../Scripts/k8s_helper.py
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

Update eNb
    Sleep    2s
    Update eNb Request
    Integer  response status  200
    String   response body enb enbType    HOME_ENB

prepare logs for tests
    Remove log files
    Save logs

E2M Logs - Verify Update
    ${result}    log_scripts.verify_log_message   ${EXECDIR}/${e2mgr_log_filename}  ${update_enb_log_message}
    Should Be Equal As Strings    ${result}      True

Redis Monitor Logs - Verify Publish
    Redis Monitor Logs - Verify Publish To Manipulation Channel    ${enb_ran_name}    UPDATED

[Teardown]
    Stop Redis Monitor