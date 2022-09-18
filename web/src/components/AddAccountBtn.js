import React, { useEffect, useState } from 'react'
import logger from '../logger'
import { usePlaidLink } from 'react-plaid-link'
import { BsPlus } from 'react-icons/bs'
import { Button } from '@chakra-ui/react'
import { colors } from '../theme'

export default function AddAccountBtn(props) {
  const [linkToken, setLinkToken] = useState(null)

  useEffect(() => {
    fetchLinkToken()
  }, [])

  // link_token is required to start linking a bank account
  const fetchLinkToken = async () => {
    await fetch(`${process.env.REACT_APP_API_URL}/link_token`, {
      method: 'GET',
      credentials: 'include',
    })
      .then(async (res) => {
        if (!res) return
        const data = await res.json().catch((err) => logger(err))
        setLinkToken(data?.link_token)
      })
      .catch((err) => {
        logger('error fetching link_token', err)
      })
  }

  /*
   * Upon linking success, plaid api will return a public_token which will be used
   * to get a permanent access_token for the user's specific linked bank account.
   */
  const onLinkAccountSuccess = async (public_token) => {
    await fetch(`${process.env.REACT_APP_API_URL}/items`, {
      method: 'POST',
      credentials: 'include',
      body: JSON.stringify({ public_token: public_token }),
    })
      .then((res) => {
        if (res.status === 200) props.onSuccess()
      })
      .catch((e) => {
        logger('error setting access token', e)
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
      color={'white'}
      _hover={{
        bg: colors.primaryFaded,
      }}
    >
      Add account
    </Button>
  )
}
