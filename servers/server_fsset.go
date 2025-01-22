//go:build allowfsset

package servers

import "github.com/pufferpanel/pufferpanel/v3/files"

func (p *Server) SetFileServer(fs files.FileServer) {
	p.fileServer = fs
}
