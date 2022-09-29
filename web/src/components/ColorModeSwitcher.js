import React from 'react'
import {
  useColorMode,
  useColorModeValue,
  Flex,
  Text,
  Switch,
  Spacer,
} from '@chakra-ui/react'
import { FaMoon, FaSun } from 'react-icons/fa'

export const ColorModeSwitcher = (props) => {
  const { toggleColorMode, colorMode } = useColorMode()
  const SwitchIcon = useColorModeValue(FaMoon, FaSun)

  console.log(colorMode)

  return (
    <Flex
      alignItems={'center'}
      justifyContent={'space-between'}
      minH={'60px'}
      p={5}
      pl={'32px'}
      color={'#DCDCE2'}
      fontWeight="500"
    >
      <SwitchIcon />
      <Text ml={'17px'}>Theme</Text>
      <Spacer />
      <Switch
        isChecked={colorMode === 'dark' ? true : false}
        onChange={toggleColorMode}
      />
    </Flex>
  )
}
