curl http://localhost:8080/games/poll \
    --include \
    --header "Content-Type: application/json" \
    --request "GET" \
    --data '{"user_id":0, "game_id":0, "token":"baf"}'

curl http://localhost:8080/games/play/card \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"card":{"suit": 0, "value": 14}, "game_auth":{"user_id":0, "game_id":0, "token":"baf"}}'
