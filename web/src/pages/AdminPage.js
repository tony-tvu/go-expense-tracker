import React, { useState, useEffect } from "react"
import { useNavigate } from "react-router-dom"
import { useLoginStatus } from "../hooks/useLoginStatus"

export default function AdminPage() {
  const navigate = useNavigate()
  const isLoggedIn = useLoginStatus()
  if (!isLoggedIn) navigate("/login")

  const [refreshToken, setRefreshToken] = useState("")
  const [accessToken, setAccessToken] = useState("")
  const [refreshTokenExp, setRefreshTokenExp] = useState("")
  const [accessTokenExp, setAccessTokenExp] = useState("")
  const [now, setNow] = useState("")

  function fetchTokenExpirations() {}

  return (
    <div style={{ marginLeft: "30px" }}>
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
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;now:{" "}
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
