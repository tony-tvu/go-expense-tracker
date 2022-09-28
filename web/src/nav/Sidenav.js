import React from 'react'
import {
  IconButton,
  Box,
  CloseButton,
  Flex,
  Icon,
  useColorModeValue,
  Link,
  Drawer,
  DrawerContent,
  Text,
  useDisclosure,
  chakra,
  Divider,
  Spacer,
} from '@chakra-ui/react'
import {
  FiHome,
  FiTrendingUp,
  FiCompass,
  FiStar,
  FiSettings,
  FiMenu,
} from 'react-icons/fi'
import { FaCat, FaThLarge, FaPen } from 'react-icons/fa'
import { ColorModeSwitcher } from '../ColorModeSwitcher'
import { colors } from '../theme'
import { appName } from '../commons'
import { Link as RouterLink } from 'react-router-dom'

const CFcat = chakra(FaCat)
const CFThLarge = chakra(FaThLarge)
const CFPen = chakra(FaPen)
const fontColor = 'gray.100'

export default function Sidenav({
  children,
  isLoggedIn,
  registrationEnabled,
  isAdmin,
}) {
  const { isOpen, onOpen, onClose } = useDisclosure()
  return (
    <Box minH="100vh" bg={useColorModeValue('gray.100', '#1E1E1E')}>
      <SidebarContent
        onClose={() => onClose}
        isLoggedIn={isLoggedIn}
        registrationEnabled={registrationEnabled}
        isAdmin={isAdmin}
        display={{ base: 'none', md: 'block' }}
      />
      <Drawer
        autoFocus={false}
        isOpen={isOpen}
        placement="left"
        onClose={onClose}
        returnFocusOnClose={false}
        onOverlayClick={onClose}
        size="full"
      >
        <DrawerContent>
          <SidebarContent onClose={onClose} />
        </DrawerContent>
      </Drawer>
      {/* mobilenav */}
      <MobileNav display={{ base: 'flex', md: 'none' }} onOpen={onOpen} />
      <Box ml={{ base: 0, md: 60 }} p="4">
        {children}
      </Box>
    </Box>
  )
}

const SidebarContent = ({
  onClose,
  isLoggedIn,
  registrationEnabled,
  ...rest
}) => {
  return (
    <Box
      bg={'#252526'}
      borderRight="1px"
      borderColor={'#1E1E1E'}
      w={{ base: 'full', md: 60 }}
      pos="fixed"
      color={fontColor}
      h="full"
      {...rest}
    >
      <Flex h="20" alignItems="center" mx="8" justifyContent="space-between">
        <RouterLink to="/">
          <CFcat size={'25px'} color={colors.primary} />
        </RouterLink>

        <RouterLink to="/">
          <Text
            ml={3}
            fontSize="2xl"
            fontFamily="monospace"
            fontWeight="bold"
            color={fontColor}
          >
            {appName()}
          </Text>
        </RouterLink>
        <Spacer />
        <CloseButton display={{ base: 'flex', md: 'none' }} onClick={onClose} />
      </Flex>

      {isLoggedIn && (
        <>
          <NavItem icon={CFThLarge}>Dashboard</NavItem>
          <Divider mt={5} borderColor={'#464646'} />
        </>
      )}

      {!isLoggedIn && registrationEnabled && (
        <>
          <RouterLink to="/register">
            <NavItem color={fontColor} icon={CFPen}>
              Register
            </NavItem>
          </RouterLink>
          <Divider mt={5} borderColor={'#464646'} />
        </>
      )}

      <Text pt={3} pl={3} color={'#79797C'} fontWeight={'600'} fontSize={'sm'}>
        PREFERENCES
      </Text>
      <ColorModeSwitcher justifySelf="flex-end" color={fontColor} />
    </Box>
  )
}

const NavItem = ({ icon, children, ...rest }) => {
  return (
    <Link
      href="#"
      style={{ textDecoration: 'none' }}
      _focus={{ boxShadow: 'none' }}
    >
      <Flex
        align="center"
        p="4"
        mx="4"
        borderRadius="lg"
        role="group"
        cursor="pointer"
        _hover={{
          bg: colors.primary,
          color: 'white',
        }}
        {...rest}
      >
        {icon && (
          <Icon
            mr="4"
            fontSize="16"
            _groupHover={{
              color: 'white',
            }}
            as={icon}
          />
        )}
        {children}
      </Flex>
    </Link>
  )
}

const MobileNav = ({ onOpen, ...rest }) => {
  return (
    <Flex
      ml={{ base: 0, md: 60 }}
      px={{ base: 4, md: 24 }}
      height="20"
      alignItems="center"
      bg={useColorModeValue('white', 'gray.900')}
      borderBottomWidth="1px"
      borderBottomColor={useColorModeValue('gray.200', 'gray.700')}
      justifyContent="flex-start"
      {...rest}
    >
      <IconButton
        variant="outline"
        onClick={onOpen}
        aria-label="open menu"
        icon={<FiMenu />}
      />

      <Text fontSize="2xl" ml="8" fontFamily="monospace" fontWeight="bold">
        Logo
      </Text>
    </Flex>
  )
}
