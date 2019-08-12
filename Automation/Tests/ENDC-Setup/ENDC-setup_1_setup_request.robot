*** Settings ***
Suite Setup     Start dockers
Resource   ../Resource/resource.robot
Library     OperatingSystem
Library    Collections
Library     REST        ${url}




*** Test Cases ***
Post Request setup node b endc-setup
    #${file}=     Get Binary File    ${PATH}
    #${file}=     Evaluate       json.loads($file)   json
    #Set To Dictionary   ${file}     ranName=test2
    Set Headers     ${header}
    POST        /v1/nodeb/endc-setup    ${json}
    Integer     response status       200



*** Keywords ***
Start dockers
     Run And Return Rc And Output    ${run_script}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    4



