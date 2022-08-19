import { useEffect, useState } from 'react'
import { usePlaidLink } from 'react-plaid-link'
import Button from 'plaid-threads/Button'
import axios from 'axios'
import {
  Box,
  Grid,
  GridItem,
  Text,
  Center,
  Square,
  Flex,
} from '@chakra-ui/react'
import { useNavigate } from 'react-router-dom'

function Homepage() {
  const [linkToken, setLinkToken] = useState(null)

  const onSuccess = (public_token) => {
    console.log(public_token)
  }

  const config = {
    token: linkToken,
    onSuccess,
  }

  const { open, ready } = usePlaidLink(config)

  // load link_token
  useEffect(() => {
    fetchLinkToken()
  }, [])

  const navigate = useNavigate()

  async function fetchLinkToken() {
    const response = await fetch(
      `${process.env.REACT_APP_API_URL}/api/create_link_token`,
      {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      }
    )
    if (!response.ok) {
      console.error('error fetching link_token')
      return
    }
    const data = await response.json()
    setLinkToken(data.link_token)
  }

  return (
    <div>
      <Flex color="white" minH="85vh">
        <Button type="button" large onClick={() => open()} disabled={!ready}>
          Connect Account
        </Button>
        {/* <Center w="500px" bg="green.500">
          <Text>Box 1</Text>
        </Center> */}
      </Flex>
    </div>
  )
}

export default Homepage
