import React, { useEffect, useState } from 'react'
import {
  Flex,
  Input,
  Button,
  InputGroup,
  Stack,
  InputLeftElement,
  chakra,
  Box,
  FormControl,
  InputRightElement,
  useColorModeValue,
  useToast,
} from '@chakra-ui/react'
import { FaUserAlt, FaLock, FaCat } from 'react-icons/fa'
import { colors } from '../theme'
import { useNavigate } from 'react-router-dom'
import logger from '../logger'
import Sidenav from '../nav/Sidenav'

const CFaUserAlt = chakra(FaUserAlt)
const CFaLock = chakra(FaLock)
const CFcat = chakra(FaCat)

export default function Login() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const handleShowClick = () => setShowPassword(!showPassword)
  const navigate = useNavigate()
  const toast = useToast()

  const logoColor = useColorModeValue('black', colors.primary)

  useEffect(() => {
    document.title = 'Login'
    fetch(`${process.env.REACT_APP_API_URL}/logged_in`, {
      method: 'GET',
      credentials: 'include',
    })
      .then(async (res) => {
        if (!res) return
        const data = await res.json().catch((err) => logger(err))
        if (data && data.logged_in) {
          navigate('/')
        }
      })
      .catch((err) => {
        logger('error verifying login state', err)
      })
  }, [navigate])

  async function handleSubmit(e) {
    e.preventDefault()
    await fetch(`${process.env.REACT_APP_API_URL}/login`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username: username, password: password }),
    })
      .then((res) => {
        if (res.status === 200) navigate('/')
        if (res.status === 404 || res.status === 403) {
          toast({
            title: 'Login failed',
            description: 'Email or password is incorrect',
            status: 'error',
            position: 'top-right',
            duration: 5000,
            isClosable: true,
          })
        }
        if (res.status === 429) {
          toast({
            title: 'Too many login attemps!',
            description: 'Try again in 1 minute',
            status: 'error',
            position: 'top-right',
            duration: 5000,
            isClosable: true,
          })
        }
      })
      .catch((e) => {
        logger('error logging in', e)
      })
  }

  return (
    <Sidenav>
      <Flex flexDirection="column">
        <Stack flexDir="column" mt="15%" alignItems="center">
          <CFcat width={'10vh'} size={'100px'} color={logoColor} />
          <Box
            minW={{ base: '90%', md: '468px' }}
            backgroundColor={'whiteAlpha.800'}
          >
            <form onSubmit={handleSubmit}>
              <Stack
                spacing={4}
                p="1rem"
                backgroundColor={useColorModeValue('whiteAlpha.800', '#252526')}
                boxShadow={'2xl'}
                borderWidth="1px"
              >
                <FormControl>
                  <InputGroup>
                    <InputLeftElement
                      pointerEvents="none"
                      children={<CFaUserAlt color={'gray.500'} />}
                    />
                    <Input
                      type="username"
                      placeholder="username"
                      _placeholder={{ color: 'gray.500' }}
                      borderColor={useColorModeValue('gray.300', 'gray.600')}
                      _hover={{
                        borderColor: 'gray.500',
                      }}
                      onChange={(event) => setUsername(event.target.value)}
                      color={useColorModeValue('black', 'white')}
                      bg={useColorModeValue('whiteAlpha.800', '#252526')}
                    />
                  </InputGroup>
                </FormControl>
                <FormControl>
                  <InputGroup>
                    <InputLeftElement
                      pointerEvents="none"
                      color={'gray.500'}
                      children={<CFaLock color={'gray.500'} />}
                    />
                    <Input
                      onChange={(event) => setPassword(event.target.value)}
                      type={showPassword ? 'text' : 'password'}
                      placeholder="password"
                      _placeholder={{ color: 'gray.500' }}
                      borderColor={useColorModeValue('gray.300', 'gray.600')}
                      _hover={{
                        borderColor: 'gray.500',
                      }}
                      color={useColorModeValue('black', 'white')}
                      bg={useColorModeValue('whiteAlpha.800', '#252526')}
                    />
                    <InputRightElement width="4.5rem">
                      <Button
                        h="1.75rem"
                        size="sm"
                        onClick={handleShowClick}
                        backgroundColor={useColorModeValue(
                          'gray.300',
                          'gray.900'
                        )}
                        color={useColorModeValue('black', 'white')}
                        _hover={{
                          backgroundColor: useColorModeValue(
                            'gray.400',
                            'gray.700'
                          ),
                        }}
                      >
                        {showPassword ? 'Hide' : 'Show'}
                      </Button>
                    </InputRightElement>
                  </InputGroup>
                </FormControl>
                <Button
                  type="submit"
                  variant="solid"
                  bg={colors.primary}
                  width="full"
                  color={'white'}
                  _hover={{
                    bg: colors.primaryFaded,
                  }}
                >
                  Login
                </Button>
              </Stack>
            </form>
          </Box>
        </Stack>
      </Flex>
    </Sidenav>
  )
}
