import React, { useEffect, useState } from 'react'
import ChatHistory from './ChatHistory';
import ChatInput from './ChatInput';

const ChatPage = ({user}) => {
  const [chatHistory, setChatHistory] = useState([]);
  const [userList, setUserList] = useState([]);
  const [message, setMessage] = useState('');
  const [selectedUserID, setSelectedUserID] = useState('');
  const [userID, setUserID] = useState('');
  const webSocketConnection = new WebSocket(`ws://localhost:8080/ws/${user}`);

  

  useEffect(() => {
    const callback = (msg) => {
          setChatHistory(prevChatHistory => [...prevChatHistory, msg]);
    }
    subscribeToSocket(callback);
  }, [])


  // useEffect(() => {
  //   connect((msg) => {
  //     console.log("New Message from useeffect?");
  //     setChatHistory(prevChatHistory => [...prevChatHistory, msg]);
  //     console.log(chatHistory);
  //   });
  // }, []);


  const subscribeToSocket = (callback) => {
    if (webSocketConnection === null) {
      return;
    }

    webSocketConnection.onopen = () => {
      console.log("Successfully Connected");
    };
  
    webSocketConnection.onmessage = (event) => {
      console.log("here in the onmessage")
      try {
        const socketPayload = JSON.parse(event.data);
        switch (socketPayload.eventname) {
          case 'register':
          case 'disconnect':
            if (!socketPayload.eventpayload) {
              return;
            }
            const userInitPayload = socketPayload.eventpayload;
            setUserList(userInitPayload.users);
            setUserID(userID === null ? userInitPayload.userid : userID);
            break;
          case 'message response':
            if (!socketPayload.eventpayload) {
              return;
            }
            const payload = socketPayload.eventpayload;
            const sentBy = payload.username ? payload.username : 'Unnamed';
            const message = payload.message;
            setMessage(`${sentBy}: ${message}`)

            break;
          default:
            break;
        }
        // callback(event.data)
      } catch (error) {
        console.log(error)
      }
    };
  
    webSocketConnection.onclose = (event) => {
      setMessage('Connected closed');
      setUserList([]);
    };
  
    webSocketConnection.onerror = (error) => {
      console.log("Error: ", error);
    };  
  }

  // const send = (e) => {
  //   if (e.keyCode === 13) {
  //     console.log('Here in the send block')
  //     webSocketConnection.send(JSON.stringify({
  //       EventName: 'message',
  //       EventPayload: {
  //         Message: e.target.value
  //       }
  //     }))
  //     e.target.value = "";
  //   }
  // }

  const handleSend = (event) => {
    try  {
      if (event.key === 'Enter') {
        if (!this.webSocketConnection || !event.target.value) { 
          return false;
        }

        webSocketConnection.send(JSON.stringify({
          EventName: 'message',
          EventPayload: {
            Message: event.target.value
          }
        }));

        event.target.value = "";

      }
    } catch(error) {
      console.log("Error when handling send", error)
    }
  }



  const setNewUserToChat = (event) => {
    if (event.target && event.target.value) {
        if (event.target.value === "select-user") {
            alert("Select a user to chat");
            return;
        }
        setSelectedUserID(event.target.value)
    }
}


  if (!user) {
    return (
      <div>Must be signed in to use the chat functions!</div>
    )
  }


  return (
    <section className="hero is-warning">
        <div className="hero-body">
            <p className="title">SaadChat</p>
        </div>
        <div className='hero-body'>
          {/* <ChatHistory messages={chatHistory} user={user} /> */}
          {userList.length === 0 ? <h1>Nobody has joined yet</h1> : 
            <select onChange={setNewUserToChat}>
                    <option value={'select-user'} className="username-list">Select User</option>
                    {
                        userList.map(user => {
                            if (user.userID !== userID) {
                                return (
                                    <option value={user.userID} className="username-list">
                                        {user.username}
                                    </option>
                                )
                            }
                        })
                    }
            </select>
          }
          <div>
            <div>
              {message}
            </div>
            <input type="text" id="message-text" onKeyPress={handleSend} />
          </div>
          {/* <ChatInput send={send}/> */}
        </div>
    </section>
  )
}

export default ChatPage