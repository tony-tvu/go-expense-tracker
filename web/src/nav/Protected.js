import React, { useEffect, useState } from 'react'
import { Navigate } from 'react-router-dom'
import logger from '../logger'
import Sidenav from './Sidenav'

export default function Protected({ adminOnly, current, children }) {
  const [isLoggedIn, setIsLoggedIn] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch(`${process.env.REACT_APP_API_URL}/logged_in`, {
      method: 'GET',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
    })
      .then(async (res) => {
        if (!res) return
        const data = await res.json().catch((err) => logger(err))
        if (data && data.logged_in) {
          setIsLoggedIn(true)
        } else {
          setIsLoggedIn(false)
        }
        setLoading(false)
      })
      .catch((err) => {
        logger('error verifying login state', err)
      })
  }, [])

  if (loading) {
    return null
  }

  if (!isLoggedIn) {
    return <Navigate to="/login" replace />
  }

  return (
    <Sidenav isLoggedIn={isLoggedIn} current={current}>
      {children}
    </Sidenav>
  )
}
