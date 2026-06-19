import { createBrowserRouter } from 'react-router-dom'
// Pages
import App from '../App.jsx'
import Home from '../pages/Home.jsx'
import Gameplay from '../pages/Gameplay.jsx'
import Dashboard from '../pages/Dashboard.jsx'

export const router = createBrowserRouter([
    {
        path: "/",
        element: <App />
    },
    {
        path: "/home",
        element: <Home />
    },
    {
        path: "/gameplay",
        element: <Gameplay />
    },
    {
        path: "/dashboard",
        element: <Dashboard />
    }
])