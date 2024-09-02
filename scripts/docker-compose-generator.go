package main

import (
	"fmt"
	"os"
	"strconv"
)

func generateDockerCompose(filename string, numClients int) {
	// Base structure of the docker-compose file
	composeStr := `name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: /server
    volumes:
      - ./server/config.ini:/config.ini
    networks:
      - testing_net
`

	// Generate client services
	for i := 1; i <= numClients; i++ {
		clientStr := fmt.Sprintf(`
  client%d:
    container_name: client%d
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=%d
    volumes:
      - ./client/config.yaml:/config.yaml
    networks:
      - testing_net
    depends_on:
      - server
`, i, i, i)
		composeStr += clientStr
	}

	// Add network configuration
	composeStr += `
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
`

	// Write the docker-compose file
	err := os.WriteFile(filename, []byte(composeStr), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ./generar-compose.sh <output_filename> <number_of_clients>")
		os.Exit(1)
	}

	outputFilename := os.Args[1]
	numClients, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Error converting number of clients to integer: %v\n", err)
		os.Exit(1)
	}

	generateDockerCompose(outputFilename, numClients)
}
