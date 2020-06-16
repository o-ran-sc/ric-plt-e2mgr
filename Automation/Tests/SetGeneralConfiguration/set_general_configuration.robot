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



*** Settings ***
Suite Setup   Prepare Enviorment
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Library     OperatingSystem
Library     REST        ${url}




*** Test Cases ***

prepare logs for tests
    Remove log files
    Save logs

Set General Configuration
    Sleep  2s
    Set General Configuration request
    Integer  response status  200
    String   response body enableRic    false

Verify e2mgr logs - Third retry to retrieve from db
   ${result}    find_error_script.find_error     ${EXECDIR}  ${e2mgr_log_filename}   ${save_general_configuration}
   Should Be Equal As Strings    ${result}      True