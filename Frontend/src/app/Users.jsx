'use client';
 
import React from 'react';
import Money from './Money.jsx';
import { get_card_src, game_type_to_string } from './helper';

export default function Users({ game, window_size, cur_index }) {

    if(!game.users) {
        return (<div className="flex justify-between"></div>)
    }

    const reorderedUsers = [
        ...game.users.slice((cur_index - 1) % game.users.length).map((user, index) => ({ user, index: (cur_index - 1 + index) % game.users.length })),
        ...game.users.slice(0, (cur_index - 1) % game.users.length).map((user, index) => ({ user, index: index })),
    ];

  return (
    <div className="flex justify-between">
        {reorderedUsers.map((item, index) => (
            <div key={index} className="flex flex-col">
                { game.state.att_index == parseInt(item.index) &&
                    <div>
                        <div className="text-3xl text-black underline">{item.user.name}</div>
                        <Money amount={item.user.money}></Money>
                        <div>{game_type_to_string(game)}</div>
                    </div>
                }
                { game.state.att_index != parseInt(item.index) &&
                    <div>
                        <div className="text-3xl text-black">{item.user.name}</div>
                        <Money amount={item.user.money}></Money>
                    </div>
                }

                { game.hlas && game.hlas[parseInt(item.index)]?.hand?.map((item, index) => (
                    <div key={index}>
                        <img className="" 
                            src={get_card_src(item)} width={window_size.width/25} id={item.suit + " " + item.spec}>
                        </img>
                    </div>
                ))}
            </div>
        ))}
    </div>
  )
}


