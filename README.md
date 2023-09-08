# Grupo15-Laboratorio-1

# Integrantes
* Dante Aspee - Rol 
* Vicente Gaete - Rol 202004604-2
* Bernardo Pinninghoff - Rol 201973543-8

# Getting started

## Nuestras VMs

Máquina - Contraseña

- VM1: dist057 - RQqxqq2H4W2U

- VM2: dist058 - SmcyWJG4EhNJ

- VM3: dist059 - ad6AwejY22VW

- VM4: dist060 - TbfeSr3ZTwMX

## Procesos a implementar

### S. Central

* Leer parametros_de_inicio.txt
* Generar cantidad de llaves al azar en intervalo segun txt, registrar hora
* Notificar con comunicacion **sincrona/gRPC** a los S.Regionales el # de llaves
* Recibir solicitudes de llaves de S.Regionales con comunicacion **asincrona/RabbitMQ**
* Registrar S.Regional y # de llaves entregadas y procesar cantidad de llaves de S.Central
* Escribir en archivo 
* Enviar respuesta de forma **sincrona/gRPC** a S.Regionales
* Repetir segun iteraciones indicadas en archivo txt

### S. Regional

* Recibir notificacion de S.Central de # de llaves disp
* Leer su propio parametros_de_inicio.txt
* Generar # de solicitud de llaves al azar en intervalo segun txt
* Enviar solicitud de llaves con comunicacion **sincrona/gRPC** al S.Central
* Recibir respuesta de llaves registrados por S.Central y restarlo para proxima solicitud
* Continuar segun iteraciones indicadas por S.Central

### Cola Rabbit

* 

## Distribución VMs y Containers

* VM1: container-america , container-central
* VM2: container-asia
* VM3: container-europa , container-rabbitmq
* VM4: container-oceania
