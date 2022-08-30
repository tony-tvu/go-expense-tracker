import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'

export default function AdminPage() {
  const [refreshToken, setRefreshToken] = useState('')
  const [accessToken, setAccessToken] = useState('')
  const [refreshTokenExp, setRefreshTokenExp] = useState('')
  const [accessTokenExp, setAccessTokenExp] = useState('')
  const [now, setNow] = useState('')

  const navigate = useNavigate()

  useEffect(() => {
    setRefreshToken(localStorage.getItem('user-refresh-token'))
    setAccessToken(localStorage.getItem('user-access-token'))
  }, [])

  function fetchTokenExpirations() {

  }

  return (
    <div style={{ marginLeft: '30px' }}>
      <h1>
        <b>refresh token: </b>
        <span>{refreshToken}</span>
      </h1>
      <h1>
        <b>access token: </b> <span>{accessToken}</span>
      </h1>
      <br />
      <br />
      <h1>
        <b>refresh expiration: </b>
        <span>{refreshTokenExp}</span>
      </h1>
      <h1>
        <b>access expiration: </b>
        <span>{accessTokenExp}</span>
      </h1>
      <h1>
        <b>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;now:{' '}
        </b>
        <span>{now}</span>
      </h1>
      <button
        onClick={() => {
          fetchTokenExpirations()
          setNow(new Date().toString())
        }}
      >
        FETCH TOKENS
      </button>
    </div>
  )
}
