import { Link } from "react-router-dom"

export default function Header() {
    return (
        <header className="border-b border-gray-300 text-gray-500 w-full flex justify-center items-center px-6 py-4">
            <Link href="/">
                <h1 className="text-3xl font-bold">Wordle</h1>
            </Link>
        </header>
    )
}