robot##############################################################################
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
Suite Setup   Prepare Enviorment
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Resource   ../Resource/scripts_variables.robot
Library    ../Scripts/find_error_script.py
Library    OperatingSystem
Library    Collections


*** Test Cases ***

Stop E2T
    stop_e2
    Sleep  3s

Prepare logs for tests
    Remove log files
    Save logs

Verify Is Dead Message Printed
    ${result}    find_error_script.find_error     ${EXECDIR}    ${e2mgr_log_filename}  ${e2_is_dead_message_printed}
    Should Be Equal As Strings    ${result}      True

Start E2T
    start_e2