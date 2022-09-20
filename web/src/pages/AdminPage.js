import React, { useCallback, useEffect } from 'react'
import {
  FormControl,
  FormHelperText,
  FormLabel,
  Input,
  VStack,
  Text,
  useColorModeValue,
  Container,
  Divider,
  Switch,
  Stack,
  Select,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  NumberIncrementStepper,
  NumberDecrementStepper,
} from '@chakra-ui/react'
import logger from '../logger'

export default function AdminPage() {
  const stackBgColor = useColorModeValue('white', 'gray.900')

  const getConfigs = useCallback(async () => {
    await fetch(`${process.env.REACT_APP_API_URL}/configs`, {
      method: 'GET',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
    })
      .then(async (res) => {
        if (!res) return
        const data = await res.json().catch((err) => logger(err))
        console.log(data)
      })
      .catch((err) => {
        logger('error getting items', err)
      })
  }, [])

  useEffect(() => {
    getConfigs()
  }, [getConfigs])

  return (
    <VStack>
      <Container maxW="container.md" mt={3}>
        <FormControl bg={stackBgColor} p={5}>
          <FormLabel fontSize="xl">App configuration</FormLabel>
          <Divider mb={5} />
          <FormLabel mt={5}>Access Token Expiration</FormLabel>
          <FormLabel fontSize={'xs'} color={'gray.500'}>
            Default: 900 seconds
          </FormLabel>
          <NumberInput
            defaultValue={15}
            min={1}
            max={2592000}
            onChange={(value) => console.log(value)}
          >
            <NumberInputField />
            <NumberInputStepper>
              <NumberIncrementStepper />
              <NumberDecrementStepper />
            </NumberInputStepper>
          </NumberInput>

          <FormLabel mt={5}>Refresh Token Expiration</FormLabel>
          <FormLabel fontSize={'xs'} color={'gray.500'}>
            Default: 3600 seconds
          </FormLabel>
          <NumberInput
            defaultValue={15}
            min={1}
            max={2592000}
            onChange={(value) => console.log(value)}
          >
            <NumberInputField />
            <NumberInputStepper>
              <NumberIncrementStepper />
              <NumberDecrementStepper />
            </NumberInputStepper>
          </NumberInput>


          <FormLabel mt={5}>Quota</FormLabel>
          <Select
            defaultValue={false}
            onChange={(event) => console.log(event.target.value)}
          >
            <option value={true}>Enabled</option>
            <option value={false}>Disabled</option>
          </Select>


          <FormLabel mt={5}>Quota Limit</FormLabel>
          <NumberInput
            defaultValue={15}
            min={1}
            max={2592000}
            onChange={(value) => console.log(value)}
          >
            <NumberInputField />
            <NumberInputStepper>
              <NumberIncrementStepper />
              <NumberDecrementStepper />
            </NumberInputStepper>
          </NumberInput>





          <FormLabel mt={5}>Scheduled Tasks</FormLabel>
          <Select
            defaultValue={false}
            onChange={(event) => console.log(event.target.value)}
          >
            <option value={true}>Enabled</option>
            <option value={false}>Disabled</option>
          </Select>

          <FormLabel mt={5}>Scheduled Tasks Interval</FormLabel>
          <FormLabel fontSize={'xs'} color={'gray.500'}>
            Default: 60 seconds
          </FormLabel>
          <NumberInput
            isDisabled={true}
            defaultValue={15}
            min={1}
            max={2592000}
            onChange={(value) => console.log(value)}
          >
            <NumberInputField />
            <NumberInputStepper>
              <NumberIncrementStepper />
              <NumberDecrementStepper />
            </NumberInputStepper>
          </NumberInput>

          <FormLabel mt={5}>User Registration</FormLabel>
          <Select
            defaultValue={false}
            onChange={(event) => console.log(event.target.value)}
          >
            <option value={true}>Enabled</option>
            <option value={false}>Disabled</option>
          </Select>
        </FormControl>
      </Container>
    </VStack>
  )
}
