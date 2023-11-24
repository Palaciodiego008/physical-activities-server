# Aplicación de Detección de Actividades Sospechosas

Esta es una aplicación desarrollada en Go (Golang) para la detección de actividades físicas sospechosas basadas en archivos CSV de datos. La aplicación incluye tanto el backend en Go como el frontend en React.js.

## Contenido

1. [Estructura del Proyecto](#estructura-del-proyecto)
2. [Tecnologías Utilizadas](#tecnologías-utilizadas)
3. [Instrucciones de Uso](#instrucciones-de-uso)
4. [Detalles del Backend](#detalles-del-backend)
5. [Detalles del Frontend](#detalles-del-frontend)
6. [Contribuir](#contribuir)
7. [Licencia](#licencia)

## Estructura del Proyecto

```bash         # Código fuente del backend (Go)
│   ├── main.go      # Punto de entrada principal
│   ├── uploads      # Directorio para almacenar archivos cargados
│   ├── go.mod       # Archivo de configuración de módulos Go
│   └── go.sum       # Archivo de resumen de módulos Go
```

## Tecnologías Utilizadas

- **Backend (Go)**:
  - [Gorilla Mux](https://github.com/gorilla/mux): Enrutador HTTP para Go.
  - [Handlers](https://github.com/gorilla/handlers): Middleware CORS para Go.

- **Frontend (React.js)**:
  - [Axios](https://github.com/axios/axios): Cliente HTTP para realizar peticiones a la API.
  - [React Dropzone](https://github.com/react-dropzone/react-dropzone): Componente de carga de archivos.
  - [Tailwind CSS](https://tailwindcss.com/): Marco de diseño CSS de utilidad de bajo nivel.

## Instrucciones de Uso

1. Clona este repositorio:
2. Instala las dependencias del frontend y backend:

```bash
# En el directorio frontend/
npm install

# En el directorio backend/
go get
```

3. Inicia el frontend y el backend (asegúrate de tener ambos entornos en ejecución):

```bash
# En el directorio frontend/
npm start

# En el directorio backend/
go run main.go
```

4. Abre tu navegador y visita `http://localhost:3000` para acceder a la aplicación.

## Detalles del Backend

### Estructura del Código

- **`main.go`**: Punto de entrada principal.
- **`uploads/`**: Directorio para almacenar archivos cargados.
- **`go.mod` y `go.sum`**: Archivos de configuración de módulos Go.

### Endpoints

- **`/upload` (POST)**: Maneja la carga de archivos CSV.
- **`/activities` (GET)**: Retorna las actividades sospechosas basadas en el archivo más reciente.

## Detalles del Frontend

### Estructura del Código

- **`public/` y `src/`**: Directorios estándar de un proyecto React.js.

### Componentes Principales

- **`FileUpload`**: Componente para cargar archivos CSV.
- **`SuspiciousActivities`**: Componente para mostrar las actividades sospechosas.

## Licencia
Este proyecto está bajo la Licencia MIT. Consulta el archivo [LICENSE](LICENSE) para más detalles.