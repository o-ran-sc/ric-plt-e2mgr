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
Suite Setup   Prepare Enviorment    ${True}
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Library    ../Scripts/find_error_script.py
Library     ../Scripts/e2mdbscripts.py
Library     ../Scripts/log_scripts.py
Library     OperatingSystem
Library    Collections
Library     REST      ${url}

*** Variables ***
${url}  ${e2mgr_address}

*** Test Cases ***

Prepare Redis Monitor Log
    Start Redis Monitor

Setup Ran and verify it's CONNECTED and associated
    Get Request node b gnb
    Integer  response status  200
    String   response body ranName    ${ranname}
    String   response body connectionStatus    CONNECTED
    String   response body associatedE2tInstanceAddress  ${e2t_alpha_address}

Stop simulator
   Stop Simulator

Verify connection status is DISCONNECTED and RAN is not associated with E2T instance
    Sleep    30s
    GET      ${getNodeb}
    Integer  response status  200
    String   response body ranName    ${ranname}
    Missing  response body associatedE2tInstanceAddress
    String   response body connectionStatus    DISCONNECTED

prepare logs for tests
    Remove log files
    Save logs

Verify E2T instance is NOT associated with RAN
   ${result}    e2mdbscripts.verify_ran_is_associated_with_e2t_instance     ${ranname}  ${e2t_alpha_address}
   Should Be True    ${result} == False

Redis Monitor Logs - Verify Publish
    Redis Monitor Logs - Verify Publish To Connection Status Channel   ${ran_name}    DISCONNECTED

[Teardown]    Run Keywords
              Start Simulator
              AND wait until keyword succeeds  1 min    10 sec    Validate Required Dockers
              AND Stop Redis Monitor