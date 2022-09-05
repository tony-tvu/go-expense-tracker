import React, { useState } from "react"
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
} from "@chakra-ui/react"
import { FaUserAlt, FaLock, FaCat } from "react-icons/fa"
import { ColorModeSwitcher } from "../ColorModeSwitcher"
import { colors } from "../theme"
import { Link as RouterLink } from "react-router-dom"
import { APP_NAME } from "../configs"
import logger from "../logger"
import { useNavigate } from "react-router-dom"
import { useLoginStatus } from "../hooks/useLoginStatus"

const CFaUserAlt = chakra(FaUserAlt)
const CFaLock = chakra(FaLock)
const CFcat = chakra(FaCat)

export default function Login() {
  const navigate = useNavigate()
  const isLoggedIn = useLoginStatus()
  if (isLoggedIn) navigate("/")

  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [showPassword, setShowPassword] = useState(false)
  const handleShowClick = () => setShowPassword(!showPassword)

  async function handleSubmit(e) {
    e.preventDefault()
    await fetch(`${process.env.REACT_APP_API_URL}/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify({ username: username, password: password }),
    })
      .then((res) => {
        if (res.status === 200) navigate("/")
      })
      .catch((e) => {
        logger("error setting access token", e)
      })
  }

  return (
    <Flex
      flexDirection="column"
      width="100wh"
      height="100vh"
      backgroundColor={useColorModeValue(colors.bgLight, colors.bgDark)}
      alignItems="center"
    >
      <Box bg="gray.800" w="100%" color="white">
        <Flex
          minH={"50px"}
          bg={"gray.800"}
          align={"center"}
          pl={"2vw"}
          pr={"2vw"}
          borderBottom={1}
          borderStyle={"solid"}
          borderColor={"gray.600"}
        >
          <Flex flex={{ base: 1 }}>
            <RouterLink to="/">
              <Text
                fontSize="xl"
                as="b"
                fontFamily={"heading"}
                color={"whiteAlpha.800"}
              >
                {APP_NAME}
              </Text>
            </RouterLink>
          </Flex>
          <ColorModeSwitcher justifySelf="flex-end" />
        </Flex>
      </Box>

      <Stack
        flexDir="column"
        mt="25vh"
        justifyContent="center"
        alignItems="center"
      >
        <CFcat
          width={"10vh"}
          size={"100px"}
          color={useColorModeValue("black", colors.primary)}
        />

        <Box
          minW={{ base: "90%", md: "468px" }}
          backgroundColor={"whiteAlpha.800"}
        >
          <form onSubmit={handleSubmit}>
            <Stack
              spacing={4}
              p="1rem"
              backgroundColor={"whiteAlpha.800"}
              boxShadow="md"
            >
              <FormControl>
                <InputGroup>
                  <InputLeftElement
                    pointerEvents="none"
                    children={<CFaUserAlt color={"gray.500"} />}
                  />
                  <Input
                    type="username"
                    placeholder="username"
                    _placeholder={{ color: "gray.500" }}
                    borderColor={"gray.200"}
                    _hover={{
                      borderColor: "gray.500",
                    }}
                    color={"black"}
                    onChange={(event) => setUsername(event.target.value)}
                    bg={"whiteAlpha.800"}
                  />
                </InputGroup>
              </FormControl>
              <FormControl>
                <InputGroup>
                  <InputLeftElement
                    pointerEvents="none"
                    color={"gray.500"}
                    children={<CFaLock color={"gray.500"} />}
                  />
                  <Input
                    onChange={(event) => setPassword(event.target.value)}
                    type={showPassword ? "text" : "password"}
                    placeholder="password"
                    _placeholder={{ color: "gray.500" }}
                    borderColor={"gray.200"}
                    _hover={{
                      borderColor: "gray.500",
                    }}
                    color={"black"}
                    bg={"white"}
                  />
                  <InputRightElement width="4.5rem">
                    <Button
                      h="1.75rem"
                      size="sm"
                      onClick={handleShowClick}
                      backgroundColor={"gray.200"}
                      _hover={{
                        backgroundColor: "gray.300",
                      }}
                      color={"black"}
                    >
                      {showPassword ? "Hide" : "Show"}
                    </Button>
                  </InputRightElement>
                </InputGroup>
                <FormHelperText textAlign="right">
                  <Link color={"gray.600"}>forgot password?</Link>
                </FormHelperText>
              </FormControl>
              <Button
                borderRadius={0}
                type="submit"
                variant="solid"
                backgroundColor={colors.primary}
                width="full"
                color={"black"}
                _hover={{
                  bg: "pink.300",
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
