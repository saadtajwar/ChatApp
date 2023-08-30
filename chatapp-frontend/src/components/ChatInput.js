import React, { useState } from 'react'

const ChatInput = ({send}) => {

    return (
        <div>
            <input onKeyDown={send}/>
        </div>
    )
}

export default ChatInput