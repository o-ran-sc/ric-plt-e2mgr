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
Library     ../Scripts/k8s_helper.py
Resource   ../Resource/resource.robot
Library     OperatingSystem
Library     Process
Variables  ../Scripts/variables.py

*** Keywords ***
Get Request nodeb
    [Arguments]    ${nodeb_name}=${ranName}
    Sleep    1s
    GET      ${getNodeb}/${nodeb_name}

Update Gnb request
    Sleep  1s
    PUT    ${update_gnb_url}   ${update_gnb_body}

Add eNb Request
    Sleep  1s
    POST    ${enb_url}   ${add_enb_request_body}

Delete eNb Request
    Sleep  1s
    DELETE    ${enb_url}/${enb_ran_name}

Update eNb Request
    Sleep  1s
    PUT    ${enb_url}/${enb_ran_name}   ${update_enb_request_body}

Set General Configuration request
    Sleep  1s
    PUT    ${set_general_configuration}   ${set_general_configuration_body}

Update Gnb request not valid
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

Prepare Enviorment
     [Arguments]     ${need_to_restart_pods}=${False}     ${set_new_timestamp}=${True}
     Init logs
     Flush And Populate DB    ${set_new_timestamp}
     Run keyword if  ${need_to_restart_pods}==${True}   Restart RM and GNB Simulator
     Wait until keyword succeeds  2 min    10 sec    Validate Required Dockers

Restart RM and GNB Simulator
    Restart routing manager
    Wait until keyword succeeds  2 min    10 sec    Validate Required Dockers
    Restart simulator


Init logs
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

Validate Required Dockers
    [Arguments]    ${required_number_of_dockers}=${pods_number}
    Log To Console  Validating all required dockers are up
    ${result}=  Run And Return Rc And Output     ${verify_all_pods_are_ready_command}
    Should Be Equal As Integers    ${result[1]}    ${required_number_of_dockers}

Start E2
     Log to Console  Starting E2Term
     Run And Return Rc And Output    ${start_e2}
     Sleep  5s

Stop E2
     Log to Console  Stopping E2Term
     Run And Return Rc And Output    ${stop_e2}
     Sleep  5s

Start E2 Manager
     Log to Console  Starting E2Mgr
     Run And Return Rc And Output    ${start_e2mgr}
     Sleep  5s

Stop E2 Manager
     Log to Console  Stopping E2Mgr
     Run And Return Rc And Output    ${stop_e2mgr}
     Sleep  5s

Start Dbass
     Log to Console  Starting redis
     Run And Return Rc And Output    ${dbass_start}
     Sleep  5s

Stop Dbass
     Log to Console  Stopping redis
     Run And Return Rc And Output    ${dbass_stop}
     Sleep  5s

Stop Simulator
    log to console  Stopping gnbe2 simulator
    Run And Return Rc And Output    ${stop_simu}
    Sleep  90s

Start Simulator
    log to console  Starting gnbe2 simulator
    Run And Return Rc And Output    ${start_simu}

Restart simulator
   Log to Console  Restarting gnbe2 simulator
   Stop Simulator
   Start Simulator

Start Routing Manager
    Log to Console  Starting routing manager
    Run And Return Rc And Output    ${start_routing_manager}
    Sleep  5s

Stop Routing Manager
    Log to Console  Stopping routing manager
    Run And Return Rc And Output    ${stop_routing_manager}
    Sleep  5s

Restart Routing Manager
    Log to Console  Restarting routing manager
    Stop Routing Manager
    Start Routing Manager

Start e2adapter
     Log to Console  Starting e2adapter
     Run And Return Rc And Output    ${start_e2adapter}
     Sleep  5s

Stop e2adapter
     Log to Console  Stopping e2adapter
     Run And Return Rc And Output    ${stop_e2adapter}
     Sleep  5s

Restart e2adapter
    Log to Console  Restarting e2adapter
    Stop e2adapter
    Start e2adapter

Flush And Populate DB
    [Arguments]    ${set_new_timestamp}=${True}
    Log To Console  Flushing and populating DB
    ${flush}=  cleanup_db.flush    ${set_new_timestamp}
    Sleep  2s
    Should Be Equal As Strings  ${flush}  True

Stop All Pods Except Simulator
    Stop E2 Manager
    Stop Dbass
    Stop E2
    Stop Routing Manager
    Stop e2adapter

Send eNB Setup Request
    Log To Console  Sending eNB setup request form e2adapter
    Restart e2adapter
    Wait until keyword succeeds  2 min    3 sec    Validate Required Dockers
    ${e2adapter_pod} =    Run And Return Rc And Output   kubectl get pods -n ricplt | /bin/grep e2adapter | /bin/grep Running | awk '{{print $1}}'
    ${send_enb_setup}    Evaluate    "kubectl -n ricplt exec -it ${e2adapter_pod[1]} cli send-e2setup-req 10.0.2.15"
    Run And Return Rc And Output    ${send_enb_setup}

Start Redis Monitor
    Log To Console  Starting redis monitor log
    ${redis_monitor_log_filename}      Evaluate      "redis_monitor.${SUITE NAME}.log".replace(" ","-")
    Set Suite Variable  ${redis_monitor_log_filename}
    Remove File  ${EXECDIR}/${redis_monitor_log_filename}
    Start Process    kubectl -n ricplt exec -it statefulset-ricplt-dbaas-server-0 redis-cli MONITOR>${EXECDIR}/${redis_monitor_log_filename}  shell=yes

Stop Redis Monitor
    Log To Console  Stopping redis monitor log
    log_scripts.kill_redis_monitor_root_process


Redis Monitor Logs - Verify Publish To Manipulation Channel
    [Arguments]       ${ran_name}    ${event}
    Log To Console  Verify Publish To Manipulation Channel
    Sleep    3s
    ${result}=  log_scripts.verify_redis_monitor_manipulation_message   ${EXECDIR}/${redis_monitor_log_filename}  ${ran_name}    ${event}
    Should Be Equal As Strings    ${result}      True

Redis Monitor Logs - Verify Publish To Connection Status Channel
    [Arguments]       ${ran_name}    ${event}
    Log To Console    Verify Publish To Connection Status Channel
    Sleep    3s
    ${result}=  log_scripts.verify_redis_monitor_connection_status_message   ${EXECDIR}/${redis_monitor_log_filename}  ${ran_name}    ${event}
    Should Be Equal As Strings    ${result}      True

Redis Monitor Logs - Verify NOT Published To Manipulation Channel
    [Arguments]       ${ran_name}    ${event}
    Log To Console  Verify NOT Published To Manipulation Channel
    Sleep    3s
    ${result}=  log_scripts.verify_redis_monitor_manipulation_message   ${EXECDIR}/${redis_monitor_log_filename}  ${ran_name}    ${event}
    Should Be Equal As Strings    ${result}      False