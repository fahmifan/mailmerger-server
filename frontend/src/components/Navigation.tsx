import {ReactElement, useState} from 'react';
import { Link } from "react-router-dom";

type NavItem = {
    url: string,
    text: string,
}

type Prop = {
    items: NavItem[],
}

function Navigation(prop: Prop): ReactElement {
    return (
        <nav>
            <ol className="navs">
                {prop.items.map(it => (
                    <li className='nav'><Link to={it.url}>{it.text}</Link></li>
                ))}
            </ol>
        </nav>
    )
}

export default Navigation