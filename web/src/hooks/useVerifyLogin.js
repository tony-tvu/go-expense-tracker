import { useEffect, useState } from 'react'
import logger from '../logger'

export function useVerifyLogin() {
  const [isLoggedIn, setIsLoggedIn] = useState(false)

  useEffect(() => {
    fetch(`${process.env.REACT_APP_API_URL}/logged_in`, {
      method: 'GET',
      credentials: 'include',
    })
      .then((res) => {
        if (res.status === 200) {
          setIsLoggedIn(true)
        }
      })
      .catch((err) => {
        logger('error verifying login state', err)
      })
  }, [])

  return isLoggedIn
}
