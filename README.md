# Grupo15-Laboratorio-1

# Integrantes
* Dante Aspee - Rol 202073524-7
* Vicente Gaete - Rol 202004604-2
* Bernardo Pinninghoff - Rol 201973543-8

# Getting started

## Comandos docker

### RabbitMQ

Armar imagen: docker build -t \[nombre-imagen-rabbit\]:latest \[path/a/dockerfile\]

Armar contenedor con la imagen: docker run -d --name \[nombre-container\] -p 5673:5673 -p 15673:15673 \[nombre-imagen-rabbit\]

## Nuestras VMs

Máquina - Contraseña

- VM1: dist057 - RQqxqq2H4W2U

- VM2: dist058 - SmcyWJG4EhNJ

- VM3: dist059 - ad6AwejY22VW

- VM4: dist060 - TbfeSr3ZTwMX

## Procesos a implementar

### S. Central

* ~~Leer parametros_de_inicio.txt~~ Listo
* ~~Generar cantidad de llaves al azar en intervalo segun txt, registrar hora~~ Listo
* ~~Notificar con comunicacion **sincrona/gRPC** a los S.Regionales el # de llaves~~ Basicamente listo pero falta probarlo con los servidores regionales para saber si funciona/falta algo
* ~~Recibir solicitudes de llaves de S.Regionales con comunicacion **asincrona/RabbitMQ**~~ Idem que punto anterior
* ~~Registrar S.Regional y # de llaves entregadas y procesar cantidad de llaves de S.Central~~ Listo
* ~~Escribir en archivo~~ Listo
* ~~Enviar respuesta de forma **sincrona/gRPC** a S.Regionales~~ Listo pero sin testear
* ~~Repetir segun iteraciones indicadas en archivo txt~~ Listo
* ~~Contenedor Docker~~
* Testing

### S. Regional

* ~~Recibir notificacion de S.Central de # de llaves disp~~  Ready
* ~~Leer su propio parametros_de_inicio.txt~~ Ready
* ~~Generar # de solicitud de llaves al azar en intervalo segun txt~~ Ready
* ~~Enviar solicitud de llaves con comunicacion **asincrona/rabbitMQ** al S.Central~~ Falta testeo
* Recibir respuesta de llaves registrados por S.Central y restarlo para proxima solicitud
* Continuar segun iteraciones indicadas por S.Central


## Distribución VMs y Containers

* VM1: container-america , container-central
* VM2: container-asia
* VM3: container-europa , container-rabbitmq
* VM4: container-oceania
