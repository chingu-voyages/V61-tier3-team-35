import { Delete } from "lucide-react";


const statuses =
{
    correct: "text-white border-none bg-lime-500",
    present: "text-white border-none bg-yellow-400",
    absent: "text-white border-none bg-gray-400",
}

export default function Key({ text, onClick, className, status }) {
    return (
        <button className={`${statuses[status]} ${className} bg-gray-200 hover:bg-gray-300 focus:bg-gray-300 md:px-5 sm:px-4 px-3 md:py-3 sm:py-3 py-3 md:text-lg sm:text-sm text-xs rounded-md font-extrabold text-gray-700 uppercase duration-300 transition-colors`} onClick={onClick}>
            {text === "Backspace" && <Delete />}
            {text !== "Backspace" && text}
        </button>
    )
}