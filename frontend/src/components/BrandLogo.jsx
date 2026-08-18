import "./BrandLogo.css";

function BrandLogo({ variant = "light" }) {
  return (
    <div className={`brand-logo brand-logo--${variant}`}>
      <span className="brand-icon">IAM</span>
      <span className="brand-name">IAM System</span>
    </div>
  );
}

export default BrandLogo;
