*** Settings ***
Suite Setup     Start dockers
Resource   ../Resource/resource.robot
Library     OperatingSystem
Library     REST      ${url}





*** Test Cases ***
Post Request setup node b x-2
    #${file}=     Get Binary File    ${PATH}
    #${file}=     Evaluate       json.loads($file)   json
    Set Headers     ${header}
    POST        /v1/nodeb/x2-setup    ${json}
    Integer     response status       200


*** Keywords ***
Start dockers
     Run And Return Rc And Output    ${run_script}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    4





