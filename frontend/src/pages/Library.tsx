import { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { textbookAPI, queryAPI } from "../services/api";
import type { Textbook, QueryResponse } from "../types";

export default function Library() {
  const [textbooks, setTextbooks] = useState<Textbook[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [uploadFile, setUploadFile] = useState<File | null>(null);
  const [uploadTitle, setUploadTitle] = useState("");
  const [isUploading, setIsUploading] = useState(false);
  const [selectedFolder, setSelectedFolder] = useState<Textbook | null>(null);
  const [question, setQuestion] = useState("");
  const [isQuerying, setIsQuerying] = useState(false);
  const [conversationHistory, setConversationHistory] = useState<
    Array<{ question: string; answer: QueryResponse }>
  >([]);
  const [showUploadModal, setShowUploadModal] = useState(false);
  const navigate = useNavigate();
  const chatContainerRef = useRef<HTMLDivElement>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  // Load textbooks when page loads
  useEffect(() => {
    loadTextbooks();
  }, []);

  // Auto-scroll to bottom when conversation updates
  useEffect(() => {
    if (chatContainerRef.current) {
      chatContainerRef.current.scrollTop =
        chatContainerRef.current.scrollHeight;
    }
  }, [conversationHistory]);

  const loadTextbooks = async () => {
    try {
      const data = await textbookAPI.list();
      setTextbooks(data);
    } catch (err: any) {
      setError("Failed to load textbooks");
      if (err.response?.status === 401) {
        // Token expired, redirect to login
        navigate("/login");
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleUpload = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!uploadFile) return;

    setIsUploading(true);
    setError("");

    try {
      await textbookAPI.upload(uploadFile, uploadTitle);
      // Reload textbooks
      await loadTextbooks();
      // Reset form and close modal
      setUploadFile(null);
      setUploadTitle("");
      setShowUploadModal(false);
    } catch (err: any) {
      setError(err.response?.data || "Upload failed");
    } finally {
      setIsUploading(false);
    }
  };

  const handleDelete = async (id: number, title: string) => {
    if (!confirm(`Delete "${title}"? This cannot be undone.`)) return;

    try {
      await textbookAPI.delete(id);
      // Close folder if it's the one being deleted
      if (selectedFolder?.id === id) {
        setSelectedFolder(null);
        setConversationHistory([]);
      }
      // Reload textbooks
      await loadTextbooks();
    } catch (err: any) {
      setError("Failed to delete textbook");
    }
  };

  const handleOpenFolder = (textbook: Textbook) => {
    setSelectedFolder(textbook);
    setConversationHistory([]);
    setQuestion("");
  };

  const handleCloseFolder = () => {
    setSelectedFolder(null);
    setConversationHistory([]);
    setQuestion("");
  };

  const handleAskQuestion = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!question.trim() || !selectedFolder) return;

    const currentQuestion = question.trim();
    setQuestion(""); // Clear input immediately
    setIsQuerying(true);
    setError("");

    // Reset textarea height
    if (textareaRef.current) {
      textareaRef.current.style.height = "auto";
    }

    try {
      const response = await queryAPI.ask({
        question: currentQuestion,
        textbook_id: selectedFolder.id,
      });

      // Add to conversation history
      setConversationHistory((prev) => [
        ...prev,
        { question: currentQuestion, answer: response },
      ]);
    } catch (err: any) {
      setError(err.response?.data || "Failed to get answer");
      // Restore question on error
      setQuestion(currentQuestion);
    } finally {
      setIsQuerying(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    // Submit on Enter (without Shift)
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      if (question.trim() && selectedFolder?.processed && !isQuerying) {
        handleAskQuestion(e as any);
      }
    }
    // Allow Shift+Enter for new line (default behavior)
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    navigate("/login");
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 flex items-center justify-center">
        <div className="text-xl text-gray-400">Loading...</div>
      </div>
    );
  }

  // Folder detail view - ChatGPT style
  if (selectedFolder) {
    return (
      <div className="flex flex-col h-screen bg-gray-900">
        {/* Header */}
        <header className="flex-shrink-0 border-b border-gray-800">
          <div className="max-w-4xl mx-auto px-4 py-4 flex justify-between items-center">
            <div className="flex items-center gap-4">
              <button
                onClick={handleCloseFolder}
                className="text-gray-400 hover:text-white transition"
              >
                ← Back
              </button>
              <div className="flex items-center gap-2">
                <div className="text-xl">📁</div>
                <h1 className="text-lg font-semibold text-white">
                  {selectedFolder.title}
                </h1>
              </div>
            </div>
            <button
              onClick={handleLogout}
              className="text-gray-400 hover:text-white transition text-sm"
            >
              Logout
            </button>
          </div>
        </header>

        {/* Main Content Area */}
        <div className="flex-1 flex flex-col overflow-hidden">
          {/* Chat Area - Scrollable */}
          <div ref={chatContainerRef} className="flex-1 overflow-y-auto">
            <div className="max-w-3xl mx-auto px-4">
              {/* Status Warning */}
              {!selectedFolder.processed && (
                <div className="bg-yellow-900/30 border border-yellow-700/50 text-yellow-300 px-6 py-4 rounded-lg mb-8 text-center mt-8">
                  This textbook is still being processed. Please check back in a
                  few minutes.
                </div>
              )}

              {/* Conversation History */}
              {conversationHistory.length > 0 && (
                <div className="space-y-8 py-8">
                  {conversationHistory.map((conversation, convIndex) => (
                    <div key={convIndex} className="space-y-8">
                      {/* User Question */}
                      <div className="flex gap-4">
                        <div className="w-8 h-8 rounded-full bg-blue-600 flex-shrink-0 flex items-center justify-center text-white font-semibold">
                          U
                        </div>
                        <div className="flex-1 pt-1">
                          <p className="text-gray-100 leading-relaxed whitespace-pre-wrap">
                            {conversation.question}
                          </p>
                        </div>
                      </div>

                      {/* AI Response */}
                      <div className="flex gap-4">
                        <div className="w-8 h-8 rounded-full bg-gray-700 flex-shrink-0 flex items-center justify-center text-white">
                          📚
                        </div>
                        <div className="flex-1 pt-1">
                          <p className="text-gray-100 leading-relaxed whitespace-pre-wrap mb-6">
                            {conversation.answer.answer}
                          </p>

                          {/* Sources */}
                          {conversation.answer.sources &&
                            conversation.answer.sources.length > 0 && (
                              <div className="mt-6">
                                <p className="text-sm text-gray-400 font-semibold mb-3">
                                  Sources
                                </p>
                                <div className="space-y-2">
                                  {conversation.answer.sources.map(
                                    (source, index) => (
                                      <div
                                        key={index}
                                        className="bg-gray-800/50 border border-gray-700 rounded-lg p-4"
                                      >
                                        <div className="text-blue-400 text-sm font-medium mb-2">
                                          Page {source.page_number}
                                        </div>
                                        <p className="text-gray-400 text-sm">
                                          "{source.content}"
                                        </p>
                                      </div>
                                    )
                                  )}
                                </div>
                              </div>
                            )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}

              {/* Error Message */}
              {error && (
                <div className="bg-red-900/30 border border-red-700/50 text-red-300 px-6 py-4 rounded-lg text-center mt-8">
                  {error}
                </div>
              )}
            </div>
          </div>

          {/* Input Area - Transitions from center to bottom */}
          <div
            className={`flex-shrink-0 bg-gray-900 transition-all duration-500 ease-in-out ${
              conversationHistory.length > 0
                ? "border-t border-gray-800"
                : "flex-1 flex flex-col items-center justify-center"
            }`}
          >
            <div
              className={`w-full max-w-3xl px-4 mx-auto transition-all duration-500 ${
                conversationHistory.length > 0 ? "py-4" : ""
              }`}
            >
              {/* Welcome Message - Only show when no conversation */}
              {conversationHistory.length === 0 && (
                <div className="text-center mb-12">
                  <div className="text-7xl mb-6">📁</div>
                  <h2 className="text-3xl font-semibold text-white mb-4">
                    {selectedFolder.title}
                  </h2>
                  <p className="text-gray-400 text-lg mb-12">
                    Ask me anything about this textbook
                  </p>
                </div>
              )}

              {/* Input Form */}
              <form
                onSubmit={handleAskQuestion}
                className="flex items-end gap-3"
              >
                <textarea
                  ref={textareaRef}
                  value={question}
                  onChange={(e) => setQuestion(e.target.value)}
                  onKeyDown={handleKeyDown}
                  placeholder="Ask anything..."
                  disabled={!selectedFolder.processed}
                  rows={1}
                  className="flex-1 px-5 py-3 bg-gray-800 border border-gray-700 text-white rounded-3xl focus:outline-none focus:border-gray-600 transition placeholder-gray-500 disabled:opacity-50 resize-none overflow-hidden"
                  style={{
                    minHeight: "48px",
                    maxHeight: "200px",
                    height: "auto",
                  }}
                  onInput={(e) => {
                    const target = e.target as HTMLTextAreaElement;
                    target.style.height = "auto";
                    target.style.height =
                      Math.min(target.scrollHeight, 200) + "px";
                  }}
                />
                <button
                  type="submit"
                  disabled={
                    isQuerying || !question.trim() || !selectedFolder.processed
                  }
                  className="w-10 h-10 flex items-center justify-center bg-white text-gray-900 rounded-full hover:bg-gray-100 disabled:bg-gray-700 disabled:text-gray-500 disabled:opacity-50 transition flex-shrink-0 mb-[2px]"
                >
                  {isQuerying ? (
                    <div className="w-4 h-4 border-2 border-gray-900 border-t-transparent rounded-full animate-spin"></div>
                  ) : (
                    <span className="text-lg font-bold">↑</span>
                  )}
                </button>
              </form>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Library view with folders
  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900">
      {/* Header */}
      <header className="bg-gray-900/50 backdrop-blur-sm border-b border-gray-800">
        <div className="max-w-7xl mx-auto px-4 py-6 flex justify-between items-center">
          <div className="flex items-center gap-2">
            <h1 className="text-3xl font-bold text-white">Lexra Dashboard</h1>
          </div>
          <button
            onClick={handleLogout}
            className="px-4 py-2 bg-gray-800 text-gray-300 rounded-lg hover:bg-gray-700 transition border border-gray-700"
          >
            Logout
          </button>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 py-8">
        {/* Error Message */}
        {error && (
          <div className="bg-red-900/50 border border-red-700 text-red-300 px-4 py-3 rounded-lg mb-6">
            {error}
          </div>
        )}

        {/* Folders Grid */}
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
          {/* Add Folder Button */}
          <button
            onClick={() => setShowUploadModal(true)}
            className="group relative bg-gray-800/30 backdrop-blur-sm border-2 border-dashed border-gray-600 rounded-lg shadow-lg p-6 hover:border-blue-500/50 transition cursor-pointer"
          >
            <div className="flex flex-col items-center">
              <div className="text-6xl mb-3 group-hover:scale-110 transition text-white">
                +
              </div>
              <h3 className="text-sm font-semibold text-gray-400 text-center group-hover:text-blue-400 transition">
                Add Folder
              </h3>
            </div>
          </button>

          {/* Existing Folders */}
          {textbooks.map((textbook) => (
            <div
              key={textbook.id}
              className="group relative bg-gray-800/50 backdrop-blur-sm border border-gray-700 rounded-lg shadow-lg p-6 hover:border-blue-500/50 transition cursor-pointer"
              onClick={() => handleOpenFolder(textbook)}
            >
              {/* Folder Icon */}
              <div className="flex flex-col items-center">
                <div className="text-6xl mb-3 group-hover:scale-110 transition">
                  📁
                </div>
                <h3 className="text-sm font-semibold text-white text-center mb-2 line-clamp-2">
                  {textbook.title}
                </h3>

                {/* Status Badge */}
                <div className="text-xs mb-3">
                  {textbook.processed ? (
                    <span className="text-green-400 font-semibold">
                      {/* Empty when ready */}
                    </span>
                  ) : (
                    <span className="text-yellow-400 font-semibold">
                      ⏳ Processing
                    </span>
                  )}
                </div>

                {/* Delete Button */}
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    handleDelete(textbook.id, textbook.title);
                  }}
                  className="w-full px-3 py-1 bg-red-900/50 text-red-300 text-xs rounded hover:bg-red-800/50 transition border border-red-700"
                >
                  Delete
                </button>
              </div>
            </div>
          ))}
        </div>

        {/* Empty State */}
        {textbooks.length === 0 && (
          <div className="text-center py-40">
            <div className="text-6xl mb-4">📚</div>
            <p className="text-gray-400 text-lg mb-2">
              No textbooks yet. Click the + button to add one!
            </p>
            <p className="text-gray-500 text-sm">
              Each textbook will appear as a folder you can click to ask
              questions.
            </p>
          </div>
        )}
      </main>

      {/* Upload Modal */}
      {showUploadModal && (
        <div className="fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800/95 backdrop-blur-sm border border-gray-700 rounded-lg shadow-2xl p-8 w-full max-w-md">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-2xl font-semibold text-white">
                Upload New Textbook
              </h2>
              <button
                onClick={() => {
                  setShowUploadModal(false);
                  setUploadFile(null);
                  setUploadTitle("");
                  setError("");
                }}
                className="text-gray-400 hover:text-white text-2xl"
              >
                ×
              </button>
            </div>

            <form onSubmit={handleUpload} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Textbook Title
                </label>
                <input
                  type="text"
                  value={uploadTitle}
                  onChange={(e) => setUploadTitle(e.target.value)}
                  placeholder="e.g., Biology 101"
                  required
                  className="w-full px-4 py-2 bg-gray-900/50 border border-gray-600 text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none placeholder-gray-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  PDF File
                </label>
                <input
                  type="file"
                  accept=".pdf"
                  onChange={(e) => setUploadFile(e.target.files?.[0] || null)}
                  required
                  className="w-full text-gray-300 file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:bg-blue-600 file:text-white file:font-semibold hover:file:bg-blue-500 file:cursor-pointer"
                />
              </div>

              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => {
                    setShowUploadModal(false);
                    setUploadFile(null);
                    setUploadTitle("");
                    setError("");
                  }}
                  className="flex-1 px-6 py-2 bg-gray-700 text-gray-300 rounded-lg hover:bg-gray-600 transition"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isUploading || !uploadFile}
                  className="flex-1 px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-500 disabled:bg-gray-600 transition shadow-lg shadow-blue-500/20"
                >
                  {isUploading ? "Uploading..." : "Upload"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
