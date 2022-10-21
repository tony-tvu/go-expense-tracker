import React, { useEffect, useState } from 'react'
import {
  FormControl,
  FormLabel,
  VStack,
  useColorModeValue,
  Container,
  Divider,
  Button,
  Input,
  InputGroup,
  InputRightElement,
  useToast,
} from '@chakra-ui/react'
import logger from '../logger'
import { colors } from '../theme'
import Sidenav from '../nav/Sidenav'
import { useNavigate } from 'react-router-dom'

export default function RegisterUser() {
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const handleShowClick = () => setShowPassword(!showPassword)
  const stackBgColor = useColorModeValue('white', 'gray.900')

  const navigate = useNavigate()
  const toast = useToast()

  useEffect(() => {
    document.title = 'Register'
  }, [])

  async function handleSubmit(e) {
    e.preventDefault()
    await fetch(`${process.env.REACT_APP_API_URL}/register`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: username,
        email: email,
        password: password,
      }),
    })
      .then(async (res) => {
        if (!res) return
        const data = await res.json().catch((err) => logger(err))
        if (res.status === 200) {
          toast({
            title: 'Success!',
            description: 'New user created',
            status: 'success',
            position: 'top-right',
            duration: 5000,
            isClosable: true,
          })
          navigate('/login')
        }
        if (res.status !== 200) {
          toast({
            title: 'Registration failed',
            description: data.error,
            status: 'error',
            position: 'top-right',
            duration: 5000,
            isClosable: true,
          })
        }
      })
      .catch((e) => {
        logger('error registering new user', e)
      })
  }

  return (
    <Sidenav>
      <VStack>
        <Container maxW="container.md" mt={3}>
          <FormControl bg={stackBgColor} p={5}>
            <FormLabel fontSize="xl">New User</FormLabel>
            <Divider mb={5} />

            <FormLabel mt={5}>Username</FormLabel>
            <Input
              type="text"
              value={username}
              onChange={(event) => setUsername(event.target.value)}
            />

            <FormLabel mt={5}>Email</FormLabel>
            <Input
              type="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
            />

            <FormLabel mt={5}>Password</FormLabel>
            <InputGroup>
              <Input
                onChange={(event) => setPassword(event.target.value)}
                type={showPassword ? 'text' : 'password'}
                _placeholder={{ color: 'gray.500' }}
                _hover={{
                  borderColor: 'gray.500',
                }}
                color={useColorModeValue('black', 'white')}
              />
              <InputRightElement width="4.5rem">
                <Button
                  h="1.75rem"
                  size="sm"
                  onClick={handleShowClick}
                  backgroundColor={useColorModeValue('gray.300', 'gray.900')}
                  color={useColorModeValue('black', 'white')}
                  _hover={{
                    backgroundColor: useColorModeValue('gray.400', 'gray.700'),
                  }}
                >
                  {showPassword ? 'Hide' : 'Show'}
                </Button>
              </InputRightElement>
            </InputGroup>
            <Button
              mt={5}
              onClick={handleSubmit}
              type="submit"
              variant="solid"
              bg={colors.primary}
              color={'white'}
              _hover={{
                bg: colors.primaryFaded,
              }}
            >
              Register
            </Button>
          </FormControl>
        </Container>
      </VStack>
    </Sidenav>
  )
}
