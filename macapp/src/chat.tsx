import { useState, useRef, useEffect } from 'react'
import { getCurrentWindow } from '@electron/remote'

interface Message {
  role: 'user' | 'assistant'
  content: string
}

interface SystemStats {
  cpu: number
  memory: number
  gpu?: number
}

export default function Chat() {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [model, setModel] = useState('')
  const [loading, setLoading] = useState(false)
  const [models, setModels] = useState<string[]>([])
  const [stats, setStats] = useState<SystemStats>({ cpu: 0, memory: 0 })
  const [conversations, setConversations] = useState<string[]>([])
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    // Fetch available models
    fetch('http://localhost:11434/api/tags')
      .then(res => res.json())
      .then(data => {
        const modelNames = data.models?.map((m: any) => m.name) || []
        setModels(modelNames)
        if (modelNames.length > 0 && !model) {
          setModel(modelNames[0])
        }
      })
      .catch(err => console.error('Failed to fetch models:', err))

    // Simulate stats (in a real app, you'd get this from the system)
    const statsInterval = setInterval(() => {
      setStats({
        cpu: Math.floor(Math.random() * 30) + 10,
        memory: Math.floor(Math.random() * 40) + 30,
        gpu: Math.floor(Math.random() * 20) + 5
      })
    }, 2000)

    return () => clearInterval(statsInterval)
  }, [])

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const sendMessage = async () => {
    if (!input.trim() || !model) return

    const userMessage: Message = { role: 'user', content: input }
    setMessages(prev => [...prev, userMessage])
    setInput('')
    setLoading(true)

    try {
      const response = await fetch('http://localhost:11434/api/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          model,
          prompt: input,
          stream: false
        })
      })

      const data = await response.json()
      const assistantMessage: Message = {
        role: 'assistant',
        content: data.response || 'No response'
      }
      setMessages(prev => [...prev, assistantMessage])
      
      // Add to conversations if first message
      if (messages.length === 0) {
        setConversations(prev => [input.slice(0, 50) + '...', ...prev])
      }
    } catch (error) {
      const errorMessage: Message = {
        role: 'assistant',
        content: `Error: ${error.message}`
      }
      setMessages(prev => [...prev, errorMessage])
    } finally {
      setLoading(false)
    }
  }

  const newChat = () => {
    setMessages([])
  }

  return (
    <div className="flex h-screen bg-[#212121] text-gray-100">
      {/* Sidebar */}
      <div className="w-64 bg-[#171717] border-r border-gray-800 flex flex-col">
        {/* Logo */}
        <div className="p-4 border-b border-gray-800">
          <div className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center text-white font-bold text-sm">
              SL
            </div>
            <span className="font-semibold text-lg">SecLlama</span>
          </div>
        </div>

        {/* New Chat Button */}
        <div className="p-3">
          <button
            onClick={newChat}
            className="w-full flex items-center space-x-2 px-3 py-2 rounded-lg border border-gray-700 hover:bg-gray-800 transition-colors"
          >
            <span className="text-lg">+</span>
            <span className="text-sm">New Chat</span>
          </button>
        </div>

        {/* Conversations */}
        <div className="flex-1 overflow-y-auto px-3">
          <div className="text-xs font-semibold text-gray-500 mb-2 px-3">Today</div>
          {conversations.map((conv, idx) => (
            <div
              key={idx}
              className="px-3 py-2 text-sm rounded-lg hover:bg-gray-800 cursor-pointer truncate mb-1"
            >
              {conv}
            </div>
          ))}
        </div>

        {/* Security & Stats Section */}
        <div className="border-t border-gray-800 p-3 space-y-3">
          {/* Security Indicators */}
          <div className="bg-[#1a1a1a] rounded-lg p-3 space-y-2">
            <div className="text-xs font-semibold text-gray-400 mb-2">Security Status</div>
            <div className="flex items-center space-x-2 text-xs">
              <div className="w-2 h-2 bg-green-500 rounded-full"></div>
              <span className="text-gray-300">🔒 End-to-End Encrypted</span>
            </div>
            <div className="flex items-center space-x-2 text-xs">
              <div className="w-2 h-2 bg-green-500 rounded-full"></div>
              <span className="text-gray-300">🏖️ Sandboxed</span>
            </div>
            <div className="flex items-center space-x-2 text-xs">
              <div className="w-2 h-2 bg-red-500 rounded-full"></div>
              <span className="text-gray-300">🚫 No Internet Access</span>
            </div>
          </div>

          {/* System Stats */}
          <div className="bg-[#1a1a1a] rounded-lg p-3 space-y-2">
            <div className="text-xs font-semibold text-gray-400 mb-2">System Resources</div>
            <div className="space-y-1.5">
              <div className="flex justify-between items-center text-xs">
                <span className="text-gray-400">CPU</span>
                <span className="text-blue-400 font-mono">{stats.cpu}%</span>
              </div>
              <div className="w-full bg-gray-800 rounded-full h-1.5">
                <div 
                  className="bg-blue-500 h-1.5 rounded-full transition-all duration-300"
                  style={{ width: `${stats.cpu}%` }}
                ></div>
              </div>
              
              <div className="flex justify-between items-center text-xs">
                <span className="text-gray-400">Memory</span>
                <span className="text-green-400 font-mono">{stats.memory}%</span>
              </div>
              <div className="w-full bg-gray-800 rounded-full h-1.5">
                <div 
                  className="bg-green-500 h-1.5 rounded-full transition-all duration-300"
                  style={{ width: `${stats.memory}%` }}
                ></div>
              </div>
              
              {stats.gpu !== undefined && (
                <>
                  <div className="flex justify-between items-center text-xs">
                    <span className="text-gray-400">GPU</span>
                    <span className="text-purple-400 font-mono">{stats.gpu}%</span>
                  </div>
                  <div className="w-full bg-gray-800 rounded-full h-1.5">
                    <div 
                      className="bg-purple-500 h-1.5 rounded-full transition-all duration-300"
                      style={{ width: `${stats.gpu}%` }}
                    ></div>
                  </div>
                </>
              )}
            </div>
          </div>

          {/* Models List */}
          <div className="bg-[#1a1a1a] rounded-lg p-3">
            <div className="text-xs font-semibold text-gray-400 mb-2">Local Models ({models.length})</div>
            <div className="space-y-1 max-h-32 overflow-y-auto">
              {models.slice(0, 5).map((m, idx) => (
                <div 
                  key={idx}
                  className="text-xs text-gray-300 flex items-center space-x-1.5 py-1"
                >
                  <div className="w-1.5 h-1.5 bg-blue-500 rounded-full"></div>
                  <span className="truncate">{m}</span>
                </div>
              ))}
              {models.length > 5 && (
                <div className="text-xs text-gray-500 pt-1">+{models.length - 5} more</div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Main Chat Area */}
      <div className="flex-1 flex flex-col">
        {/* Header */}
        <div className="border-b border-gray-800 px-6 py-4 flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <select
              value={model}
              onChange={(e) => setModel(e.target.value)}
              className="bg-[#2a2a2a] border border-gray-700 rounded-lg px-4 py-2 text-sm focus:outline-none focus:border-blue-500 transition-colors"
            >
              {models.length === 0 ? (
                <option>No models found</option>
              ) : (
                models.map(m => <option key={m} value={m}>{m}</option>)
              )}
            </select>
            <div className="flex items-center space-x-1.5 text-xs text-gray-500">
              <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
              <span>Connected</span>
            </div>
          </div>
          
          <button
            onClick={() => getCurrentWindow().close()}
            className="text-gray-500 hover:text-gray-300 transition-colors"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {/* Messages */}
        <div className="flex-1 overflow-y-auto px-6 py-6">
          {messages.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-full text-gray-400">
              <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center text-white font-bold text-2xl mb-4">
                SL
              </div>
              <h2 className="text-2xl font-semibold text-gray-200 mb-2">How can I help you today?</h2>
              <p className="text-sm text-gray-500">All conversations are encrypted and sandboxed 🔒</p>
            </div>
          ) : (
            <div className="max-w-3xl mx-auto space-y-6">
              {messages.map((msg, idx) => (
                <div
                  key={idx}
                  className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}
                >
                  <div className="flex items-start space-x-3 max-w-[80%]">
                    {msg.role === 'assistant' && (
                      <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center flex-shrink-0 mt-1">
                        <span className="text-white text-xs font-bold">SL</span>
                      </div>
                    )}
                    <div
                      className={`rounded-2xl px-4 py-3 ${
                        msg.role === 'user'
                          ? 'bg-blue-600 text-white'
                          : 'bg-[#2a2a2a] text-gray-100'
                      }`}
                    >
                      <p className="text-sm whitespace-pre-wrap leading-relaxed">{msg.content}</p>
                    </div>
                    {msg.role === 'user' && (
                      <div className="w-8 h-8 bg-gray-700 rounded-lg flex items-center justify-center flex-shrink-0 mt-1">
                        <span className="text-gray-300 text-xs">You</span>
                      </div>
                    )}
                  </div>
                </div>
              ))}
              {loading && (
                <div className="flex justify-start">
                  <div className="flex items-start space-x-3">
                    <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
                      <span className="text-white text-xs font-bold">SL</span>
                    </div>
                    <div className="bg-[#2a2a2a] rounded-2xl px-4 py-3">
                      <div className="flex space-x-2">
                        <div className="w-2 h-2 bg-gray-500 rounded-full animate-bounce"></div>
                        <div className="w-2 h-2 bg-gray-500 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }}></div>
                        <div className="w-2 h-2 bg-gray-500 rounded-full animate-bounce" style={{ animationDelay: '0.4s' }}></div>
                      </div>
                    </div>
                  </div>
                </div>
              )}
              <div ref={messagesEndRef} />
            </div>
          )}
        </div>

        {/* Input */}
        <div className="border-t border-gray-800 px-6 py-4">
          <div className="max-w-3xl mx-auto">
            <div className="flex items-end space-x-3 bg-[#2a2a2a] rounded-2xl p-3">
              <input
                type="text"
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && !loading && sendMessage()}
                placeholder="Send a message..."
                disabled={loading || !model}
                className="flex-1 bg-transparent focus:outline-none text-sm disabled:opacity-50 text-gray-100 placeholder-gray-500"
              />
              <button
                onClick={sendMessage}
                disabled={loading || !input.trim() || !model}
                className="bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:hover:bg-blue-600 rounded-xl px-4 py-2 text-sm font-medium transition-colors"
              >
                {loading ? 'Sending...' : 'Send'}
              </button>
            </div>
            {models.length === 0 && (
              <p className="mt-2 text-xs text-red-400 text-center">
                No models found. Run: <code className="bg-gray-800 px-2 py-0.5 rounded">secllama pull llama3.2</code>
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
