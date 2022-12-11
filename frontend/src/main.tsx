import React from 'react'
import {createRoot} from 'react-dom/client'
import { createHashRouter, RouterProvider } from "react-router-dom";

import Home from './Home';
import ListCampaign, { Campaign, campaignLoader } from './Campaign';
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
            <ListCampaign />
        </>
    },
    {
        path: "/campaigns/:campaignID",
        element: <>
            <Nav />
            <Campaign />
        </>,
        loader: campaignLoader,
    },
    {
        path: "/templates",
        element: <>
            <Nav />
            <h1>Templates</h1>
        </>
    }
])

const App = (props: any) => <>
    {props.children}
</>

root.render(
    <React.StrictMode>
        <App root={root}>
            <RouterProvider router={router} />
        </App>
    </React.StrictMode>
)
