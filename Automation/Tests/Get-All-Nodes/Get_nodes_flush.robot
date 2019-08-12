*** Settings ***
Test Teardown       Flush Redis


*** Variables ***
${file}     ${CURDIR}/flush.py

*** Keywords ***
Flush Redis
    ${result}=  Run Process    python3.6    ${file}
    Should Be Equal As Strings      ${result.stdout}        Flush Success