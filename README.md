# Grupo15-Laboratorio-1

# Integrantes
* Dante Aspee - Rol 202073524-7
* Vicente Gaete - Rol 202004604-2
* Bernardo Pinninghoff - Rol 201973543-8

# Getting started

## Comandos docker

Para instrucciones sobre como armar cada contenedor con la makefile, usar el siguiente comando:
`make help`

## Nuestras VMs

Máquina - Contraseña

- VM1: dist057 - RQqxqq2H4W2U

- VM2: dist058 - SmcyWJG4EhNJ

- VM3: dist059 - ad6AwejY22VW

- VM4: dist060 - TbfeSr3ZTwMX

## Distribución VMs y Containers

* VM1: container-america , container-central
* VM2: container-asia
* VM3: container-europa , container-rabbitmq
* VM4: container-oceania

## Consideraciones

No nos dio el tiempo, por lo tanto hay errores que no pudimos solucionar. Estos son:

- Por algun motivo, prints en la consola que antes funcionaban sin problema dejaron de aparecer y no sabemos como ni por que.

- El servidor central se queda pegado en su segunda iteracion, antes de procesar los mensajes consumidos por la cola de rabbit. Sospechamos que se debe a que a la cola no le llega ningun mensaje de servidores regionales despues de la primera iteracion, y como esta vacia se queda esperando que le lleguen mensajes.

- Los servidores regionales no restan a los usuarios registrados de su total. Lo hubieramos solucionado pero nos dimos cuenta muy tarde.
