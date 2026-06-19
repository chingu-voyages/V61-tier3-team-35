import { Link } from "react-router-dom"

export default function Footer() {
    return (
        <footer className="border-t border-gray-300 px-6 py-4 flex justify-center">
            <Link to="https://github.com/chingu-voyages/V61-tier3-team-35" className="text-gray-600 hover:text-gray-800 focus:text-gray-800 duration-200 transition-colors">Our Team's Repository: V61-tier-3-team-35
            </Link>
        </footer>
    )
}