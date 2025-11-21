import React from 'react'
import { redirect } from 'next/navigation'

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
    const res = await fetch(url); // TODO: put api config stuff

    let data
    try {
        data = await res.json();
    } catch {
        data = null;
    }

    if (!res.ok) {
        console.error("Failed to fetch applications:", data);
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