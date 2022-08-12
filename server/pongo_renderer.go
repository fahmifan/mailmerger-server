package server

import (
	"io"
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
}
type PongoRendererConfig struct {
	BaseURL      string
	RootDir      string
	DebugEnabled bool
}

// NewPongoRenderer ..
func NewPongoRenderer(cfg *PongoRendererConfig) *PongoRenderer {
	return &PongoRenderer{cfg}
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
