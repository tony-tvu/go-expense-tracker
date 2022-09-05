import { useState, useEffect } from "react"
import logger from "../logger"

export function useLoginStatus() {
  const [isLoggedIn, setIsLoggedIn] = useState(false)

  useEffect(() => {
    fetch(`${process.env.REACT_APP_API_URL}/logged_in`, {
      method: "GET",
      credentials: "include",
    })
      .then((res) => {
        if (res.status === 200) {
          setIsLoggedIn(true)
        } else {
          setIsLoggedIn(false)
        }
      })
      .catch((err) => {
        logger("error verifying login state", err)
      })
  }, [])

  return isLoggedIn
}
