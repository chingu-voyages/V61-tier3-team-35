"use client"
import { useState, useEffect, useRef } from "react"
import Board from "../components/Board/Board"
import Keyboard from "../components/Keyboard/Keyboard"
// using keyboard rows to check if pressed key should is valid
import keyboardRows from "../components/Keyboard/keyboardRows"
import mockStatus from "../components/mockStatus"

// Win / lose modal
import WinModal from "../components/game/WinModal"
import LoseModal from "../components/game/LoseModal"
import ErrorModal from "../components/game/ErrorModal"
import { ContactRound } from "lucide-react"

// Header
import Header from "../components/Header"

// Empty tile 
const emptyTile = {
    letter: "",
    status: ""
}

// Creates empty 6x5 board with tiles.
const emptyBoard = Array.from({ length: 6 }, () =>
    Array.from({ length: 5 }, () => ({ ...emptyTile }))
)

const InitialKeyboardStatuses = Object.fromEntries(
    "abcdefghijklmnopqrstuvwxyz"
        .split("")
        .map(letter => [letter, ""])
)

export default function () {
    const [keyValue, setKeyValue] = useState("")
    const [keyActive, setKeyActive] = useState(false)
    const [board, setBoard] = useState(emptyBoard)
    const [keyboardStatuses, setKeyboardStatuses] = useState(InitialKeyboardStatuses)
    // Curr row and col
    const [currCol, setCurrCol] = useState(0)
    const [currRow, setCurrRow] = useState(0)
    const [gameStatus, setGameStatus] = useState("playing")
    const [showWinModal, setShowWinModal] = useState(false)
    const [showLoseModal, setShowLoseModal] = useState(false)
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [error, setError] = useState("")

    const API_BASE_URL = "https://wordle-grqh.onrender.com"


    useEffect(() => {
        if (!error) return;

        const timer = setTimeout(() => {
            setError("")
        }, 1500);

        // clear time out
        return () => clearTimeout(timer);
    }, [error])

    useEffect(() => {

        const handleKeyDown = (event) => {
            handleKeyPress(event.key)
            setKeyActive(true)
        }

        window.addEventListener("keydown", handleKeyDown)


        return () => {
            window.removeEventListener("keydown", handleKeyDown)
        }
    }, [currRow, currCol, board])


    const handleKeyPress = (key) => {
        if (gameStatus !== "playing" || isSubmitting) return

        setKeyValue(key)

        // col - letters
        // rows - guesses

        if (key === "Backspace") {
            handleBackspace()
        }
        else if (key === "Enter") {
            submitGuess()
        }
        else {
            handleLetter(key)
        }

    }

    // Handle letter

    const handleLetter = (key) => {

        const newBoard = [...board]
        newBoard[currRow] = [...newBoard[currRow]]

        if (keyboardRows.flat().includes(key) && currCol < 5) {

            // Move the cursor forward
            setCurrCol(prev => prev + 1)

            // Add the letter to the board
            newBoard[currRow][currCol] = {
                letter: key,
                status: ""
            }
            setBoard(newBoard)
        }
    }


    // Backspace
    const handleBackspace = () => {

        const newBoard = [...board]
        newBoard[currRow] = [...newBoard[currRow]]

        if (currCol > 0) {
            // Clear the  tile
            newBoard[currRow][currCol - 1] = {
                letter: "",
                status: ""
            }
            setBoard(newBoard)

            // Move the cursor to previous tile
            setCurrCol(prev => prev - 1)
        }
    }

    // Submit Guess
    const submitGuess = async () => {

        if (isSubmitting) { return }

        if (currCol === 5 && currRow < 6) {
            setIsSubmitting(true)
            // Check the status of letters in guess word - api call
            const success = await checkGuess()

            if (!success) {
                setIsSubmitting(false)
                return;
            }

            // Set col to 0 and move to next row
            setCurrCol(0)
            setCurrRow(prev => prev + 1)
            // when checkGuess has run, set IsSubmitting to false
            setIsSubmitting(false)
        }

    }


    // Check guess
    const checkGuess = async () => {

        const newBoard = [...board]
        newBoard[currRow] = [...newBoard[currRow]]

        const guess = newBoard[currRow];
        let guessWord = ""
        for (let i = 0; i < guess.length; i++) {
            guessWord += newBoard[currRow][i].letter
        }

        const response = await fetch(`${API_BASE_URL}/api/guess`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                guess: guessWord,
            }),
            credentials: "include"
        })

        const data = await response.json()

        console.log("data is:", data)


        if (data?.error) {
            if (data.error === "word is not in the accepted word list") {
                setError("Word Not Found")
            }
            return false;
        }
        else {
            setError("")
        }

        const guessStatus = data?.feedback;
        const isCorrect = data?.is_correct;

        const newKeyboardStatuses = { ...keyboardStatuses }

        for (let i = 0; i < guess.length; i++) {

            const letter = guess[i].letter;

            // Checking if keyboard status is higher than result status for relevant color updates
            let keyStatus = newKeyboardStatuses[letter]
            let resultStatus = guessStatus[i].status

            // if key status weight is not higher than resultStatus weight, change key status.

            if (!newKeyboardStatuses[letter] ||
                !checkKeyStatusWeight(keyStatus, resultStatus)) {
                newKeyboardStatuses[letter] = guessStatus[i].status
            }

            // Update gameboard status
            newBoard[currRow][i] = {
                ...newBoard[currRow][i],
                status: guessStatus[i].status
            }

        }
        setBoard(newBoard)
        setKeyboardStatuses(newKeyboardStatuses)

        if (isCorrect) {
            setGameStatus("won")
            setShowWinModal(true)
        }
        else if (!isCorrect) {
            if (currRow > 4) {
                setGameStatus("lost")
                setShowLoseModal(true)


            }
        }

        return true;

    }

    // Checking highest status to update keyboard key's statuses
    const checkKeyStatusWeight = (keyStatus, resultStatus) => {

        const statuses = {
            absent: 0,
            present: 1,
            correct: 2
        }

        const keyStatusWeight = statuses[keyStatus]
        const resultStatusWeight = statuses[resultStatus]

        // returns true or false
        return keyStatusWeight > resultStatusWeight

    }

    // Get new word
    const getNewWord = async () => {
        resetGame()


        const response = await fetch(`${API_BASE_URL}/api/daily-word`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        })
        const data = await response.json();

    }

    // Reset Game
    const resetGame = () => {
        setBoard(emptyBoard)
        setKeyboardStatuses(InitialKeyboardStatuses)
        setShowLoseModal(false)
        setShowWinModal(false)
        setCurrCol(0)
        setCurrRow(0)
        setGameStatus("playing")
    }


    return (
        <div className="w-full relative">
            {/* overlay */}
            {(showLoseModal || showWinModal) && (<div className="bg-white/70 absolute inset-0 z-10"></div>)}
            <Header className="w-full"
                action={
                    <button className="px-4 py-1 text-white bg-black rounded-full hover:bg-gray-700 transition-colors duration-200 z-10" onClick={resetGame}>
                        Restart the game
                    </button>
                }
            />
            <main className={`flex flex-col items-center z-0 min-h-[82vh] justify-center`}>
                {showWinModal && (<WinModal onClose={() => { setShowWinModal(false) }} newGame={getNewWord}/>)}
                {showLoseModal && (<LoseModal onClose={() => { setShowLoseModal(false) }} newGame={getNewWord}/>)}
                <div className={`flex flex-col items-center justify-center ${gameStatus === "playing" ? "md:gap-12 gap-10" : "gap-0"} md:-mt-4`}>
                    <Board board={board} />
                    {gameStatus === "won" && (
                        <div className="bg-gray-200 rounded-full p-1 px-2 text-sm font-semibold my-2">You Win! 🏆</div>
                    )}
                    {gameStatus === "lost" && (
                        <div className="bg-gray-200 rounded-full p-1 px-2 text-sm font-semibold my-2">You Lose! 🥲</div>
                    )}
                    {error && (<ErrorModal error={error} />)}
                    <Keyboard activeKey={keyValue} keyboardStatuses={keyboardStatuses} handleKeyPress={handleKeyPress} />
                </div>
            </main>
        </div>

    )
}