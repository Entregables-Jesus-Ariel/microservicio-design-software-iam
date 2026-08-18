import { useNavigate, Link } from "react-router-dom";
import { logoutUser } from "../services/authService";
import BrandLogo from "../components/BrandLogo";
import "./Dashboard.css";

function getRolesFromToken() {
  const token = localStorage.getItem("accessToken");
  if (!token) return [];
  try {
    const payload = JSON.parse(atob(token.split(".")[1]));
    return payload.roles || [];
  } catch (e) {
    return [];
  }
}

function Dashboard() {
  const navigate = useNavigate();
  const roles = getRolesFromToken();
  const isAdmin = roles.includes("ADMIN");

  const handleLogout = async () => {
    try {
      await logoutUser();
    } catch (error) {
      console.error("Error logging out", error);
    } finally {
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
      navigate("/login");
    }
  };

  return (
    <div className="dashboard-page">
      <header className="dashboard-header">
        <BrandLogo variant="light" />
        <nav className="dashboard-nav">
          {isAdmin && (
            <Link to="/admin" className="admin-link">
              Admin Dashboard
            </Link>
          )}
          <button onClick={handleLogout} className="logout-button">
            Cerrar Sesión
          </button>
        </nav>
      </header>
      <main className="dashboard-main">
        <h1>Bienvenido al Sistema IAM</h1>
        <p>Has iniciado sesión exitosamente.</p>
        <div className="roles-info">
          <h3>Tus roles:</h3>
          {roles.length > 0 ? (
            <ul>
              {roles.map((r, idx) => <li key={idx}>{r}</li>)}
            </ul>
          ) : (
            <p>No tienes roles asignados.</p>
          )}
        </div>
      </main>
    </div>
  );
}

export default Dashboard;
