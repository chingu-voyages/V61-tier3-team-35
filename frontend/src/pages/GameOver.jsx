import Header from "../components/Header"
import { useNavigate } from "react-router-dom"

export default function GameOver() {

    const API_BASE_URL = "https://wordle-grqh.onrender.com"
    const navigate = useNavigate()

    const newGame = async () => {

        const response = await fetch(`${API_BASE_URL}/api/practice/new-game`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        })

        navigate("/")
    }

    return (
        <div>
            <Header />
            <div className="min-h-[83vh] flex flex-col items-center justify-center gap-5 text-center px-6">
                <h1 className="text-4xl font-bold">
                    Game Over
                </h1>

                <p className="max-w-md text-lg text-gray-600">
                    You've already completed today's challenge. Come back tomorrow for a new word!
                </p>

                <p className="text-gray-500">
                    Want more practice? Try the unlimited mode and keep playing as much as you'd like.
                </p>

                <button onClick={newGame} className="border mt-2 px-4 py-3 rounded-md hover:bg-gray-700 focus:bg-gray-700 focus:text-white hover:text-white transition-colors duration-200">
                    Play Unlimited
                </button>

            </div>

        </div>
    )
}