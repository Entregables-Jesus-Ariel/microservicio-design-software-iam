import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { getUsers, getRoles, assignRole, revokeRole } from "../services/rbacService";
import BrandLogo from "../components/BrandLogo";
import "./AdminDashboard.css";

function AdminDashboard() {
  const [users, setUsers] = useState([]);
  const [roles, setRoles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [selectedRole, setSelectedRole] = useState({});

  const fetchData = async () => {
    setLoading(true);
    try {
      const [usersData, rolesData] = await Promise.all([getUsers(), getRoles()]);
      setUsers(usersData || []);
      setRoles(rolesData || []);
    } catch (err) {
      setError(err.message || "Error al cargar datos");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleAssignRole = async (userId) => {
    const roleId = selectedRole[userId];
    if (!roleId) return;
    try {
      await assignRole(userId, roleId);
      // Actualizar localmente o refetch
      fetchData();
    } catch (err) {
      alert(err.message || "Error al asignar rol");
    }
  };

  const handleRevokeRole = async (userId, roleId) => {
    if (!window.confirm("¿Seguro que deseas revocar este rol?")) return;
    try {
      await revokeRole(userId, roleId);
      fetchData();
    } catch (err) {
      alert(err.message || "Error al revocar rol");
    }
  };

  const handleRoleSelectChange = (userId, value) => {
    setSelectedRole({ ...selectedRole, [userId]: value });
  };

  if (loading) return <div className="admin-loading">Cargando panel de administración...</div>;

  return (
    <div className="admin-dashboard">
      <header className="admin-header">
        <BrandLogo variant="light" />
        <nav style={{ display: 'flex', gap: '1rem' }}>
          <Link to="/admin/audit" className="nav-link">Ver Auditoría de Accesos</Link>
          <Link to="/" className="nav-link">Volver al Dashboard</Link>
        </nav>
      </header>

      <main className="admin-main">
        <h1>Administración de Usuarios y Roles</h1>
        {error && <div className="admin-error">{error}</div>}

        <div className="table-container">
          <table className="users-table">
            <thead>
              <tr>
                <th>Usuario</th>
                <th>Email</th>
                <th>Roles Actuales</th>
                <th>Acciones</th>
              </tr>
            </thead>
            <tbody>
              {users.map((item) => (
                <tr key={item.User.ID}>
                  <td>{item.User.FirstName} {item.User.LastName}</td>
                  <td>{item.User.Email}</td>
                  <td>
                    <div className="roles-list">
                      {(item.Roles || []).map((r, idx) => (
                        <span key={idx} className="role-chip">
                          {r.Name}
                          <button 
                            className="role-remove" 
                            onClick={() => {
                              // Buscar el ID real del rol para revocar
                              const roleObj = roles.find(ro => ro.Name === r.Name);
                              if (roleObj) handleRevokeRole(item.User.ID, roleObj.ID);
                            }}
                          >
                            ×
                          </button>
                        </span>
                      ))}
                    </div>
                  </td>
                  <td>
                    <div className="assign-role-group">
                      <select 
                        value={selectedRole[item.User.ID] || ""} 
                        onChange={(e) => handleRoleSelectChange(item.User.ID, e.target.value)}
                      >
                        <option value="">Seleccionar rol...</option>
                        {roles.map((r) => (
                          <option key={r.ID} value={r.ID}>{r.Name}</option>
                        ))}
                      </select>
                      <button 
                        className="assign-button"
                        onClick={() => handleAssignRole(item.User.ID)}
                        disabled={!selectedRole[item.User.ID]}
                      >
                        Asignar
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </main>
    </div>
  );
}

export default AdminDashboard;
