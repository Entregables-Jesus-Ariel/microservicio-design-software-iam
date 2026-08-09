import axios from "axios";

const API_BASE_URL = "http://localhost:8082/api";

export async function registerUser({ email, password, firstName, lastName }) {
  try {
    const response = await axios.post(`${API_BASE_URL}/auth/register`, {
      email,
      password,
      firstName,
      lastName,
    });
    return response.data;
  } catch (error) {
    if (error.response && error.response.data && error.response.data.error) {
      throw new Error(error.response.data.error);
    }
    throw new Error("No se pudo conectar con el servidor. Intenta de nuevo.");
  }
}

export async function loginUser({ email, password }) {
  try {
    const response = await axios.post(`${API_BASE_URL}/auth/login`, {
      email,
      password,
    });
    return response.data;
  } catch (error) {
    if (error.response && error.response.data && error.response.data.error) {
      throw new Error(error.response.data.error);
    }
    throw new Error("No se pudo conectar con el servidor. Intenta de nuevo.");
  }
}
