import React, { useEffect, useState } from 'react'
import {connect, sendMsg} from '../api'
import ChatHistory from './ChatHistory';
import ChatInput from './ChatInput';

const ChatPage = () => {
  const [chatHistory, setChatHistory] = useState([]);

  useEffect(() => {
    connect((msg) => {
      console.log("New Message from useeffect?");
      setChatHistory(prevChatHistory => [...prevChatHistory, msg]);
      console.log(chatHistory);
    });
  }, []);


  const send = (e) => {
    if (e.keyCode === 13) {
      sendMsg(e.target.value);
      e.target.value = "";
    }
  }

  return (
    <section className="hero is-warning">
        <div className="hero-body">
            <p className="title">SaadChat</p>
        </div>
        <div className='hero-body'>
          <ChatHistory messages={chatHistory} />
          <ChatInput send={send}/>
        </div>
    </section>
  )
}

export default ChatPage