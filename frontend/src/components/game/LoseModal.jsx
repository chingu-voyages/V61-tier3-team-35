import Button from "../Button"
import { X } from "lucide-react"
import { useEffect } from "react"

export default function LoseModal({ newGame, onClose }) {

    useEffect(() => {
        const handleKeyDown = (event) => {
            if (event.key === "Enter") {
                newGame()
            }
        }

        window.addEventListener("keydown", handleKeyDown)

        return () => {
            window.removeEventListener("keydown", handleKeyDown)
        }
    }, [newGame])

    return (
        <article className="absolute z-20 bg-white top-1/4 w-80 rounded-md overflow-hidden shadow-2xl">
            <div className="w-full bg-green-100 font-bold text-center text-gray-700 py-3">
                <button onClick={onClose} className="absolute right-2 text-gray-500" >
                    <X />
                </button>
                You Lose! 🥲
            </div>
            <div className="flex flex-col items-center justify-center gap-6 px-5 py-8">
                <div className="flex flex-col gap-2">The answer was: <span className="uppercase font-bold border border-dashed border-gray-400 p-2 text-center text-gray-700 bg-gray-100">apple</span></div>
                <p className="text-center font-semibold text-gray-600">Better luck next time! Click the button below to begin a new game.</p>
                <div className="flex flex-col">
                    <Button onClick={newGame} text="new game" />
                    <span className="text-xs text-gray-700 pt-1">or Press Enter to play again</span>
                </div>
            </div>
        </article>
    )
}