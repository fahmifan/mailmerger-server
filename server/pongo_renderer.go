package server

import (
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/flosch/pongo2"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

var _ echo.Renderer = (*PongoRenderer)(nil)

// PongoRenderer implements echo.Renderer
type PongoRenderer struct {
	*PongoRendererConfig
	globalData map[string]interface{}
}
type PongoRendererConfig struct {
	BaseURL      string
	RootDir      string
	DebugEnabled bool
}

// NewPongoRenderer ..
func NewPongoRenderer(cfg *PongoRendererConfig) *PongoRenderer {
	pr := &PongoRenderer{PongoRendererConfig: cfg}
	pr.loadGlobalData()
	return pr
}

func (r *PongoRenderer) loadGlobalData() {
	if r.globalData == nil {
		r.globalData = make(map[string]interface{})
	}

	nav := make(map[string]interface{})
	navPath := path.Join(r.RootDir, "data/navigation.json")
	bt, err := os.ReadFile(navPath)
	if err != nil {
		return
	}

	if err = json.Unmarshal(bt, &nav); err != nil {
		return
	}

	r.globalData["navigation"] = nav
}

// Render implement echo.Renderer
func (r *PongoRenderer) Render(w io.Writer, name string, data interface{}, ec echo.Context) (err error) {
	var ctx pongo2.Context
	if data != nil {
		switch val := data.(type) {
		case pongo2.Context:
			ctx = pongo2.Context(val)
		case echo.Map:
			ctx = pongo2.Context(val)
		case map[string]interface{}:
			ctx = pongo2.Context(val)
		default:
			ctx = make(pongo2.Context)
		}
	} else {
		ctx = make(pongo2.Context)
	}

	tpl, err := r.getTemplate(name)
	if err != nil {
		log.Err(err).Msg("getTemplate")
		return err
	}

	ctx[csrf.TemplateTag] = csrf.TemplateField(ec.Request())
	ctx["baseURL"] = r.BaseURL
	ctx["debugEnabled"] = r.DebugEnabled
	ctx["reverse"] = ec.Echo().Reverse

	for k, v := range r.globalData {
		ctx[k] = v
	}

	if err = tpl.ExecuteWriter(ctx, w); err != nil {
		log.Err(err).Msg("exec writer")
		return err
	}
	return nil
}

func (r *PongoRenderer) getTemplate(name string) (tpl *pongo2.Template, err error) {
	name = path.Join(r.RootDir, name)
	if r.DebugEnabled {
		tpl, err = pongo2.FromFile(name)
	} else {
		tpl, err = pongo2.FromCache(name)
	}

	return
}
