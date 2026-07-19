import { Link } from "react-router-dom"
// import SelectMode from "../components/SelectMode"

export default function Header({ newGame, className, action }) {
    return (
        <header className={`w-full md:px-22 px-6 md:pt-6 pt-6 ${className}`}>
            <div className="flex justify-between items-center">
                <Link href="/" className="h-10 z-10">
                    <img src="logo/logo-dark.svg" alt="Wordle Logo" className="object-cover w-full h-full" />
                </Link>
                {/* <SelectMode /> */}
                {action}
            </div>
        </header>
    )
}