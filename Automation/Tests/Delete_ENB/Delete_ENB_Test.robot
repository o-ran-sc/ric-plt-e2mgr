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
Suite Setup   Prepare Enviorment  ${False}
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Library     OperatingSystem
Library     ../Scripts/log_scripts.py
Library     ../Scripts/k8s_helper.py
Library     REST        ${url}


*** Variables ***
${url}  ${e2mgr_address}

*** Test Cases ***

Prepare Redis Monitor Log
    Start Redis Monitor

Add eNB
    Sleep  2s
    Add eNb Request

Delete eNb
    Sleep    2s
    Delete eNb Request
    Integer  response status  204


prepare logs for tests
    Remove log files
    Save logs

E2M Logs - Verify Deletion
    ${result}    log_scripts.verify_log_message   ${EXECDIR}/${e2mgr_log_filename}  ${delete_enb_log_message}
    Should Be Equal As Strings    ${result}      True

Redis Monitor Logs - Verify Publish
    Redis Monitor Logs - Verify Publish To Manipulation Channel    ${enb_ran_name}    DELETED

[Teardown]
    Stop Redis Monitor