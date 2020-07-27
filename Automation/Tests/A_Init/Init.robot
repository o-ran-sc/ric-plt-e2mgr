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
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot

*** Test Cases ***

Init Environment And Validate
    Stop All Pods Except Simulator
    Restart simulator
    Wait until keyword succeeds  2 min    10 sec    Validate Required Dockers    1

    Start E2 Manager
    Start Dbass
    Wait until keyword succeeds  2 min    10 sec    Validate Required Dockers    3

    Start Routing Manager
    Wait until keyword succeeds  2 min    10 sec    Validate Required Dockers    4

    Start E2
    Wait until keyword succeeds  2 min    10 sec    Validate Required Dockers





