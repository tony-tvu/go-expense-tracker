import React, { useEffect, useState } from 'react'

import {
  Text,
  Spacer,
  VStack,
  Center,
  HStack,
  Spinner,
  Container,
  useColorModeValue,
} from '@chakra-ui/react'
import CreateRuleBtn from '../components/CreateRuleBtn'
import { IoIosWarning } from 'react-icons/io'
import logger from '../logger'
import DeleteRuleBtn from '../components/DeleteRuleBtn'

export default function Rules() {
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(true)

  const stackBgColor = useColorModeValue('white', 'gray.900')
  const bgColor = useColorModeValue('white', '#252526')

  function onSuccess() {
    setData([])
    setLoading(true)
  }

  useEffect(() => {
    document.title = 'Rules'
    if (loading) {
      fetch(`${process.env.REACT_APP_API_URL}/rules`, {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      })
        .then(async (res) => {
          if (!res) return
          const resData = await res.json().catch((err) => logger(err))
          if (res.status === 200 && resData.rules) {
            setData(resData.rules)
          }
          setLoading(false)
        })
        .catch((err) => {
          logger('error getting rules', err)
        })
    }
  }, [loading])

  function renderRules() {
    if (data.length === 0 && !loading) {
      return null
    }

    return data.map((rule) => {
      return (
        <HStack
          key={rule.id}
          width={'100%'}
          borderWidth="1px"
          borderRadius="lg"
          bg={stackBgColor}
          boxShadow={'2xl'}
          p={3}
          mb={5}
        >
          <Text fontSize="xl" as="b">
            "{rule.substring}" ={' '}
            {rule.category.charAt(0).toUpperCase() + rule.category.slice(1)}
          </Text>
          <Spacer />

          <DeleteRuleBtn rule={rule} onSuccess={onSuccess} />
        </HStack>
      )
    })
  }

  return (
    <VStack>
      <Container maxW="container.md" centerContent bg={bgColor} mb={5}>
        <HStack alignItems="end" width="100%" mt={5} mb={5}>
          <Text fontSize="3xl" as="b" pl={1}>
            Rules
          </Text>
          <Spacer />
          <CreateRuleBtn onSuccess={onSuccess} />
        </HStack>
        <VStack alignItems="start" width="100%" mb={10}>
          <Text fontSize="md" pl={1} mb={5}>
            Rules provide a quick and easy way to categorize transactions. When
            creating a rule, you can specify a substring and the category that
            it belongs to. For instance, a substring of "Bananas Market" and
            category of "groceries" will make all transactions with that
            substring change to the "groceries" category. More specifically, a
            transaction with the name "Super Fruits and Bananas Market" will be
            categorized as "groceries".
          </Text>
          <HStack>
            <IoIosWarning size={'60px'} color={'red'} />
            <Text>
              Be careful! If you specify a short and generic substring, such as
              "a", then ALL transactions with the letter "a" will be changed to
              the category you specified for that rule.
            </Text>
          </HStack>
        </VStack>
      </Container>

      <Container maxW="container.md" centerContent minH={'100px'} p={0}>
        {loading ? (
          <Center pt={10}>
            <Spinner
              thickness="4px"
              speed="0.65s"
              emptyColor="gray.200"
              color="blue.500"
              size="xl"
            />
          </Center>
        ) : (
          renderRules()
        )}
      </Container>
    </VStack>
  )
}
