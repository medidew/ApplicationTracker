import React from "react"
import LoginForm from "../components/LoginForm"
import HomeButton from "../components/HomeButton"

const API_BASE = process.env.NEXT_PUBLIC_API_BASE;

const LoginPage = async () => {
    
    return (
        <>
            <div className="flex flex-col items-center justify-start py-10">
                <HomeButton/>
                <div className="text-2xl font-bold mt-4">Login to ApplicationTracker</div>
            </div>
            <div className="flex flex-col justify-center items-center font-bold font-sans">
                <LoginForm/>
            </div>
        </>
    )
}

export default LoginPage