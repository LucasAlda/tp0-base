echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
go run ./scripts/docker-compose-generator.go $1 $2
