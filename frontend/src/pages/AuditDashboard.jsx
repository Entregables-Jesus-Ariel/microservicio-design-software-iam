import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { getLoginAudits } from "../services/auditService";
import BrandLogo from "../components/BrandLogo";
import "./AuditDashboard.css";

function AuditDashboard() {
  const [logs, setLogs] = useState([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const limit = 20;

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const data = await getLoginAudits(page, limit);
        setLogs(data.Items || []);
        setTotal(data.Total || 0);
      } catch (err) {
        setError(err.message || "Error al cargar la auditoría");
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [page]);

  const totalPages = Math.ceil(total / limit);

  const getOutcomeStyle = (outcome) => {
    switch (outcome) {
      case "SUCCESS": return "outcome-success";
      case "INVALID_PASSWORD": return "outcome-warning";
      case "ACCOUNT_LOCKED": return "outcome-danger";
      case "USER_NOT_FOUND": return "outcome-gray";
      default: return "outcome-gray";
    }
  };

  return (
    <div className="audit-dashboard">
      <header className="audit-header">
        <BrandLogo variant="light" />
        <nav className="audit-nav">
          <Link to="/admin" className="nav-link">← Usuarios y Roles</Link>
          <Link to="/" className="nav-link">Dashboard Principal</Link>
        </nav>
      </header>

      <main className="audit-main">
        <h1>Auditoría de Accesos</h1>
        <p className="subtitle">Historial de intentos de inicio de sesión en el sistema.</p>
        
        {error && <div className="audit-error">{error}</div>}

        <div className="table-container">
          <table className="audit-table">
            <thead>
              <tr>
                <th>Fecha y Hora</th>
                <th>Email Intentado</th>
                <th>Resultado</th>
                <th>ID Usuario</th>
              </tr>
            </thead>
            <tbody>
              {loading && logs.length === 0 ? (
                <tr>
                  <td colSpan="4" className="text-center">Cargando registros...</td>
                </tr>
              ) : logs.length === 0 ? (
                <tr>
                  <td colSpan="4" className="text-center">No hay registros de auditoría.</td>
                </tr>
              ) : (
                logs.map((log) => (
                  <tr key={log.ID}>
                    <td>{new Date(log.AttemptedAt).toLocaleString()}</td>
                    <td>{log.EmailAttempted}</td>
                    <td>
                      <span className={`outcome-chip ${getOutcomeStyle(log.Outcome)}`}>
                        {log.Outcome}
                      </span>
                    </td>
                    <td className="user-id-cell" title={log.UserID || "N/A"}>
                      {log.UserID ? log.UserID.split("-")[0] + "..." : "-"}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {totalPages > 1 && (
          <div className="pagination">
            <button 
              disabled={page === 1 || loading} 
              onClick={() => setPage(page - 1)}
            >
              Anterior
            </button>
            <span>Página {page} de {totalPages}</span>
            <button 
              disabled={page === totalPages || loading} 
              onClick={() => setPage(page + 1)}
            >
              Siguiente
            </button>
          </div>
        )}
      </main>
    </div>
  );
}

export default AuditDashboard;
