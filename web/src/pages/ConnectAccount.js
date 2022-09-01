import { useEffect, useState } from "react"
import { usePlaidLink } from "react-plaid-link"
import { Button } from "@chakra-ui/react"
import {
  Box,
  Grid,
  GridItem,
  Text,
  Center,
  Square,
  Flex,
} from "@chakra-ui/react"
import { useNavigate } from "react-router-dom"
import logger from "../logger"

export default function ConnectAccount() {
  const [linkToken, setLinkToken] = useState(null)
  const navigate = useNavigate()

  useEffect(() => {
    fetch(`${process.env.REACT_APP_API_URL}/ping`, {
      method: "GET",
      credentials: "include",
    })
      .then((res) => {
        if (res.status !== 200) navigate("/login")
        else fetchLinkToken()
      })
      .catch((err) => {
        logger("error pinging server", err)
      })
  }, [navigate])

  const fetchLinkToken = async () => {
    const response = await fetch(
      `${process.env.REACT_APP_API_URL}/create_link_token`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
      }
    ).catch((err) => {
      logger("error fetching link_token", err)
    })

    if (!response) return
    const data = await response.json().catch((err) => logger(err))
    setLinkToken(data?.link_token)
  }

  const onSuccess = async (public_token) => {
    await fetch(`${process.env.REACT_APP_API_URL}/set_access_token`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Plaid-Public-Token": public_token,
      },
    }).catch((e) => {
      logger("error setting access token", e)
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
        <Button
          type="button"
          onClick={() => open()}
          disabled={!ready}
          colorScheme="teal"
          size="md"
        >
          Connect Account
        </Button>
        <Center w="500px" bg="green.500">
          <Text>Box 1</Text>
        </Center>
      </Flex>
    </div>
  )
}
