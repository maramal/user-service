# Servicio de Usuarios
Este servicio proporciona acciones CRUD y autenticación de usuarios mediante JWT token.

# CORS
No implementado. Para modificar el comportamiento actual (aceptar todos los orígenes) editar el archivo `middlewares/cors.go`.

# Variables de entorno

| Variable               | Descripción                                                              |
|------------------------|--------------------------------------------------------------------------|
| APP_PORT               | Puerto del servidor HTTP                                                 |
| MONGO_URI              | Dirección URL de la base de datos                                        |
| SESSION_SECRET_KEY     | Clave secreta utilizada para la creación del token. Debe ser de al menos 32 de caracteres de largo                     |
| ACCESS_TOKEN_DURATION  | Duración del token de la sesión, por ejemplo: 15m (15 minutos)           |
| REFRESH_TOKEN_DURATION | Duración del token de refresco de la sesión, por ejemplo: 24h (24 horas) |
| APM_APPNAME            | Nombre de la aplicación de APM (New Relic)                               |
| APM_LICENSE            | Licencia del APM (New Relic)                                             |

## Entorno de desarrollo
Para poder hacer uso de las variables de entorno en desarrollo debemos crear un archivo `app.env` en la raíz del directorio. Cualquier variable que se deba agregar, también debe ser agregada en el archivo `utils/config.go` para que estas puedan ser accedidas desde la aplicación.

## Entorno de producción
Para configurar las variables de entorno se deben exponer en el contenedor mediante el parámetro `-e VARIABLE=VALOR`. Para más información, leer la [documentación oficial de Docker](https://docs.docker.com/engine/reference/commandline/run/).

# Migraciones
Para crear un superadministrador, ejecutar el script `scripts/create-superadmin.go` con los argumentos siguientes:

| Argumento | Descripción                                                          |
|-----------|----------------------------------------------------------------------|
| firstname | Primer nombre del super administrador                                |
| lastname  | Apellido del super administrador                                     |
| email     | Correo electrónico del super administrador (utilizado para el login) |
| password  | Contraseña del super administrador                                   |

## Ejemplo
```cmd
go run scripts/create-superadmin.go --firstname Juan --lastname Perez --email juanp@mail.com --password 123456
```

# Documentación
La documentación de la API está integrada y es accesible en la ruta `/docs/index.html`, que nos proporciona una interfaz gráfica "amigable" para realizar pruebas e identificar modelos y estructuras.

Para crear y/o actualizar la documentación necesitamos tener instalado la herramienta CLI de [Swagger](https://github.com/swaggo/swag#getting-started), luego escribimos `swag init --parseDependency --parseInternal --parseDepth 1` en la raíz de nuestro directorio.

> **Nota 1:** La documentación debería ser actualizada cada vez que se modifican los handlers.

> **Nota 2:** En caso de obtener error al documentar código, cambiar el `--parseDepth` de `1` a `2` (por ejemplo).

# Autenticación
En la ruta `/login` podemos ingresar con el usuario creado en las migraciones iniciales. Luego, en el modal que se abre al presionar en el botón "Authorize" ingresamos "Bearer <ACCESS TOKEN>" y podremos hacer las consultas a los endpoints que sean restringidos:

![](https://i.imgur.com/udxMlex.png)

# Administración
Para restringir el acceso a administradores y superadministradores a ciertas rutas se debe utilizar el middleware `middlewares.AdminMiddleware()` en cada ruta / grupo de rutas. Por ejemplo:

```go
userRouter := router.Group("/users")
userRouter.Use(middlewares.AdminMiddleware())
```

# Endpoints
La ruta base de la API es `/api`, por lo que las llamadas a los endpoints deben hacerse a `<host>:<puerto>/api/<endpoint>/<id>` donde `<id>` es el ID del recurso.

![](https://i.imgur.com/TeCfZoh.png)

Cada llamada a la API debe contener en la cabecera `Authorization` el token de acceso (`access_token`) obtenido luego del proceso de login de usuario.

# Deployment

## Desarrollo
Para realizar tareas de desarrollo y/o pruebas en diferentes entornos, recomiendo las configuraciones siguientes:

### Levantar servidor de pruebas local
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

### Levantar servidor en contenedor Docker
1. Construir la imagen del proyecto con `docker build -t user-service:1.0 .` en la raíz del directorio.
2. Inicializar el contenedor con `docker run -p 8080:8080 --env APP_PORT=8080 --env MONGO_URI=mongodb+srv://<usuario>:<password>@cluster0.ejvre.mongodb.net/myFirstDatabase?retryWrites=true&w=majority --env APM_APPNAME=userservice --env APM_LICENSE=0d0bf600800495f7f39304bf09e5342a642cNRAL user-service`