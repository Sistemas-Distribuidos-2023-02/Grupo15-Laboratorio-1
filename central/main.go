package main

import (
	"context"
	"log"
	"net"
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
	"math/rand"
	"time"

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

func keygen(minKey, maxKey int) []int {
	rand.Seed(time.Now().UnixNano())

	n := rand.Intn(maxKey - minKey + 1) + minKey

	keys := make([]int, n)
	for i := 0; i < n; i++ {
		keys[i] = minKey + rand.Intn(n)
	}
	return keys
}

func main() {
	filePath := "./parametros_de_inicio.txt"

	minKey, maxKey, err := startupParameters(filePath)
	if err != nil {
		fmt.Printf("Error reading startup_parameters: %v\n", err)
		return
	}

	keys := keygen(minKey, maxKey)
	fmt.Printf("Keys: %v\n", keys)
}
