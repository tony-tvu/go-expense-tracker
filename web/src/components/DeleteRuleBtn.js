import React from 'react'
import { Button } from '@chakra-ui/react'
import logger from '../logger'
import { useNavigate } from 'react-router-dom'

export default function DeleteRuleBtn({ rule, onSuccess }) {
  const navigate = useNavigate()

  async function deleteAccount() {
    await fetch(`${process.env.REACT_APP_API_URL}/rules/${rule.id}`, {
      method: 'DELETE',
      credentials: 'include',
    })
      .then((res) => {
        if (res.status === 401) {
          navigate('/login')
        }
        if (res.status === 200) {
          onSuccess()
        }
      })
      .catch((e) => {
        logger('error deleting account', e)
      })
  }

  return (
    <Button
      variant="solid"
      justifyContent="space-between"
      fontWeight="normal"
      fontSize="md"
      as="b"
      onClick={async () => await deleteAccount()}
      color={'white'}
      bg={'#E63E3F'}
      _hover={{
        bg: '#AA2E2F',
      }}
    >
      Delete
    </Button>
  )
}
