"use client";

import { useRouter, useSearchParams } from 'next/navigation';
import React, { useState, useEffect } from 'react';
import Money from '../Money.jsx';
import { generate_user_auth } from '../helper.js';

export default function MainPage() {
    const router = useRouter();
    const search_params = useSearchParams();
    var game_id = parseInt(search_params.get("game_id"));
    const [game, set_game] = useState({});
    const [user_id, set_user_id] = useState(0);
    const [token, set_token] = useState('');

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
                set_game(g);
            }
        }
    }

    const get_card_src = (card_val) => {
        return `/cards/${card_val.suit}-${card_val.value}.png`
    }

    useEffect(() => {
        set_user_id(parseInt(localStorage.getItem('user_id')));
        set_token(localStorage.getItem('token'));

        fetch_game_data();
    }, []);

    const start_game = async () => {
        const response = await fetch(`http://localhost:8080/games/start`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(generate_user_auth(4, game_id, token))
        });
    
        // var parsed_response = await response.json();

        if (response.ok) {
            fetch_game_data();
        }
    }

    return(
        <div className="bg-background h-screen" style={{ padding: '2rem', textAlign: 'center' }}>
            <div className="flex h-full flex-col justify-between">
                <div className="flex justify-between">
                    {game.users?.map((item, index) => (
                        <div key={index} className="flex flex-col">
                            <div className="text-3xl text-black">{item.name}</div>
                            <Money amount={item.money}></Money>
                        </div>
                    ))}         
                </div>

                {game.running == false &&
                    <div>
                        <button onClick={start_game}>Start game</button>
                    </div>
                }

                <div className="flex justify-between">
                    <div className="flex">
                        {game.hands && game.hands[0].hand.map((item, index) => (
                            <img key={index} src={get_card_src(item)} width={window.innerWidth/20} style={{position: 'absolute', left: `${index * 4}%`, bottom:'5%'}}></img>
                        ))}
                    </div>  
                </div>
            </div>
        </div>
    );
}