import React from 'react'
import {
  IconButton,
  Box,
  CloseButton,
  Flex,
  Icon,
  useColorModeValue,
  Drawer,
  DrawerContent,
  Text,
  useDisclosure,
  chakra,
  Divider,
  Spacer,
} from '@chakra-ui/react'
import { FiSettings, FiMenu } from 'react-icons/fi'
import {
  FaCat,
  FaThLarge,
  FaPen,
  FaMugHot,
  FaMoneyBill,
  FaListUl,
} from 'react-icons/fa'
import { ColorModeSwitcher } from '../components/ColorModeSwitcher'
import { colors } from '../theme'
import { appName } from '../commons'
import { Link as RouterLink } from 'react-router-dom'

const CFcat = chakra(FaCat)
const textColor = '#DCDCE2'
const hoverBgColor = '#303031'
const navBgColor = '#252526'

export default function Sidenav({
  children,
  isLoggedIn,
  registrationEnabled,
  isAdmin,
  current,
}) {
  const { isOpen, onOpen, onClose } = useDisclosure()
  return (
    <Box minH="100vh" bg={useColorModeValue('gray.100', '#1E1E1E')}>
      <SidebarContent
        onClose={() => onClose}
        isLoggedIn={isLoggedIn}
        registrationEnabled={registrationEnabled}
        isAdmin={isAdmin}
        current={current}
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
          <SidebarContent
            isLoggedIn={isLoggedIn}
            registrationEnabled={registrationEnabled}
            isAdmin={isAdmin}
            current={current}
            onClose={onClose}
          />
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
  isAdmin,
  current,
  ...rest
}) => {
  return (
    <Box
      bg={navBgColor}
      borderRight="1px"
      borderColor={'#1E1E1E'}
      w={{ base: 'full', md: 60 }}
      pos="fixed"
      color={textColor}
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
            color={textColor}
          >
            {appName()}
          </Text>
        </RouterLink>
        <Spacer />
        <CloseButton display={{ base: 'flex', md: 'none' }} onClick={onClose} />
      </Flex>

      {isLoggedIn && (
        <>
          <Text p={3} color={'#79797C'} fontWeight={'600'} fontSize={'sm'}>
            MANAGE
          </Text>
          <NavItem
            to="/"
            icon={FaThLarge}
            bgColor={current === 'dashboard' ? hoverBgColor : navBgColor}
            iconColor={current === 'dashboard' ? colors.primary : textColor}
          >
            Dashboard
          </NavItem>
          <NavItem
            to="/expenses"
            icon={FaMoneyBill}
            bgColor={current === 'expenses' ? hoverBgColor : navBgColor}
            iconColor={current === 'expenses' ? colors.primary : textColor}
          >
            Expenses
          </NavItem>

          <Divider mt={5} borderColor={'#464646'} />
        </>
      )}

      {!isLoggedIn && registrationEnabled && (
        <>
          <NavItem to="/register" color={textColor} icon={FaPen}>
            Register
          </NavItem>
          <Divider mt={5} borderColor={'#464646'} />
        </>
      )}

      <Text pl={3} pt={3} color={'#79797C'} fontWeight={'600'} fontSize={'sm'}>
        PREFERENCES
      </Text>
      <ColorModeSwitcher justifySelf="flex-end" color={textColor} />

      {isLoggedIn && (
        <>
          <NavItem
            to="/accounts"
            icon={FaListUl}
            bgColor={current === 'linked_accounts' ? hoverBgColor : navBgColor}
            iconColor={
              current === 'linked_accounts' ? colors.primary : textColor
            }
          >
            Accounts
          </NavItem>
          {isAdmin && (
            <NavItem
              to="/admin"
              icon={FaMugHot}
              bgColor={current === 'admin' ? hoverBgColor : navBgColor}
              iconColor={current === 'admin' ? colors.primary : textColor}
            >
              Admin
            </NavItem>
          )}
          <NavItem
            to="/settings"
            icon={FiSettings}
            bgColor={current === 'settings' ? hoverBgColor : navBgColor}
            iconColor={current === 'settings' ? colors.primary : textColor}
          >
            Settings
          </NavItem>
        </>
      )}
    </Box>
  )
}

const NavItem = ({ icon, children, bgColor, iconColor, to, ...rest }) => {
  return (
    <RouterLink
      to={to}
      style={{ textDecoration: 'none' }}
      _focus={{ boxShadow: 'none' }}
      color={textColor}
      fontWeight="500"
      _hover={{
        color: textColor,
      }}
    >
      <Flex
        align="center"
        p="4"
        mx="4"
        mb={1}
        borderRadius="4px"
        role="group"
        cursor="pointer"
        bg={bgColor}
        _hover={{
          bg: hoverBgColor,
          color: textColor,
        }}
        {...rest}
      >
        {icon && (
          <Icon
            mr="4"
            fontSize="16"
            color={iconColor}
            _groupHover={{
              color: colors.primary,
            }}
            as={icon}
          />
        )}
        {children}
      </Flex>
    </RouterLink>
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
        mr={5}
      />

      <RouterLink to="/">
        <CFcat size={'25px'} color={colors.primary} />
      </RouterLink>

      <RouterLink to="/">
        <Text
          ml={3}
          fontSize="2xl"
          fontFamily="monospace"
          fontWeight="bold"
          color={useColorModeValue('black', textColor)}
        >
          {appName()}
        </Text>
      </RouterLink>
    </Flex>
  )
}
