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

import shlex
import sys

#file = 'e2mgr.log'

path = '/'

#file_path = sys.argv[1] + path + file

f = open('/home/ubuntu/PycharmProjects/Oran/e2mgr.log','r')

found_message_10370 = False
found_message_10371 = False

for l in f:
    #t = shlex.split(l)
    #for i in t:
        #m = i.split(",")

        for l in f:

            if l.find('MType: 10370') > 0 and l.find('Meid: \\"test1\\"') > 0:
                found_message_10370 = True
            elif l.find('MType: 10371') > 0 and l.find('Meid: \\"test1\\"') > 0:
                found_message_10371 = True
            if found_message_10370 and found_message_10371:
                break


if found_message_10370 and found_message_10371:
    print("Found All Configuration Update logs")
else:
    print("Didn't find any Configuration Update logs ")



