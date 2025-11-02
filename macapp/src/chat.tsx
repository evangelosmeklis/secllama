import { useState, useRef, useEffect } from 'react'
import { getCurrentWindow } from '@electron/remote'

interface Message {
  role: 'user' | 'assistant'
  content: string
}

export default function Chat() {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [model, setModel] = useState('llama3.2')
  const [loading, setLoading] = useState(false)
  const [models, setModels] = useState<string[]>([])
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

  const clearChat = () => {
    setMessages([])
  }

  return (
    <div className="flex h-screen flex-col bg-gray-50">
      {/* Header */}
      <div className="drag flex items-center justify-between border-b bg-white px-4 py-3">
        <div className="flex items-center space-x-3">
          <h1 className="text-lg font-semibold text-gray-900">SecLlama Chat 🔒</h1>
          <select
            value={model}
            onChange={(e) => setModel(e.target.value)}
            className="no-drag rounded-md border border-gray-300 px-2 py-1 text-sm focus:border-black focus:outline-none"
          >
            {models.length === 0 ? (
              <option>No models found</option>
            ) : (
              models.map(m => <option key={m} value={m}>{m}</option>)
            )}
          </select>
        </div>
        <div className="flex space-x-2">
          <button
            onClick={clearChat}
            className="no-drag rounded-md px-3 py-1 text-sm text-gray-600 hover:bg-gray-100"
          >
            Clear
          </button>
          <button
            onClick={() => getCurrentWindow().close()}
            className="no-drag rounded-md px-3 py-1 text-sm text-gray-600 hover:bg-gray-100"
          >
            Close
          </button>
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4">
        {messages.length === 0 ? (
          <div className="flex h-full items-center justify-center text-gray-400">
            <div className="text-center">
              <p className="text-lg">Start a conversation</p>
              <p className="mt-2 text-sm">All messages are encrypted end-to-end 🔒</p>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            {messages.map((msg, idx) => (
              <div
                key={idx}
                className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}
              >
                <div
                  className={`max-w-[70%] rounded-lg px-4 py-2 ${
                    msg.role === 'user'
                      ? 'bg-black text-white'
                      : 'bg-white text-gray-900 shadow-sm'
                  }`}
                >
                  <p className="whitespace-pre-wrap text-sm">{msg.content}</p>
                </div>
              </div>
            ))}
            {loading && (
              <div className="flex justify-start">
                <div className="max-w-[70%] rounded-lg bg-white px-4 py-2 shadow-sm">
                  <p className="text-sm text-gray-400">Thinking...</p>
                </div>
              </div>
            )}
            <div ref={messagesEndRef} />
          </div>
        )}
      </div>

      {/* Input */}
      <div className="border-t bg-white p-4">
        <div className="flex space-x-2">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && !loading && sendMessage()}
            placeholder="Type your message..."
            disabled={loading || !model}
            className="no-drag flex-1 rounded-lg border border-gray-300 px-4 py-2 focus:border-black focus:outline-none disabled:bg-gray-100"
          />
          <button
            onClick={sendMessage}
            disabled={loading || !input.trim() || !model}
            className="no-drag rounded-lg bg-black px-6 py-2 text-white hover:brightness-110 disabled:opacity-50"
          >
            Send
          </button>
        </div>
        {models.length === 0 && (
          <p className="mt-2 text-xs text-red-600">
            No models found. Run: secllama pull llama3.2
          </p>
        )}
      </div>
    </div>
  )
}

