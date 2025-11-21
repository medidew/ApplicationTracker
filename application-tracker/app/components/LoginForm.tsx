'use client'; // enables events and reactivity
import React from 'react'

const LoginForm = () => {

    async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
        event.preventDefault();

        const formData = new FormData(event.currentTarget);
        const API_BASE = process.env.NEXT_PUBLIC_API_BASE;
        const url = API_BASE + "/login"

        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'multipart/form-data',
                //'Access-Control-Allow-Origin': 'http://localhost:3000',
            },
            body: formData,
        });

        if (response.ok) {
            // Handle successful login (e.g., redirect to applications page)
            console.log('Login successful');
        } else {
            // Handle login failure (e.g., show error message)
            console.error('Login failed');
        }
    }

    return (
        <>
            <form onSubmit={handleSubmit}>
                <input type="text" name="username" placeholder="Username" required />
                <input type="password" name="password" placeholder="Password" required />
                <button type="submit">Login</button>
            </form>
        </>
    )
}

export default LoginForm