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
Suite Setup   Prepare Enviorment    ${True}
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Resource    ../Resource/scripts_variables.robot
Library     OperatingSystem
Library     ../Scripts/find_rmr_message.py
Library     ../Scripts/cleanup_db.py
Library     ../Scripts/e2t_db_script.py

*** Test Cases ***

Test New E2T Send Init
    Stop E2
    Stop Routing Manager
    Flush And Populate DB  ${False}
    Start Routing Manager
    Start E2
    wait until keyword succeeds  1 min    10 sec    Validate Required Dockers

Prepare Logs For Tests
    Remove log files
    Save logs

E2M Logs - Verify RMR Message
    ${result}    find_rmr_message.verify_logs   ${EXECDIR}   ${e2mgr_log_filename}  ${E2_INIT_message_type}    ${None}
    Should Be Equal As Strings    ${result}      True

Verify E2T keys in DB
    ${result}=    e2t_db_script.verify_e2t_addresses_key
    Should Be Equal As Strings  ${result}    True

    ${result}=    e2t_db_script.verify_e2t_instance_key
    Should Be Equal As Strings  ${result}    True



