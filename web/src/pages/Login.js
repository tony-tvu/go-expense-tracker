import React, { useState } from "react"
import {
  Flex,
  Heading,
  Input,
  Button,
  InputGroup,
  Stack,
  InputLeftElement,
  chakra,
  Box,
  Link,
  Avatar,
  FormControl,
  FormHelperText,
  InputRightElement,
  Text,
  useBreakpointValue,
  useColorModeValue,
} from "@chakra-ui/react"
import { FaUserAlt, FaLock, FaCat } from "react-icons/fa"
import { ColorModeSwitcher } from "../ColorModeSwitcher"
import { colors } from "../theme"
import { Link as RouterLink } from "react-router-dom"
import { APP_NAME } from "../configs"

const CFaUserAlt = chakra(FaUserAlt)
const CFaLock = chakra(FaLock)
const CFcat = chakra(FaCat)

export default function Login() {
  const [showPassword, setShowPassword] = useState(false)

  const handleShowClick = () => setShowPassword(!showPassword)

  return (
    <Flex
      flexDirection="column"
      width="100wh"
      height="100vh"
      backgroundColor={useColorModeValue(colors.bgLight, colors.bgDark)}
      alignItems="center"
    >
      <Box bg={colors.gray.dark} w="100%" color="white">
        <Flex
          minH={"50px"}
          bg={colors.gray.dark}
          color={colors.white}
          align={"center"}
          pl={"2vw"}
          pr={"2vw"}
          borderBottom={1}
          borderStyle={"solid"}
          borderColor={colors.gray.medium}
        >
          <Flex flex={{ base: 1 }}>
            <RouterLink to="/">
              <Text
                fontSize="xl"
                as="b"
                fontFamily={"heading"}
                color={colors.white.extra}
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
          color={useColorModeValue(colors.black, colors.primary)}
        />

        <Box
          minW={{ base: "90%", md: "468px" }}
          backgroundColor={colors.white.light}
        >
          <form>
            <Stack
              spacing={4}
              p="1rem"
              backgroundColor={colors.white.light}
              boxShadow="md"
            >
              <FormControl>
                <InputGroup>
                  <InputLeftElement
                    pointerEvents="none"
                    children={<CFaUserAlt color={colors.gray.light} />}
                  />
                  <Input
                    type="email"
                    placeholder="email address"
                    _placeholder={{ color: colors.gray.light }}
                    borderColor={colors.gray.extraLight}
                    _hover={{
                      borderColor: colors.gray.light,
                    }}
                    color={colors.black}
                  />
                </InputGroup>
              </FormControl>
              <FormControl>
                <InputGroup>
                  <InputLeftElement
                    pointerEvents="none"
                    color={colors.gray.light}
                    children={<CFaLock color={colors.gray.light} />}
                  />
                  <Input
                    type={showPassword ? "text" : "password"}
                    placeholder="password"
                    _placeholder={{ color: colors.gray.light }}
                    borderColor={colors.gray.extraLight}
                    _hover={{
                      borderColor: colors.gray.light,
                    }}
                    color={colors.black}
                  />
                  <InputRightElement width="4.5rem">
                    <Button
                      h="1.75rem"
                      size="sm"
                      onClick={handleShowClick}
                      backgroundColor={colors.gray.extraLight}
                      _hover={{
                        backgroundColor: "gray.300",
                      }}
                      color={colors.black}
                    >
                      {showPassword ? "Hide" : "Show"}
                    </Button>
                  </InputRightElement>
                </InputGroup>
                <FormHelperText textAlign="right">
                  <Link color={colors.gray.medium}>forgot password?</Link>
                </FormHelperText>
              </FormControl>
              <Button
                borderRadius={0}
                type="submit"
                variant="solid"
                backgroundColor={colors.primary}
                width="full"
                color={colors.black}
                _hover={{
                  bg: colors.pink.medium,
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
