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
Resource    ../Resource/scripts_variables.robot
Library     String
Library     OperatingSystem
Library     Process
Resource   ../Resource/Keywords.robot
Library     ../Scripts/find_rmr_message.py
Library     ../Scripts/find_error_script.py
Test Teardown  Start Dbass with 4 dockers



*** Test Cases ***
Verify logs - Reset Sent by simulator
    ${Reset}=   Grep File  ./gnb.log  ResetRequest has been sent
    Should Be Equal     ${Reset}     gnbe2_simu: ResetRequest has been sent

Verify logs for restart received
    ${result}    find_rmr_message.verify_logs     ${EXECDIR}  ${e2mgr_log_filename}  ${RIC_X2_RESET_REQ_message_type}    ${Meid_test1}
    Should Be Equal As Strings    ${result}      True

Verify for error on retrying
    ${result}    find_error_script.find_error    ${EXECDIR}     ${e2mgr_log_filename}   ${failed_to_retrieve_nodeb_message}
    Should Be Equal As Strings    ${result}      True


*** Keywords ***
Start Dbass with 4 dockers
     Run And Return Rc And Output    ${dbass_remove}
     Run And Return Rc And Output    ${dbass_start}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    ${docker_number-1}
     Sleep  5s
