import { Link } from "react-router-dom"

export default function Header({newGame}) {
    return (
        <header className="w-full px-6 py-4">
            <div className="flex justify-between">
                <Link href="/" className="h-10">
                    <img src="logo/logo-dark.svg" alt="Wordle Logo" className="object-cover w-full h-full" />
                </Link>
                <button className="px-4 py-1 text-white bg-black rounded-full hover:bg-gray-700 transition-colors duration-200" onClick={newGame}>
                    Restart the game
                </button>
            </div>
        </header>
    )
}