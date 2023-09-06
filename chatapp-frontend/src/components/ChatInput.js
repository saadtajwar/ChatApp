import React, { useState } from 'react'

const ChatInput = ({handleSend}) => {

    return (
        <div>
            <input onKeyDown={handleSend}/>
        </div>
    )
}

export default ChatInput