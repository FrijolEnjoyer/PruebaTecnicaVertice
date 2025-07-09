# 🚀 API Prueba Técnica Vértice

Este proyecto es una API REST desarrollada en **Go (Golang)** con el framework **Gin** que gestiona usuarios, productos y órdenes.
Incluye autenticación JWT, documentación Swagger, y se despliega fácilmente con **Docker Compose**.

---

## ✅ Tecnologías Usadas

* **Golang** (Gin, Gorm)
* **MySQL**
* **Swagger** (Documentación API)
* **Docker & Docker Compose**

---

## 📂 Estructura del Proyecto

```
Api/
🌟
🔹🔹 cmd/                     # Punto de entrada (main.go)
🔹🔹 docs/                    # Documentación generada por Swagger
🔹🔹 handler/                 # Controladores (Users, Products, Orders)
🔹🔹 models/                  # Modelos (ORM)
🔹🔹 repo/                    # Repositorios (acceso a la DB)
🔹🔹 services/                # Servicios (lógica de negocio)
🔹🔹 utils/                   # Utilidades (JWT, Hash, etc.)
🔹🔹 server/                  # Configuración del servidor y rutas
🌟
🔹🔹 docker-compose.yml       # Configuración de Docker Compose
🔹🔹 Dockerfile               # Imagen Docker multi-stage
🔹🔹 .env.example             # Archivo de ejemplo de variables de entorno
🌟
🔹🔹 README.md                # Este archivo
```

---

## ⚙️ Requisitos Previos

* Tener instalado:

  * [Docker](https://docs.docker.com/get-docker/)
  * [Docker Compose](https://docs.docker.com/compose/install/)

---

## 📋 Configuración de Variables de Entorno

1. Crea un archivo `.env`:

```bash
cp .env.example .env
```

2. Completa el archivo `.env` con las variables necesarias:

```env
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_PORT=
DB_HOST=

SECRET_KEY=
TIME_TOKEN=
TIME_REFRESH_TOKEN=
```

---

## 🚀 Cómo Levantar el Proyecto

### 1️⃣ Clona el repositorio:

```bash
git clone https://github.com/FrijolEnjoyer/PruebaTecnicaVertice.git 
cd tu-repositorio
```

### 2️⃣ Levanta los contenedores:

```bash
docker-compose up -d --build
```

* Esto construye las imágenes y levanta los servicios:

  * **API**: [http://localhost:8080](http://localhost:8080)
  * **MySQL**: localhost:3306

### 3️⃣ Verifica que estén corriendo:

```bash
docker-compose ps
```

### 4️⃣ (Opcional) Ver Logs:

```bash
docker-compose logs -f
```

---

## 📝 Documentación Swagger (API Docs)

Cuando el proyecto esté corriendo, accede a:

```
http://localhost:8080/swagger/index.html
```

Ahí encontrarás todas las rutas documentadas.

---

## 🥪 Ejecutar Tests 

Ejecuta el siguiente script:

chmod +x run_tests.sh
./run_tests.sh
---

## 📛 Detener los Servicios

```bash
docker-compose down
```

---

## ✅ Servicios Incluidos

| Servicio | Descripción                | Puerto |
| -------- | -------------------------- | ------ |
| app      | Aplicación principal (API) | 8080   |
| db       | Base de datos MySQL        | 3306   |

---

## ✅ Rutas Principales (Resumen)

* **POST /api/auth/register** → Crear usuario
* **POST /api/auth/login** → Login usuario (JWT)
* **GET /api/auth/me** → Obtener usuario logueado *(JWT Requerido)*
* **CRUD Productos y Órdenes** → Protegidos con JWT

---

## 🔒 Seguridad

* Autenticación con JWT (Bearer Token).
* Todas las rutas protegidas requieren el header:

```bash
Authorization: Bearer {token}
```

---
