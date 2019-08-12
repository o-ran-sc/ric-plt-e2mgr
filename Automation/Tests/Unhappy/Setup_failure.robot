*** Settings ***
Suite Setup     Start dockers
Resource   ../Resource/resource.robot
Library     OperatingSystem
Library     REST        ${url}




*** Test Cases ***
Post Request setup node b x2-setup - setup failure
    #${file}=     Get Binary File    ${PATH}
    #${file}=     Evaluate       json.loads($file)   json
    Set Headers     ${header}
    POST        /v1/nodeb/x2-setup   ${json}
    Sleep    1s
    POST        /v1/nodeb/x2-setup   ${json}
    Sleep    1s
    GET      /v1/nodeb/test1
    Integer    response status       200
    String     response body connectionStatus     CONNECTED_SETUP_FAILED
    String     response body failureType     X2_SETUP_FAILURE
    String     response body setupFailure networkLayerCause       HO_TARGET_NOT_ALLOWED




*** Keywords ***
Start dockers
     Run And Return Rc And Output    ${run_script}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    4