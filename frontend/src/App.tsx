import { useEffect, useState } from 'react'
import './index.css'

interface HealthResponse {
  status: string
  service: string
}

function App() {
  const [status, setStatus] = useState<string>('loading')
  const [error, setError] = useState<string>('')

  useEffect(() => {
    fetch('/api/health')
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`)
        return res.json()
      })
      .then((data: HealthResponse) => {
        setStatus(data.status)
      })
      .catch((err: Error) => {
        setError(err.message)
      })
  }, [])

  if (error) {
    return (
      <div className="app">
        <h1>Project Daedalus</h1>
        <p className="error">Failed to reach backend: {error}</p>
      </div>
    )
  }

  return (
    <div className="app">
      <h1>Project Daedalus</h1>
      <p className="status">status: {status}</p>
    </div>
  )
}

export default App
