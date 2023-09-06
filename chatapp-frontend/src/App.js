import { useEffect, useState } from 'react';
import ChatPage from './components/ChatPage';
import HomePage from './components/HomePage';
import LoginPage from './components/LoginPage'
import RegisterPage from './components/RegisterPage'
import {BrowserRouter as Router, Routes, Route} from 'react-router-dom'


const App = () => {
  const [user, setUser] = useState(null);

  useEffect(() => {
    const loggedUser = window.localStorage.getItem('loggedUser');
    if (loggedUser) {
      setUser(JSON.parse(loggedUser));
    }
  }, [])



  return (
    <Router>
      <Routes>
        <Route path='/register' element={<RegisterPage user={user} />} />
        <Route path='/login' element={<LoginPage user={user} setUser={setUser} />}  />
        <Route path='/chat' element={user === null ? <>Login required</> : <ChatPage user={user} />} />
        <Route path='/' element={<HomePage />} />
      </Routes>
    </Router>
  );
}

export default App;
