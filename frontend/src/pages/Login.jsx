import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { loginUser } from "../services/authService";
import "./Login.css";
import BrandLogo from "../components/BrandLogo";
import { Lock } from "lucide-react";

function Login() {
  const navigate = useNavigate();
  const [form, setForm] = useState({ email: "", password: "" });
  const [error, setError] = useState("");
  const [isLocked, setIsLocked] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setIsLocked(false);
    setLoading(true);
    try {
      const data = await loginUser(form);
      localStorage.setItem("accessToken", data.accessToken);
      localStorage.setItem("refreshToken", data.refreshToken);
      navigate("/");
    } catch (err) {
      setError(err.message);
      setIsLocked(Boolean(err.isLocked));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-page">
      <BrandLogo variant="light" />
      <div className="login-card">
        <h2>Iniciar sesión</h2>
        <form onSubmit={handleSubmit}>
          <div className="login-field">
            <label>Email</label>
            <input
              type="email"
              name="email"
              autoComplete="off"
              value={form.email}
              onChange={handleChange}
              required
            />
          </div>
          <div className="login-field">
            <label>Contraseña</label>
            <input
              type="password"
              name="password"
              autoComplete="current-password"
              value={form.password}
              onChange={handleChange}
              required
            />
          </div>

          {error && (
            <div className={isLocked ? "login-locked" : "login-error"}>
              {isLocked && <span className="login-locked-icon">🔒</span>}
              {error}
            </div>
          )}

          <button type="submit" className="login-button" disabled={loading}>
            {loading ? "Ingresando..." : "Ingresar"}
          </button>
        </form>
        <div className="login-footer">
          <div><Link to="/forgot-password">¿Olvidaste tu contraseña?</Link></div>
          <div style={{marginTop: "0.5rem"}}>¿No tienes cuenta? <Link to="/register">Regístrate</Link></div>
        </div>
      </div>
    </div>
  );
}

export default Login;
