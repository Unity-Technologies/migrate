package source

import (
	"bytes"
	"io"
	text "text/template"
)

// Templater is a Driver decorator that allows for template parameters to be passed in.
type Templater struct {
	Driver Driver
	params map[string]interface{}
}

func NewTemplater(d Driver, params map[string]interface{}) *Templater {
	return &Templater{
		params: params,
		Driver: d,
	}
}

func (t *Templater) Open(url string) (Driver, error) {
	return t.Driver.Open(url)
}

func (t *Templater) Close() error {
	return t.Driver.Close()
}

func (t *Templater) First() (version uint, err error) {
	return t.Driver.First()
}

func (t *Templater) Prev(version uint) (prevVersion uint, err error) {
	return t.Driver.Prev(version)
}

func (t *Templater) Next(version uint) (nextVersion uint, err error) {
	return t.Driver.Next(version)
}

func (t *Templater) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	r, identifier, err = t.Driver.ReadUp(version)
	if err != nil {
		return
	}
	return t.template(r, identifier)
}

func (t *Templater) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	r, identifier, err = t.Driver.ReadDown(version)
	if err != nil {
		return
	}
	return t.template(r, identifier)
}

func (t *Templater) template(r io.ReadCloser, identifier string) (io.ReadCloser, string, error) {
	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, r)
	if err != nil {
		return nil, "", err
	}
	// Close the original reader here, because we're going to replace it with a templatized version.
	r.Close()

	tpl, err := text.New(identifier).Parse(buffer.String())
	if err != nil {
		return nil, "", err
	}

	buffer.Reset()
	err = tpl.Execute(&buffer, t.params)
	if err != nil {
		return nil, "", err
	}
	return io.NopCloser(bytes.NewReader(buffer.Bytes())), identifier, nil
}
