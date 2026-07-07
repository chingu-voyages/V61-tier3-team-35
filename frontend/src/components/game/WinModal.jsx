import { useEffect, useRef } from "react"
import Button from "../Button"
import { X } from "lucide-react"

export default function WinModal({ newGame, onClose }) {

    const newGameButtonRef = useRef(null);
    const closeButtonRef = useRef(null);

    useEffect(() => {
        // Bring tab focus to the close and new game button for accessibility
        closeButtonRef.current.focus()
        newGameButtonRef.current.focus()
    }, [])

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
        <dialog aria-labelledby="win-title" open className="absolute z-20 bg-white top-1/4 left-1/4 w-80 rounded-md overflow-hidden shadow-2xl">
            <div className="w-full bg-green-100 font-bold text-center text-gray-700 py-3">
                <button ref={closeButtonRef} onClick={onClose} className="absolute right-2 text-gray-500" aria-label="close modal" >
                    <X aria-hidden />
                </button>
                <h1 id="win-title">You Win! 🏆</h1>
            </div>
            <div className="flex flex-col items-center justify-center gap-6 px-5 py-8">
                <p className="text-center font-semibold text-gray-600">Congrats on winning! Click the button below to begin a new game.</p>
                <div className="flex flex-col">
                    <Button ref={newGameButtonRef} onClick={newGame} text="new game" />
                    <span className="text-xs text-gray-700 pt-1">or Press Enter to play again</span>
                </div>

            </div>
        </dialog>
    )
}