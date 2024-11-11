export const generate_user_auth = (user_id, game_id, token) => {
    var data = {}
    data.user_id = user_id
    data.game_id = game_id
    data.token = token 

    return data
}