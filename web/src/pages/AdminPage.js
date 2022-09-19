import React, { useCallback, useEffect } from 'react'
import {
  FormControl,
  FormHelperText,
  FormLabel,
  Input,
  VStack,
  Text,
  useColorModeValue,
  Container,
} from '@chakra-ui/react'
import { useVerifyAdmin } from '../hooks/useVerifyAdmin'
import logger from '../logger'
import { useNavigate } from 'react-router-dom'

export default function AdminPage() {
  const navigate = useNavigate()
  const isAdmin = useVerifyAdmin()
  if (!isAdmin) navigate('/login')

  const stackBgColor = useColorModeValue('white', 'gray.900')

  const getConfigs = useCallback(async () => {
    await fetch(`${process.env.REACT_APP_API_URL}/configs`, {
      method: 'GET',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
    })
      .then(async (res) => {
        if (!res) return
        const data = await res.json().catch((err) => logger(err))
        console.log(data)
      })
      .catch((err) => {
        logger('error getting items', err)
      })
  }, [])

  useEffect(() => {
    getConfigs()
  }, [getConfigs])

  return (
    <VStack>
      <Container maxW="container.md" centerContent>
        <Text fontSize="2xl" as="b" mb={5} mt={5}>
          Configuration
        </Text>

        <FormControl
          borderWidth="1px"
          borderRadius="lg"
          bg={stackBgColor}
          boxShadow={'2xl'}
          p={5}
          mb={5}
        >
          <FormLabel>Email address</FormLabel>
          <Input type="email" />
          <FormHelperText>We'll never share your email.</FormHelperText>
        </FormControl>
      </Container>
    </VStack>
  )
}
