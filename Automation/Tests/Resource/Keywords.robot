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

*** Settings ***
Documentation   Keywords file
Library     ../Scripts/cleanup_db.py
Resource   ../Resource/resource.robot
Library     OperatingSystem





*** Keywords ***
Post Request setup node b x-2
    Set Headers     ${header}
    POST        /v1/nodeb/x2-setup    ${json}


Put Request Resource Status Start
    Set Headers     ${header}
    PUT        /v1/general/resourcestatus  ${resource_status_start_json}


Put Request Resource Status Stop
    Set Headers     ${header}
    PUT        /v1/general/resourcestatus  ${resource_status_stop_json}

Get Request node b enb test1
    Sleep    1s
    GET      /v1/nodeb/test1


Get Request node b enb test2
    Sleep    1s
    GET      /v1/nodeb/test2


Remove log files
    Remove File  ${EXECDIR}/${gnb_log_filename}
    Remove File  ${EXECDIR}/${e2mgr_log_filename}
    Remove File  ${EXECDIR}/${e2t_log_filename}
    Remove File  ${EXECDIR}/${rsm_log_filename}
    Remove File  ${EXECDIR}/${e2e_simu_log_filename}
    Remove File  ${EXECDIR}/${e2adapter_log_filename}

Save logs
    Sleep   1s
    Run     ${Save_sim_log}
    Run     ${Save_e2mgr_log}
    Run     ${Save_e2t_log}
    Run     ${Save_rsm_log}
    Run     ${Save_e2e_simu_log}
    Run     ${Save_e2adapter_log}


Post Request setup node b endc-setup
    Set Headers     ${header}
    POST        /v1/nodeb/endc-setup    ${endcjson}

Stop Simulator
    Run And Return Rc And Output    ${stop_simu}


Prepare Simulator For Load Information
     Run And Return Rc And Output    ${stop_simu}
     Run And Return Rc And Output    ${docker_Remove}
     ${flush}  cleanup_db.flush
     Should Be Equal As Strings  ${flush}  True
     Run And Return Rc And Output    ${run_simu_load}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    ${docker_number}

Prepare Enviorment
     ${starting_timestamp}    Evaluate   datetime.datetime.now(datetime.timezone.utc).isoformat("T")   modules=datetime 
     ${e2t_log_filename}      Evaluate      "e2t.${SUITE NAME}.log".replace(" ","-")
     ${e2mgr_log_filename}    Evaluate      "e2mgr.${SUITE NAME}.log".replace(" ","-")
     ${gnb_log_filename}      Evaluate      "gnb.${SUITE NAME}.log".replace(" ","-")
     ${rsm_log_filename}      Evaluate      "rsm.${SUITE NAME}.log".replace(" ","-")
     ${e2e_simu_log_filename}      Evaluate      "e2esimu.${SUITE NAME}.log".replace(" ","-")
     ${e2adapter_log_filename}    Evaluate  "e2adapter.${SUITE NAME}.log".replace(" ","-")
     ${Save_sim_log}          Evaluate   'docker logs --since ${starting_timestamp} gnbe2_simu > ${gnb_log_filename}'
     ${Save_e2mgr_log}        Evaluate   'docker logs --since ${starting_timestamp} e2mgr > ${e2mgr_log_filename}'
     ${Save_e2t_log}          Evaluate   'docker logs --since ${starting_timestamp} e2 > ${e2t_log_filename}'
     ${Save_rsm_log}          Evaluate   'docker logs --since ${starting_timestamp} rsm > ${rsm_log_filename}'
     ${Save_e2e_simu_log}     Evaluate   'docker logs --since ${starting_timestamp} e2esimu > ${e2e_simu_log_filename}'
     ${Save_e2adapter_log}    Evaluate   'docker logs --since ${starting_timestamp} e2adapter > ${e2adapter_log_filename}'
     Set Suite Variable  ${e2t_log_filename}  
     Set Suite Variable  ${e2mgr_log_filename}  
     Set Suite Variable  ${gnb_log_filename}   
     Set Suite Variable  ${rsm_log_filename}  
     Set Suite Variable  ${e2e_simu_log_filename} 
     Set Suite Variable  ${e2adapter_log_filename} 
     Set Suite Variable  ${Save_sim_log}
     Set Suite Variable  ${Save_e2mgr_log}
     Set Suite Variable  ${Save_e2t_log}
     Set Suite Variable  ${Save_rsm_log}
     Set Suite Variable  ${Save_e2e_simu_log}
     Set Suite Variable  ${Save_e2adapter_log}
	
     ${flush}  cleanup_db.flush
     Should Be Equal As Strings  ${flush}  True
     Run And Return Rc And Output    ${stop_simu}
     Run And Return Rc And Output    ${stop_e2e_simu}
     Run And Return Rc And Output    ${docker_Remove}
     Run And Return Rc And Output    ${docker_remove_e2e_simu}
     Run And Return Rc And Output    ${run_simu_regular}
     Run And Return Rc And Output    ${run_e2e_simu_regular}
     Run And Return Rc And Output    ${restart_e2adapter}
     Sleep  2s
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    ${docker_number}

Start E2
     Run And Return Rc And Output    ${start_e2}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    ${docker_number}
     Sleep  2s

Start Dbass
     Run And Return Rc And Output    ${dbass_remove}
     Run And Return Rc And Output    ${dbass_start}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    ${docker_number}

Stop Dbass
     Run And Return Rc And Output    ${dbass_stop}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    ${docker_number-1}





