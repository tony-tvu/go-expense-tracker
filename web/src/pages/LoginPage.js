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
  Link,
  FormControl,
  FormHelperText,
  InputRightElement,
  Text,
  useColorModeValue,
} from '@chakra-ui/react'
import { FaUserAlt, FaLock, FaCat } from 'react-icons/fa'
import { ColorModeSwitcher } from '../ColorModeSwitcher'
import { colors } from '../theme'
import { Link as RouterLink } from 'react-router-dom'
import { useNavigate } from 'react-router-dom'
import logger from '../logger'

const CFaUserAlt = chakra(FaUserAlt)
const CFaLock = chakra(FaLock)
const CFcat = chakra(FaCat)

export default function LoginPage() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [registrationEnabled, setRegistrationEnabled] = useState(false)
  const handleShowClick = () => setShowPassword(!showPassword)
  const navigate = useNavigate()

  const bgColor = useColorModeValue(colors.bgLight, colors.bgDark)
  const logoColor = useColorModeValue('black', colors.primary)

  useEffect(() => {
    Promise.all([
      fetch(`${process.env.REACT_APP_API_URL}/logged_in`, {
        method: 'GET',
        credentials: 'include',
      })
        .then((res) => {
          if (res.status === 200) {
            navigate('/')
          }
        })
        .catch((err) => {
          logger('error verifying login state', err)
        }),
      fetch(`${process.env.REACT_APP_API_URL}/registration_enabled`, {
        method: 'GET',
      })
        .then(async (res) => {
          if (!res) return
          const data = await res.json().catch((err) => logger(err))
          if (data && data.registration_enabled) {
            setRegistrationEnabled(true)
          }
        })
        .catch((err) => {
          logger('error getting registration_enabled', err)
        }),
    ])
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
      })
      .catch((e) => {
        logger('error setting access token', e)
      })
  }

  return (
    <Flex
      flexDirection="column"
      width="100wh"
      height="100vh"
      backgroundColor={bgColor}
      alignItems="center"
    >
      <Box bg="gray.800" w="100%" color="white">
        <Flex
          minH={'50px'}
          bg={'gray.800'}
          align={'center'}
          pl={'2vw'}
          pr={'2vw'}
          borderBottom={1}
          borderStyle={'solid'}
          borderColor={'gray.600'}
        >
          <Flex flex={{ base: 1 }}>
            <RouterLink to="/">
              <Text fontSize="xl" as="b" fontFamily={'heading'} color={'white'}>
                {process.env.REACT_APP_NAME}
              </Text>
            </RouterLink>
          </Flex>
          {registrationEnabled && (
            <Link
              px={2}
              py={1}
              rounded={'md'}
              _hover={{
                textDecoration: 'none',
                bg: 'gray.700',
              }}
              href={'/'}
              color={'white'}
            >
              Register
            </Link>
          )}

          <ColorModeSwitcher justifySelf="flex-end" />
        </Flex>
      </Box>
      <Stack
        flexDir="column"
        mt="25vh"
        justifyContent="center"
        alignItems="center"
      >
        <CFcat width={'10vh'} size={'100px'} color={logoColor} />
        <Box
          minW={{ base: '90%', md: '468px' }}
          backgroundColor={'whiteAlpha.800'}
        >
          <form onSubmit={handleSubmit}>
            <Stack
              spacing={4}
              p="1rem"
              backgroundColor={useColorModeValue('whiteAlpha.800', 'gray.900')}
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
                    bg={useColorModeValue('whiteAlpha.800', 'gray.900')}
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
                    bg={useColorModeValue('whiteAlpha.800', 'gray.900')}
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
                <FormHelperText textAlign="right">
                  <Link color={useColorModeValue('black', 'whiteAlpha.800')}>
                    forgot password?
                  </Link>
                </FormHelperText>
              </FormControl>
              <Button
                borderRadius={0}
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
  )
}
