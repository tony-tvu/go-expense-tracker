import { Flex, useColorModeValue } from "@chakra-ui/react"
import React from "react"
import Navbar from "../components/Navbar"
import { useVerifyLogin } from "../hooks/useVerifyLogin"
import { colors } from "../theme"

export default function HomePage() {
  useVerifyLogin()

  return (
    <>
      <Navbar />
      <Flex
        flexDirection="column"
        width="100wh"
        height="100vh"
        backgroundColor={useColorModeValue(colors.bgLight, colors.bgDark)}
        alignItems="center"
      >
        HOMEPAGE
      </Flex>
    </>
  )
}
