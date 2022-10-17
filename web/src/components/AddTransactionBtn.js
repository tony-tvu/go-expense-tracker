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
} from '@chakra-ui/react'
import logger from '../logger'
import { BsPlus } from 'react-icons/bs'
import { useNavigate } from 'react-router-dom'
import { colors } from '../theme'
import DatePicker from './Datepicker'

export default function AddTransactionBtn({ onSuccess, icon }) {
  const [loading, setLoading] = useState(false)
  const [date, setDate] = useState(new Date())
  const [name, setName] = useState(null)
  const [category, setCategory] = useState('bills')
  const [amount, setAmount] = useState(0)
  const { isOpen, onOpen, onClose } = useDisclosure()
  const cancelRef = React.useRef()
  const navigate = useNavigate()
  const toast = useToast()

  async function createTransaction() {
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
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        date: date.toUTCString(),
        name: name,
        category: category,
        amount: amount,
      }),
    })
      .then((res) => {
        if (res.status === 401) {
          navigate('/login')
        }
        if (res.status === 200) {
          onClose()
          onSuccess()
          toast({
            title: 'Success!',
            description: 'New transaction saved',
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
        logger('error creating transaction', e)
      })
  }

  return (
    <>
      <Button
        leftIcon={icon ? icon : <BsPlus />}
        type="button"
        variant="solid"
        onClick={onOpen}
        bg={colors.primary}
        color={'white'}
        _hover={{
          bg: colors.primaryFaded,
        }}
      >
        Add Transaction
      </Button>

      <AlertDialog
        isOpen={isOpen}
        leastDestructiveRef={cancelRef}
        onClose={onClose}
      >
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              New Transaction
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
                <NumberInput>
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
                  createTransaction()
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
