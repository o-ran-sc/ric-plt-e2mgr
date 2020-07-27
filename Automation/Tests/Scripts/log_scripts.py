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

import subprocess


def verify_log_message(file_path, message):

    file = open(file_path, 'r')

    for line in file:

        if line.find(message) > 0:
            return True

    return False


def verify_redis_monitor_manipulation_message(file_path, ran_name, event):
    message = "\"PUBLISH\" \"{e2Manager},RAN_MANIPULATION\" \"" + ran_name + "_" + event + "\""
    return verify_log_message(file_path, message)

def verify_redis_monitor_connection_status_message(file_path, ran_name, event):
    message = "\"PUBLISH\" \"{e2Manager},RAN_CONNECTION_STATUS_CHANGE\" \"" + ran_name + "_" + event + "\""
    return verify_log_message(file_path, message)


def kill_redis_monitor_root_process():
    kill_command = "for pid in $(pidof redis-cli); do sudo kill -9 $pid; done"
    return subprocess.check_output(["/bin/bash", "-c", kill_command], universal_newlines=True)
