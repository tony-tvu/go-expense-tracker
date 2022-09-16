import React, { useEffect, useState } from "react"
import logger from "../logger"
import { useQuery, gql, useMutation } from "@apollo/client"
import { usePlaidLink } from "react-plaid-link"
import { BsPlus } from "react-icons/bs"
import { Button } from "@chakra-ui/react"
import { colors } from "../theme"

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

export default function AddAccountBtn(props) {
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
    })
      .then((res) => {
        if (!res.errors) {
          props.onSuccess()
        }
      })
      .catch((e) => {
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
    <Button
      leftIcon={<BsPlus />}
      type="button"
      variant="solid"
      onClick={() => openLinkingPopup()}
      disabled={!isReadyToLinkAccount}
      bg={colors.primary}
      color={"white"}
      _hover={{
        bg: colors.primaryFaded,
      }}
    >
      Add account
    </Button>
  )
}
