import React from 'react'

const ChatHistory = ({messages}) => {
    // console.log(messages);
  return (
    <div>
        <ul>
            {messages.map(msg =>(
                <li key={msg.timeStamp}>{msg.data}</li>
            ))}
        </ul>
    </div>
  )
}

export default ChatHistory