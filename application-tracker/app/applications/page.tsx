import React from "react"
import { redirect } from "next/navigation"
import { cookies } from "next/headers"
import HomeButton from "../components/HomeButton"

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
            "Cookie": browser_cookies.toString() // Forward cookies to the API
        },
    });

    let data
    try {
        data = await response.json();
    } catch {
        data = null;
    }

    if (!response.ok) {
        console.error("Failed to fetch applications:", response);
        if (response.status === 401) {
            console.log("User not authenticated, redirecting to login.");
            redirect("/login")
        }
    }

    let applications: Application[] = data;

    return (
        <>
            <div className="flex flex-col items-center justify-start py-10">
                <HomeButton/>
            </div>
            <div className="flex flex-col justify-center items-center font-bold font-sans">
                <h1>Applications</h1>
                <ul>
                    {applications.map((app, index) => (
                        <li key={index}>
                            <h2>{app.company} - {app.role}</h2>
                            <p>Status: {app.status}</p>
                            <h3>Notes:</h3>
                            <ul>
                                {app.notes.map((note, noteIndex) => (
                                    <li key={noteIndex}>{note}</li>
                                ))}
                            </ul>
                        </li>
                    ))}
                </ul>
            </div>
        </>
    )
}

export default ApplicationsPage