curl http://localhost:8080/games/poll \
    --include \
    --header "Content-Type: application/json" \
    --request "GET" \
    --data '{"user_id":1, "game_id":0, "token":"baf"}'