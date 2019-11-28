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
Resource   ../Resource/Keywords.robot
Library     String
Library     OperatingSystem
Library     Process
Library     ../Scripts/find_rmr_message.py
Library     ../Scripts/find_error_script.py



*** Test Cases ***
Verify logs - Reset Sent by e2adapter
    ${result}    find_error_script.find_error  ${EXECDIR}  ${e2adapter_log_filename}  ${E2ADAPTER_Setup_Resp}
    Should Be Equal As Strings    ${result}      True

Verify logs - e2mgr logs - messege sent
    ${result}    find_rmr_message.verify_logs  ${EXECDIR}  ${e2mgr_log_filename}  ${RIC_X2_RESET_REQ_message_type}  ${Meid_test1}
    Should Be Equal As Strings    ${result}      True

Verify logs - e2mgr logs - messege received
    ${result}    find_rmr_message.verify_logs  ${EXECDIR}  ${e2mgr_log_filename}  ${RIC_X2_RESET_RESP_message_type}  ${Meid_test1}
    Should Be Equal As Strings    ${result}      True

RAN Restarted messege sent
    ${result}    find_rmr_message.verify_logs  ${EXECDIR}  ${e2mgr_log_filename}  ${RAN_RESTARTED_message_type}  ${Meid_test1}
    Should Be Equal As Strings    ${result}      True

RSM RESOURCE STATUS REQUEST message not sent
    ${result}    find_rmr_message.verify_logs     ${EXECDIR}    ${rsm_log_filename}  ${RIC_RES_STATUS_REQ_message_type_successfully_sent}    ${RAN_NAME_test2}
    Should Be Equal As Strings    ${result}      False