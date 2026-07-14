package page

type Meta struct {
	Title        string
	Description  string
	Keywords     []string
	CanonicalURL string
	OGImage      string
	NoIndex      bool
}

type Page struct {
	Template string
	Layout   string
	Meta     Meta
	Data     any
}

func New(template string, data any) *Page {
	return &Page{
		Template: template,
		Data:     data,
	}
}

func (p *Page) WithLayout(layout string) *Page {
	p.Layout = layout
	return p
}

func (p *Page) WithMeta(meta Meta) *Page {
	p.Meta = meta
	return p
}
