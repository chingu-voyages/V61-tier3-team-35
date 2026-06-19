import { createBrowserRouter } from 'react-router-dom'
import MainLayout from '../layouts/MainLayout.jsx'
// Pages
import { Children } from 'react'
import Game from '../pages/Game.jsx'

export const router = createBrowserRouter([
    {
        element: <MainLayout />,
        children: [
            {
                path: "/",
                element: <Game />
            }
        ]
    }
]);