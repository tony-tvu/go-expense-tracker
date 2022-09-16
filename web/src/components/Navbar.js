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
  useColorMode,
  Stack,
  chakra,
} from "@chakra-ui/react"
import { Link as RouterLink } from "react-router-dom"
import { HamburgerIcon, CloseIcon } from "@chakra-ui/icons"
import { useNavigate } from "react-router-dom"
import logger from "../logger"
import { gql, useMutation } from "@apollo/client"
import { ColorModeSwitcher } from "../ColorModeSwitcher"
import { FaCat } from "react-icons/fa"
import { colors } from "../theme"

const CFcat = chakra(FaCat)

const mutation = gql`
  mutation {
    logout
  }
`

export default function Navbar() {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const { colorMode, toggleColorMode } = useColorMode()
  const navigate = useNavigate()
  const linkBgColor = useColorModeValue("gray.200", "gray.700")

  const [logout] = useMutation(mutation)

  function handleLogout() {
    logout({})
      .then((res) => {
        if (!res.errors) {
          navigate("/login")
        }
      })
      .catch((err) => {
        logger(err)
      })
  }

  return (
    <>
      <Box>
        <Flex
          bg={"gray.800"}
          h={"50px"}
          pl={"2vw"}
          pr={"2vw"}
          alignItems={"center"}
          justifyContent={"space-between"}
          borderBottom={1}
          borderStyle={"solid"}
          borderColor={"gray.600"}
        >
          <IconButton
            size={"md"}
            icon={isOpen ? <CloseIcon /> : <HamburgerIcon />}
            aria-label={"Open Menu"}
            display={{ md: "none" }}
            onClick={isOpen ? onClose : onOpen}
            bg={"gray.700"}
            color={"white"}
            _hover={{
              borderColor: "gray.500",
            }}
          />
          <HStack spacing={8} alignItems={"center"}>
            <RouterLink to="/">
              <CFcat size={"30px"} color={colors.primary} />
            </RouterLink>
            <HStack
              as={"nav"}
              spacing={4}
              display={{ base: "none", md: "flex" }}
            >
              <Link
                px={2}
                py={1}
                rounded={"md"}
                _hover={{
                  textDecoration: "none",
                  bg: "gray.700",
                }}
                href={"/"}
                color={"white"}
              >
                Overview
              </Link>
            </HStack>
          </HStack>

          <Flex alignItems={"center"}>
            <ColorModeSwitcher justifySelf="flex-end" color="white" />
            <Menu>
              <MenuButton
                ml={"20px"}
                as={Button}
                rounded={"full"}
                variant={"link"}
                cursor={"pointer"}
                minW={0}
              >
                <Avatar size={"sm"} bg={colors.primary} />
              </MenuButton>
              <MenuList>
                <MenuItem onClick={() => navigate("/accounts")}>
                  Accounts
                </MenuItem>
                <MenuItem>Settings</MenuItem>
                <MenuDivider />
                <MenuItem onClick={() => handleLogout()}>Logout</MenuItem>
              </MenuList>
            </Menu>
          </Flex>
        </Flex>

        {isOpen ? (
          <Box pb={4} display={{ md: "none" }}>
            <Stack as={"nav"} spacing={4}>
              <Link
                px={2}
                py={1}
                rounded={"md"}
                _hover={{
                  textDecoration: "none",
                  bg: linkBgColor,
                }}
                href={"/"}
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
