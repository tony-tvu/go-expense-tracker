import React, { useEffect, useState } from 'react'
import logger from '../logger'
import { BsPlus } from 'react-icons/bs'
import { Button } from '@chakra-ui/react'
import { colors } from '../theme'
import { useNavigate } from 'react-router-dom'
import { loadScript } from '../util'

// Ensure the Teller Connect script is loaded
// Returns the `window.TellerConnect` object once it exists
function loadTellerConnect() {
  return new Promise((resolve) => {
    function check() {
      if (window.TellerConnect) {
        return resolve(window.TellerConnect)
      }
      loadScript('teller-script', 'https://cdn.teller.io/connect/connect.js')
      setTimeout(check, 100)
    }
    check()
  })
}

export default function AddAccountBtn({ onSuccess }) {
  const [tellerApi, setTellerApi] = useState(null)
  const navigate = useNavigate()

  useEffect(() => {
    loadTellerConnect().then((tellerApi) => {
      setTellerApi(tellerApi)
    })
  }, [])

  async function createEnrollment(enrollment) {
    await fetch(`${process.env.REACT_APP_API_URL}/enrollments`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        access_token: enrollment.accessToken,
        enrollment_id: enrollment.enrollment.id,
        institution: enrollment.enrollment.institution.name,
      }),
    })
      .then((res) => {
        if (res.status === 401) navigate('/login')
        if (res.status === 200) onSuccess()
      })
      .catch((e) => {
        logger('error saving access token', e)
      })
  }
  return (
    <Button
      leftIcon={<BsPlus />}
      type="button"
      variant="solid"
      onClick={async () => {
        const res = await tellerApi.setup({
          applicationId: process.env.REACT_APP_TELLER_APPLICATION_ID,
          environment: process.env.REACT_APP_TELLER_ENV,
          onSuccess: async (enrollment) => {
            await createEnrollment(enrollment)
          },
        })
        res.open()
      }}
      bg={colors.primary}
      color={'white'}
      _hover={{
        bg: colors.primaryFaded,
      }}
    >
      Add account
    </Button>
  )
}
