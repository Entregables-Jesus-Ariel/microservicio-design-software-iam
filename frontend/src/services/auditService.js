import { api } from "./authService";

export async function getLoginAudits(page = 1, limit = 20) {
  const response = await api.get(`/admin/audit/logins?page=${page}&limit=${limit}`);
  return response.data;
}
