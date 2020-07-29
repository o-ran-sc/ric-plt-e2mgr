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
import k8s_helper

e2mgr_ip = k8s_helper.extract_service_ip("e2mgr-http")
e2mgr_address = "http://" + e2mgr_ip + ":3800"

e2t_alpha_ip = k8s_helper.extract_service_ip("e2term-rmr-alpha")
e2t_alpha_address = e2t_alpha_ip + ":38000"

e2adapter_pod_name = k8s_helper.extract_pod_name("e2adapter")

