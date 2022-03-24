# Servicio de Usuarios
Este servicio proporciona acciones CRUD y autenticación de usuarios mediante JWT token.

## Variables de entorno

| Variable               | Descripción                                                              |
|------------------------|--------------------------------------------------------------------------|
| APP_PORT               | Puerto del servidor HTTP                                                 |
| DATABASE_USERNAME      | Nombre de usuario de la conexión a la base de datos                      |
| DATABASE_PASSWORD      | Contraseña de usuario de la conexión a la base de datos                  |
| DATABASE_HOSTNAME      | Nombre del servidor de base de datos                                     |
| DATABASE_NAME          | Nombre de la base de datos                                               |
| SESSION_SECRET_KEY     | Clave secreta de la sesión                                               |
| ACCESS_TOKEN_DURATION  | Duración del token de la sesión, por ejemplo: 15m (15 minutos)           |
| REFRESH_TOKEN_DURATION | Duración del token de refresco de la sesión, por ejemplo: 24h (24 horas) |

### Entorno de desarrollo
Para poder hacer uso de las variables de entorno en desarrollo debemos crear un archivo `app.env` en la raíz del directorio. Cualquier variable que se deba agregar, también debe ser agregada en el archivo `utils/config.go` para que estas puedan ser accedidas desde la aplicación.

### Entorno de producción
Para configurar las variables de entorno se deben exponer en el contenedor mediante el parámetro `-e VARIABLE=VALOR`. Para más información, leer la [documentación oficial de Docker](https://docs.docker.com/engine/reference/commandline/run/).

## Migraciones
Al ejecutar por primera vez el servidor, en caso de no tener configurado la base de datos, podemos correr el comando `go run cmd/migration.go` desde la raíz del directorio para crear la base de datos, las tablas de usuarios y sesiones, así como crear un usuario nuevo, cuyos datos de acceso son otorgados en la CLI.

## Documentación
La documentación de la API está integrada y es accesible en la ruta `/docs/index.html`, que nos proporciona una interfaz gráfica "amigable" para realizar pruebas e identificar modelos y estructuras.

Para crear la documentación necesitamos tener instalado la herramienta CLI de [Swagger](https://github.com/swaggo/swag#getting-started), luego escribimos `swag init --parseDependency --parseInternal --parseDepth 1` en la raíz de nuestro directorio.

### Autenticación
En la ruta `/login` podemos ingresar con el usuario creado en las migraciones iniciales. Luego, en el modal que se abre al presionar en el botón "Authorize" ingresamos "Bearer <ACCESS TOKEN>" y podremos hacer las consultas a los endpoints que sean restringidos:

![](https://i.imgur.com/udxMlex.png)

## Endpoints
La ruta base de la API es `/api`, por lo que las llamadas a los endpoints deben hacerse a `<host>:<puerto>/api/<endpoint>/<id>` donde `<id>` es el ID del recurso.

![](https://i.imgur.com/TeCfZoh.png)

Cada llamada a la API debe contener en la cabecera `Authorization` el token de acceso (`access_token`) obtenido luego del proceso de login de usuario.

## Levantar servidor de pruebas local
1. Instalar Go en nuestro equipo y agregar el `GOPATH` a nuestras variables de entorno.
2. Clonar el repositorio usando Git: `git clone https://github.com/maramal/user-service.git` en la ruta `src` de Go de nuestro usuario:
    * En Windows: `%GOPATH%\src\github.com\maramal`
    * En Linux o Mac: `$GOPATH/src/github.com/maramal`
3. Levantar un servidor MySQL.
4. Crear un archivo `app.env` en la raíz del directorio con las variables de entorno necesarias (Ver **Variables de entorno**) y sus respectivos valores.
5. Correr el comando de migraciones (Ver **Migraciones**) y copiar el nombre de usuario y contraseña que aparecen en el resultado de las migraciones en el CLI.
6. Levantar el servidor con el comando `go run main.go` en la raíz del directorio del servicio.
7. Crear la documentación mediante el comando de creación de documentación (Ver **Documentación**).
8. Abrir el navegador y acceder a la documentación (Ver **Documentación**).