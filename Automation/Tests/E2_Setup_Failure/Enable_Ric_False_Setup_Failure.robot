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
Suite Setup   Prepare Enviorment
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Library     OperatingSystem
Library     ../Scripts/find_rmr_message.py
Library     ../Scripts/e2mdbscripts.py
Library     REST        ${url}
Suite Teardown  Flush And Populate DB

*** Variables ***
${url}  ${e2mgr_address}

*** Test Cases ***

Disable ric and restart simulator
    ${result}    e2mdbscripts.set_enable_ric_false
    Restart simulator
    wait until keyword succeeds  2 min    10 sec    Validate Required Dockers

prepare logs for tests
    Remove log files
    Save logs

E2M Logs - Verify RMR Message
    ${result}    find_rmr_message.verify_logs   ${EXECDIR}   ${e2mgr_log_filename}  ${Setup_failure_message_type}    ${None}
    Should Be Equal As Strings    ${result}      True