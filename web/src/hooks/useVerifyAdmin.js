import { useEffect, useState } from 'react'
import logger from '../logger'

export function useVerifyAdmin() {
  const [isAdmin, setIsAdmin] = useState(false)

  useEffect(() => {
    fetch(`${process.env.REACT_APP_API_URL}/is_admin`, {
      method: 'GET',
      credentials: 'include',
    })
      .then((res) => {
        if (res.status === 200) {
          setIsAdmin(true)
        }
      })
      .catch((err) => {
        logger('error verifying is_admin', err)
      })
  }, [])

  return isAdmin
}
