import { createBrowserRouter } from 'react-router-dom'
import MainLayout from '../layouts/MainLayout.jsx'

// Pages
import { Children } from 'react'
import Game from '../pages/Game.jsx'
import GameOver from '../pages/GameOver.jsx'

export const router = createBrowserRouter([
    {
        element: <MainLayout />,
        children: [
            {
                path: "/",
                element: <Game />
            },
            {
                path: "/game-over",
                element: <GameOver />
            },
        ]
    }
]);