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

Para agregar el graceful shutdown tanto en el cliente como en el servidor se utilizó el paquete `os/signal` para capturar las señales de interrupción y terminar los loops de los respectivos servicios.

En el servidor se utiliza el context para que una vez este es cancelado, cierre el listener y finalice el loop de aceptar nuevas conexiones y el programa finalice correctamente una vez que la conexion actual termina.

En el cliente se reemplaza el time.Sleep por un select entre el contexto y un time.After, de forma que si el contexto se cancela, el select case se ejecuta y proceso termina correctamente, caso contrario se queda esperando que el tiempo de sleep termine y se vuelve a ejecutar el loop.

En ambos casos se utiliza defer para asegurar que las conexiones se cierren una vez terminado su closure, en el cliente se cierra cuando termina la funcion `sendMessage` y en el servidor se cierra cuando termina el handle de la nueva conexion.

## Parte 2: Repaso de Comunicaciones

### Ejercicio N°5:

La implementacion del protocolo se encuentra en la carpeta `shared/protocol` y es utilizada tanto por el cliente como por el servidor.

El protocolo cuenta con 2 secciones:

- Messages: Encargada de la serializacion y deserializacion de los datos.
- Network: Encargada del Read y Write de los datos en el socket.

#### Messages

Los mensajes se definen como estructuras que implementan la interfaz `Message` y tienen que implementar los siguientes metodos:

- `GetMessageType() MessageType`: Devuelve el tipo de mensaje.
- `Encode() string`: Devuelve el mensaje serializado.
- `Decode(data string) error`: Decodifica un string en un mensaje.

`MessageType` es un enum que se utiliza para identificar el tipo de mensaje.

#### Network

La capa de network se encarga de la lectura y escritura de los mensajes en el socket. Para esto se utiliza la funcion `Send` para enviar un mensaje y `Receive` para recibir un mensaje.

El payload de los mensajes cuenta con 3 partes:

- `size`: 4 bytes que indican el tamaño del mensaje. Es necesario ya que el mensaje puede ser de un tamaño variable y este dato es usado para evitar short-reads y short-writes.
- `messageType`: 4 bytes que indican el tipo de mensaje. Es necesario ya que el mensaje puede ser de un tipo variable y este dato es usado para deserializar el mensaje correctamente.
- `data`: El mensaje serializado. Es el string que se le pasa a la funcion `Decode` de un struct que implemente la interfaz `Message` para deserializar el mensaje.

`Send` recibe una conexion y un mensaje que implemente la interfaz `Message`, serializa el mensaje y lo envía a través del socket utilizando el protocolo definido anteriormente.

`Receive` recibe una conexion y devuelve un struct llamado `ReceivedMessage` que contiene el tipo de mensaje, size y string recibido. Este struct es lo que devuelve la funcion `Receive` y en la logica del servidor se utiliza el `MessageType` para ejecutar el `Decode` del struct correspondiente.

#### Mensajes implementados

- `MessageBet`
- `MessageBetAck`

```go
type MessageBet struct {
	FirstName string
	LastName  string
	Document  string
	Birthdate string
	Number    string
}

type MessageBetAck struct {
	MessageType int32
}
```

#### Funcionamiento del cliente

El cliente lee del archivo de configuracion los datos de la apuesta y los utiliza para crear un mensaje `MessageBet` que es enviado al servidor con `Send`.

El servidor recibe el `ReceivedMessage` con el mensaje de apuesta utilizando `Receive`, verifica que el MessageType sea `MessageTypeBet` y lo decodifica utilizando el metodo `Decode` de `MessageBet`. Una vez decodificado, crea el `Bet` y guarda los datos en el csv para, finalmente, mandar un mensaje de confirmacion `MessageBetAck` con la funcion `Send`.

El cliente recibe el mensaje de confirmacion utilizando `Receive` y decodificando luego de revisar que sea del tipo `MessageTypeAck`, lo imprime por pantalla y finaliza.

#### Configuracion del cliente

La apuesta del cliente se define en el archivo de configuracion `client/config.yaml` agregando un apartado `bet` y se utiliza para crear el mensaje `MessageBet` que es enviado al servidor.

```yml
bet:
  firstName: "john"
  lastName: "doe"
  document: "43000000"
  birthdate: "2002-01-04"
  number: "1000"
```

Tambien se puede pisar cualquier valor con las variables de entorno:

```bash
export CLI_NOMBRE="juan"
export CLI_APELLIDO="perez"
export CLI_DOCUMENTO="43000000"
export CLI_NACIMIENTO="2002-01-05"
export CLI_NUMERO="1001"
```

### Ejercicio N°6:

Para implementar el envio de multiples apuestas en batch no se tuvo que alterar tanto la logica del cliente como del servidor, mostrando una buena encapsulamiento de la capa de network. Para esto se creó un mensaje `MessageBetBatch` que contiene un slice de `MessageBet` y un `MessageAllBetsSent` que es enviado al servidor cuando el cliente termina de enviar todas las apuestas.

El servidor ahora tiene un loop infinito que acepta nuevas conexiones y las agrega a la lista de clientes. Este loop utiliza `Receive` para obtener un `ReceivedMessage` del cliente, y es esta abstraccion la que permite que el servidor
haga un switch sobre el MessageType y se encarge de procesar los batchs de apuestas y los `MessageAllBetsSent` dependiendo del caso.

Por cada `MessageBetBatch` que llega al servidor, este es decodificado y los bets son agregados a una lista de apuestas. Cuando se llega al final del batch, se procesan las apuestas, se hace el append al archivo csv y se envía un `MessageBetAck` al cliente.

Cuando el servidor recibe un `MessageAllBetsSent` este es procesado y se desconecta el cliente para proceder con el sigiuente.

Con el siguiente comando se puede validar un `wc` al archivo `bets.csv` que se encuentra en el contenedor del servidor:

```bash
docker exec server wc bets.csv
```

### Ejercicio N°7:

## Parte 3: Repaso de Concurrencia

### Ejercicio N°8:
