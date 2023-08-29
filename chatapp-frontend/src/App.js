import ChatPage from './components/ChatPage';
import HomePage from './components/HomePage';
import LoginPage from './components/LoginPage'
import RegisterPage from './components/RegisterPage'
import {BrowserRouter as Router, Routes, Route} from 'react-router-dom'


const App = () => {
  return (
    <Router>
      <Routes>
        <Route path='/register' element={<RegisterPage />} />
        <Route path='/login' element={<LoginPage />} />
        <Route path='/chat' element={<ChatPage />} />
        <Route path='/' element={<HomePage />} />
      </Routes>
    </Router>
  );
}

export default App;
