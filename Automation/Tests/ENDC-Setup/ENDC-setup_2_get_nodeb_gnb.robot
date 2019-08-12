*** Settings ***
Resource   ../Resource/resource.robot
#Suite Setup     Start dockers
Library     OperatingSystem
Library    Collections
Library     REST        ${url}



*** Test Cases ***
Get request gnb
    Sleep    1s
    GET      /v1/nodeb/test1
    Integer  response status  200
    String   response body ranName    test1
    String   response body ip    10.0.2.15
    String   response body connectionStatus    CONNECTED
    Integer  response body port     5577
    String   response body nodeType     GNB
    Integer  response body gnb servedNrCells 0 servedNrCellInformation nrPci  99
    String   response body gnb servedNrCells 0 servedNrCellInformation cellId  02f829:0007ab0120
    String   response body gnb servedNrCells 0 servedNrCellInformation servedPlmns 0   "02f829"
    String   response body gnb servedNrCells 0 servedNrCellInformation nrMode  FDD
    String   response body gnb servedNrCells 0 servedNrCellInformation choiceNrMode fdd ulFreqInfo nrArFcn   100
    Integer  response body gnb servedNrCells 0 servedNrCellInformation choiceNrMode fdd ulFreqInfo frequencyBands 0 nrFrequencyBand       9
    Integer  response body gnb servedNrCells 0 servedNrCellInformation choiceNrMode fdd ulFreqInfo frequencyBands 0 supportedSulBands 0    9
    String   response body gnb servedNrCells 0 servedNrCellInformation choiceNrMode fdd dlFreqInfo nrArFcn   100
    Integer  response body gnb servedNrCells 0 servedNrCellInformation choiceNrMode fdd dlFreqInfo frequencyBands 0 nrFrequencyBand        9
    Integer  response body gnb servedNrCells 0 servedNrCellInformation choiceNrMode fdd dlFreqInfo frequencyBands 0 supportedSulBands 0     9
    String   response body gnb servedNrCells 0 servedNrCellInformation choiceNrMode fdd ulTransmissionBandwidth nrscs   SCS15
    String  response body gnb servedNrCells 0 servedNrCellInformation choiceNrMode fdd ulTransmissionBandwidth ncnrb   NRB11
    String  response body gnb servedNrCells 0 servedNrCellInformation choiceNrMode fdd dlTransmissionBandwidth nrscs   SCS15
    String  response body gnb servedNrCells 0 servedNrCellInformation choiceNrMode fdd dlTransmissionBandwidth ncnrb   NRB11




#*** Keywords ***
#Start dockers
     #Run And Return Rc And Output    ${run_script}
     #${result}=  Run And Return Rc And Output     ${docker_command}
     #Should Be Equal As Integers    ${result[1]}    4

