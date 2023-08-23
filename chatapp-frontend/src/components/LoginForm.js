import React, {useState} from 'react'

const LoginForm = ({user, setUser}) => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');

    if (user) {
        return (
            <div>Already logged in!</div>
        )
    }

    const handleLogin = async (e) => {
        try {
            e.preventDefault();
            const attemptedUser = {
                username,
                password
            }
            // const loggedUser = await axios.post(attemptedUser to the route for logging in)
            // setUser(loggedUser);
            // window.localStorage.setItem('loggedUser', JSON.stringify(loggedUser));
            setUsername('');
            setPassword('');
        } catch (error) {
            console.log('Registration failed');
        }
    }

    return (
        <div>
            Login Here
            <form onSubmit={handleLogin}>
                Username: <input value={username} onChange={(e)=>setUsername(e.target.value)} />
                Password: <input type="password" value={password} onChange={(e)=>setPassword(e.target.value)} />
                <button type="submit">Login</button>
            </form>
        </div>
    )
}

export default LoginForm