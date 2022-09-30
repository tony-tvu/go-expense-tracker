import React, { useEffect, useState } from 'react'
import { Navigate } from 'react-router-dom'
import logger from '../logger'
import Navbar from './Navbar'
import Sidenav from './Sidenav'

export default function Protected({ adminOnly, current, children }) {
  const [isLoggedIn, setIsLoggedIn] = useState(null)
  const [isAdmin, setIsAdmin] = useState(null)
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
        if (data && data.is_admin) {
          setIsAdmin(true)
        } else {
          setIsAdmin(false)
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

  if (adminOnly && !isAdmin) {
    return <Navigate to="/" replace />
  }

  return (
    <Sidenav isLoggedIn={isLoggedIn} isAdmin={isAdmin} current={current}>
      {children}
    </Sidenav>
  )
}
