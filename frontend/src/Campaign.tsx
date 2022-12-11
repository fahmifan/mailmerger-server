import {ReactElement, useEffect, useState, Component} from 'react';
import { Link, useLoaderData } from "react-router-dom"
import AceEditor from "react-ace";
import "ace-builds/src-noconflict/mode-plain_text";
import "ace-builds/src-noconflict/theme-solarized_light";
import "ace-builds/src-noconflict/ext-language_tools";

import * as gotype from '../wailsjs/go/models';
import * as goApp from '../wailsjs/go/main/App';
import * as runtime from '../wailsjs/runtime/runtime';

function ListCampaign(): ReactElement {
    let [campaigns, setCampaigns] = useState<gotype.service.Campaign[]>([])

    useEffect(() => {
        findCampaigns()
    }, [campaigns])


    const findCampaigns = async () => {
        const res = await goApp.ListCampaigns()
        setCampaigns(res)
    }

    const handleEdit = (campaignID: string) => {
        runtime.LogDebug("edit >>> "+campaignID)
    }

    return <main>
        <h1>Campaigns</h1>
        <a href="">
            <button>New</button>
        </a>
        <table>
            <thead>
                <th>Campaign</th>
                <th>Action</th>
            </thead>
            <tbody>
                <tr>
                    {campaigns.map(campaign => (
                        <>
                            <td> 
                                <Link to={`/campaigns/${campaign.id}`}>
                                { campaign.name }
                                </Link>
                            </td>
                            <td style={{"display": "flex"}}>
                                <button onClick={() => handleEdit(campaign.id)}>Edit</button>
                                <button>Delete</button>
                            </td>
                        </>
                    ))}
                </tr>
            </tbody>
        </table>
    </main>
}

export async function campaignLoader({ params }: any): Promise<gotype.service.Campaign> {
    return goApp.ShowCampaign(params.campaignID)
}

export function Campaign(): ReactElement {
    const campaignx: any = useLoaderData();
    const campaign: gotype.service.Campaign = campaignx;
    const [content, setContent] = useState(campaign.body)
    const [renderedContent, setRenderedContent] = useState('')

    function onChange(value: any) {
        setContent(value)
    }

    useEffect(() => {
        renderContent()
    }, 
    [content])

    async function renderContent() {
        let res: string
        try {
            res = await goApp.CreateRenderedTemplate(campaign.template?.id, content ? content : '')
            setRenderedContent(res)
        } catch (err: any) {
            runtime.LogError(err)
        }
    }

    return <main>
        <h1>Campaign</h1>
        <Link to={`campaigns/${campaign.id}/edit`}>Edit</Link>
        <p>ID: {campaign.id}</p>
        <p>Name: {campaign.name}</p>
        <p>File: {campaign.file.fileName}</p>
        <p>Template: {campaign.template ? campaign.template.name : <i>No Template</i>}</p>

        <AceEditor
            onChange={onChange}
            placeholder="Write your content"
            mode="plain_text"
            theme="solarized_light"
            name="template editor"
            fontSize={16}
            showPrintMargin={true}
            showGutter={true}
            width="100%"
            height='300px'
            highlightActiveLine={true}
            value={content}
            readOnly={true}
            setOptions={{
                enableBasicAutocompletion: true,
                enableLiveAutocompletion: true,
                enableSnippets: true,
                showLineNumbers: true,
                tabSize: 4,
            }}
        />

        <br />
        <iframe 
                id="rendered"
                srcDoc={renderedContent}
                style={{width: "100%", height: "100%", minHeight: "300px"}} 
        />

        {/*
        <label for="subject">Subject</label>
        <div  style="border: 2px solid gray; border-radius: 4px; margin-top: 8px; margin-bottom: 8px; padding-left: 8px">
            <pre id="subject">{{ campaign.Subject }}</pre>
        </div>

        <label for="body">Body</label>
        <div id="editor">{{campaign.Body}}</div>

        <label for="rendered">Rendered</label>
        <div style="width: 100%; min-width: 500px; height: 100%">
            <iframe 
                id="rendered"
                src="{{ reverse("render-ondemand", campaign.ID) }}?body={{campaign.Body}}&templateID={{campaign.Template.ID}}" 
                frameborder="0"
                style="width: 100%; height: 100%; min-height: 300px;" 
            ></iframe>
        </div>

        <h2>Events</h2>
        <form action="{{ reverse("events-create") }}" method="post">
            {{ csrfField | safe }}
            
            <input type="hidden" name="campaign_id" value="{{ campaign.ID }}">
            <input type="submit" value="Sent Blast Emails">
        </form>
        {% if campaign.IsNoEvent() %}
            <p><i>No event</i></p>
        {% else %}
        <div style="max-height: 200px; overflow: scroll;">
            <table>
                <thead>
                    <th>Status</th>
                    <th>Detail</th>
                    <th>Date</th>
                </thead>
                <tbody style="flex-direction: row-reverse">
                {% for event in campaign.Events %}
                    <tr>
                        <td>{{event.Status}}</td>
                        <td>
                            {% if event.Detail == "" %}
                            -
                            {% else %}
                            {{event.Detail}}
                            {% endif %}
                        </td>
                        <td>{{event.CreatedAt | date:"2006-02-01 15:04"}}</td>
                    </tr>
                {% endfor %}
                </tbody>
            </table>
        </div>
        {% endif %}
        */}
    </main> 
}

export default ListCampaign