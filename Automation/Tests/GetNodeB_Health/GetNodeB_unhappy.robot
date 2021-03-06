##############################################################################
#
#   Copyright (c) 2020 Samsung Electronics Co., Ltd. All Rights Reserved.
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
${url}                      ${e2mgr_address}
${empty_list_nodeb_body}    {}
${invalid_list_nodeb_body}    {"ranList" :["abcd"]}




*** Test Cases ***

Get nodeb health empty list
    Sleep    2s
    Get request nodeb health             request_body=${empty_list_nodeb_body}
    Integer  response status             404
    Integer  response body errorCode     406
    String   response body errorMessage  No RAN in Connected State

Get nodeb health invalid RAN
    Sleep    2s
    Get request nodeb health             request_body=${invalid_list_nodeb_body}
    Integer  response status             404
    Integer  response body errorCode     406
    String   response body errorMessage  No RAN in Connected State

Prepare Logs For Tests
    Remove log files
    Save logs

