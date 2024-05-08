package mediatype

import (
	"mime"
	"path/filepath"
)

const (
	// JupyterNotebookMediaType is the local IANA compliant mimetype to indicate a python notebook
	// using this because it is mentioned here: https://discourse.jupyter.org/t/i-cant-download-my-notebook-as-ipynb-anymore-it-saves-as-json/7043/4
	// Another common media type is application/x-ipynb+json
	JupyterNotebookMediaType = "application/x.jupyter.notebook+json"

	// JupyterNotebookExtension indicates an interactive python notebook
	JupyterNotebookExtension = ".ipynb"
)

// DetermineType read the file extension of the file, and returns a type to be included when creating a public artifact
func DetermineType(path string) string {
	fExt := filepath.Ext(path)
	if fExt == JupyterNotebookExtension {
		return JupyterNotebookMediaType
	}
	// TODO we could also consider using http.DetectContentType(buf) to use the contents of the file to determine the content type
	return mime.TypeByExtension(fExt)
}
