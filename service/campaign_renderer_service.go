package service

import (
	"bytes"
	"io"

	"github.com/flosch/pongo2"
)

type CampaignRenderer struct {
	body     string
	template *Template
}

func NewCampaignRenderer(body string, template *Template) CampaignRenderer {
	return CampaignRenderer{body: body, template: template}
}

func (c CampaignRenderer) Render() (io.Reader, error) {
	res := bytes.NewBuffer(nil)

	pongoTpl, err := pongo2.FromString(c.template.HTML)
	if err != nil {
		return nil, err
	}

	err = pongoTpl.ExecuteWriter(pongo2.Context{"body": c.body}, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
