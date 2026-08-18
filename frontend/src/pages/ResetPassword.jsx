import { useState, useEffect } from "react";
import { Link, useSearchParams, useNavigate } from "react-router-dom";
import { resetPassword } from "../services/authService";
import "./ResetPassword.css";
import BrandLogo from "../components/BrandLogo";

function ResetPassword() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const token = searchParams.get("token");

  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!token) {
      setError("Token inválido o no proporcionado.");
    }
  }, [token]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!token) return;

    if (newPassword !== confirmPassword) {
      setError("Las contraseñas no coinciden.");
      return;
    }

    if (newPassword.length < 8) {
      setError("La contraseña debe tener al menos 8 caracteres.");
      return;
    }

    setError("");
    setLoading(true);
    
    try {
      await resetPassword(token, newPassword);
      setSuccess(true);
      setTimeout(() => {
        navigate("/login");
      }, 3000);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="reset-password-page">
      <BrandLogo variant="light" />
      <div className="reset-password-card">
        <h2>Nueva Contraseña</h2>
        <p className="subtitle">Crea una nueva contraseña segura para tu cuenta.</p>

        {success ? (
          <div className="success-message">
            Contraseña actualizada con éxito. Serás redirigido al inicio de sesión en unos segundos...
            <div className="footer-links" style={{ marginTop: '1rem' }}>
              <Link to="/login">Ir a Iniciar Sesión ahora</Link>
            </div>
          </div>
        ) : (
          <form onSubmit={handleSubmit}>
            <div className="form-field">
              <label>Nueva Contraseña</label>
              <input
                type="password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                required
                disabled={!token}
              />
            </div>
            
            <div className="form-field">
              <label>Confirmar Contraseña</label>
              <input
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
                disabled={!token}
              />
            </div>

            {error && <div className="error-message">{error}</div>}

            <button type="submit" className="submit-button" disabled={loading || !token}>
              {loading ? "Actualizando..." : "Restablecer contraseña"}
            </button>
          </form>
        )}
      </div>
    </div>
  );
}

export default ResetPassword;
