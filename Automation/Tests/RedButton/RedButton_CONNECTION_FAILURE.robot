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
Library     OperatingSystem
Library    Collections
Library     REST      ${url}


*** Test Cases ***

Execute X2 setup twice and verify its connection status is CONNECTED_SETUP_FAILED
    Execute setup twice and verify connected setup failed

Execute Shutdown
    Execute Shutdown

Verify nodeb's connection status is SHUT_DOWN and it's NOT associated to an e2t instance
   Verify shutdown for enb

