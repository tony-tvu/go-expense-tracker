import { useEffect, useState } from 'react'
import { usePlaidLink } from 'react-plaid-link'
import Button from 'plaid-threads/Button'
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
import logger from '../logger'

function Homepage() {
  const [linkToken, setLinkToken] = useState(null)

  // fetch link_token on page load
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
    ).catch(e => {
      logger('error fetching link_token', e)
    })
    if (!response) return
    const data = await response.json()
    setLinkToken(data.link_token)
  }

  const onSuccess = async (public_token) => {
    await fetch(
      `${process.env.REACT_APP_API_URL}/api/set_access_token`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Public-Token': public_token
        },
      }
    ).catch(e => {
      logger('error setting access token', e)
    })
  }

  const config = {
    token: linkToken,
    onSuccess,
  }
  const { open, ready } = usePlaidLink(config)

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
