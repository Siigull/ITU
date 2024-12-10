curl http://localhost:8080/games/poll \
    --include \
    --header "Content-Type: application/json" \
    --request "GET" \
    --data '{"user_id":2, "game_id":0, "token":"baf"}'

curl http://localhost:8080/games/play/choice \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"choices":[], "game_auth":{"user_id":2, "game_id":0, "token":"baf"}}'