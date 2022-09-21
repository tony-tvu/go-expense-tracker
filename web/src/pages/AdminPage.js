import React, { useEffect, useState } from 'react'
import {
  FormControl,
  FormLabel,
  VStack,
  Text,
  useColorModeValue,
  Container,
  Divider,
  Select,
  chakra,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  NumberIncrementStepper,
  NumberDecrementStepper,
  Tooltip,
  HStack,
  Spinner,
  Center,
  Button,
} from '@chakra-ui/react'
import logger from '../logger'
import { FaInfoCircle } from 'react-icons/fa'
import { colors } from '../theme'
const CFInfoCircle = chakra(FaInfoCircle)

export default function AdminPage() {
  const [loading, setLoading] = useState(true)
  const [accessTokenExp, setAccessTokenExp] = useState(0)
  const [refreshTokenExp, setRefreshTokenExp] = useState(0)
  const [quotaEnabled, setQuoteEnabled] = useState(true)
  const [quotaLimit, setQuotaLimit] = useState(0)
  const [tasksEnabled, setTasksEnabled] = useState(true)
  const [tasksInterval, setTasksInterval] = useState(0)
  const [registrationEnabled, setRegistrationEnabled] = useState(false)

  const stackBgColor = useColorModeValue('white', 'gray.900')
  const tooltipBg = useColorModeValue('white', 'gray.900')
  const tooltipColor = useColorModeValue('black', 'white')

  useEffect(() => {
    if (loading) {
      fetch(`${process.env.REACT_APP_API_URL}/configs`, {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      })
        .then(async (res) => {
          if (!res) return
          const data = await res.json().catch((err) => logger(err))
          setAccessTokenExp(data.access_token_exp)
          setRefreshTokenExp(data.refresh_token_exp)
          setQuoteEnabled(data.quota_enabled)
          setQuotaLimit(data.quota_limit)
          setTasksEnabled(data.tasks_enabled)
          setTasksInterval(data.tasks_interval)
          setRegistrationEnabled(data.registration_enabled)
          setLoading(false)
        })
        .catch((err) => {
          logger('error getting items', err)
        })
    }
  }, [loading])

  if (loading) {
    return (
      <Center pt={10}>
        <Spinner
          thickness="4px"
          speed="0.65s"
          emptyColor="gray.200"
          color="blue.500"
          size="xl"
        />
      </Center>
    )
  }

  async function handleUpdate() {
    await fetch(`${process.env.REACT_APP_API_URL}/configs`, {
      method: 'PUT',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        "access_token_exp": Number(accessTokenExp),
        "refresh_token_exp": Number(refreshTokenExp),
        "quota_enabled": String(quotaEnabled) === "true" ? true : false,
        "quota_limit": Number(quotaLimit),
        "tasks_enabled": String(tasksEnabled) === "true" ? true : false,
        "tasks_interval": Number(tasksInterval),
        "registration_enabled": String(registrationEnabled) === "true" ? true : false,
      }),
    })
      .then(async (res) => {
        if (!res) return
        const data = await res.json().catch((err) => logger(err))
        console.log(data)
      })
      .catch((err) => {
        logger('error getting items', err)
      })
  }

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
            value={accessTokenExp}
            min={1}
            max={2592000}
            onChange={(value) => setAccessTokenExp(value)}
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
            value={refreshTokenExp}
            min={1}
            max={2592000}
            onChange={(value) => setRefreshTokenExp(value)}
          >
            <NumberInputField />
            <NumberInputStepper>
              <NumberIncrementStepper />
              <NumberDecrementStepper />
            </NumberInputStepper>
          </NumberInput>

          <FormLabel mt={5}>
            <HStack>
              <Text>Quota</Text>
              <Tooltip
                label="Quotas can be set to control the number of times a user is able to link a new account"
                fontSize="md"
                bg={tooltipBg}
                color={tooltipColor}
                borderWidth="1px"
                boxShadow={'2xl'}
                borderRadius="lg"
                p={5}
              >
                <span>
                  <CFInfoCircle />
                </span>
              </Tooltip>
            </HStack>
          </FormLabel>

          <Select
            value={quotaEnabled}
            onChange={(event) => setQuoteEnabled(event.target.value)}
          >
            <option value={true}>Enabled</option>
            <option value={false}>Disabled</option>
          </Select>

          <FormLabel mt={5}>Quota Limit</FormLabel>
          <NumberInput
            value={quotaLimit}
            min={1}
            max={2592000}
            onChange={(value) => setQuotaLimit(value)}
          >
            <NumberInputField />
            <NumberInputStepper>
              <NumberIncrementStepper />
              <NumberDecrementStepper />
            </NumberInputStepper>
          </NumberInput>

          <FormLabel mt={5}>
            <HStack>
              <Text>Scheduled Tasks</Text>
              <Tooltip
                label="Tasks must be enabled to refresh transactions and account values"
                fontSize="md"
                bg={tooltipBg}
                color={tooltipColor}
                borderWidth="1px"
                boxShadow={'2xl'}
                borderRadius="lg"
                p={5}
              >
                <span>
                  <CFInfoCircle />
                </span>
              </Tooltip>
            </HStack>
          </FormLabel>
          <Select
            value={tasksEnabled}
            onChange={(event) => setTasksEnabled(event.target.value)}
          >
            <option value={true}>Enabled</option>
            <option value={false}>Disabled</option>
          </Select>

          <FormLabel mt={5}>Scheduled Tasks Interval</FormLabel>
          <FormLabel fontSize={'xs'} color={'gray.500'}>
            Default: 60 seconds
          </FormLabel>
          <NumberInput
            value={tasksInterval}
            min={1}
            max={2592000}
            onChange={(value) => setTasksInterval(value)}
          >
            <NumberInputField />
            <NumberInputStepper>
              <NumberIncrementStepper />
              <NumberDecrementStepper />
            </NumberInputStepper>
          </NumberInput>

          <FormLabel mt={5}>
            <HStack>
              <Text>User Registration</Text>
              <Tooltip
                label="Allow new users to create an account from the login page"
                fontSize="md"
                bg={tooltipBg}
                color={tooltipColor}
                borderWidth="1px"
                boxShadow={'2xl'}
                borderRadius="lg"
                p={5}
              >
                <span>
                  <CFInfoCircle />
                </span>
              </Tooltip>
            </HStack>
          </FormLabel>
          <Select
            value={registrationEnabled}
            onChange={(event) => setRegistrationEnabled(event.target.value)}
          >
            <option value={true}>Enabled</option>
            <option value={false}>Disabled</option>
          </Select>
          <Button
            mt={5}
            onClick={handleUpdate}
            type="submit"
            variant="solid"
            bg={colors.primary}
            color={'white'}
            _hover={{
              bg: colors.primaryFaded,
            }}
          >
            Update
          </Button>
        </FormControl>
      </Container>
    </VStack>
  )
}
