import React, { useState } from 'react'
import {
  Button,
  useDisclosure,
  AlertDialog,
  AlertDialogOverlay,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogBody,
  AlertDialogFooter,
  FormControl,
  FormLabel,
  Input,
  Select,
  Spinner,
  useToast,
  NumberInputField,
  NumberInput,
  IconButton,
  useColorModeValue,
} from '@chakra-ui/react'
import logger from '../logger'
import { useNavigate } from 'react-router-dom'
import { colors } from '../theme'
import DatePicker from './Datepicker'
import { FaPencilAlt } from 'react-icons/fa'

export default function EditTransactionBtn({
  onSuccess,
  forceRefresh,
  transaction,
}) {
  const [loading, setLoading] = useState(false)
  const [date, setDate] = useState(new Date(transaction.date))
  const [name, setName] = useState(transaction.name)
  const [category, setCategory] = useState(transaction.category)
  const [amount, setAmount] = useState(transaction.amount)
  const { isOpen, onOpen, onClose } = useDisclosure()
  const cancelRef = React.useRef()
  const navigate = useNavigate()
  const toast = useToast()

  const bgColor = useColorModeValue('white', '#252526')

  async function updateTransaction() {
    if (date > new Date()) {
      toast({
        title: 'Invalid Date',
        description: 'Cannot be in the future',
        status: 'error',
        position: 'top-right',
        duration: 5000,
        isClosable: true,
      })
      setLoading(false)
      return
    } else if (!name) {
      toast({
        title: 'Name required',
        description: 'Cannot be blank',
        status: 'error',
        position: 'top-right',
        duration: 5000,
        isClosable: true,
      })
      setLoading(false)
      return
    } else if (amount === 0) {
      toast({
        title: 'Amount required',
        description: 'Must be a positive or negative number',
        status: 'error',
        position: 'top-right',
        duration: 5000,
        isClosable: true,
      })
      setLoading(false)
      return
    }

    await fetch(`${process.env.REACT_APP_API_URL}/transactions`, {
      method: 'PATCH',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        transaction_id: transaction.transaction_id,
        date: date.toUTCString(),
        name: name,
        category: category,
        amount: amount.toString(),
      }),
    })
      .then((res) => {
        if (res.status === 401) {
          navigate('/login')
        }
        if (res.status === 200) {
          onClose()
          forceRefresh()
          onSuccess()
          toast({
            title: 'Success!',
            description: 'Transaction updated',
            status: 'success',
            position: 'top-right',
            duration: 5000,
            isClosable: true,
          })
        }
        if (res.status !== 200) {
          toast({
            title: 'Something went wrong',
            description: 'Please try again later',
            status: 'error',
            position: 'top-right',
            duration: 5000,
            isClosable: true,
          })
        }
        setLoading(false)
      })
      .catch((e) => {
        logger('error updating transaction', e)
      })
  }

  return (
    <>
      <IconButton
        type="button"
        onClick={onOpen}
        bg={bgColor}
        icon={<FaPencilAlt />}
      />
      <AlertDialog
        isOpen={isOpen}
        leastDestructiveRef={cancelRef}
        onClose={onClose}
      >
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Updated Transaction
            </AlertDialogHeader>

            <AlertDialogBody>
              <FormControl>
                <FormLabel>Date</FormLabel>
                <DatePicker
                  selected={date}
                  onChange={(date) => setDate(date)}
                />
                <FormLabel mt={3}>Name</FormLabel>
                <Input
                  defaultValue={name}
                  onChange={(event) => setName(event.target.value)}
                  mb={3}
                />
                <FormLabel>Category</FormLabel>
                <Select
                  defaultValue={category}
                  onChange={async (event) => {
                    setCategory(event.target.value)
                  }}
                  mb={3}
                >
                  <option value={'bills'}>Bills</option>
                  <option value={'entertainment'}>Entertainment</option>
                  <option value={'groceries'}>Groceries</option>
                  <option value={'ignore'}>Ignore</option>
                  <option value={'income'}>Income</option>
                  <option value={'restaurant'}>Restaurant</option>
                  <option value={'transportation'}>Transportation</option>
                  <option value={'vacation'}>Vacation</option>
                  <option value={'uncategorized'}>Uncategorized</option>
                </Select>
                <FormLabel>Amount</FormLabel>
                <NumberInput defaultValue={transaction.amount}>
                  <NumberInputField
                    onChange={(event) => setAmount(event.target.value)}
                    mb={3}
                  />
                </NumberInput>
              </FormControl>
            </AlertDialogBody>
            <AlertDialogFooter>
              {loading && (
                <Spinner
                  thickness="4px"
                  speed="0.65s"
                  emptyColor="gray.200"
                  color="blue.500"
                  size="md"
                  mr={5}
                />
              )}
              <Button
                ref={cancelRef}
                onClick={() => {
                  onClose()
                  setLoading(false)
                }}
              >
                Cancel
              </Button>
              <Button
                bg={colors.primary}
                color={'white'}
                _hover={{
                  bg: colors.primaryFaded,
                }}
                onClick={() => {
                  setLoading(true)
                  updateTransaction()
                }}
                ml={3}
                disabled={loading}
              >
                Submit
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </>
  )
}
