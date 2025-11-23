import React from 'react'
import { redirect } from 'next/navigation'
import { cookies } from 'next/headers'

const API_BASE = process.env.NEXT_PUBLIC_API_BASE;

interface Application {
    company: string;
    role: string;
    status: number;
    notes: string[];
    username: string;
}

const ApplicationsPage = async () => {
    const url = API_BASE + "/applications"
    const browser_cookies = await cookies();
    
    const response = await fetch(url, {
        method: "GET",
        credentials: "include",
        headers: {
            'Cookie': browser_cookies.toString() // Forward cookies to the API
        },
    }); // TODO: put api config stuff

    let data
    try {
        data = await response.json();
    } catch {
        data = null;
    }

    if (!response.ok) {
        console.error("Failed to fetch applications:", response);
        //redirect("/login")
    }

    let applications: Application[] = data;

    return (
        <>
            <h1>Applications</h1>
            <ul>
                {applications.map(application => application.company)}
            </ul>
        </>
    )
}

export default ApplicationsPage