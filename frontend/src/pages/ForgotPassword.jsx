import { useState } from "react";
import { Link } from "react-router-dom";
import { forgotPassword } from "../services/authService";
import "./ForgotPassword.css";
import BrandLogo from "../components/BrandLogo";

function ForgotPassword() {
  const [email, setEmail] = useState("");
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [devToken, setDevToken] = useState(""); // Solo para entorno de desarrollo

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setMessage("");
    setDevToken("");
    setLoading(true);
    
    try {
      const data = await forgotPassword(email);
      setMessage("Se ha enviado un correo con las instrucciones para restablecer tu contraseña.");
      
      // Mostrar token de forma temporal para pruebas de desarrollo
      if (data && data.resetToken) {
        setDevToken(data.resetToken);
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="forgot-password-page">
      <BrandLogo variant="light" />
      <div className="forgot-password-card">
        <h2>Recuperar Contraseña</h2>
        <p className="subtitle">
          Ingresa tu dirección de correo electrónico y te enviaremos un enlace para restablecer tu contraseña.
        </p>

        {message ? (
          <div className="success-message">
            {message}
            {devToken && (
              <div className="dev-token-box">
                <p><strong>Dev Only:</strong> Copia este enlace para restablecer:</p>
                <Link to={`/reset-password?token=${devToken}`}>Ir a restablecer contraseña</Link>
              </div>
            )}
          </div>
        ) : (
          <form onSubmit={handleSubmit}>
            <div className="form-field">
              <label>Email</label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                placeholder="tu@email.com"
              />
            </div>

            {error && <div className="error-message">{error}</div>}

            <button type="submit" className="submit-button" disabled={loading}>
              {loading ? "Enviando..." : "Enviar enlace"}
            </button>
          </form>
        )}
        
        <div className="footer-links">
          <Link to="/login">Volver al inicio de sesión</Link>
        </div>
      </div>
    </div>
  );
}

export default ForgotPassword;
