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
} from "@chakra-ui/react"
import { Link as RouterLink } from "react-router-dom"
import { HamburgerIcon, CloseIcon } from "@chakra-ui/icons"
import { MoonIcon, SunIcon } from "@chakra-ui/icons"
import { APP_NAME } from "../configs"
import { useNavigate } from "react-router-dom"
import logger from "../logger"
import { gql, useMutation } from "@apollo/client"

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
      <Box bg={useColorModeValue("gray.100", "gray.900")} px={4}>
        <Flex h={16} alignItems={"center"} justifyContent={"space-between"}>
          <IconButton
            size={"md"}
            icon={isOpen ? <CloseIcon /> : <HamburgerIcon />}
            aria-label={"Open Menu"}
            display={{ md: "none" }}
            onClick={isOpen ? onClose : onOpen}
          />
          <HStack spacing={8} alignItems={"center"}>
            <RouterLink to="/">
              <Box>{APP_NAME}</Box>
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
                  bg: linkBgColor,
                }}
                href={"/"}
              >
                Navlink 1
              </Link>
            </HStack>
          </HStack>

          <Flex alignItems={"center"}>
            <Button onClick={toggleColorMode}>
              {colorMode === "light" ? <MoonIcon /> : <SunIcon />}
            </Button>
            <Menu>
              <MenuButton
                ml={"20px"}
                as={Button}
                rounded={"full"}
                variant={"link"}
                cursor={"pointer"}
                minW={0}
              >
                <Avatar size={"sm"} bg="teal.500" />
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
                Navlink 1
              </Link>
            </Stack>
          </Box>
        ) : null}
      </Box>
    </>
  )
}
