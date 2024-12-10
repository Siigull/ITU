curl http://localhost:8080/games/poll \
    --include \
    --header "Content-Type: application/json" \
    --request "GET" \
    --data '{"user_id":0, "game_id":0, "token":"baf"}'

curl http://localhost:8080/games/play/card \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"card":{"suit": 3, "value": 7}, "game_auth":{"user_id":0, "game_id":0, "token":"baf"}}'

curl http://localhost:8080/games/play/card \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"card":{"suit": 2, "value": 12}, "game_auth":{"user_id":0, "game_id":0, "token":"baf"}}'
