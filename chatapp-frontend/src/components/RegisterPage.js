import React, {useState} from 'react'
import axios from 'axios'
import { useNavigate } from "react-router-dom";


const RegisterPage = ({user}) => {
    // const [name, setName] = useState('');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const navigate = useNavigate();

    if (user) {
        return (
            <div>Already logged in</div>
        )
    }

    const handleRegister = async (event) => {
        event.preventDefault();
        console.log("Got into this func");
        try {
            const userToRegister = {
                username,
                // name,
                password
            }
            const newUser = await axios.post("http://localhost:8080/users", userToRegister);
            // console.log("New user added:", newUser);
            // setName('');
            setUsername('');
            setPassword('');
            navigate('/login')
        } catch (error) {
            console.log('Registration failed');
        }
    }

    return (
        <div className="login-page">
         <div className="card">
          <form onSubmit={handleRegister}>
            <div className="title">Sign Up</div>
            <input
              type="text"
              name="username"
              placeholder="Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
            <input
              type="password"
              name="secret"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            <button type="submit">SIGN UP</button>
          </form>
          <style>{`
            .login-page { width: 100vw; height: 100vh; padding-top: 6vw; background: linear-gradient(180deg, rgba(117,84,160,1) 7%, rgba(117,84,160,1) 17%, rgba(106,95,168,1) 29%, rgba(99,103,174,1) 44%, rgba(87,116,184,1) 66%, rgba(70,135,198,1) 83%, rgba(44,163,219,1) 96%, rgba(22,188,237,1) 100%, rgba(0,212,255,1) 100%); }
            .card { width: 200px; position: relative; left: calc(50vw - 100px); text-align: center; }
            .title { padding-top: 32px; font-size: 22px; color: blue; font-weight: 700; }
            input { width: calc(100% - 16px); margin-top: 12px; padding: 8px; background-color: #e6f7ff; outline: none; border: 1px solid #e6f7ff; }
            button { margin-top: 12px; width: 100%; padding: 8px; }
        `}</style>
            </div>
        </div>
    )
}

export default RegisterPage