import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { handleCallback } from '../auth'

export default function CallbackPage() {
  const navigate = useNavigate()

  useEffect(() => {
    handleCallback()
      .catch(err => console.error('Auth callback failed:', err))
      .finally(() => navigate('/booking', { replace: true }))
  }, [navigate])

  return null
}
