import { useState } from 'react'
import { usePlaidLink } from 'react-plaid-link'
import Button from 'plaid-threads/Button'
import {
  Box,
  Grid,
  GridItem,
  Flex,
  Text,
  Center,
  Square,
} from '@chakra-ui/react'
import GoogleLoginBtn from '../components/GoogleLoginBtn'

function Homepage() {
  const [linkToken, setLinkToken] = useState('')

  function fetchLinkToken() {
    axios
      .request({
        method: 'GET',
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

  return (
    <div>
      <Flex color="white" minH="85vh">
        <Button type="button" large onClick={() => open()} disabled={!ready}>
          Launch Link
        </Button>
        <Center w="500px" bg="green.500">
          <Text>Box 1</Text>
        </Center>
      </Flex>
    </div>
  )
}

export default Homepage
