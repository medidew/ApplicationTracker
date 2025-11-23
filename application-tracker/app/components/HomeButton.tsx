import React from 'react'
import Image from "next/image";

const HomeButton = () => {
    return (
        <div>
            <a href="/">
            <Image
                src="/steam_pfp.jpg"
                alt="Steam Profile Picture"
                width={60}
                height={60}
            />
            </a>
        </div>
    )
}

export default HomeButton