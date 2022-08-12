import React, { useEffect } from 'react'
import axios from 'axios'
import { useToast } from '@chakra-ui/react'
import { useNavigate } from 'react-router-dom'

const GoogleLanding = () => {
  const toast = useToast()
  const navigate = useNavigate()

  useEffect(() => {
    let currentURL = window.location.href
    const params = new URLSearchParams(currentURL)
    const idToken = params.get('id_token')
    if (idToken !== null) {
      axios
        .request({
          method: 'POST',
          url: `${process.env.REACT_APP_API_URL}/auth/refresh_token`,
          headers: { Authorization: idToken },
        })
        .then(res => {
          if (res.status === 200) {
            localStorage.setItem(
              'user-access-token',
              res.headers['user-access-token']
            )
            localStorage.setItem(
              'user-refresh-token',
              res.headers['user-refresh-token']
            )
            navigate('/admin')
          }
        })
        .catch(() => {
          // toast({
          //   title: 'Failed to Load',
          //   description: 'Something went wrong on our side!',
          //   status: 'error',
          //   duration: 10,
          //   isClosable: false,
          //   position: 'top',
          // })
        })
    }
  }, [])
  return <div>Google Landing Page</div>
}

export default GoogleLanding
