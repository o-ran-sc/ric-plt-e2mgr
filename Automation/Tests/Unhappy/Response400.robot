*** Settings ***
Suite Setup     Start dockers
Resource   ../Resource/resource.robot
Library     OperatingSystem
Library     REST      ${url}



*** Test Cases ***
Post Request setup node b endc-setup - 400 validation of fields
    Set Headers     ${header}
    POST        /v1/nodeb/endc-setup
    Integer    response status   400


*** Keywords ***
Start dockers
     Run And Return Rc And Output    ${run_script}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    4