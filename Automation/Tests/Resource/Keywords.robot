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
Documentation   Keywords file
Library     ../Scripts/cleanup_db.py
Resource   ../Resource/resource.robot
Library     OperatingSystem

*** Keywords ***
Get Request node b gnb
    Sleep    1s
    GET      ${getNodeb}

Update Ran request
    Sleep  1s
    PUT    ${update_gnb_url}   ${update_gnb_body}

Set General Configuration request
    Sleep  1s
    PUT    ${set_general_configuration}   ${set_general_configuration_body}

Update Ran request not valid
    Sleep  1s
    PUT    ${update_gnb_url}   ${update_gnb_body_notvalid}

Remove log files
    Remove File  ${EXECDIR}/${gnb_log_filename}
    Remove File  ${EXECDIR}/${e2mgr_log_filename}
    Remove File  ${EXECDIR}/${e2t_log_filename}

Save logs
    Sleep   1s
    Run     ${Save_sim_log}
    Run     ${Save_e2mgr_log}
    Run     ${Save_e2t_log}

Stop Simulator
    Run And Return Rc And Output    ${stop_simu}

Prepare Enviorment
     [Arguments]    ${prepare_logs_only}=${False}
     prepare enviorment step 1 - init logs
     run keyword if  ${prepare_logs_only}==${False}    Prepare Enviorment Step 2
     wait until keyword succeeds  1 min    10 sec    Validate Required Dockers

Prepare Enviorment Step 1 - Init logs
    ${starting_timestamp}    Evaluate   datetime.datetime.now(datetime.timezone.utc).isoformat("T")   modules=datetime
    ${e2t_log_filename}      Evaluate      "e2t.${SUITE NAME}.log".replace(" ","-")
    ${e2mgr_log_filename}    Evaluate      "e2mgr.${SUITE NAME}.log".replace(" ","-")
    ${gnb_log_filename}      Evaluate      "gnb.${SUITE NAME}.log".replace(" ","-")
    ${Save_sim_log}          Evaluate  "kubectl -n ricplt logs --since-time=${starting_timestamp} $(${gnbe2_sim_pod}) > ${gnb_log_filename}"
    ${Save_e2mgr_log}        Evaluate   "kubectl -n ricplt logs --since-time=${starting_timestamp} $(${e2mgr_pod}) > ${e2mgr_log_filename}"
    ${Save_e2t_log}          Evaluate   "kubectl -n ricplt logs --since-time=${starting_timestamp} $(${e2term_pod}) > ${e2t_log_filename}"

    Set Suite Variable  ${e2t_log_filename}
    Set Suite Variable  ${e2mgr_log_filename}
    Set Suite Variable  ${gnb_log_filename}
    Set Suite Variable  ${Save_sim_log}
    Set Suite Variable  ${Save_e2mgr_log}
    Set Suite Variable  ${Save_e2t_log}

Prepare Enviorment Step 2
    Flush And Populate DB
    Restart simulator

Validate Required Dockers
    [Arguments]    ${required_number_of_dockers}=${docker_number}
    Log To Console  Validating all required dockers are up
    ${result}=  Run And Return Rc And Output     ${docker_command}
    Should Be Equal As Integers    ${result[1]}    ${required_number_of_dockers}

Start E2
     Log to Console  Starting E2Term
     Run And Return Rc And Output    ${start_e2}
     sleep  5s

Stop E2
     Log to Console  Stopping E2Term
     Run And Return Rc And Output    ${stop_e2}
     Sleep  2s
     Validate Required Dockers  ${docker_number-1}

Start Dbass
     Log to Console  Starting redis
     Run And Return Rc And Output    ${dbass_remove}
     Run And Return Rc And Output    ${dbass_start}
     Validate Required Dockers

Stop Dbass
     Log to Console  Stopping redis
     Run And Return Rc And Output    ${dbass_stop}
     Validate Required Dockers  required_number_of_dockers = ${docker_number-1}

Restart simulator
    Log to Console  Restarting gnbe2 simulator
    Run And Return Rc And Output    ${restart_simu}

Start Routing Manager
    Log to Console  Starting routing manager
    Run And Return Rc And Output    ${start_routing_manager}
    sleep  5s

Stop Routing Manager
    Log to Console  Stopping routing manager
    Run And Return Rc And Output    ${stop_routing_manager}
    Sleep  2s

Restart simulator with less docker
    Log to Console  Restarting gnbe2 simulator
    Run And Return Rc And Output    ${restart_simu}
    Validate Required Dockers  required_number_of_dockers = ${docker_number-1}

Flush And Populate DB
    [Arguments]     ${insert_e2_details}=True
    Log To Console  Flushing and populating DB
    ${flush}  cleanup_db.flush   ${insert_e2_details}
    Sleep  2s
    Should Be Equal As Strings  ${flush}  True



#To be deleted:
Start RoutingManager Simulator
    Run And Return Rc And Output    ${start_routing_manager}

Stop RoutingManager Simulator
    Run And Return Rc And Output    ${stop_routing_manager}
