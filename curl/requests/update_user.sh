#!/bin/bash
payload=@../payloads/update_user_1.json

curl --header "Content-Type: application/json"  \
    --header "x-access-token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOiI1ZTkzMzUxOTJhNDdjZGNhMzU2MTUxOWIiLCJleHAiOjE1ODY3MDcyNzl9.0gDAXmiShsgkNofi-QapUQfuKAeIA8gTcwvTXZW_W1A"   \
    --request POST \
    --data $payload \
    http://localhost:5000/api/user/update

