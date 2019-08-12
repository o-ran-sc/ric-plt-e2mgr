*** Settings ***
Suite Setup     Start dockers
Resource   ../Resource/resource.robot
Library     OperatingSystem
Library     REST      ${url}


*** Test Cases ***
Get Request node b gnb - resource not found 404
    GET      /v1/nodeb/test5
    Integer   response status   404
    Integer   response body errorCode   404
    String   response body errorMessage   "Resource not found"


*** Keywords ***
Start dockers
     Run And Return Rc And Output    ${run_script}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    4