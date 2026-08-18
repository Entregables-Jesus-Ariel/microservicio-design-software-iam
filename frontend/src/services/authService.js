import axios from "axios";

const API_BASE_URL = "http://localhost:8082/api";

const ERROR_MESSAGES = {
  "invalid email or password": "Email o contraseña incorrectos.",
  "account is temporarily locked": "Tu cuenta está bloqueada temporalmente por varios intentos fallidos. Intenta de nuevo en unos minutos.",
  "email is already registered": "Ese email ya está registrado.",
};

function translateError(rawMessage) {
  return ERROR_MESSAGES[rawMessage] || "No se pudo conectar con el servidor. Intenta de nuevo.";
}

export const api = axios.create({
  baseURL: API_BASE_URL,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("accessToken");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    
    // Si obtenemos un 401 y no es ya un reintento (evitar loops infinitos)
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      const refreshToken = localStorage.getItem("refreshToken");
      
      if (refreshToken) {
        try {
          const res = await axios.post(`${API_BASE_URL}/auth/refresh`, { refreshToken });
          localStorage.setItem("accessToken", res.data.accessToken);
          localStorage.setItem("refreshToken", res.data.refreshToken);
          
          originalRequest.headers.Authorization = `Bearer ${res.data.accessToken}`;
          return api(originalRequest);
        } catch (err) {
          // Si el refresh también falla, forzamos logout
          localStorage.removeItem("accessToken");
          localStorage.removeItem("refreshToken");
          window.location.href = "/login";
        }
      }
    }
    return Promise.reject(error);
  }
);

export async function registerUser({ email, password, firstName, lastName }) {
  try {
    const response = await api.post("/auth/register", {
      email,
      password,
      firstName,
      lastName,
    });
    return response.data;
  } catch (error) {
    const rawMessage = error.response?.data?.error;
    throw new Error(rawMessage ? translateError(rawMessage) : translateError());
  }
}

export async function loginUser({ email, password }) {
  try {
    const response = await api.post("/auth/login", {
      email,
      password,
    });
    return response.data;
  } catch (error) {
    const rawMessage = error.response?.data?.error;
    const isLocked = rawMessage === "account is temporarily locked";
    const err = new Error(rawMessage ? translateError(rawMessage) : translateError());
    err.isLocked = isLocked;
    throw err;
  }
}

export async function logoutUser() {
  try {
    const refreshToken = localStorage.getItem("refreshToken");
    if (refreshToken) {
      await api.post("/auth/logout", { refreshToken });
    }
  } catch (err) {
    console.error("Error al cerrar sesión", err);
  } finally {
    localStorage.removeItem("accessToken");
    localStorage.removeItem("refreshToken");
    window.location.href = "/login";
  }
}

export async function forgotPassword(email) {
  try {
    const response = await api.post("/auth/forgot-password", { email });
    return response.data;
  } catch (error) {
    const rawMessage = error.response?.data?.error;
    throw new Error(rawMessage ? translateError(rawMessage) : translateError());
  }
}

export async function resetPassword(token, newPassword) {
  try {
    const response = await api.post("/auth/reset-password", { token, newPassword });
    return response.data;
  } catch (error) {
    const rawMessage = error.response?.data?.error;
    throw new Error(rawMessage ? translateError(rawMessage) : translateError());
  }
}
