import { extendTheme } from "@chakra-ui/react"

export const colors = {
  primary: "pink.400",
  bgLight: "gray.300",
  bgDark: "gray.800",
  black: "black",
  pink: {
    light: "pink.50",
    medium: "pink.300",
    dark: "pink.400",
  },
  gray: {
    dark: "gray.800",
    medium: "gray.600",
    light: "gray.500",
    extraLight: "gray.200",
  },
  white: {
    light: "whiteAlpha.800",
    extra: "white",
  },
}

export const extendedTheme = extendTheme({
  components: {
    Input: {
      baseStyle: {
        field: {
          _autofill: {
            textFillColor: "#000000",
            boxShadow: "0 0 0px 1000px #ffffff inset",
            transition: "background-color 5000s ease-in-out 0s",
          },
        },
      },
    },
  },
})
