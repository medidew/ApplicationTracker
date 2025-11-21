import React from 'react'
import LoginForm from '../components/LoginForm';

const API_BASE = process.env.NEXT_PUBLIC_API_BASE;

const LoginPage = async () => {
    
    return (
        <>
            <LoginForm/>
        </>
    )
}

export default LoginPage