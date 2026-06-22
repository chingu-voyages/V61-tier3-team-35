import { Link } from "react-router-dom"

export default function Header() {
    return (
        <header className="border-b border-gray-300 text-gray-500 mb-10 w-full flex justify-center items-center px-6 py-4">
            <Link href="/">
                <h1 className="text-3xl font-bold">Wordle</h1>
            </Link>
            {/* <bzutton className="rounded-md px-5 py-3 text-white bg-gray-800 hover:bg-gray-300 hover:text-black focus:bg-gray-300 focus:text-black transition-colors duration-300">Sign Out</button> */}
        </header>
    )
}