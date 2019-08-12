*** Settings ***
Documentation    Resource file


*** Variables ***
${url}   http://localhost:3800
${json}    {"ranIp": "10.0.2.15","ranPort": 5577,"ranName":"test1"}
${header}  {"Content-Type": "application/json"}
${run_script}      /home/ubuntu/run.sh
${docker_command}  docker ps | grep 1.0 | wc --lines
