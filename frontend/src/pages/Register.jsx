import { useState } from "react";
import { Link } from "react-router-dom";
import { registerUser } from "../services/authService";
import "./Register.css";
import BrandLogo from "../components/BrandLogo";
import { Check, Circle } from "lucide-react";

function Register() {
  const [form, setForm] = useState({
    email: "",
    password: "",
    firstName: "",
    lastName: "",
  });
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);

  const passwordChecks = {
    length: form.password.length >= 8,
    upper: /[A-Z]/.test(form.password),
    number: /[0-9]/.test(form.password),
  };
  const isPasswordValid = Object.values(passwordChecks).every(Boolean);

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess(false);

    if (!isPasswordValid) {
      setError("La contraseña no cumple los requisitos mínimos.");
      return;
    }

    setLoading(true);
    try {
      await registerUser(form);
      setSuccess(true);
      setForm({ email: "", password: "", firstName: "", lastName: "" });
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="register-page">
      <BrandLogo variant="dark" />
      <div className="register-card">
        <h2>Crear cuenta</h2>
        <form onSubmit={handleSubmit}>
          <div className="register-field">
            <label>Nombre</label>
            <input
              type="text"
              name="firstName"
              autoComplete="off"
              value={form.firstName}
              onChange={handleChange}
              required
            />
          </div>
          <div className="register-field">
            <label>Apellido</label>
            <input
              type="text"
              name="lastName"
              autoComplete="off"
              value={form.lastName}
              onChange={handleChange}
              required
            />
          </div>
          <div className="register-field">
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
          <div className="register-field">
            <label>Contraseña</label>
            <input
              type="password"
              name="password"
              autoComplete="new-password"
              value={form.password}
              onChange={handleChange}
              required
            />
            {form.password.length > 0 && (
              <ul className="password-rules">
                <li className={passwordChecks.length ? "valid" : ""}>
                  <span className="icon">{passwordChecks.length ? <Check size={14} /> : <Circle size={6} fill="currentColor" />}</span>
                  Mínimo 8 caracteres
                </li>
                <li className={passwordChecks.upper ? "valid" : ""}>
                  <span className="icon">{passwordChecks.upper ? <Check size={14} /> : <Circle size={6} fill="currentColor" />}</span>
                  Al menos una mayúscula
                </li>
                <li className={passwordChecks.number ? "valid" : ""}>
                  <span className="icon">{passwordChecks.number ? <Check size={14} /> : <Circle size={6} fill="currentColor" />}</span>
                  Al menos un número
                </li>
              </ul>
            )}
          </div>

          {error && <div className="register-error">{error}</div>}
          {success && (
            <div className="register-success">¡Usuario registrado correctamente!</div>
          )}

          <button type="submit" className="register-button" disabled={loading}>
            {loading ? "Registrando..." : "Registrarse"}
          </button>
        </form>
        <div className="register-footer">
          ¿Ya tienes cuenta? <Link to="/login">Inicia sesión</Link>
        </div>
      </div>
    </div>
  );
}

export default Register;
