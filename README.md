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

Para permitir editar los archivos de configuración dentro del contenedor sin necesidad de volver a generar la imagen, se utilizó volumes para montar los archivos en el contenedor.

Estos volúmenes se encuentran definidos en el docker compose y montan en los contenedores indicados el archivo de configuración ubicado en la raíz del proyecto, de manera que no hay que mover los archivos de lugar y es transparente para el usuario.

### Ejercicio N°3:

Para validar el funcionamiento del servidor se utilizó el script `validar-echo-server.sh` que se encuentra en la raíz del proyecto. Este script se ejecuta en el contenedor de alpine y se encarga de probar el servidor con netcat.

El script valida que el servidor en la network del `docker-compose-dev.yaml` escuchando en `server:12345` responda con "ping" cuando se le envía "ping". Para esto se utiliza el contenedor de alpine como un cliente netcat que se conecta al servidor para evitar que el usuario tenga que instalar netcat en su sistema.

```bash
./validar-echo-server.sh
```

### Ejercicio N°4:

## Parte 2: Repaso de Comunicaciones

### Ejercicio N°5:

### Ejercicio N°6:

### Ejercicio N°7:

## Parte 3: Repaso de Concurrencia

### Ejercicio N°8:
