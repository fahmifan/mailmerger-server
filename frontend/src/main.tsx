import React from 'react'
import {createRoot} from 'react-dom/client'
import { createHashRouter, RouterProvider } from "react-router-dom";

import Home from './Home';
import Navigation from "./components/Navigation";

import './style.css'
import navigationData from "./data/navigation.json"

const Nav = () => <Navigation items={navigationData} />
const container = document.getElementById('root')
const root = createRoot(container!)

const router = createHashRouter([
    {
        path: "/",
        element: <>
            <Nav />
            <Home />
        </>
    },
    {
        path: "/campaigns",
        element: <>
            <Nav />
            <h1>Campaigns</h1>
        </>
    },
    {
        path: "/templates",
        element: <>
            <Nav />
            <h1>Templates</h1>
        </>
    }
])

root.render(
    <React.StrictMode>
        <RouterProvider router={router} />
    </React.StrictMode>
)
