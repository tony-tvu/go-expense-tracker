import { useEffect, useState } from 'react'
import logger from '../logger'

export function useRegistration() {
  const [registrationEnabled, setRegistrationEnabled] = useState(false)

  useEffect(() => {
    fetch(`${process.env.REACT_APP_API_URL}/registration_enabled`, {
      method: 'GET',
      credentials: 'include',
    })
      .then(async (res) => {
        if (!res) return
        const data = await res.json().catch((err) => logger(err))
        if (data && data.registration_enabled) {
          setRegistrationEnabled(true)
        }
      })
      .catch((err) => {
        logger('error getting registration_enabled', err)
      })
  }, [])

  return registrationEnabled
}
