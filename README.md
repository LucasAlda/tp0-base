# TP0: Docker + Comunicaciones + Concurrencia

## Previo a los ejercicios

Una vez leído el enunciado, tome la decision de migrar el server a go para poder practicar el lenguaje y su manejo de la concurrencia, que no solo me resultaba interesante sino que es distinto a lo que hemos utilizado en los anteriores materias (como el de python que ya había sido usado en Redes).

Para la migracion se mantuvo la arquitectura original, respetando archivos, funciones, tests y logica en general para que no sea afectado el desarrollo de los ejercicios. Las librerias externas se mantuvieron similares a las usadas en el cliente (como viper).

## Parte 1: Introducción a Docker

### Ejercicio N°1:

Para generar los archivos de docker-compose.yaml se utilizó el script `generar-compose.sh` que se encuentra en la raíz del proyecto. Este script recibe como parámetros el archivo de configuración y la cantidad de clientes a levantar.

Cuando se ejecuta el script de bash, se corre el go encontrado en `scripts/docker-compose-generator.go`, por lo que es necesario tener instalado go en el sistema host.

```bash
./generar-compose.sh docker-compose-dev.yaml 5
```

### Ejercicio N°2:

### Ejercicio N°3:

### Ejercicio N°4:

## Parte 2: Repaso de Comunicaciones

### Ejercicio N°5:

### Ejercicio N°6:

### Ejercicio N°7:

## Parte 3: Repaso de Concurrencia

### Ejercicio N°8:
