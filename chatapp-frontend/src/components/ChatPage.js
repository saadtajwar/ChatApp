import React, { useEffect, useMemo, useState } from 'react'
import ChatHistory from './ChatHistory';
import ChatInput from './ChatInput';

const ChatPage = ({user}) => {
  const [chatHistory, setChatHistory] = useState([]);
  const [userList, setUserList] = useState([]);
  const [message, setMessage] = useState('');
  const [selectedUserID, setSelectedUserID] = useState('');
  const [userID, setUserID] = useState('');
  const [webSocketConnection, setWebSocketConnection] = useState(null);

  useEffect(() => {
    if (userID !== '') return;
    const setConnection = () => {
      const webSocketConnection = new WebSocket(`ws://localhost:8080/ws/${user}`)
      setWebSocketConnection(webSocketConnection);
      console.log("Set the ws conn, here is what it is ", webSocketConnection);
      webSocketConnection.onopen = () => {
        console.log("Successfully Connected");
      };

      webSocketConnection.onmessage = (event) => {
        console.log("here in the onmessage")
        try {
          const socketPayload = JSON.parse(event.data);
          console.log('payload', socketPayload);
          switch (socketPayload.eventname) {
            case 'register':
              console.log("Here in the register");
              if (!socketPayload.eventpayload) {
                return;
              }
              const userInitPayload = socketPayload.eventpayload;
              console.log('userinitpayload: ', userInitPayload);
              if (userInitPayload.username === user) {
                console.log("made it in here to assign the userID");
                setUserID(userInitPayload.userid);
              }
              // const newUser = {username: userInitPayload.username, userid: userInitPayload.id};
              setUserList(userInitPayload.users);
              // if (userList !== userInitPayload) setUserList(userInitPayload.users);
              // if (!userID) setUserID(userInitPayload.userid);
              console.log('The new userID: ', userID);
              console.log('The new userlist: ', userList);
              break;
            case 'disconnect':
              if (!socketPayload.eventpayload) {
                return;
              }
              const newUserList = userList.filter(u => u.userid !== socketPayload.eventpayload.userid);
              console.log('here in disconnect');
              // setUserList(newUserList);
              break;
            case 'message response':
              console.log("here in the messages response");
              if (!socketPayload.eventpayload) {
                return;
              }
              const payload = socketPayload.eventpayload;
              const sentBy = payload.username ? payload.username : 'Unnamed';
              const message = payload.message;
              setMessage(`${sentBy}: ${message}`);
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
        // setUserList([]);
      };
    
      webSocketConnection.onerror = (error) => {
        console.log("Error for the websocket: ", error);
      };  

    }

    console.log("Here in the useeffect");
    setConnection();
    // subscribeToSocket();

    // const callback = (msg) => {
    //       setChatHistory(prevChatHistory => [...prevChatHistory, msg]);
    // }
  }, [])


  console.log('userlist', userList);
  console.log("userid", userID)

  const handleSend = (event) => {
    try  {
      if (event.keyCode === 13) {
        webSocketConnection.send(JSON.stringify({
          eventname: 'message',
          eventpayload: {
            userid: selectedUserID,
            message: event.target.value,
            username: user
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
        setSelectedUserID(event.target.value);
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
                        userList.map(userListUser => {
                            if (userListUser.userid !== userID) {
                                    return <option value={userListUser.userid} className="username-list">
                                        {userListUser.username}
                                    </option>
                            }
                        })
                    }
            </select>
          }
          <div>
            <div>
              {message}
            </div>
            <input onKeyDown={handleSend} />
          </div>
        </div>
    </section>
  )
}

export default ChatPage