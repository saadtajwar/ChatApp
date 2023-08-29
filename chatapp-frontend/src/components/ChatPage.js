import React, { useEffect, useState } from 'react'
import {connect, sendMsg} from '../api'
import ChatHistory from './ChatHistory';

const ChatPage = () => {
  const [chatHistory, setChatHistory] = useState([]);

  useEffect(() => {
    connect((msg) => {
      console.log("New Message from useeffect?");
      setChatHistory(prevChatHistory => [...prevChatHistory, msg]);
      console.log(chatHistory);
    });
  }, []);


  const send = () => {
    sendMsg("Goodbye");
    console.log("Send function in chatpage called with hello");
  }

  return (
    <section className="hero is-warning">
        <div className="hero-body">
            <p className="title">SaadChat</p>
        </div>
        <div className='hero-body'>
          <ChatHistory messages={chatHistory} />
          <button onClick={send}>Hit</button>
        </div>
    </section>
  )
}

export default ChatPage