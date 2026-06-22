"use client"
import { useState, useEffect } from "react"
import Board from "../components/Board/Board"
import Keyboard from "../components/Keyboard/Keyboard"
// using keyboard rows to check if pressed key should is valid
import keyboardRows from "../components/Keyboard/keyboardRows"
import mockStatus from "../components/mockStatus"

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

        if (gameStatus !== "complete") {

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

        else {
            console.log("round complete")
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
                status:""
            }
            setBoard(newBoard)

            // Move the cursor to previous tile
            setCurrCol(prev => prev - 1)
        }
    }

    // Submit Guess
    const submitGuess = () => {

        if (currCol === 5 && currRow < 6) {
            // Check the status of letters in guess word - api call
            checkGuess()

            // Set col to 0 and move to next row
            setCurrCol(0)
            setCurrRow(prev => prev + 1)
        }

        else if (currRow > 5) {
            setGameStatus("complete")
        }

    }


    // Check guess
    const checkGuess = () => {

        const newBoard = [...board]
        newBoard[currRow] = [...newBoard[currRow]]

        const guess = newBoard[currRow];

        // Looping through the guess

        for (let i = 0; i < guess.length; i++) {
            // does the word contain words in mock?

            const letter = newBoard[currRow][i].letter

            // Update gameboard status
            newBoard[currRow][i] = {
                ...newBoard[currRow][i],
                status: mockStatus[letter] ?? "absent"
            }

            if (mockStatus[letter]) {
                keyboardStatuses[letter] = mockStatus[letter]
            }
            else {
                keyboardStatuses[letter] = "absent"
            }



        }
        setBoard(newBoard)

    }



    return (
        <div className="flex flex-col items-center">
            <Board board={board} />
            <Keyboard activeKey={keyValue} keyboardStatuses={keyboardStatuses} handleKeyPress={handleKeyPress} />
        </div>
    )
}