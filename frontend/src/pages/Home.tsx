import { useNavigate } from "react-router-dom";
import { useEffect, useRef, useState } from "react";

export default function Home() {
  const navigate = useNavigate();
  const [isVisible, setIsVisible] = useState(false);
  const featuresRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        setIsVisible(entry.isIntersecting);
      },
      {
        threshold: 0.1,
        rootMargin: "-100px",
      }
    );

    if (featuresRef.current) {
      observer.observe(featuresRef.current);
    }

    return () => {
      if (featuresRef.current) {
        observer.unobserve(featuresRef.current);
      }
    };
  }, []);

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-900 overflow-x-hidden">
      {/* Navigation Bar - Matching page background */}
      <nav className="bg-slate-900/50 backdrop-blur-sm border-b border-gray-800">
        <div className="max-w-7xl mx-auto px-6 py-4 flex justify-between items-center">
          <div className="flex items-center gap-3">
            <div className="static-glow text-3xl">📚</div>
            <h1 className="static-glow text-2xl font-bold bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
              Lexra
            </h1>
          </div>
          <div className="flex gap-4">
            <button
              onClick={() => navigate("/login")}
              className="px-6 py-2.5 text-purple-300 font-medium hover:text-blue-200 transition-colors duration-200"
            >
              Log in →
            </button>
            <button
              onClick={() => navigate("/register")}
              className="px-6 py-2.5 bg-gradient-to-r from-purple-600 to-blue-900 text-white rounded-lg font-medium hover:from-blue-500 hover:to-blue-400 transition-all duration-200 shadow-lg shadow-blue-500/20"
            >
              Register
            </button>
          </div>
        </div>
      </nav>

      {/* Hero Section - Always Visible */}
      <main className="pt-24">
        <div className="max-w-7xl mx-auto px-6 py-20 text-center">
          <div className="animate-fade-in">
            <h2 className="text-6xl md:text-7xl font-bold mb-6 leading-tight">
              <span className="static-glow text-white">Learn smarter.</span>
              <br />
              <span className="neon-glow bg-gradient-to-r md:text-9xl from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
                Lexra
              </span>
            </h2>
            <p className="text-xl md:text-1xl text-slate-300 mb-10 max-w-2xl mx-auto leading-relaxed">
              Upload your textbooks and ask questions. Get instant, accurate
              answers with AI-powered insights and page references.
            </p>

            {/* Get Started Button - Only button visible */}
            <div className="flex justify-center items-center mb-16">
              <button
                onClick={() => navigate("/register")}
                className="px-6 py-3 bg-gradient-to-r from-purple-600 to-blue-900 text-white text-lg font-semibold rounded-lg hover:from-blue-500 hover:to-blue-400 transition-all duration-200 shadow-2xl shadow-blue-500/30 hover:shadow-blue-500/50 hover:scale-105"
              >
                Get started
              </button>
            </div>
          </div>
        </div>

        {/* Features Section - Scroll-Triggered Animation - Stacked Vertically */}
        <div className="min-h-screen flex items-center justify-center py-20">
          <div
            ref={featuresRef}
            className={`w-full px-6 transition-all duration-1000 ease-out ${
              isVisible
                ? "opacity-100 translate-y-0"
                : "opacity-0 translate-y-20"
            }`}
          >
            <div className="text-center mb-16">
              <h3 className="static-glow text-5xl font-bold text-white mb-6">
                Everything you need to succeed
              </h3>
              <p className="text-xl text-slate-400 max-w-2xl mx-auto">
                Our AI-powered platform makes learning interactive and efficient
              </p>
            </div>

            {/* Feature boxes stacked vertically, full width */}
            <div className="flex flex-col gap-8 max-w-6xl mx-auto">
              {/* Feature 1 - Upload PDFs */}
              <div
                className={`bg-slate-900/60 backdrop-blur-md border border-slate-700/50 rounded-2xl shadow-2xl p-12 hover:border-blue-500/50 hover:shadow-blue-500/20 transition-all duration-500 ${
                  isVisible
                    ? "opacity-100 translate-y-0"
                    : "opacity-0 translate-y-20"
                }`}
                style={{ transitionDelay: "100ms" }}
              >
                <div className="grid md:grid-cols-2 gap-8 items-center">
                  <div className="text-left">
                    <div className="text-6xl mb-6">📤</div>
                    <h4 className="text-3xl font-bold mb-4 text-white">
                      Upload PDFs
                    </h4>
                    <p className="text-slate-400 text-lg leading-relaxed">
                      Upload your textbook PDFs and let our AI process them into
                      searchable, interactive content instantly.
                    </p>
                  </div>
                  {/* GIF Placeholder for Upload PDFs */}
                  <div className="rounded-xl overflow-hidden shadow-xl border border-slate-600/50 bg-slate-800/50">
                    <img
                      src="/assets/images/upload.gif"
                      alt="Upload Demo"
                      className="w-full h-auto object-cover"
                      style={{ maxHeight: "300px" }}
                    />
                    {/* Fallback */}
                    <div className="min-h-[250px] flex items-center justify-center text-slate-600">
                      <div className="text-center">
                        <div className="text-5xl mb-3">📤</div>
                        <p className="text-sm">upload.gif placeholder</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              {/* Feature 2 - Ask Questions */}
              <div
                className={`bg-slate-900/60 backdrop-blur-md border border-slate-700/50 rounded-2xl shadow-2xl p-12 hover:border-purple-500/50 hover:shadow-purple-500/20 transition-all duration-500 ${
                  isVisible
                    ? "opacity-100 translate-y-0"
                    : "opacity-0 translate-y-20"
                }`}
                style={{ transitionDelay: "200ms" }}
              >
                <div className="grid md:grid-cols-2 gap-8 items-center">
                  {/* GIF Placeholder for Ask Questions */}
                  <div className="rounded-xl overflow-hidden shadow-xl border border-slate-600/50 bg-slate-800/50 order-2 md:order-1">
                    <img
                      src="/assets/images/questions.gif"
                      alt="Questions Demo"
                      className="w-full h-auto object-cover"
                      style={{ maxHeight: "300px" }}
                    />
                    {/* Fallback */}
                    <div className="min-h-[250px] flex items-center justify-center text-slate-600">
                      <div className="text-center">
                        <div className="text-5xl mb-3">💡</div>
                        <p className="text-sm">questions.gif placeholder</p>
                      </div>
                    </div>
                  </div>
                  <div className="text-left order-1 md:order-2">
                    <div className="text-6xl mb-6">💡</div>
                    <h4 className="text-3xl font-bold mb-4 text-white">
                      Ask Questions
                    </h4>
                    <p className="text-slate-400 text-lg leading-relaxed">
                      Ask any question about your textbook and get detailed,
                      contextual answers powered by advanced AI.
                    </p>
                  </div>
                </div>
              </div>

              {/* Feature 3 - Page References */}
              <div
                className={`bg-slate-900/60 backdrop-blur-md border border-slate-700/50 rounded-2xl shadow-2xl p-12 hover:border-pink-500/50 hover:shadow-pink-500/20 transition-all duration-500 ${
                  isVisible
                    ? "opacity-100 translate-y-0"
                    : "opacity-0 translate-y-20"
                }`}
                style={{ transitionDelay: "300ms" }}
              >
                <div className="grid md:grid-cols-2 gap-8 items-center">
                  <div className="text-left">
                    <div className="text-6xl mb-6">🔍</div>
                    <h4 className="text-3xl font-bold mb-4 text-white">
                      Page References
                    </h4>
                    <p className="text-slate-400 text-lg leading-relaxed">
                      Every answer includes exact page numbers and relevant
                      excerpts from your textbook for easy verification.
                    </p>
                  </div>
                  {/* GIF Placeholder for Page References */}
                  <div className="rounded-xl overflow-hidden shadow-xl border border-slate-600/50 bg-slate-800/50">
                    <img
                      src="/assets/images/references.gif"
                      alt="References Demo"
                      className="w-full h-auto object-cover"
                      style={{ maxHeight: "300px" }}
                    />
                    {/* Fallback */}
                    <div className="min-h-[250px] flex items-center justify-center text-slate-600">
                      <div className="text-center">
                        <div className="text-5xl mb-3">🔍</div>
                        <p className="text-sm">references.gif placeholder</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Additional CTA in Features Section */}
            <div className="text-center mt-16">
              <p className="text-slate-400 text-lg mb-6">
                Work smarter. Learn smarter.
              </p>
              <button
                onClick={() => navigate("/register")}
                className="px-6 py-3 bg-gradient-to-r from-purple-600 to-blue-900 text-white text-lg font-semibold rounded-lg hover:from-purple-500 hover:to-pink-500 transition-all duration-200 shadow-2xl shadow-purple-500/30 hover:shadow-purple-500/50 hover:scale-105"
              >
                Start Now
              </button>
            </div>
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="bg-slate-900/80 backdrop-blur-md border-t border-slate-800/50 py-8 mt-20">
        <div className="max-w-7xl mx-auto px-6 text-center">
          <p className="text-slate-400">
            © 2025 Lexra. Learn smarter with AI-powered textbooks.
          </p>
        </div>
      </footer>

      <style>{`
        @keyframes fade-in {
          from {
            opacity: 0;
            transform: translateY(20px);
          }
          to {
            opacity: .7;
            transform: translateY(0);
          }
        }

        .animate-fade-in {
          animation: fade-in 0.8s ease-out;
        }

        .static-glow {
          text-shadow: 0 0 10px rgba(255, 255, 255, 0.3),
                       0 0 20px rgba(147, 51, 234, 0.2),
                       0 0 30px rgba(59, 130, 246, 0.2);
        }

        .neon-glow {
          filter: drop-shadow(0 0 10px rgba(147, 51, 234, 0.7))
                  drop-shadow(0 0 25px rgba(59, 130, 246, 0.5))
                  drop-shadow(0 0 40px rgba(236, 72, 153, 0.4));
          animation: neon-pulse 4s ease-in-out infinite;
        }

        @keyframes neon-pulse {
          0%, 100% {
            filter: drop-shadow(0 0 10px rgba(147, 51, 234, 0.7))
                    drop-shadow(0 0 25px rgba(59, 130, 246, 0.5))
                    drop-shadow(0 0 40px rgba(236, 72, 153, 0.4));
          }
          50% {
            filter: drop-shadow(0 0 10px rgba(147, 51, 234, 0.9))
                    drop-shadow(0 0 25px rgba(59, 130, 246, 0.7))
                    drop-shadow(0 0 40px rgba(236, 72, 153, 0.6));
          }
        }
      `}</style>
    </div>
  );
}
