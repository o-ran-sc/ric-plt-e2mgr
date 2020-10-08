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
Library     ../Scripts/log_scripts.py
Library     REST        ${url}


*** Variables ***
${url}  ${e2mgr_address}

*** Test Cases ***
[Setup]
    Prepare Enviorment  ${False}

Add eNB
    Sleep  2s
    Add eNb Request    ${add_enb_type_ng_request_body}
    Integer  response status  400

prepare logs for tests
    Remove log files
    Save logs