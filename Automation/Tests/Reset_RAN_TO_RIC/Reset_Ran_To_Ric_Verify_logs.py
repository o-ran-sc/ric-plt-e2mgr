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


def verify(directory):

    file = 'e2mgr.log'

    path = '/'

    file_path = directory + path + file

    f = open(file_path,'r')

    found_message_10070 = False
    found_message_10071 = False

    for l in f:
        if l.find('MType: 10070') > 0 and l.find('Meid: \\"test1\\"') > 0:
            found_message_10070 = True
        elif l.find('MType: 10071') > 0 and l.find('Meid: \\"test1\\"') > 0:
            found_message_10071 = True

        if found_message_10070 and found_message_10071:
            break

    if found_message_10070 and found_message_10071:
        return True
    else:
        return False



