package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	// Importa tu módulo personalizado para que Caddy lo pueda registrar.
	// Asegúrate de que el path coincida con la estructura de tu proyecto Go.
	_ "github.com/your-username/my-caddy-dns-check/module"
)

func main() {
	// Permite que Caddy se ejecute y use la funcionalidad normal, incluyendo la recarga.
	caddycmd.EnableReload()
	caddycmd.Main()
}
