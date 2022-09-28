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
  useColorModeValue,
  Stack,
  chakra,
  Text,
  Spacer,
} from '@chakra-ui/react'
import { Link as RouterLink } from 'react-router-dom'
import { HamburgerIcon, CloseIcon } from '@chakra-ui/icons'
import logger from '../logger'
import { ColorModeSwitcher } from '../ColorModeSwitcher'
import { FaCat } from 'react-icons/fa'
import { colors } from '../theme'
import { useNavigate } from 'react-router-dom'

const CFcat = chakra(FaCat)

export default function Navbar({ isLoggedIn, registrationEnabled, isAdmin }) {
  const navigate = useNavigate()
  const { isOpen, onOpen, onClose } = useDisclosure()
  const linkBgColor = useColorModeValue('gray.200', 'gray.700')
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

          <HStack spacing={8}>
            {isLoggedIn ? (
              <RouterLink to="/">
                <CFcat size={'30px'} ml={3} color={colors.primary} />
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
                <Link
                  px={2}
                  py={1}
                  rounded={'md'}
                  _hover={{
                    textDecoration: 'none',
                    bg: 'gray.700',
                  }}
                  href={'/expenses'}
                  color={'white'}
                >
                  Expenses
                </Link>
              </HStack>
            )}
          </HStack>
          <Spacer />
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
                  ml={5}
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
              <Link
                px={2}
                py={1}
                rounded={'md'}
                _hover={{
                  textDecoration: 'none',
                  bg: linkBgColor,
                }}
                href={'/expenses'}
              >
                Expenses
              </Link>
            </Stack>
          </Box>
        ) : null}
      </Box>
    </>
  )
}
