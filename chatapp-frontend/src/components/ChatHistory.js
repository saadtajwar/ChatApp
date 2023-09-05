import React from 'react'

const ChatHistory = ({messages, user}) => {
    console.log(messages);
  return (
    <div>
        <ul>
          <li>Test</li>
            {/* {messages.map(msg =>(
                <li key={msg.timeStamp}>{user}: {JSON.parse(msg.data).body}</li>
            ))} */}
        </ul>
    </div>
  )
}

export default ChatHistory