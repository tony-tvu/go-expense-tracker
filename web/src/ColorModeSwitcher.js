import React from "react"
import { useColorMode, useColorModeValue, IconButton } from "@chakra-ui/react"
import { FaMoon, FaSun } from "react-icons/fa"
import { colors } from "./theme"

export const ColorModeSwitcher = (props) => {
  const { toggleColorMode } = useColorMode()
  const text = useColorModeValue("dark", "light")
  const SwitchIcon = useColorModeValue(FaMoon, FaSun)

  return (
    <IconButton
      _hover={{
        backgroundColor: useColorModeValue(
          colors.gray.light,
          colors.gray.light
        ),
        color: colors.white.light,
      }}
      size="md"
      fontSize="lg"
      aria-label={`Switch to ${text} mode`}
      variant="ghost"
      color="current"
      marginLeft="2"
      onClick={toggleColorMode}
      icon={<SwitchIcon />}
      {...props}
    />
  )
}
