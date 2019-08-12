*** Settings ***
Suite Setup  Flush Redis
Library      Process
Resource   ../Resource/resource.robot
Library     OperatingSystem
Library     REST      ${url}


*** Variables ***
${file}     ${CURDIR}/addtoredis.py
${flush_file}     ${CURDIR}/flush.py

*** Test Cases ***
Add nodes to redis db
    ${result}=  Run Process    python3.6    ${file}
    Should Be Equal As Strings      ${result.stdout}        Insert successfully to Redis


Get all node ids
    ${result}=  GET     v1/nodeb-ids
    #Output
    Integer  response status   200
    String   response body 0 inventoryName  test1
    String   response body 0 globalNbId plmnId   02f829
    String   response body 0 globalNbId nbId     007a80
    String   response body 1 inventoryName  test2
    String   response body 1 globalNbId plmnId   03f829
    String   response body 1 globalNbId nbId     001234



*** Keywords ***
Flush Redis
    ${result}=  Run Process    python3.6    ${flush_file}
    Should Be Equal As Strings      ${result.stdout}        Flush Success







