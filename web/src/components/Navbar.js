import {
  Box,
  Flex,
  Avatar,
  HStack,
  Link,
  IconButton,
  Button,
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
  MenuDivider,
  useDisclosure,
  Text,
  useColorModeValue,
  Stack,
  chakra,
} from '@chakra-ui/react'
import { Link as RouterLink } from 'react-router-dom'
import { HamburgerIcon, CloseIcon } from '@chakra-ui/icons'
import { useNavigate } from 'react-router-dom'
import logger from '../logger'
import { ColorModeSwitcher } from '../ColorModeSwitcher'
import { FaCat } from 'react-icons/fa'
import { colors } from '../theme'
import { useLocation } from 'react-router-dom'
import { useVerifyAdmin } from '../hooks/useVerifyAdmin'
import { useVerifyLogin } from '../hooks/useVerifyLogin'
import { useRegistration } from '../hooks/useRegistration'

const CFcat = chakra(FaCat)

export default function Navbar() {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const navigate = useNavigate()
  const location = useLocation()
  const linkBgColor = useColorModeValue('gray.200', 'gray.700')

  const isAdmin = useVerifyAdmin()
  const isLoggedIn = useVerifyLogin()
  const registrationEnabled = useRegistration()

  // const checkIsLoggedIn = useCallback(async () => {
  //   await fetch(`${process.env.REACT_APP_API_URL}/logged_in`, {
  //     method: 'GET',
  //     credentials: 'include',
  //   })
  //     .then((res) => {
  //       if (res.status === 200) {
  //         setIsLoggedIn(true)
  //       }
  //       if (
  //         res.status === 200 &&
  //         (location.pathname === '/login/' || location.pathname === '/login')
  //       ) {
  //         navigate('/')
  //       }
  //       if (res.status !== 200) {
  //         navigate('/login')
  //       }
  //     })
  //     .catch((err) => {
  //       logger('error verifying login state', err)
  //     })
  // }, [navigate, location.pathname])

  // const checkIsRegistrationEnabled = useCallback(async () => {
  //   await fetch(`${process.env.REACT_APP_API_URL}/registration_enabled`, {
  //     method: 'GET',
  //     credentials: 'include',
  //   })
  //     .then(async (res) => {
  //       if (!res) return
  //       const data = await res.json().catch((err) => logger(err))
  //       if (data && data.registration_enabled) {
  //         setRegistrationEnabled(true)
  //       }
  //     })
  //     .catch((err) => {
  //       logger('error getting registration_enabled', err)
  //     })
  // }, [])

  // useEffect(() => {
  //   Promise.all([
  //     checkIsLoggedIn(),
  //     checkIsRegistrationEnabled(),
  //   ])
  // }, [checkIsLoggedIn, checkIsRegistrationEnabled])

  function logout() {
    fetch(`${process.env.REACT_APP_API_URL}/logout`, {
      method: 'POST',
      credentials: 'include',
    })
      .then((res) => {
        if (res.status === 200) {
          navigate('/login')
        }
      })
      .catch((err) => {
        logger('error logging out', err)
      })
  }

  return (
    <>
      <Box>
        <Flex
          bg={'gray.800'}
          h={'50px'}
          pl={'2vw'}
          pr={'2vw'}
          alignItems={'center'}
          justifyContent={'space-between'}
          borderBottom={1}
          borderStyle={'solid'}
          borderColor={'gray.600'}
        >
          {isLoggedIn && (
            <IconButton
              size={'md'}
              icon={isOpen ? <CloseIcon /> : <HamburgerIcon />}
              aria-label={'Open Menu'}
              display={{ md: 'none' }}
              onClick={isOpen ? onClose : onOpen}
              bg={'gray.700'}
              color={'white'}
              _hover={{
                borderColor: 'gray.500',
              }}
            />
          )}

          <HStack spacing={8} alignItems={'center'}>
            {isLoggedIn ? (
              <RouterLink to="/">
                <CFcat size={'30px'} color={colors.primary} />
              </RouterLink>
            ) : (
              <RouterLink to="/login">
                <Text
                  fontSize="xl"
                  as="b"
                  fontFamily={'heading'}
                  color={'white'}
                >
                  {process.env.REACT_APP_NAME}
                </Text>
              </RouterLink>
            )}

            {isLoggedIn && (
              <HStack
                as={'nav'}
                spacing={4}
                display={{ base: 'none', md: 'flex' }}
              >
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
                  Overview
                </Link>
              </HStack>
            )}
          </HStack>

          <Flex alignItems={'center'}>
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
            <ColorModeSwitcher justifySelf="flex-end" color="white" />

            {isLoggedIn && (
              <Menu>
                <MenuButton
                  ml={'20px'}
                  as={Button}
                  rounded={'full'}
                  variant={'link'}
                  cursor={'pointer'}
                  minW={0}
                >
                  <Avatar size={'sm'} bg={colors.primary} />
                </MenuButton>
                <MenuList>
                  <MenuItem onClick={() => navigate('/accounts')}>
                    Accounts
                  </MenuItem>
                  {isAdmin && (
                    <MenuItem onClick={() => navigate('/admin')}>
                      Admin
                    </MenuItem>
                  )}
                  <MenuDivider />
                  <MenuItem onClick={logout}>Logout</MenuItem>
                </MenuList>
              </Menu>
            )}
          </Flex>
        </Flex>

        {isOpen ? (
          <Box pb={4} display={{ md: 'none' }}>
            <Stack as={'nav'} spacing={4}>
              <Link
                px={2}
                py={1}
                rounded={'md'}
                _hover={{
                  textDecoration: 'none',
                  bg: linkBgColor,
                }}
                href={'/'}
              >
                Overview
              </Link>
            </Stack>
          </Box>
        ) : null}
      </Box>
    </>
  )
}
