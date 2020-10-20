// package main. Šeit uzrakstīšu vispārīgās problēmas, un ierosinājumus uzalbojumiem :)
// 1. Mums pipeline tiek izpildīts `go vet ./...` un `golint ./...`, šajā projektā ir lint kļūdas.
// 2. Go source kodam ir viens pareizais formāts, kuru var iegūt, izpildot `go fmt ./...`
// 3. Redzu, ka izmanto kaut kādu interneta resursu kā datubāzi. Ierosinu šādu resursu pielietošanai izmantot docker.
package main

import (
	"goTestProj/API/Initialization"
)

func main() {
	Initialization.Initialize()
}