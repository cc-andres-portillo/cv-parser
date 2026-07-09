package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	"github.com/cc-andres-portillo/cv-parser/internal/adapters/storages"
	"github.com/cc-andres-portillo/cv-parser/internal/core/ports"
)

func main() {
	data, err := ioutil.ReadFile("sample_good.zip")
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(data)
	svc := ports.NewCVService(storage.NewDocumentExtractorAdapter(), storage.NewMockDatabaseAdapter())
	res, err := svc.ParseCV(context.Background(), reader, int64(len(data)), ".docx")
	fmt.Println("err:", err)
	fmt.Printf("res: %+v\n", res)
}
