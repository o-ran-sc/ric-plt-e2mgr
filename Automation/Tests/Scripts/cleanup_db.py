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
import config
import redis
import time
import k8s_helper


def flush(set_new_timestamp):
    c = config.redis_ip_address

    p = config.redis_ip_port

    r = redis.Redis(host=c, port=p, db=0, )

    e2t_ip = k8s_helper.extract_service_ip("e2term-rmr-alpha")
    et2_address = e2t_ip + ":38000"

    r.flushall()
    r.set("{e2Manager},GENERAL", "{\"enableRic\":true}")
    r.set("{e2Manager},E2TAddresses", "[\"{}\"]".format(et2_address))

    timestamp = str(int((time.time() + 2) * 1000000000)) if set_new_timestamp else str(
        int((time.time() - 300) * 1000000000))
    r.set("{{e2Manager}},E2TInstance:{}".format(et2_address),
          "{{\"address\":\"{}\",\"associatedRanList\":[],\"keepAliveTimestamp\":".format(et2_address) + timestamp +
          ",\"state\":\"ACTIVE\",\"deletionTimeStamp\":0}")

    return True
