import { extendTheme } from "@chakra-ui/react"

export const colors = {
  primary: "purple.600",
  primaryFaded: "purple.500",
  bgLight: "gray.300",
  bgDark: "gray.800",
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
