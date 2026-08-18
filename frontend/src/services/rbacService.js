import { api } from "./authService";

export async function getUsers() {
  const response = await api.get("/admin/users");
  return response.data;
}

export async function getRoles() {
  const response = await api.get("/admin/roles");
  return response.data;
}

export async function assignRole(userId, roleId) {
  const response = await api.post(`/admin/users/${userId}/roles`, { roleId });
  return response.data;
}

export async function revokeRole(userId, roleId) {
  const response = await api.delete(`/admin/users/${userId}/roles/${roleId}`);
  return response.data;
}
