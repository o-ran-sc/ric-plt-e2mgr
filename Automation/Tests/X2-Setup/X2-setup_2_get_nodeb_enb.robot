*** Settings ***
#Suite Setup    Start dockers
Resource   ../Resource/resource.robot
Library     OperatingSystem
Library     REST      ${url}





*** Test Cases ***
Get Request node b enb
    Sleep    1s
    GET      /v1/nodeb/test1
    Integer  response status  200
    String   response body ranName    test1
    String   response body ip    10.0.2.15
    Integer  response body port     5577
    String   response body connectionStatus    CONNECTED
    String   response body nodeType     ENB
    String   response body enb enbType     MACRO_ENB
    Integer  response body enb servedCells 0 pci  99
    String   response body enb servedCells 0 cellId   02f829:0007ab00
    String   response body enb servedCells 0 tac    0102
    String   response body enb servedCells 0 broadcastPlmns 0   "02f829"
    Integer  response body enb servedCells 0 choiceEutraMode fdd ulearFcn    1
    Integer  response body enb servedCells 0 choiceEutraMode fdd dlearFcn    1
    String   response body enb servedCells 0 choiceEutraMode fdd ulTransmissionBandwidth   BW50
    String   response body enb servedCells 0 choiceEutraMode fdd dlTransmissionBandwidth   BW50



*** Keywords ***
#Start dockers
     #Run And Return Rc And Output    ${run_script}
     #${result}=  Run And Return Rc And Output     ${docker_command}
     #Should Be Equal As Integers    ${result[1]}    4





