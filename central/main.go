package main

import (
	"context"
	"log"
	"net"
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"

	"google.golang.org/grpc"
)

type server struct {}

func startupParameters(filePath string) (minKey, maxKey int, err error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, 0, err
	}

	cotas := strings.Split(string(content), "-")
	if len(cotas) != 2 {
		return 0, 0, fmt.Errorf("invalid file format in parametros_de_inicio.txt")
	}

	minKey, err = strconv.Atoi(strings.TrimSpace(cotas[0]))
	if err != nil {
		return 0, 0, err
	}

	maxKey, err = strconv.Atoi(strings.TrimSpace(cotas[1]))
	if err != nil {
		return 0, 0, err
	}

	return minKey, maxKey, nil
}

func main() {

}
