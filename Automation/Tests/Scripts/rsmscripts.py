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

import config
import redis
import cleanup_db


def getRedisClientDecodeResponse():

    c = config.redis_ip_address

    p = config.redis_ip_port

    return redis.Redis(host=c, port=p, db=0, decode_responses=True)


def verify_rsm_ran_info():

    r = getRedisClientDecodeResponse()
    
    value = "{\"ranName\":\"test1\",\"enb1MeasurementId\":1,\"enb2MeasurementId\":0,\"action\":\"start\",\"actionStatus\":false}"

    if r.get("{rsm},RAN:test1") == value:
        return True
    else:
        return False