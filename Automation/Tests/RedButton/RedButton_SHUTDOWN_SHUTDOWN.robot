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
Suite Setup   Prepare Enviorment
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Resource    red_button_keywords.robot
Resource   ../Resource/scripts_variables.robot
Library    ../Scripts/find_error_script.py
Library     OperatingSystem
Library    Collections
Library     REST      ${url}

*** Test Cases ***

Verify nodeb connection status is CONNECTED and it's associated to an e2t instance
   Verify connected and associated

Execute Shutdown
   Execute Shutdown

prepare logs for tests
    Remove log files
    Save logs

Verify e2mgr logs - Set and Publish Disconnect True
  ${result}    find_error_script.find_error     ${EXECDIR}  ${e2mgr_log_filename}    ${set_and_publish_disconnect}
   Should Be Equal As Strings    ${result}      True

Verify nodeb's connection status is SHUT_DOWN and it's NOT associated to an e2t instance
   Verify shutdown for gnb

Verify E2T instance has no associated RANs
   Verify E2T instance has no associated RANs

Execute second Shutdown
   Execute Shutdown

Verify again nodeb's connection status is SHUT_DOWN and it's NOT associated to an e2t instance
   Verify shutdown for gnb

Verify again E2T instance has no associated RANs
   Verify E2T instance has no associated RANs
