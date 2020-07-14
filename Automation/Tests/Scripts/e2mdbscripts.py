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
import json
import redis
import variables


def get_redis_client_decode_response():
    c = config.redis_ip_address
    p = config.redis_ip_port
    return redis.Redis(host=c, port=p, db=0, decode_responses=True)


def verify_ran_is_associated_with_e2t_instance(ran_name, e2t_address):
    r = get_redis_client_decode_response()
    e2t_instance_json = r.get("{e2Manager},E2TInstance:" + e2t_address)

    if e2t_instance_json is None:
        return False

    e2t_instance_dic = json.loads(e2t_instance_json)
    assoc_ran_list = e2t_instance_dic.get("associatedRanList")
    if assoc_ran_list is None:
        return False
    else:
        return ran_name in assoc_ran_list


def verify_e2t_instance_has_no_associated_rans(e2t_address):
    r = get_redis_client_decode_response()
    e2t_instance_json = r.get("{e2Manager},E2TInstance:" + e2t_address)
    e2t_instance_dic = json.loads(e2t_instance_json)
    assoc_ran_list = e2t_instance_dic.get("associatedRanList")
    return not assoc_ran_list


def verify_e2t_instance_exists_in_addresses(e2t_address):
    r = get_redis_client_decode_response()
    e2t_addresses_json = r.get("{e2Manager},E2TAddresses")
    e2t_addresses = json.loads(e2t_addresses_json)
    return e2t_address in e2t_addresses


def verify_e2t_instance_key_exists(e2t_address):
    r = get_redis_client_decode_response()
    return r.exists("{e2Manager},E2TInstance:" + e2t_address)


def populate_e2t_instances_in_e2m_db_for_get_e2t_instances_tc():
    r = get_redis_client_decode_response()
    r.set("{e2Manager},E2TAddresses", "[\"e2t.att.com:38000\"]")
    r.set("{e2Manager},E2TInstance:e2t.att.com:38000",
          "{\"address\":\"e2t.att.com:38000\",\"associatedRanList\":[\"test1\",\"test2\",\"test3\"],"
          "\"keepAliveTimestamp\":1577619310484022369,\"state\":\"ACTIVE\"}")
    return True


def verify_e2t_addresses_for_e2t_initialization_tc():
    r = get_redis_client_decode_response()

    value = "[\"{}\"]".format(variables.e2t_alpha_address)

    return r.get("{e2Manager},E2TAddresses") == value


def verify_e2t_instance_for_e2t_initialization_tc():
    r = get_redis_client_decode_response()

    e2_address = "\"address\":\"{}\"".format(variables.e2t_alpha_address)
    e2_associated_ran_list = "\"associatedRanList\":[]"
    e2_state = "\"state\":\"ACTIVE\""

    e2_db_instance = r.get("{{e2Manager}},E2TInstance:{}".format(variables.e2t_alpha_address))

    if e2_db_instance.find(e2_address) < 0:
        return False
    if e2_db_instance.find(e2_associated_ran_list) < 0:
        return False
    if e2_db_instance.find(e2_state) < 0:
        return False

    return True


def set_enable_ric_false():
    r = get_redis_client_decode_response()
    r.set("{e2Manager},GENERAL", "{\"enableRic\":false}")
    return True

