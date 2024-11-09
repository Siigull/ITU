curl http://localhost:8080/join/kys \
    --include \
    --header "Content-Type: application/json" \
    --request "GET"

curl http://localhost:8080/join/kys \
    --include \
    --header "Content-Type: application/json" \
    --request "GET"

curl http://localhost:8080/games/join \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"user_id":0, "game_id":0, "token":"baf"}'

curl http://localhost:8080/games/join \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"user_id":1, "game_id":0, "token":"baf"}'

curl http://localhost:8080/games/join \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"user_id":2, "game_id":0, "token":"baf"}'
    