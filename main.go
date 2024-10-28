package main

import (
    "crypto/rand"
    "encoding/base64"
    "net/http"

    "github.com/gin-gonic/gin"
)

func generateRandomToken() string {
    bytes := make([]byte, 32)
	
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(bytes)
}

type User struct {
    ID int `json:"id"`
    Name string `json:"name"`
    Token string `json:"token"`
    Money int `json:"money"`
}

var users = []User{{ID:0, Name:"I'm everlasting", Token:generateRandomToken(), Money:69420},}
	
type GameStateEnum int

const (
    STATE_CHOICE GameStateEnum = iota 
    STATE_BETTING
    STATE_GAME
    STATE_END
)

type GameTypeEnum int 

const (
    TYPE_GAME GameTypeEnum = iota
    TYPE_SEVEN
    TYPE_HUNDRED
    TYPE_HUNDRED_SEVEN
)

type ChoiceStateEnum int 

const (
    STATE_CHOOSING_TRUMP ChoiceStateEnum = iota
    STATE_CHOOSING_GAME
    STATE_BETS
)

type ChoiceState struct {
    state ChoiceStateEnum
    game_flek bool
    seven_flek bool
    hundred_flek bool
    betting_round int
}

type GameStateAll struct {
    game_state GameStateEnum
    choice_state ChoiceState
    game_type GameTypeEnum

    cur_player_index int 
    attacker_index int
} 

type Card struct {
    Suit         int `json:"suit"`
    Value        int `json:"value"`
    SpecialValue int `json:"special_value"`
}

var default_deck = []Card{
    {Suit: 0, Value: 7, SpecialValue: 7}, {Suit: 0, Value: 8, SpecialValue: 8}, {Suit: 0, Value: 9, SpecialValue: 9}, {Suit: 0, Value: 10, SpecialValue: 14},
    {Suit: 0, Value: 11, SpecialValue: 11}, {Suit: 0, Value: 12, SpecialValue: 12}, {Suit: 0, Value: 13, SpecialValue: 13}, {Suit: 0, Value: 14, SpecialValue: 15},

    {Suit: 1, Value: 7, SpecialValue: 7}, {Suit: 1, Value: 8, SpecialValue: 8}, {Suit: 1, Value: 9, SpecialValue: 9}, {Suit: 1, Value: 10, SpecialValue: 14},
    {Suit: 1, Value: 11, SpecialValue: 11}, {Suit: 1, Value: 12, SpecialValue: 12}, {Suit: 1, Value: 13, SpecialValue: 13}, {Suit: 1, Value: 14, SpecialValue: 15},

    {Suit: 2, Value: 7, SpecialValue: 7}, {Suit: 2, Value: 8, SpecialValue: 8}, {Suit: 2, Value: 9, SpecialValue: 9}, {Suit: 2, Value: 10, SpecialValue: 14},
    {Suit: 2, Value: 11, SpecialValue: 11}, {Suit: 2, Value: 12, SpecialValue: 12}, {Suit: 2, Value: 13, SpecialValue: 13}, {Suit: 2, Value: 14, SpecialValue: 15},

    {Suit: 3, Value: 7, SpecialValue: 7}, {Suit: 3, Value: 8, SpecialValue: 8}, {Suit: 3, Value: 9, SpecialValue: 9}, {Suit: 3, Value: 10, SpecialValue: 14},
    {Suit: 3, Value: 11, SpecialValue: 11}, {Suit: 3, Value: 12, SpecialValue: 12}, {Suit: 3, Value: 13, SpecialValue: 13}, {Suit: 3, Value: 14, SpecialValue: 15},
}

type UsersDeck struct {
    User_ID int `json:"id"`
    Hand []Card `json:"hand"`
}

type Game struct {
    Users []*User `json:"users"`
    Hands []UserDeck `json:"hands"`
    ID int `json:"id"`
    Running bool `json:"running"`
    Game_type string `json:"game_type"`
    state GameStateAll
}

func (*Game) init(self *Game) {
    self.Hands = []Hands{}
    var temp_deck []Card
    copy(temp_deck, default_deck)

    rand.Shuffle(len(temp_deck), func(i, j int) {
        temp_deck[i], temp_deck[j] = temp_deck[j], temp_deck[i]
    })

    var deal_arr = [][]int{{5,12}, {12,22}, {22,32}}
    for i, _ := range self.users {
        var user_i = (i + self.cur_player_index) % 3
        append(self.Hands, UserDeck{User_ID:self.users[user_i], Hand:temp_deck[deal_arr[i][0]:deal_arr[i][1]]})
    }

    self.attacker_index = (self.attacker_index + 1) % 3
    self.cur_player_index = 0
    self.state.game_state = StateChoice
    self.state.choice_state.state = STATE_CHOOSING_TRUMP
    self.state.choice_state.game_flek = false
    self.state.choice_state.seven_flek = false
    self.state.choice_state.hundred_flek = false
    self.state.choice_state.betting_round = 0 
}

func (*Game) next_player(self *Game) {
    self.cur_player_index = (self.cur_player_index + 1) % 3
}

func (*Game) get_choices(self *Game) {
    state = self.state

    switch state.game_state{
        case STATE_CHOICE:
            switch state.choice_state {
                case STATE_CHOOSING_TRUMP:
                    return "card"

                case STATE_CHOOSING_GAME:
                    return []string{"game", "seven", "hundred", "hundred seven"}

                case STATE_BETS:
                    var bets = []string{}
                    if(state.choice_state.game_flek == false) {
                        bets = append(bets, "on game");
                    }
                    if(state.choice_state.seven_flek == false) {
                        bets = append(bets, "on seven");
                    }
                    if(state.choice_state.hundred_flek == false) {
                        bets = append(bets, "on hundred");
                    }

                    if(len(bets) == 0) {
                        return "next"
                    }

                    return bets
            }
                

        case STATE_BETTING:

        case STATE_GAME:
            return "card"

        case STATE_END:
            return "next game"
    }
}

func (*Game) next_state(self *Game) {
    state = self.state
    
    switch state.game_state {
        case StateChoice:
            switch state.choice_state {
                case STATE_CHOOSING_TRUMP:
                    state.choice_state = STATE_CHOOSING_GAME

                case STATE_CHOOSING_GAME:
                    self.next_player()

                case STATE_BETS:


                    self.next_player()
                    if(self.cur_player_index == self.attacker_index) {
                        switch self.state.game_type {
                            case TYPE_GAME:
                                if(self.state.choice_state.game_flek == false) {
                                    self.state.game_state = STATE_GAME
                                }
                            case TYPE_SEVEN:
                                if(self.state.choice_state.game_flek == false &&
                                   self.state.choice_state.seven_flek == false) {

                                    self.state.game_state = STATE_GAME
                                }
                            case TYPE_HUNDRED:
                                if(self.state.choice_state.game_flek == false &&
                                   self.state.choice_state.hundred_flek == false) {
                                     
                                     self.state.game_state = STATE_GAME
                                }
                            case TYPE_HUNDRED_SEVEN:
                                if(self.state.choice_state.game_flek == false &&
                                   self.state.choice_state.hundred_flek == false) {
                                      
                                      self.state.game_state = STATE_GAME
                                }
                        }
                        if(self.state.choice_state.game_flek == false &&
                           ) {

                        }
                    }
            }
        case StateBetting:

        case StateGame:

        case StateEnd:

    }
}

var games = []Game{
    {Users: []*User{}, ID: 0},
}

func list_games(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, games)
}

func list_users(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, users)
}

func find_game(id int) *Game {
    for _, g := range games {
        if g.ID == id {
            return &g
        }
    }
    return nil;
}

func find_user(id int) *User {
    for _, u := range users {
        if u.ID == id {
            return &u
        }
    }
    return nil;
}

type GameAuth struct {
    User_id int `json:"user_id"`
    Game_id int `json:"game_id"`
    User_token string `json:"token"`
}

func start_game(c *gin.Context) {
    var data GameAuth

    if err := c.BindJSON(&data); err != nil {
        return
    }

    game := find_game(data.Game_id)
    if (game == nil) {
        c.String(400, "Game not found");
        return;
    }

    user := find_user(data.User_id)
    if (user == nil) {
        c.String(400, "User not found");
        return;
    }

    if (game.Users[0].ID != user.ID) {
        c.String(400, "You are not the game owner");
        return;
    }

    if (len(game.Users) != 3) {
        c.String(http.StatusConflict, "Not enough players")
        return
    }

    if (game.Running != true) {
        game.Running = true;
    } else {
        c.String(http.StatusConflict, "Game is already running");
        return;
    }

    game.attacker_index = -1
    game.init()
    
    for i, g := range games {
        if g.ID == game.ID {
            games[i] = *game
            break
        }
    }
}

func join_game(c *gin.Context) {
    var data GameAuth

    if err := c.BindJSON(&data); err != nil {
        return
    }

    game := find_game(data.Game_id)
    if (game == nil) {
        c.String(400, "Game not found");
        return;
    }

    user := find_user(data.User_id)
    if (user == nil) {
        c.String(400, "User not found");
        return;
    }

    if (len(game.Users) >= 3) {
        c.String(http.StatusConflict, "Game is full")
        return
    }

    for _, u := range game.Users {
        if u.ID == user.ID {
            c.String(http.StatusConflict, "User already in the game")
            return
        }
    }

    game.Users = append(game.Users, user)
    for i, g := range games {
        if g.ID == game.ID {
            games[i] = *game
            break
        }
    }
}

type UserName struct {
    user_name string `json:"user"`
}

func new_user(c *gin.Context) {
    name := c.Param("name")

    last_id := users[len(users) - 1].ID

    user := User{ID: last_id + 1, Name: name, Token: generateRandomToken(), Money: 0}

    users = append(users, user)

    c.IndentedJSON(http.StatusOK, user)
}

type PollData struct {
    User_choices []string `json:"choices"`
    Game_data Game `json:"game"`
}

func poll_game(c *gin.Context) {
    var data GameAuth

    if err := c.BindJSON(&data); err != nil {
        return
    }

    game := find_game(data.Game_id)
    if (game == nil) {
        c.String(400, "Game not found");
        return;
    }

    user := find_user(data.User_id)
    if (user == nil) {
        c.String(400, "User not found");
        return;
    }

    found = false;
    for _, u := range game.Users {
        if u.ID == user.ID {
            found = true;
            break
        }
    }

    if(found == false) {
        c.String(http.StatusConflict, "You are not in this game")
        return
    }

    var PollData return_val = return_val{Choices:nil, Game:game}

    if(game.Users[cur_player_index].ID == user.ID) {
        return_val.User_choices = game.get_choices();
    }

    c.IndentedJSON(http.StatusOK, return_val)
}

func (*Game) play_choice(self *Game, choices []string) string {
    state = self.state

    switch state.game_state{
        case STATE_CHOICE:
            switch state.choice_state {
                case STATE_CHOOSING_TRUMP:
                    return "You have to choose a card"

                case STATE_CHOOSING_GAME:
                    switch choices[0] {
                        case "game":
                            state.game_type = TYPE_GAME
                        case "seven":
                            state.game_type = TYPE_SEVEN
                        case "hundred":
                            state.game_type = TYPE_HUNDRED
                        case "hundred seven":
                            state.game_type = TYPE_HUNDRED_SEVEN

                        default:
                            return "Not a valid choice" 
                    }

                case STATE_BETS:
                    for _, c := range choices {
                        switch c {
                            case "on game":
                                state.choice_state.game_flek = true;
                            case "on seven":
                                state.choice_state.seven_flek = true;
                            case "on hundred":
                                state.choice_state.hundred_flek = true;

                            default:
                                return "Not a valid choice" 
                        }
                    }
            }
            
        case STATE_BETTING:

        case STATE_GAME:
            return "card"

        case STATE_END:
            return "next game"
    }

    self.next_state()
}

func play_card(c *gin.Context) {

}

type Choice struct {
    Choices []string `json:"choices"`
    GameData GameAuth `json:"game_auth"`
}

func parse_choice(c *gin.Context) {
    var data Choice

    if err := c.BindJSON(&data); err != nil {
        return
    }

    game := find_game(data.Game_id)
    if (game == nil) {
        c.String(400, "Game not found");
        return;
    }

    user := find_user(data.User_id)
    if (user == nil) {
        c.String(400, "User not found");
        return;
    }

    player_i = -1;
    for i, u := range game.Users {
        if u.ID == user.ID {
            player_i = i;
            break
        }
    }

    if(player_i == false) {
        c.String(http.StatusConflict, "You are not in this game")
        return
    }

    if(player_i != game.cur_player_index) {
        c.String(http.StatusConflict, "It is not your turn")
        return
    }

    play_choice(data.Choices)
}

func main() {
    rand.Seed(time.Now().UnixNano())

    router := gin.Default()
    router.GET("/games", list_games)
    router.GET("/users", list_users)
    router.GET("/join/:name", new_user)

    router.POST("/games/join", join_game)
    router.POST("/games/start", start_game)
    router.GET("/games/poll", poll_game)
    router.POST("/games/play/choice", parse_choice)
    router.POST("/games/play/card", parse_card)

    router.Run("localhost:8080")
}



///////////////////////////////////////////
// package main

// import (
//     "net/http"

//     "github.com/gin-gonic/gin"
// )

// // album represents data about a record album.
// type album struct {
//     ID     string  `json:"id"`
//     Title  string  `json:"title"`
//     Artist string  `json:"artist"`
//     Price  float64 `json:"price"`
// }

// // albums slice to seed record album data.
// var albums = []album{
//     {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
//     {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
//     {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
// }

// func main() {
//     router := gin.Default()
//     router.GET("/albums", getAlbums)
//     router.GET("/albums/:id", getAlbumByID)
//     router.POST("/albums", postAlbums)

//     router.Run("localhost:8080")
// }

// // getAlbums responds with the list of all albums as JSON.
// func getAlbums(c *gin.Context) {
//     c.IndentedJSON(http.StatusOK, albums)
// }

// // postAlbums adds an album from JSON received in the request body.
// func postAlbums(c *gin.Context) {
//     var newAlbum album

//     // Call BindJSON to bind the received JSON to
//     // newAlbum.
//     if err := c.BindJSON(&newAlbum); err != nil {
//         return
//     }

//     // Add the new album to the slice.
//     albums = append(albums, newAlbum)
//     c.IndentedJSON(http.StatusCreated, newAlbum)
// }

// // getAlbumByID locates the album whose ID value matches the id
// // parameter sent by the client, then returns that album as a response.
// func getAlbumByID(c *gin.Context) {
//     id := c.Param("id")

//     // Loop through the list of albums, looking for
//     // an album whose ID value matches the parameter.
//     for _, a := range albums {
//         if a.ID == id {
//             c.IndentedJSON(http.StatusOK, a)
//             return
//         }
//     }
//     c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
// }