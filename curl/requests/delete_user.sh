#!/bin/bash
payload=@../payloads/delete_user_1.json
tk="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOiI1ZTkzMmE1ZWIwN2MzMTFkZDIxNTJkOWQiLCJleHAiOjE1ODY3MDQ1ODV9.ekS-qnKG9WyLvX0waXITR9vGbJybI_2gyWV998XeMEY"
curl --header "Content-Type: application/json"  \
     --header "x-access-token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOiI1ZTkzMzBiNDcyNzZkMWYwMDczYmM5NDUiLCJleHAiOjE1ODY3MDQ3NTR9.NafA1xOODQszC10I9dE3ZwccHkax3FxPuqJkVabcCl8"   \
    --request DELETE \
    http://localhost:5000/api/user/delete

