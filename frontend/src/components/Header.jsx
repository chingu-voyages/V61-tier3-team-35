import { Link } from "react-router-dom"

export default function Header({ newGame, className, action }) {
    return (
        <header className={`w-full md:px-22 px-6 md:pt-6 pt-6 ${className}`}>
            <div className="flex justify-between">
                <Link href="/" className="h-10 z-10">
                    <img src="logo/logo-dark.svg" alt="Wordle Logo" className="object-cover w-full h-full" />
                </Link>
                {action}
            </div>
        </header>
    )
}