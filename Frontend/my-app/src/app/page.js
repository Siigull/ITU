"use client";

import React, { useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { NodeNextRequest } from 'next/dist/server/base-http/node';

export default function NicknameForm() {
  const router = useRouter();
  const [nickname, set_nickname] = useState('');  

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (nickname.trim()) {
      const response = await fetch(`http://localhost:8080/join/${nickname}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        }
      });

      var parsed_response = await response.json()
      var token = await parsed_response.token
      var id = await parsed_response.id
      
      if(response.ok) {
        router.push(`/menu?nickname=${nickname}&token=${token}&id=${id}`);
      }
    }
  };

  return (
    <div className="bg-background h-screen flex items-center justify-center">
      <div className="flex flex-col items-center space-y-4">
        <h1 className="text-2xl font-bold">Enter your nickname</h1>
        <form onSubmit={handleSubmit}>
          <div className="flex flex-col items-center space-y-4">
            <input maxLength="20" cols="22" rows="1"
              placeholder="Enter your nickname"
              value={nickname}
              onChange={(e) => set_nickname(e.target.value)}
              required
              style={{ resize: 'none', color:'black' }}
            />
            <button type="submit">Confirm</button>
          </div>
        </form>
      </div>
    </div>
  );
}
