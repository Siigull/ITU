'use client';
 
import React from 'react';

export default function Money({ amount }) {
 
  return (
    <div>
        {parseInt(amount) < 0 &&
            <div className="text-left text-red-500">{amount}Kč</div>
        }
        {parseInt(amount) > 0 &&
            <div className="text-left text-green-500">{amount}Kč</div>
        }
        {parseInt(amount) === 0 &&
            <div className="text-left">{amount}Kč</div>
        }
    </div>
  )
}
