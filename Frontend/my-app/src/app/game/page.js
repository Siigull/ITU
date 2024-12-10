"use client";

import { useRouter, useSearchParams } from 'next/navigation';
import React, { useState, useEffect } from 'react';
import Table from '../Table.jsx';
import Users from '../Users.jsx';
import use_window_size, { generate_user_auth, game_type_to_string, butt_style, get_card_src } from '../helper.js';

export default function MainPage() {
    const router = useRouter();
    const search_params = useSearchParams();
    var game_id = parseInt(search_params.get("game_id"));
    const [game, set_game] = useState({});
    const [user_id, set_user_id] = useState(-1);
    const [user_index, set_user_index] = useState(0);
    const [token, set_token] = useState('');
    const [window_width, set_windows_width] = useState(1920);
    const [choices, set_choices] = useState({});
    const [chosen_choices, set_chosen_choices] = useState({});
    const window_size = use_window_size();
    var interval;

    const send_choices = async (choice) => {
        let temp_choices = []
        for(var ch in chosen_choices) {
            if(chosen_choices[ch]) {
                temp_choices.push(ch)
            }
        }
        
        const response = await fetch(`http://localhost:8080/games/play/choice`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ game_auth:generate_user_auth(user_id, game_id, token),
                                   choices: temp_choices
            })
        });
    
        // var parsed_response = await response.json();

        if(response.ok) {
            set_chosen_choices([])
            console.log("Hooray");
        }
    }

    const set_choice = async (choice, id) => {
        let chosen = chosen_choices;
        if(chosen[choice]) {
            chosen[choice] = !chosen[choice]
        } else {
            chosen[choice] = true;
        }
        set_chosen_choices(chosen);

        const element = document.getElementById(id);
        var color = element.style.color;
        
        if(color == 'grey') {
            color = 'black';
        } else {
            color = 'grey';
        }

        element.style.color = color;
    }

    const send_card = async (suit, value) => {
        const response = await fetch(`http://localhost:8080/games/play/card`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ game_auth:generate_user_auth(user_id, game_id, token),
                                   card:{suit:suit, value:value}
            })
        });
    
        // var parsed_response = await response.json();

        if(response.ok) {
            console.log("Hooray");
        }
    }

    const poll_game = async () => {
        const response = await fetch(`http://localhost:8080/games/poll`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(generate_user_auth(user_id, game_id, token))
        });
    
        var parsed_response = await response.json();

        set_choices({});
        if (parsed_response.choices != null) {
            var temp_choices = {card: false, choices:[]}
            for(var ch of parsed_response.choices){
                if(ch == "card") {
                    temp_choices.card = true;
                } else {
                    temp_choices.choices.push(ch);
                }
            }
            set_choices(temp_choices);
        }

        set_game(parsed_response.game)
    }

    const start_interval = () => {
        interval = setInterval(() => poll_game(), 1000);
    }

    const fetch_game_data = async () => {
        const response = await fetch(`http://localhost:8080/games`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            }
        });
    
        var parsed_response = await response.json();

        for(var g of parsed_response) {
            if(game_id == g.id) {
                for(var index in g.users) {
                    console.log(user_id)
                    console.log(parseInt(g.users[parseInt(index)].id))
                    if(g.users[parseInt(index)].id == user_id) {
                        set_user_index(index)
                    }
                }
                set_game(g);
            }
        }
    }

    useEffect(() => {
        if(token != '') {
            start_interval();
        }
        
        fetch_game_data();
    }, [token]);

    useEffect(() => {
        if (user_id == -1) {
            set_user_id(parseInt(localStorage.getItem('user_id')));
        }
        if (token == '') {
            set_token(localStorage.getItem('token'));
        }

        fetch_game_data();
    }, []);

    const start_game = async () => {
        console.log(user_id)

        const response = await fetch(`http://localhost:8080/games/start`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(generate_user_auth(user_id, game_id, token))
        });
    
        // var parsed_response = await response.json();

        if (response.ok) {
            fetch_game_data();
        }
    }

    const go_to_menu = () => {
        router.push(`/menu`);
    }

    return(
        <div className="bg-background h-screen" style={{ padding: '2rem', textAlign: 'left' }}>
            <div className="flex h-full flex-col justify-between">

                <Users game={game} window_size={window_size} cur_index={user_index}></Users>

                {game.running == false &&
                    <div>
                        <button className={butt_style} onClick={start_game}>Start game</button>
                    </div>
                }

                <div>
                    <div className="flex flex-col text-center">   
                        <div>
                            {game.running == true && choices.choices &&
                                choices.choices.map((item, index) => (
                                    <button className={butt_style + "pr-10"} key={index} style={{color:'grey'}} id={index} onClick={() => set_choice(item, index)}>{ item }</button>
                                ))
                            }
                        </div>
                        <div>
                            {game.running == true && choices.choices != undefined && choices.choices.length > 0 &&
                                <button className={butt_style}onClick={send_choices}>Potvrdit</button>
                            }
                        </div>
                        <div>
                            {choices.card == true &&
                                <h1>Play Card</h1>
                            }
                        </div>
                        
                        <Table game={game} window_size={window_size} cur_index={user_index}></Table>

                        <div className="pt-10" style={{display: 'flex', justifyContent:'center'}}>
                            {game.betting == true && 
                                <img className="object-center" 
                                    src={get_card_src(game.trump_card)} width={window_size.width/13} 
                                    id={game.trump_card.suit + " " + game.trump_card.value}>
                                </img>
                            }
                        </div>
                    </div>

                    <div className="flex text-left">
                        <button className={butt_style}onClick={go_to_menu}>Odejít</button>
                    </div>
                </div>

                <div className="flex justify-between w-full h-1/2">
                    <div className="flex flex-col w-1/3 justify-end">
                        <div className="bg-background_green relative rounded-3xl w-3/4 h-1/3 md:h-1/2 xl:h-2/3 lg:3/5 mt-8">
                            {game.talon && game.talon.map((item, index) => (
                                <div className="content-center" key={index}>
                                    <img className="h-3/4 w-auto object-cover" 
                                        src={get_card_src(item)} id={item.suit + " " + item.spec}
                                        style={{position: 'absolute', left: `${window_size.width /200 + index * 26}%`, bottom:"13%"}}
                                        onClick={() => send_card(item.suit, item.value)}>
                                    </img>
                                </div>
                            ))}
                        </div>
                    </div>
                    <div className="relative flex w-2/3">
                        {game.hands && game.hands[user_index].hand.map((item, index) => (
                            <div key={index}>
                                { game.trump_card.suit == item.suit && game.trump_card.value == item.value && game.game_type == "" &&
                                    <div>
                                        {game.betting == false &&
                                            <img className="rotate-90" 
                                                src={get_card_src(item)} width={window_size.width/13} id={item.suit + " " + item.spec}
                                                style={{position: 'absolute', left: `${index * 7}%`, bottom:'40%'}}
                                                onClick={() => send_card(item.suit, item.value)}>
                                            </img>
                                        }
                                    </div>
                                }
                                { ((game.trump_card.suit != item.suit || game.trump_card.value != item.value) || game.game_type != "") &&
                                    <img className="transform hover:-translate-y-14 transition-transform duration-100 ease-in-out" 
                                        src={get_card_src(item)} width={window_size.width/13} id={item.suit + " " + item.spec}
                                        style={{position: 'absolute', left: `${index * 7}%`, bottom:'5%'}}
                                        onClick={() => send_card(item.suit, item.value)}>
                                    </img>
                                }
                            </div>
                        ))}
                    </div>  
                </div>
            </div>
        </div>
    );
}