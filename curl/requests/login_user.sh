#!/bin/bash
payload=@../payloads/login_user_1.json

curl --header "Content-Type: application/json" \
    --request POST \
    --data $payload \
    http://localhost:5000/api/user/login

