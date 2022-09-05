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
import { useLoginStatus } from "../hooks/useLoginStatus"

export default function ConnectAccount() {
  const navigate = useNavigate()
  const isLoggedIn = useLoginStatus()
  if (!isLoggedIn) navigate("/login")

  const [linkToken, setLinkToken] = useState(null)

  useEffect(() => {
    fetchLinkToken()
  }, [])

  // link_token is required to start linking a bank account
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

  /*
   * Upon linking success, plaid api will return a public_token which will be used
   * to get a permanent access_token for the user's specific linked bank account.
   */
  const onLinkAccountSuccess = async (public_token) => {
    await fetch(`${process.env.REACT_APP_API_URL}/set_access_token`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Plaid-Public-Token": public_token,
      },
      credentials: "include",
    }).catch((e) => {
      logger("error setting access token", e)
    })
  }

  const plaidConfig = {
    token: linkToken,
    onSuccess: onLinkAccountSuccess,
  }
  const { open: openLinkingPopup, ready: isReadyToLinkAccount } =
    usePlaidLink(plaidConfig)

  return (
    <div>
      <Flex color="white" minH="85vh">
        <Button
          type="button"
          onClick={() => openLinkingPopup()}
          disabled={!isReadyToLinkAccount}
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
