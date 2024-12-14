package main

import (
	"encoding/base64"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"github.com/gin-contrib/cors"
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
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
	Money int    `json:"money"`
}

var users = []User{{ID: 0, Name: "test_user", Token: generateRandomToken(), Money: 69420}}

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
	state         ChoiceStateEnum
	game_flek     bool
	seven_flek    bool
	hundred_flek  bool
	betting_round int
}

type GameStateAll struct {
	game_state   GameStateEnum
	choice_state ChoiceState
	game_type    GameTypeEnum

	game_bets    int
	seven_bets   int
	hundred_bets int

	won_seven        bool
	Cur_player_index int `json:"cur_index"`
	Attacker_index   int `json:"att_index"`
}

type Card struct {
	Suit         int `json:"suit"`
	Value        int `json:"value"`
	SpecialValue int `json:"spec"`
}

var default_deck = []Card{
	{Suit: 0, Value: 7, SpecialValue: 7}, {Suit: 0, Value: 8, SpecialValue: 8}, {Suit: 0, Value: 9, SpecialValue: 9}, {Suit: 0, Value: 14, SpecialValue: 10},
	{Suit: 0, Value: 11, SpecialValue: 11}, {Suit: 0, Value: 12, SpecialValue: 12}, {Suit: 0, Value: 13, SpecialValue: 13}, {Suit: 0, Value: 15, SpecialValue: 14},

	{Suit: 1, Value: 7, SpecialValue: 7}, {Suit: 1, Value: 8, SpecialValue: 8}, {Suit: 1, Value: 9, SpecialValue: 9}, {Suit: 1, Value: 14, SpecialValue: 10},
	{Suit: 1, Value: 11, SpecialValue: 11}, {Suit: 1, Value: 12, SpecialValue: 12}, {Suit: 1, Value: 13, SpecialValue: 13}, {Suit: 1, Value: 15, SpecialValue: 14},

	{Suit: 2, Value: 7, SpecialValue: 7}, {Suit: 2, Value: 8, SpecialValue: 8}, {Suit: 2, Value: 9, SpecialValue: 9}, {Suit: 2, Value: 14, SpecialValue: 10},
	{Suit: 2, Value: 11, SpecialValue: 11}, {Suit: 2, Value: 12, SpecialValue: 12}, {Suit: 2, Value: 13, SpecialValue: 13}, {Suit: 2, Value: 15, SpecialValue: 14},

	{Suit: 3, Value: 7, SpecialValue: 7}, {Suit: 3, Value: 8, SpecialValue: 8}, {Suit: 3, Value: 9, SpecialValue: 9}, {Suit: 3, Value: 14, SpecialValue: 10},
	{Suit: 3, Value: 11, SpecialValue: 11}, {Suit: 3, Value: 12, SpecialValue: 12}, {Suit: 3, Value: 13, SpecialValue: 13}, {Suit: 3, Value: 15, SpecialValue: 14},
}

type UsersDeck struct {
	User_ID int    `json:"id"`
	Hand    []Card `json:"hand"`
}

type Game struct {
	Users       []*User      `json:"users"`
	Hands       []UsersDeck  `json:"hands"`
	Hlas        []UsersDeck  `json:"hlas"`
	Table       []Card       `json:"table"`
	User_scores []int        `json:"scores"`
	Talons      []UsersDeck  `json:"talons"`
	Last_talon  []Card       `json:"talon"`
	ID          int          `json:"id"`
	Running     bool         `json:"running"`
	Betting     bool         `json:"betting"`
	Game_type   string       `json:"game_type"`
	Trump_card  Card         `json:"trump_card"`
	Trump_suit  int          `json:"trump_suit"`
	State       GameStateAll `json:"state"`
	temp_deck   []Card
}

func (self *Game) init() {
	self.Hands = []UsersDeck{}
	self.temp_deck = make([]Card, len(default_deck))
	copy(self.temp_deck, default_deck)

	rand.Shuffle(len(self.temp_deck), func(i, j int) {
		self.temp_deck[i], self.temp_deck[j] = self.temp_deck[j], self.temp_deck[i]
	})

	self.Hands = make([]UsersDeck, 3)
	self.Hlas = make([]UsersDeck, 3)
	self.Talons = make([]UsersDeck, 3)
	self.User_scores = make([]int, 3)
	self.Last_talon = []Card{}

	self.State.Attacker_index = (self.State.Attacker_index + 1) % 3
	self.State.Cur_player_index = self.State.Attacker_index

	var deal_arr = [][]int{{5, 12}, {12, 22}, {22, 32}}
	for i := range self.Users {
		self.User_scores[i] = 0
		var user_i = (i + self.State.Cur_player_index) % 3
		hand_copy := make([]Card, deal_arr[i][1]-deal_arr[i][0])
		copy(hand_copy, self.temp_deck[deal_arr[i][0]:deal_arr[i][1]])
		self.Talons[user_i] = UsersDeck{User_ID: self.Users[user_i].ID, Hand: []Card{}}
		self.Hands[user_i] = UsersDeck{User_ID: self.Users[user_i].ID, Hand: hand_copy}

		sort.Slice(self.Hands[user_i].Hand, func(i, j int) bool {
			if self.Hands[user_i].Hand[i].Suit == self.Hands[user_i].Hand[j].Suit {
				return self.Hands[user_i].Hand[i].Value < self.Hands[user_i].Hand[j].Value
			}
			return self.Hands[user_i].Hand[i].Suit < self.Hands[user_i].Hand[j].Suit
		})
	}

	self.State.won_seven = false
	self.Game_type = ""
	self.Trump_card = Card{}
	self.Running = true
	self.Betting = false
	self.State.game_state = STATE_CHOICE
	self.State.choice_state.state = STATE_CHOOSING_TRUMP
	self.State.choice_state.game_flek = false
	self.State.choice_state.seven_flek = false
	self.State.choice_state.hundred_flek = false
	self.State.choice_state.betting_round = 0
	self.State.game_bets = 0
	self.State.seven_bets = 0
	self.State.hundred_bets = 0
}

func (self *Game) next_player() {
	self.State.Cur_player_index = (self.State.Cur_player_index + 1) % 3
}

func (self *Game) get_choices() []string {
	state := &self.State

	if !self.Running {
		return []string{}
	}

	switch state.game_state {
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
			if state.choice_state.game_flek == false &&
				(state.game_type == TYPE_GAME ||
					state.game_type == TYPE_SEVEN) {
				bets = append(bets, "on game")
			}
			if state.choice_state.seven_flek == false &&
				(state.game_type == TYPE_SEVEN ||
					state.game_type == TYPE_HUNDRED_SEVEN) {
				bets = append(bets, "on seven")
			}
			if state.choice_state.hundred_flek == false &&
				(state.game_type == TYPE_HUNDRED ||
					state.game_type == TYPE_HUNDRED_SEVEN) {
				bets = append(bets, "on hundred")
			}

			if len(bets) == 0 {
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
	{Users: []*User{}, ID: 0},
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
	return nil
}

func find_user(id int) *User {
	for _, u := range users {
		if u.ID == id {
			return &u
		}
	}
	return nil
}

type GameAuth struct {
	User_id    int    `json:"user_id"`
	Game_id    int    `json:"game_id"`
	User_token string `json:"token"`
}

func start_game(c *gin.Context) {
	var data GameAuth

	if err := c.BindJSON(&data); err != nil {
		return
	}

	game := find_game(data.Game_id)
	if game == nil {
		c.String(400, "Game not found")
		return
	}

	user := find_user(data.User_id)
	if user == nil {
		c.String(400, "User not found")
		return
	}

	if game.Users[0].ID != user.ID {
		c.String(400, "You are not the game owner")
		return
	}

	if len(game.Users) != 3 {
		c.String(http.StatusConflict, "Not enough players")
		return
	}

	if game.Running != true {
		game.Running = true
	} else {
		c.String(http.StatusConflict, "Game is already running")
		return
	}

	game.State.Attacker_index = -1
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
	if game == nil {
		c.String(400, "Game not found")
		return
	}

	user := find_user(data.User_id)
	if user == nil {
		c.String(400, "User not found")
		return
	}

	if len(game.Users) >= 3 {
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

func new_game(c *gin.Context) {
	last_id := games[len(games)-1].ID

	game := Game{Users: []*User{}, ID: last_id + 1}

	games = append(games, game)

	c.IndentedJSON(http.StatusOK, game)
}

func new_user(c *gin.Context) {
	name := c.Param("name")

	last_id := users[len(users)-1].ID

	user := User{ID: last_id + 1, Name: name, Token: generateRandomToken(), Money: 0}

	users = append(users, user)

	c.IndentedJSON(http.StatusOK, user)
}

type PollData struct {
	User_choices []string `json:"choices"`
	Game_data    Game     `json:"game"`
}

func poll_game(c *gin.Context) {
	var data GameAuth

	if err := c.BindJSON(&data); err != nil {
		return
	}

	game := find_game(data.Game_id)
	if game == nil {
		c.String(400, "Game not found")
		return
	}

	user := find_user(data.User_id)
	if user == nil {
		c.String(400, "User not found")
		return
	}

	found := false
	for _, u := range game.Users {
		if u.ID == user.ID {
			found = true
			break
		}
	}

	if found == false {
		c.String(http.StatusConflict, "You are not in this game")
		return
	}

	var return_val PollData = PollData{User_choices: nil, Game_data: *game}

	if game.Users[game.State.Cur_player_index].ID == user.ID {
		return_val.User_choices = game.get_choices()
	}

	c.IndentedJSON(http.StatusOK, return_val)
}

func (self *Game) no_flek() bool {
	if self.State.game_type == TYPE_GAME {
		return self.State.choice_state.game_flek == false

	} else if self.State.game_type == TYPE_SEVEN {
		return self.State.choice_state.game_flek == false &&
			self.State.choice_state.seven_flek == false

	} else if self.State.game_type == TYPE_HUNDRED {
		return self.State.choice_state.hundred_flek == false

	} else if self.State.game_type == TYPE_HUNDRED_SEVEN {
		return self.State.choice_state.hundred_flek == false &&
			self.State.choice_state.seven_flek == false
	}

	return false
}

func (self *Game) set_fleks() {
	self.State.choice_state.game_flek = true
	self.State.choice_state.seven_flek = true
	self.State.choice_state.hundred_flek = true
}

func (self *Game) play_choice(choices []string) string {
	state := &self.State

	switch state.game_state {
	case STATE_CHOICE:
		switch state.choice_state.state {
		case STATE_CHOOSING_TRUMP:
			return "You have to choose a card"
		case STATE_THROW_AWAY1:
			return "You have to choose a card"
		case STATE_THROW_AWAY2:
			return "You have to choose a card"

		case STATE_CHOOSING_GAME:
			if len(choices) != 1 {
				return "Bad number of choices"
			}

			switch choices[0] {
			case "game":
				state.game_type = TYPE_GAME
				state.game_bets += 1
				self.Game_type = "game"
			case "seven":
				state.game_type = TYPE_SEVEN
				state.game_bets += 1
				state.seven_bets += 1
				self.Game_type = "seven"
			case "hundred":
				state.game_type = TYPE_HUNDRED
				state.hundred_bets += 1
				self.Game_type = "hundred"
			case "hundred seven":
				state.game_type = TYPE_HUNDRED_SEVEN
				state.hundred_bets += 1
				state.seven_bets += 1
				self.Game_type = "hundred_seven"

			default:
				return "Not a valid choice"
			}
			self.Betting = true
			state.choice_state.state = STATE_BETS
			self.next_player()

		case STATE_BETS:
			if self.State.Attacker_index == self.State.Cur_player_index {
				game_f := !state.choice_state.game_flek
				seven_f := !state.choice_state.seven_flek
				hundred_f := !state.choice_state.hundred_flek

				self.set_fleks()

				if len(choices) == 0 {
					self.Betting = false
					self.Trump_card = Card{-1, -1, -1}
					state.game_state = STATE_GAME
					return ""
				}

				for _, c := range choices {
					switch c {
					case "on game":
						if game_f {
							state.game_bets += 1
							state.choice_state.game_flek = false
						} else {
							return "Not a valid choice"
						}
					case "on seven":
						if seven_f {
							state.seven_bets += 1
							state.choice_state.seven_flek = false
						} else {
							return "Not a valid choice"
						}
					case "on hundred":
						if hundred_f {
							state.hundred_bets += 1
							state.choice_state.hundred_flek = false
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
						state.game_bets += 1
						state.choice_state.game_flek = true
					case "on seven":
						state.seven_bets += 1
						state.choice_state.seven_flek = true
					case "on hundred":
						state.hundred_bets += 1
						state.choice_state.hundred_flek = true

					default:
						return "Not a valid choice"
					}
				}
			}

			self.next_player()
			if self.State.Attacker_index == self.State.Cur_player_index {
				if self.no_flek() {
					self.Betting = false
					self.Trump_card = Card{-1, -1, -1}
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
		if choices[0] == "next game" {
			self.init()
		}
	}

	return ""
}

func (self *Game) clear_table() {

}

func (self *Game) compare_cards(c1 Card, c2 Card) bool {
	if c1.Suit == c2.Suit {
		return c1.Value > c2.Value

	} else if c1.Suit == self.Trump_suit {
		return true
	}

	return false
}

func (self *Game) hand_eval() int {
	next_player := (self.State.Cur_player_index + 1) % 3
	max_card := self.Table[0]

	for i, _ := range self.Table {
		cur_card := self.Table[i]
		if self.compare_cards(cur_card, max_card) {
			next_player = (self.State.Cur_player_index + i + 1) % 3
			max_card = cur_card
		}
	}

	if len(self.Hands[0].Hand) == 0 &&
		(self.State.game_type == TYPE_SEVEN ||
			self.State.game_type == TYPE_HUNDRED_SEVEN) {

		if max_card.Value == 7 && max_card.Suit == self.Trump_suit &&
			next_player == self.State.Attacker_index {
			self.State.won_seven = true

		} else {
			self.State.won_seven = false
		}
	}

	return next_player
}

func (self *Game) end() bool {
	return len(self.Hands[0].Hand) == 0
}

func (self *Game) next_round() string {
	if len(self.Table) != 3 {
		return "Not enough cards for next turn"
	}

	p_won_round := self.hand_eval()

	self.Talons[p_won_round].Hand = append(self.Talons[p_won_round].Hand, self.Table...)
	self.Last_talon = make([]Card, len(self.Table))
	copy(self.Last_talon, self.Table)
	self.Table = []Card{}

	score := 0
	for _, el := range self.Last_talon {
		if el.SpecialValue == 10 || el.SpecialValue == 14 {
			score += 10
		}
	}

	self.User_scores[p_won_round] += score

	if self.end() {
		self.User_scores[p_won_round] += 10
		return "end"
	}

	self.State.Cur_player_index = p_won_round

	return ""
}

func (self *Game) attacker_draw() {
	att_index := self.State.Attacker_index
	self.Hands[att_index].Hand = append(self.Hands[att_index].Hand, self.temp_deck[0:5]...)

	sort.Slice(self.Hands[att_index].Hand, func(i, j int) bool {
		if self.Hands[att_index].Hand[i].Suit == self.Hands[att_index].Hand[j].Suit {
			return self.Hands[att_index].Hand[i].Value < self.Hands[att_index].Hand[j].Value
		}
		return self.Hands[att_index].Hand[i].Suit < self.Hands[att_index].Hand[j].Suit
	})
}

func (self *Game) can_place(card Card) string {
	card_to_compare := Card{-1, -1, -1}
	trump_card := Card{-1, -1, -1}

	if len(self.Table) == 0 {
		return ""

	} else if len(self.Table) == 1 {
		card_to_compare = self.Table[0]

	} else if len(self.Table) == 2 {
		if self.Table[0].Suit == self.Table[1].Suit {
			if self.Table[0].Value > self.Table[1].Value {
				card_to_compare = self.Table[0]
			} else {
				card_to_compare = self.Table[1]
			}

		} else if self.Table[1].Suit == self.Trump_suit {
			trump_card = self.Table[1]
			card_to_compare = self.Table[0]

		} else {
			card_to_compare = self.Table[0]
		}

	} else {
		return "Somehow there are too many cards"
	}

	has_comp_suit := false
	has_to_higher := false
	has_trump := false
	highest_trump_val := 0

	fmt.Println(card_to_compare)

	for _, el := range self.Hands[self.State.Cur_player_index].Hand {
		fmt.Println(el)
		if el.Suit == card_to_compare.Suit {
			if el.Value > card_to_compare.Value && trump_card.Suit == -1 {
				has_to_higher = true
			}

			has_comp_suit = true
		}

		if el.Suit == self.Trump_suit {
			has_trump = true
			if highest_trump_val < el.Value {
				highest_trump_val = el.Value
			}
		}
	}

	if has_comp_suit {
		if card.Suit != card_to_compare.Suit {
			return "Not correct suit"
		}

		if has_to_higher && card.Value < card_to_compare.Value {
			return "Play bigger"
		}

	} else {
		if has_trump {
			if card.Suit != self.Trump_suit {
				return "Play trump"
			}

			if highest_trump_val > trump_card.Value && trump_card.Value > card.Value {
				return "Play bigger"
			}
		}
	}

	return ""
}

func (self *Game) is_hlas(card Card) string {
	p_index := self.State.Cur_player_index

	is_king := card.Value == 13
	is_queen := card.Value == 12

	if !(is_king || is_queen) {
		return "not"
	}

	for _, el := range self.Hands[p_index].Hand {
		if is_king && el.Value == 12 && el.Suit == card.Suit {
			return "play queen"
		}
		if is_queen && el.Value == 13 && el.Suit == card.Suit {
			return "yes"
		}
	}

	return "not"
}

func (self *Game) place_card(card Card) string {
	p_index := self.State.Cur_player_index
	found := false
	found_i := 0
	for i, el := range self.Hands[p_index].Hand {
		if el.Suit == card.Suit && el.Value == card.Value {
			found = true
			found_i = i
			break
		}
	}

	if found == false {
		return "This card is not yours"
	}

	can := self.can_place(card)

	if can != "" {
		return can
	}

	hlas_res := self.is_hlas(card)

	if hlas_res == "yes" {
		self.Hlas[p_index].Hand = append(self.Hlas[p_index].Hand, card)
		if card.Suit == self.Trump_suit {
			self.User_scores[p_index] += 40

		} else {
			self.User_scores[p_index] += 20
		}

	} else if hlas_res != "not" {
		return hlas_res
	}

	self.Table = append(self.Table, card)

	self.Hands[p_index].Hand = append(self.Hands[p_index].Hand[:found_i], self.Hands[p_index].Hand[found_i+1:]...)

	if len(self.Table) == 3 {
		return "next"
	}

	return ""
}

func (self *Game) card_throw_away(card Card) string {
	found := false
	att_index := self.State.Attacker_index
	for i, el := range self.Hands[att_index].Hand {
		if el.Suit == card.Suit && el.Value == card.Value {
			found = true
			self.Hands[att_index].Hand = append(self.Hands[att_index].Hand[:i], self.Hands[att_index].Hand[i+1:]...)
			break
		}
	}

	if found == false {
		return "This card is not yours"
	}

	self.Talons[(att_index+1)%3].Hand = append(self.Talons[(att_index+1)%3].Hand, card)

	return ""
}

func (self *Game) game_eval() {
	att_index := self.State.Attacker_index

	won_seven := self.State.won_seven

	att_won := false

	att_score := 0
	def_score := 0
	for i, el := range self.User_scores {
		if i == att_index {
			att_score += el

		} else {
			def_score += el
		}
	}

	if att_score > def_score {
		att_won = true
	}

	above_hundred := (att_score - 100) / 10
	if above_hundred >= 0 {
		above_hundred += 1
	}

	def_above_hundred := (def_score-100)/10 + 1

	var att *User
	def := make([]*User, 2)
	def_i := 0

	for i, el := range self.Users {
		if i == self.State.Attacker_index {
			att = el
		} else {
			def[def_i] = el
			def_i++
		}
	}

	bet := 0

	fmt.Println(self.State.seven_bets)

	switch self.State.game_type {
	case TYPE_GAME:
		if above_hundred > 0 {
			bet += above_hundred * 2
		}
		if def_above_hundred > 0 {
			bet -= above_hundred * 2
		}

		if att_won {
			bet = int(math.Pow(2, float64(self.State.game_bets-1)))
		} else {
			bet = -int(math.Pow(2, float64(self.State.game_bets-1)))
		}

	case TYPE_SEVEN:
		if above_hundred > 0 {
			bet += above_hundred * 2
		}
		if def_above_hundred > 0 {
			bet -= above_hundred * 2
		}

		if att_won {
			bet += int(math.Pow(2, float64(self.State.game_bets-1)))
		} else {
			bet -= int(math.Pow(2, float64(self.State.game_bets-1)))
		}
		if won_seven {
			bet += 4 * int(math.Pow(2, float64(self.State.seven_bets-1)))
		} else {
			bet -= 4 * int(math.Pow(2, float64(self.State.seven_bets-1)))
		}

	case TYPE_HUNDRED:
		bet += above_hundred * 4

	case TYPE_HUNDRED_SEVEN:
		bet += above_hundred * 4
		if won_seven {
			bet += 4 * int(math.Pow(2, float64(self.State.seven_bets-1)))
		} else {
			bet -= 4 * int(math.Pow(2, float64(self.State.seven_bets-1)))
		}
	}

	// Hearts double
	if self.Trump_suit == 2 {
		bet *= 2
	}

	att.Money += bet * 2
	def[0].Money -= bet
	def[1].Money -= bet
}

func (self *Game) play_card(card Card) string {
	state := &self.State

	switch state.game_state {
	case STATE_CHOICE:
		switch state.choice_state.state {
		case STATE_CHOOSING_TRUMP:
			self.Trump_card = card
			self.Trump_suit = card.Suit
			state.choice_state.state = STATE_THROW_AWAY1
			self.attacker_draw()

		case STATE_THROW_AWAY1:
			err := self.card_throw_away(card)
			if err != "" {
				return err
			}
			state.choice_state.state = STATE_THROW_AWAY2

		case STATE_THROW_AWAY2:
			err := self.card_throw_away(card)
			if err != "" {
				return err
			}
			state.choice_state.state = STATE_CHOOSING_GAME

		case STATE_CHOOSING_GAME:
			return "You can't choose a card here"
		case STATE_BETS:
			return "You can't choose a card here"
		}

	case STATE_GAME:
		res := self.place_card(card)
		if res == "next" {
			res = self.next_round()
			if res == "end" {
				self.game_eval()
				self.State.game_state = STATE_END
				self.State.Cur_player_index = 0
				return ""
			}

			return res

		} else if res == "" {
			self.next_player()

		} else {
			return res
		}

	case STATE_END:
		return "next game"
	}

	return ""
}

type CardChoice struct {
	ChosenCard Card     `json:"card"`
	GameData   GameAuth `json:"game_auth"`
}

func parse_card(c *gin.Context) {
	var data CardChoice

	if err := c.BindJSON(&data); err != nil {
		return
	}

	game := find_game(data.GameData.Game_id)
	if game == nil {
		c.String(400, "Game not found")
		return
	}

	user := find_user(data.GameData.User_id)
	if user == nil {
		c.String(400, "User not found")
		return
	}

	player_i := -1
	for i, u := range game.Users {
		if u.ID == user.ID {
			player_i = i
			break
		}
	}

	if player_i == -1 {
		c.String(http.StatusConflict, "You are not in this game")
		return
	}

	if player_i != game.State.Cur_player_index {
		c.String(http.StatusConflict, "It is not your turn")
		return
	}

	player_hand_i := -1
	for i, el := range game.Hands {
		if el.User_ID == data.GameData.User_id {
			player_hand_i = i
		}
	}

	fmt.Println(player_hand_i)

	found := false
	card := data.ChosenCard
	for _, el := range game.Hands[player_hand_i].Hand {
		if el.Suit == card.Suit && el.Value == card.Value {
			card = el
			found = true
			break
		}
	}

	if found == false {
		c.String(http.StatusConflict, "This card is not yours")
	}

	err := game.play_card(card)
	if err != "" {
		c.String(http.StatusConflict, err)
	}
}

type StringChoice struct {
	Choices  []string `json:"choices"`
	GameData GameAuth `json:"game_auth"`
}

func parse_choice(c *gin.Context) {
	var data StringChoice

	if err := c.BindJSON(&data); err != nil {
		return
	}

	game := find_game(data.GameData.Game_id)
	if game == nil {
		c.String(400, "Game not found")
		return
	}

	user := find_user(data.GameData.User_id)
	if user == nil {
		c.String(400, "User not found")
		return
	}

	player_i := -1
	for i, u := range game.Users {
		if u.ID == user.ID {
			player_i = i
			break
		}
	}

	if player_i == -1 {
		c.String(http.StatusConflict, "You are not in this game")
		return
	}

	if player_i != game.State.Cur_player_index {
		c.String(http.StatusConflict, "It is not your turn")
		return
	}

	err := game.play_choice(data.Choices)

	if err != "" {
		c.String(http.StatusConflict, err)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/games", list_games)
	router.GET("/users", list_users)
	router.GET("/join/:name", new_user)
	router.POST("/create", new_game)

	router.POST("/games/join", join_game)
	router.POST("/games/start", start_game)
	router.POST("/games/poll", poll_game)
	router.POST("/games/play/choice", parse_choice)
	router.POST("/games/play/card", parse_card)

	router.Run("localhost:8080")
}
