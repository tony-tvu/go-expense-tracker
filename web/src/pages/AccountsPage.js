import { useEffect, useState } from "react"
import { usePlaidLink } from "react-plaid-link"
import { Button } from "@chakra-ui/react"
import { Text, Center, Flex } from "@chakra-ui/react"
import logger from "../logger"
import { useVerifyLogin } from "../hooks/useVerifyLogin"
import Navbar from "../components/Navbar"
import { useQuery, gql, useMutation } from "@apollo/client"

const query = gql`
  query {
    linkToken
  }
`

const mutation = gql`
  mutation ($input: PublicTokenInput!) {
    setAccessToken(input: $input)
  }
`

export default function Accounts() {
  useVerifyLogin()

  // link_token is required to start linking a bank account
  const [linkToken, setLinkToken] = useState(null)

  const { data } = useQuery(query, {
    fetchPolicy: "no-cache",
  })

  const [setAccessToken] = useMutation(mutation)

  useEffect(() => {
    if (data && data.linkToken) {
      setLinkToken(data.linkToken)
    }
  }, [data])

  /*
   * Upon linking success, plaid api will return a public_token which will be used
   * to get a permanent access_token for the user's specific linked bank account.
   */
  const onLinkAccountSuccess = async (public_token) => {
    setAccessToken({
      variables: {
        input: {
          publicToken: public_token,
        },
      },
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
      <Navbar />
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
