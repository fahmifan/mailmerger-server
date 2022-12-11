import {ReactElement} from 'react';
import { Link } from "react-router-dom"

function Home(): ReactElement {
    return (
        <>
            <main>
                <h1>Mail Merger</h1>
                <p>Send customize email to your csv lists</p>

                <Link to="/campaigns">
                    <button>Campaigns</button>
                </Link>
            </main>
        </>
    )
}

export default Home;
