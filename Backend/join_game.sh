curl http://localhost:8080/games/join \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"user_id":2, "game_id":2, "token":"baf"}'

# curl http://localhost:8080/games/join \
#     --include \
#     --header "Content-Type: application/json" \
#     --request "POST" \
#     --data '{"user_id":1, "game_id":0, "token":"baf"}'

# curl http://localhost:8080/games/join \
#     --include \
#     --header "Content-Type: application/json" \
#     --request "POST" \
#     --data '{"user_id":2, "game_id":0, "token":"baf"}'