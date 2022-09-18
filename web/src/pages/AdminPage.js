import { Flex, useColorModeValue } from "@chakra-ui/react"
import React from "react"
import Navbar from "../components/Navbar"
import { useVerifyAdmin } from "../hooks/useVerifyAdmin"
import { colors } from "../theme"

export default function AdminPage() {
  useVerifyAdmin()

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
        ADMIN PAGE
      </Flex>
    </>
  )
}
