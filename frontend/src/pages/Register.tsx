import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { authAPI } from "../services/api";

// Helper function to format error messages for display
const formatErrorMessage = (error: string): string => {
  // Capitalize first letter
  const formatted = error.charAt(0).toUpperCase() + error.slice(1);

  // Add period if not present
  return formatted.endsWith(".") ? formatted : formatted + ".";
};

export default function Register() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    // Validate passwords match
    if (password !== confirmPassword) {
      setError("Passwords do not match");
      return;
    }

    // Validate password length
    if (password.length < 8) {
      setError("Password must be at least 8 characters");
      return;
    }

    setIsLoading(true);

    try {
      const response = await authAPI.register({ email, password });

      // Save token to localStorage
      localStorage.setItem("token", response.token);
      localStorage.setItem("user", JSON.stringify(response.user));

      // Redirect to library
      navigate("/library");
    } catch (err: any) {
      const errorMessage =
        err.response?.data || "Registration failed. Please try again.";
      setError(formatErrorMessage(errorMessage));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 flex items-center justify-center p-4">
      <div className="relative w-full max-w-md">
        <div className="bg-gray-900/50 backdrop-blur-sm border border-gray-700 rounded-lg shadow-xl p-8 w-full">
          <div className="flex items-center justify-center gap-2 mb-2">
            <div className="text-2xl">📚</div>
            <h1 className="text-3xl font-bold text-white text-center">
              Join Lexra
            </h1>
          </div>
          <p className="text-gray-400 text-center mb-8">
            Sign up to get started
          </p>

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label
                htmlFor="email"
                className="block text-sm font-medium text-gray-300 mb-2"
              >
                Email
              </label>
              <input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="your@email.com"
                required
                className="w-full px-4 py-2 bg-gray-900/50 border border-gray-600 text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition placeholder-gray-500"
              />
            </div>

            <div>
              <label
                htmlFor="password"
                className="block text-sm font-medium text-gray-300 mb-2"
              >
                Password
              </label>
              <input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="At least 8 characters"
                required
                className="w-full px-4 py-2 bg-gray-900/50 border border-gray-600 text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition placeholder-gray-500"
              />
            </div>

            <div>
              <label
                htmlFor="confirmPassword"
                className="block text-sm font-medium text-gray-300 mb-2"
              >
                Confirm Password
              </label>
              <input
                id="confirmPassword"
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                placeholder="Re-enter your password"
                required
                className="w-full px-4 py-2 bg-gray-900/50 border border-gray-600 text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition placeholder-gray-500"
              />
            </div>

            {error && (
              <div className="bg-red-900/50 border border-red-700 text-red-300 px-4 py-3 rounded-lg">
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={isLoading}
              className="w-full bg-blue-600 hover:bg-blue-500 disabled:bg-gray-600 text-white font-semibold py-2 px-4 rounded-lg transition duration-200 shadow-lg shadow-blue-500/20"
            >
              {isLoading ? "Creating Account..." : "Create Account"}
            </button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-gray-400">
              Already have an account?{" "}
              <button
                onClick={() => navigate("/login")}
                className="text-blue-400 font-medium hover:text-blue-300"
              >
                Sign in
              </button>
            </p>
          </div>
        </div>
        {/* Back to Home Link - Bottom Left */}
        <button
          onClick={() => navigate("/")}
          className="absolute -bottom-12 left-0 text-slate-400 hover:text-blue-400 transition-colors duration-200 flex items-center gap-1 text-sm"
        >
          <span>←</span>Back to Home
        </button>
      </div>
    </div>
  );
}
