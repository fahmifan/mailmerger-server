import {
    ReactElement,
    useEffect,
    useState,
} from 'react';

import {
    Link, useLoaderData,
} from "react-router-dom"

import AceEditor from "react-ace";
import "ace-builds/src-noconflict/mode-html";
import "ace-builds/src-noconflict/theme-solarized_light";
import "ace-builds/src-noconflict/ext-language_tools";

import * as gotype from '../wailsjs/go/models';
import * as goapp from '../wailsjs/go/main/App';
import * as runtime from '../wailsjs/runtime/runtime';


export default function ListTemplate(): ReactElement {
    const [templates, setTemplates] = useState<gotype.service.Template[]>([]);

    async function findTemplates() {
        const res = await goapp.ListTemplates()
        setTemplates(res)
    }

    useEffect(() => {
        findTemplates()
    })

    return <main>
        <h1>Templates</h1>
        <a href="">
            <button>New</button>
        </a>
        <ul>
            {
                templates.map(tpl => (
                    <li>
                        <Link to={`/templates/${tpl.id}`}>{tpl.name}</Link>
                    </li>
                ))
            }
        </ul>
    </main>
}

export async function templateLoader({ params }: any): Promise<gotype.service.Template> {
    return goapp.ShowTemplate(params.templateID)
}

export function Template(): ReactElement {
    const templatex: any = useLoaderData();
    const template: gotype.service.Template = templatex;

    return <main>
        <h1>Template</h1>
        <a href="#edit">
            <button>
                Edit
            </button>
        </a>

        <p>ID: { template.id }</p>
        <p>Name: { template.name }</p>

        <label>Body</label>

        <AceEditor
            placeholder="No template"
            mode="html"
            theme="solarized_light"
            name="template editor"
            fontSize={16}
            showPrintMargin={true}
            showGutter={true}
            width="100%"
            height='500px'
            highlightActiveLine={true}
            value={template.html}
            readOnly={true}
            setOptions={{
                enableBasicAutocompletion: true,
                enableLiveAutocompletion: true,
                enableSnippets: true,
                showLineNumbers: true,
                tabSize: 4,
            }} 
        />
    </main>
}