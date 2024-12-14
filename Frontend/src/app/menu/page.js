"use client";

import { useRouter, useSearchParams } from 'next/navigation';
import React, { useState, useEffect } from 'react';
import Money from '../Money.jsx'
import { generate_user_auth, butt_style } from '../helper.js';

export default function MainPage() {
  const router = useRouter();
  const search_params = useSearchParams();
  var temp_nickname = search_params.get("nickname");
  var temp_token = search_params.get("token");
  var temp_id = parseInt(search_params.get("id"));

  const [nickname, set_nickname] = useState(temp_nickname);
  const [user_id, set_user_id] = useState(temp_id);
  const [token, set_token] = useState(temp_token);
  const [balance, set_balance] = useState(0);
  const [games, set_games] = useState([]);
  const [users, set_users] = useState([]);

  const create_game = async () => {
    const response = await fetch(`http://localhost:8080/create`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      }
    });

    var parsed_response = await response.json();

    fetch_games();
  }

  const fetch_games = async () => {
    const response = await fetch(`http://localhost:8080/games`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      }
    });

    var parsed_response = await response.json();

    set_games(parsed_response);
  }

  const fetch_balance = async () => {
    const response = await fetch(`http://localhost:8080/users`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      }
    });

    var parsed_response = await response.json();

    for(var user of parsed_response) {
      if(user.id == user_id) {
        set_balance(user.money);
        break;
      }
    }
  }

  const fetch_users = async () => {
    const response = await fetch(`http://localhost:8080/users`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      }
    });

    var parsed_response = await response.json();

    set_users(parsed_response);
  }

  const join_game = async (game_id) => {
    const response = await fetch(`http://localhost:8080/games/join`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(generate_user_auth(user_id, game_id, token))
    });

    if (response.ok || await response.text() == "Game is full") {
      localStorage.setItem('token', token);
      localStorage.setItem('user_id', user_id);
      localStorage.setItem('nickname', nickname);
      router.push(`/game?game_id=${game_id}`);
    }
  }

  useEffect(() =>{
    if(!user_id) {
       set_user_id(parseInt(localStorage.getItem('user_id')));
       set_token(localStorage.getItem('token'));
       set_nickname(localStorage.getItem('nickname'));
   }
  }, [user_id]);

  useEffect(() => {
    fetch_users();
    fetch_games();
    fetch_balance();
  }, []);

  return (
    <div className="bg-background h-screen" style={{ padding: '2rem', textAlign: 'center' }}>
      <div className="flex flex-col">
        <div className="flex flex-wrap">
          <div className="flex flex-col">
            <h1 className="text-black text-2xl">{nickname}</h1>
            <h2 className="text-white text-left">
              <Money amount={balance}></Money>
            </h2>
          </div>
        </div>
        <div className="flex flex-wrap justify-between justify-items-center pt-20">
          <div className="flex flex-col">
            <h1 className="text-black text-left text-4xl pb-4">Otevřené hry</h1>
            <table className="min-w-full text-left border-collapse bg-white shadow-md">
              <thead>
                <tr className="bg-gray-200 text-black uppercase text-sm leading-normal">
                    <td className="">ID</td>
                    <td className="pl-5">Počet hráčů</td>
                    <td className="pl-5">Zapnuta</td>
                    <td className="pl-5">Stav</td>
                    <td className="pl-5"></td>
                </tr>
              </thead>
              <tbody>
                {games.map((item, index) => (
                    <tr key={index} className="text-black border-gray-200 hover:bg-gray-100">
                        <td>{item.id}</td>
                        <td className="pl-5">{item.users.length}/3</td>
                        <td className="pl-5">{item.running ? 'Ano' : 'Ne'}</td>
                        <td className="pl-5">{item.game_type || 'Čekání'}</td>
                        <td className="pl-5 pr-5">
                          <button className="pl-1 pr-1 hover:bg-gray-500" onClick={() => join_game(item.id)}> 
                            {item.users.some(item => item.id === user_id) == true? "Jsi v ní" : item.users.length >= 3 ? 'Plno' : 'Připojit'}
                          </button>
                        </td>
                    </tr>
                ))}
              </tbody>
            </table>
            <button className={butt_style} onClick={create_game} >Nová hra</button>
            <button className={butt_style} onClick={fetch_games}>Přenačíst</button>
          </div>

          <div className="flex items-end flex-col">
            <h1 className="text-black text-left text-4xl pb-4">Aktivní uživatelé</h1>
            <table className="min-w-full text-left border-collapse bg-white shadow-md">
              <thead>
                <tr className="bg-gray-200 text-black uppercase text-sm leading-normal">
                    <td className="">Přezdívka</td>
                </tr>
              </thead>
              <tbody>
                {users.map((item, index) => (
                    <tr key={index} className="text-black border-gray-200 hover:bg-gray-100">
                        <td>{item.name}</td>
                    </tr>
                ))}
              </tbody>
            </table>
            <button className={butt_style} onClick={fetch_users}>Přenačíst</button>
          </div>
        </div>
        <h1>Welcome, {nickname}</h1>
        <p>This is the main page.</p>
      </div>
    </div>
  );
}
