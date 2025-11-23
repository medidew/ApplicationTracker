import Image from "next/image";

export default function Home() {
  return (
    <>
    <div className="flex flex-col items-center justify-start py-20">
      <Image
        src="/steam_pfp.jpg"
        alt="Steam Profile Picture"
        width={150}
        height={150}
        className="mb-4"
      />
    </div>
    <div className="flex flex-col justify-center items-center font-bold font-sans">
      <h1 className="flex text-4xl">ApplicationTracker</h1>
      <div className="p-[20]">
        <a href="/login" className="inline-block px-6 py-2 bg-blue-500 text-white rounded hover:bg-blue-700 text-center">
          Login
        </a>
      </div>
    </div>
    </>
  );
}
