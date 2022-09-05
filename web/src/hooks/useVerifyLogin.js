import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import logger from "../logger"

export function useVerifyLogin() {
  const navigate = useNavigate()

  useEffect(() => {
    fetch(`${process.env.REACT_APP_API_URL}/logged_in`, {
      method: "GET",
      credentials: "include",
    })
      .then((res) => {
        if (res.status !== 200) {
          navigate("/login")
        } 
      })
      .catch((err) => {
        logger("error verifying login state", err)
      })
  }, [navigate])
}
