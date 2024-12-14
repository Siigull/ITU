package main

import (
    "github.com/gin-contrib/cors"
    "math/rand"
    "encoding/base64"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "fmt"
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

var users = []User{{ID:0, Name:"I'm everlasting", Token:generateRandomToken(), Money:69420},
                   {ID:1, Name:"User1", Token:generateRandomToken(), Money:-69}}
	
type GameStateEnum int

const (
    STATE_CHOICE GameStateEnum = iota 
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
    STATE_THROW_AWAY1
    STATE_THROW_AWAY2
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
    specialValue int 
}

var default_deck = []Card{
    {Suit: 0, Value: 7, specialValue: 7}, {Suit: 0, Value: 8, specialValue: 8}, {Suit: 0, Value: 9, specialValue: 9}, {Suit: 0, Value: 10, specialValue: 14},
    {Suit: 0, Value: 11, specialValue: 11}, {Suit: 0, Value: 12, specialValue: 12}, {Suit: 0, Value: 13, specialValue: 13}, {Suit: 0, Value: 14, specialValue: 15},

    {Suit: 1, Value: 7, specialValue: 7}, {Suit: 1, Value: 8, specialValue: 8}, {Suit: 1, Value: 9, specialValue: 9}, {Suit: 1, Value: 10, specialValue: 14},
    {Suit: 1, Value: 11, specialValue: 11}, {Suit: 1, Value: 12, specialValue: 12}, {Suit: 1, Value: 13, specialValue: 13}, {Suit: 1, Value: 14, specialValue: 15},

    {Suit: 2, Value: 7, specialValue: 7}, {Suit: 2, Value: 8, specialValue: 8}, {Suit: 2, Value: 9, specialValue: 9}, {Suit: 2, Value: 10, specialValue: 14},
    {Suit: 2, Value: 11, specialValue: 11}, {Suit: 2, Value: 12, specialValue: 12}, {Suit: 2, Value: 13, specialValue: 13}, {Suit: 2, Value: 14, specialValue: 15},

    {Suit: 3, Value: 7, specialValue: 7}, {Suit: 3, Value: 8, specialValue: 8}, {Suit: 3, Value: 9, specialValue: 9}, {Suit: 3, Value: 10, specialValue: 14},
    {Suit: 3, Value: 11, specialValue: 11}, {Suit: 3, Value: 12, specialValue: 12}, {Suit: 3, Value: 13, specialValue: 13}, {Suit: 3, Value: 14, specialValue: 15},
}

type UsersDeck struct {
    User_ID int `json:"id"`
    Hand []Card `json:"hand"`
}

type Game struct {
    Users  []*User `json:"users"`
    Hands  []UsersDeck `json:"hands"`
    Table  []Card `json:"table"`
    Talons []UsersDeck `json:"talons"`
    ID int `json:"id"`
    Running bool `json:"running"`
    Betting bool `json:"betting"`
    Game_type string `json:"game_type"`
    Trump_card Card `json:"trump_card"`
    temp_deck []Card
    state GameStateAll
}

func (self *Game) init() {
    self.Hands = []UsersDeck{}
    self.temp_deck = make([]Card, len(default_deck))
    copy(self.temp_deck, default_deck)

    rand.Shuffle(len(self.temp_deck), func(i, j int) {
        self.temp_deck[i], self.temp_deck[j] = self.temp_deck[j], self.temp_deck[i]
    })

    self.Hands = []UsersDeck{}
    self.Talons = []UsersDeck{}

    self.state.attacker_index = (self.state.attacker_index + 1) % 3
    self.state.cur_player_index = self.state.attacker_index

    var deal_arr = [][]int{{5,12}, {12,22}, {22,32}}
    for i, _ := range self.Users {
        var user_i = (i + self.state.cur_player_index) % 3
        hand_copy := make([]Card, deal_arr[i][1] - deal_arr[i][0])
        copy(hand_copy, self.temp_deck[deal_arr[i][0]:deal_arr[i][1]])
        self.Talons = append(self.Talons, UsersDeck{User_ID:self.Users[user_i].ID, Hand:[]Card{}})
        self.Hands = append(self.Hands, UsersDeck{User_ID:self.Users[user_i].ID, Hand:hand_copy})
    }
    
    self.Trump_card = Card{}
    self.Running = true
    self.Betting = false
    self.state.game_state = STATE_CHOICE
    self.state.choice_state.state = STATE_CHOOSING_TRUMP
    self.state.choice_state.game_flek = false
    self.state.choice_state.seven_flek = false
    self.state.choice_state.hundred_flek = false
    self.state.choice_state.betting_round = 0 
}

func (self *Game) next_player() {
    self.state.cur_player_index = (self.state.cur_player_index + 1) % 3
}

func (self *Game) get_choices() []string {
    state := &self.state

    if(!self.Running) {
        return []string{}
    }

    switch state.game_state{
        case STATE_CHOICE:
            switch state.choice_state.state {
                case STATE_CHOOSING_TRUMP:
                    return []string{"card"}
                    
                case STATE_THROW_AWAY1:
                    return []string{"card"}
                
                case STATE_THROW_AWAY2:
                    return []string{"card"}

                case STATE_CHOOSING_GAME:
                    return []string{"game", "seven", "hundred", "hundred seven"}

                case STATE_BETS:
                        var bets = []string{}
                        if(state.choice_state.game_flek == false && 
                        (state.game_type == TYPE_GAME ||
                            state.game_type == TYPE_SEVEN)) {
                            bets = append(bets, "on game");
                        }
                        if(state.choice_state.seven_flek == false &&
                        (state.game_type == TYPE_SEVEN || 
                            state.game_type == TYPE_HUNDRED_SEVEN)) {
                            bets = append(bets, "on seven");
                        }
                        if(state.choice_state.hundred_flek == false &&
                        (state.game_type == TYPE_HUNDRED ||
                            state.game_type == TYPE_HUNDRED_SEVEN)) {
                            bets = append(bets, "on hundred");
                        }

                        if(len(bets) == 0) {
                            return []string{"next"}
                        }

                    return bets
            }

        case STATE_GAME:
            return []string{"card"}

        case STATE_END:
            return []string{"next game"}
    }

    return []string{"error"}
}

var games = []Game{
    {Users: []*User{&users[0], &users[1]}, ID: 0},
    {Users: []*User{}, ID: 1},
    {Users: []*User{}, ID: 2},
    {Users: []*User{}, ID: 3},
}

func list_games(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, games)
}

func list_users(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, users)
}

func find_game(id int) *Game {
    for i, g := range games {
        if g.ID == id {
            return &games[i]
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

    game.state.attacker_index = -1
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

    found := false;
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

    var return_val PollData = PollData{User_choices:nil, Game_data:*game}
    
    if(game.Users[game.state.cur_player_index].ID == user.ID) {
        return_val.User_choices = game.get_choices();
    }

    c.IndentedJSON(http.StatusOK, return_val)
}

func (self *Game) no_flek() bool {
    if (self.state.choice_state.game_flek == false &&
        self.state.choice_state.seven_flek == false && 
        self.state.choice_state.hundred_flek == false) {
        
        return true;
    }
    return false;
}

func (self* Game) clear_fleks() {
    self.state.choice_state.game_flek = false;
    self.state.choice_state.seven_flek = false;
    self.state.choice_state.hundred_flek = false; 
}

func (self *Game) play_choice(choices []string) string {
    state := &self.state

    switch state.game_state{
        case STATE_CHOICE:
            switch state.choice_state.state {
                case STATE_CHOOSING_TRUMP:
                    return "You have to choose a card"
                case STATE_THROW_AWAY1:
                    return "You have to choose a card"
                case STATE_THROW_AWAY2:
                    return "You have to choose a card"

                case STATE_CHOOSING_GAME:
                    fmt.Println(len(choices))
                    if (len(choices) != 1) {
                        return "Bad number of choices"
                    }

                    switch choices[0] {
                        case "game":
                            state.game_type = TYPE_GAME
                            self.Game_type = "game"
                        case "seven":
                            state.game_type = TYPE_SEVEN
                            self.Game_type = "seven"
                        case "hundred":
                            state.game_type = TYPE_HUNDRED
                            self.Game_type = "hundred"
                        case "hundred seven":
                            state.game_type = TYPE_HUNDRED_SEVEN
                            self.Game_type = "hundred_seven"

                        default:
                            return "Not a valid choice" 
                    }
                    self.Betting = true
                    state.choice_state.state = STATE_BETS
                    self.next_player()

                case STATE_BETS:
                    if (self.state.attacker_index == self.state.cur_player_index) {
                        game_f := state.choice_state.game_flek
                        seven_f := state.choice_state.seven_flek
                        hundred_f := state.choice_state.hundred_flek

                        self.clear_fleks()
                        
                        fmt.Println(len(choices))
                        if(len(choices) == 0) {
                            self.Betting = false
                            state.game_state = STATE_GAME
                            return ""
                        }

                        for _, c := range choices {
                            switch c {
                                case "on game":
                                    if(game_f) {
                                        state.choice_state.game_flek = true;
                                    } else {
                                        return "Not a valid choice" 
                                    }
                                case "on seven":
                                    if(seven_f) {
                                        state.choice_state.seven_flek = true;
                                    } else {
                                        return "Not a valid choice" 
                                    }
                                case "on hundred":
                                    if(hundred_f) {
                                        state.choice_state.hundred_flek = true;
                                    } else {
                                        return "Not a valid choice" 
                                    }
    
                                default:
                                    return "Not a valid choice" 
                            }
                        }
                    } else {
                        for _, c := range choices {
                            switch c {
                                case "next":
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
                    
                    self.next_player()
                    if (self.state.attacker_index == self.state.cur_player_index) {
                        if (self.no_flek()) {
                            self.Betting = false
                            state.game_state = STATE_GAME
                        }

                        state.choice_state.game_flek = !state.choice_state.game_flek
                        state.choice_state.seven_flek = !state.choice_state.seven_flek
                        state.choice_state.hundred_flek = !state.choice_state.hundred_flek
                    }
            }
            
        case STATE_GAME:
            self.Betting = false
            return "You have to choose a card"

        case STATE_END:
            if (choices[0] == "next game") {
                self.init()
            }
    }

    return ""
}

func (self *Game) clear_table() {

}

func (self *Game) place_card(card Card) bool {
    self.Table = append(self.Table, card) 

    if (len(self.Table) == 3) {
        return true;
    }
    
    return false;
}

func (self *Game) next_round() {
    // for i, el := range self.Table {
    // }
}

func (self *Game) attacker_draw() {
    self.Hands[0].Hand = append(self.Hands[0].Hand, self.temp_deck[0:5]...)
}

func (self *Game) card_throw_away(card Card) string {
    found := false
    for i, el := range self.Hands[0].Hand {
        if (el.Suit == card.Suit && el.Value == card.Value) {
            found = true
            self.Hands[0].Hand = append(self.Hands[0].Hand[:i], self.Hands[0].Hand[i+1:]...)
            break
        }
    }

    if(found == false) {
        return "This card is not yours"
    }

    self.Talons[1].Hand = append(self.Talons[1].Hand, card)

    return ""
}

func (self *Game) play_card(card Card) string {
    state := &self.state

    fmt.Println(state)

    switch state.game_state{
        case STATE_CHOICE:
            switch state.choice_state.state {
                case STATE_CHOOSING_TRUMP:
                    self.Trump_card = card;
                    state.choice_state.state = STATE_THROW_AWAY1;
                    self.attacker_draw()

                case STATE_THROW_AWAY1:
                    err := self.card_throw_away(card)
                    if(err != "") {
                        return err;
                    }
                    state.choice_state.state = STATE_THROW_AWAY2

                case STATE_THROW_AWAY2:
                    err := self.card_throw_away(card)
                    if(err != "") {
                        return err;
                    }
                    state.choice_state.state = STATE_CHOOSING_GAME 

                case STATE_CHOOSING_GAME:
                    return "You can't choose a card here"
                case STATE_BETS:
                    return "You can't choose a card here"
            }
            
        case STATE_GAME:
            if(self.place_card(card)) {
                self.next_round()
            } else {
                self.next_player()
            }

        case STATE_END:
            return "next game"
    }

    return ""
}

type CardChoice struct {
    ChosenCard Card `json:"card"`
    GameData GameAuth `json:"game_auth"`
}

func parse_card(c *gin.Context) {
    var data CardChoice

    if err := c.BindJSON(&data); err != nil {
        return
    }

    game := find_game(data.GameData.Game_id)
    if (game == nil) {
        c.String(400, "Game not found");
        return;
    }

    user := find_user(data.GameData.User_id)
    if (user == nil) {
        c.String(400, "User not found");
        return;
    }

    player_i := -1;
    for i, u := range game.Users {
        if u.ID == user.ID {
            player_i = i;
            break
        }
    }

    if(player_i == -1) {
        c.String(http.StatusConflict, "You are not in this game")
        return
    }

    if(player_i != game.state.cur_player_index) {
        c.String(http.StatusConflict, "It is not your turn")
        return
    }

    player_hand_i := -1
    for i, el := range game.Hands {
        if(el.User_ID == data.GameData.User_id) {
            player_hand_i = i
        }
    }

    found := false
    for _, el := range game.Hands[player_hand_i].Hand {
        if (el.Suit == data.ChosenCard.Suit && el.Value == data.ChosenCard.Value) {
            found = true
            break
        }
    }

    if(found == false) {
        c.String(http.StatusConflict, "This card is not yours")
    }

    err := game.play_card(data.ChosenCard)
    if(err != "") {
        c.String(http.StatusConflict, err)
    }
}

type StringChoice struct {
    Choices []string `json:"choices"`
    GameData GameAuth `json:"game_auth"`
}

func parse_choice(c *gin.Context) {
    var data StringChoice

    if err := c.BindJSON(&data); err != nil {
        return
    }

    game := find_game(data.GameData.Game_id)
    if (game == nil) {
        c.String(400, "Game not found");
        return;
    }

    user := find_user(data.GameData.User_id)
    if (user == nil) {
        c.String(400, "User not found");
        return;
    }

    player_i := -1;
    for i, u := range game.Users {
        if u.ID == user.ID {
            player_i = i;
            break
        }
    }

    if(player_i == -1) {
        c.String(http.StatusConflict, "You are not in this game")
        return
    }

    if(player_i != game.state.cur_player_index) {
        c.String(http.StatusConflict, "It is not your turn")
        return
    }

    err := game.play_choice(data.Choices)

    if(err != "") {
        c.String(http.StatusConflict, err)
    } 
}

func main() {
    rand.Seed(time.Now().UnixNano())

    router := gin.Default()
    
    router.Use(cors.Default())

    router.GET("/games", list_games)
    router.GET("/users", list_users)
    // TODO(Sigull): Add option to fetch only a single user
    router.GET("/join/:name", new_user)

    router.POST("/games/join", join_game)
    router.POST("/games/start", start_game)
    router.POST("/games/poll", poll_game)
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
