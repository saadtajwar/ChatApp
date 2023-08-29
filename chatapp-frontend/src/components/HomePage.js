import React from 'react'
import {BrowserRouter as Router, Link, Routes, Route} from 'react-router-dom'


const HomePage = () => {
  return (
    <div className="login-page">
        <div className="card">
            <div className="title">Welcome to SaadChat</div>
                <Link className="link" to='/register'>Register</Link>
                <Link className="link" to='/login'>Log In</Link>
        </div>
        <style>{`
        .login-page { width: 100vw; height: 100vh; padding-top: 6vw; background: linear-gradient(180deg, rgba(117,84,160,1) 7%, rgba(117,84,160,1) 17%, rgba(106,95,168,1) 29%, rgba(99,103,174,1) 44%, rgba(87,116,184,1) 66%, rgba(70,135,198,1) 83%, rgba(44,163,219,1) 96%, rgba(22,188,237,1) 100%, rgba(0,212,255,1) 100%); }
        .card { width: 200px; position: relative; left: calc(50vw - 100px); text-align: center; }
        .title { padding-top: 32px; font-size: 22px; color: white; font-weight: 700; }
        input { width: calc(100% - 16px); margin-top: 12px; padding: 8px; background-color: #e6f7ff; outline: none; border: 1px solid #e6f7ff; }
        .link { margin-top: 12px; width: 100%; padding: 8px; }
        `}</style>
    </div>
  )
}

export default HomePage