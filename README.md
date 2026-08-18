# IAM Microservice (Identity and Access Management)

Este proyecto consiste en la implementación completa de un sistema de Gestión de Identidad y Acceso (IAM) basado en una arquitectura de microservicios. Partiendo de una base de datos preexistente y documentación de requerimientos, se construyó tanto el backend como el frontend para cumplir con todas las funcionalidades solicitadas.

## Tecnologías Utilizadas

- **Backend**: Go (Golang)
- **Frontend**: React.js
- **Base de Datos**: PostgreSQL
- **Seguridad**: JSON Web Tokens (JWT), encriptación bcrypt.

## ¿Qué se construyó?

A partir del esquema de base de datos (`identity`, `identity_audit`, `rbac`), el proyecto se desarrolló siguiendo buenas prácticas y arquitectura limpia, logrando las siguientes funcionalidades principales:

1. **Autenticación y Autorización (Login & Registro):**
   - Registro de usuarios encriptando contraseñas usando bcrypt.
   - Inicio de sesión con validación de credenciales.
   - Emisión de tokens de acceso (JWT).

2. **Gestión de Sesiones (Refresh Tokens):**
   - Implementación de tokens opacos (Opaque Tokens) guardados en la base de datos para manejar el ciclo de vida de la sesión.
   - Endpoints para rotación y refresco seguro de los Access Tokens sin exponer las credenciales nuevamente.
   - Endpoint de cierre de sesión (Logout) que revoca el Refresh Token activamente.

3. **Recuperación de Contraseñas:**
   - Flujo de *Forgot Password* para solicitar un token temporal de recuperación.
   - Flujo de *Reset Password* validando la expiración y uso del token temporal para actualizar la contraseña de manera segura.

4. **Control de Acceso Basado en Roles (RBAC):**
   - Soporte nativo para múltiples roles.
   - Inyección dinámica de los roles autorizados dentro del payload del JWT.
   - Middlewares en el backend (`RequireAuth` y `RequireRole`) para proteger las rutas.

5. **Panel de Administración y Auditoría:**
   - Dashboard exclusivo para usuarios con el rol `ADMIN`.
   - **Gestión de Usuarios:** Interfaz para listar a todos los usuarios, visualizar sus roles actuales y asignar o revocar permisos con un solo clic de forma visual (usando etiquetas/chips).
   - **Auditoría de Accesos:** Tabla de historial paginada que lee de `identity_audit.audit_login`, permitiendo a los administradores observar quién, cuándo y con qué resultado (Éxito, Contraseña Inválida, etc.) intentó ingresar al sistema, destacando cada evento con colores semánticos.

## Estructura del Código

- `backend/`: Código fuente en Go estructurado en capas (Domain, Application Ports, Usecases, Persistence Postgres, Transport HTTP).
- `frontend/`: Aplicación React estructurada en páginas (Login, Register, Dashboard, Admin, Audit) y servicios (Axios interceptors).
- `database/`: Scripts de inicialización y esquemas de la BD en PostgreSQL.

---

## Demostración de Funcionamiento

En el siguiente video se muestra en detalle cómo quedaron las interfaces de React interactuando con los endpoints de Go, pasando por el registro, inicio de sesión, rotación de tokens, asignación de roles y la revisión de la auditoría:

[👉 **Ver Video Demostrativo aquí**](#) *(Reemplaza este `#` con el link al video real)*

---

## Cómo Ejecutar el Proyecto Localmente

### Requisitos Previos
- **Go** (v1.20+)
- **Node.js** (v18+)
- **PostgreSQL** ejecutándose localmente.

### 1. Base de Datos
Asegúrate de ejecutar los scripts SQL ubicados en la carpeta `database/01_ddl/` dentro de tu instancia de PostgreSQL para crear las tablas necesarias.

### 2. Ejecutar el Backend (Go)
Abre una terminal, sitúate en la carpeta del backend y ejecuta el servidor:
```bash
cd backend
go run cmd/server/main.go
```
El backend se ejecutará por defecto en el puerto `8080`.

### 3. Ejecutar el Frontend (React)
Abre otra terminal, sitúate en la carpeta del frontend, instala las dependencias e inicia la app web:
```bash
cd frontend
npm install
npm run dev
```
La interfaz de React estará disponible en `http://localhost:5173` (o el puerto indicado por Vite en tu consola).
