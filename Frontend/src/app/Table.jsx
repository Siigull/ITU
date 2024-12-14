'use client';
 
import React from 'react';
import { get_card_src } from './helper';

export default function Table({ game, window_size, cur_index }) {

  if(!game.table) {
      return (<div></div>)
  }

  let p_index = game.state.cur_index - game.table.length

  const tableCopy = [...game.table];

  for(let i=game.table.length; i < 3; i++) {
    tableCopy.push({"suit": 0, "value": 0, "spec": 0})
  }

  const reorderedTable1 = [
      ...tableCopy.slice(p_index % tableCopy.length).map((card, _) => (card)),
      ...tableCopy.slice(0, p_index % tableCopy.length).map((card, _) => (card)),
  ];

  console.log(p_index)
  console.log(cur_index)
  console.log("baf1")
  console.log(reorderedTable1)

  const reorderedTable = [
      ...reorderedTable1.slice((cur_index - 1) % tableCopy.length).map((card, _) => (card)),
      ...reorderedTable1.slice(0, (cur_index - 1) % tableCopy.length).map((card, _) => (card)),
  ];

  console.log(reorderedTable)

  return (
    <div>
        { reorderedTable.map((item, index) => (
            <div key={index}>
              {index == 0 &&
                <div>
                  <img className="" 
                      src={get_card_src(item)} width={window_size.width/13} id={item.suit + " " + item.spec}
                      style={{position: 'absolute', left: `${35 + index * 10}%`, bottom:'50%', 
                              transform: `rotate(108deg)`}}>
                  </img>
                </div>
              }
              {index == 1 &&
                <div>
                  <img className="" 
                      src={get_card_src(item)} width={window_size.width/13} id={item.suit + " " + item.spec}
                      style={{position: 'absolute', left: `${35 + index * 10}%`, bottom:'40%', 
                              transform: `rotate(-5deg)`}}>
                  </img>
                </div>
              }
              {index == 2 &&
                <div>
                  <img className="" 
                      src={get_card_src(item)} width={window_size.width/13} id={item.suit + " " + item.spec}
                      style={{position: 'absolute', left: `${35 + index * 10}%`, bottom:'50%', 
                              transform: `rotate(80deg)`}}>
                  </img>
                </div>
              }
              
            </div>
        ))}
    </div>
  )
}
