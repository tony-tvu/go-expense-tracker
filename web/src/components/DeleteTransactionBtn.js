import React from 'react'
import {
  AlertDialog,
  AlertDialogBody,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogOverlay,
  Button,
  useDisclosure,
  IconButton,
  useColorModeValue,
} from '@chakra-ui/react'
import logger from '../logger'
import { useNavigate } from 'react-router-dom'
import { FaTrashAlt } from 'react-icons/fa'

export default function DeleteTransactionBtn({ transaction, onSuccess }) {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const cancelRef = React.useRef()
  const navigate = useNavigate()

  const bgColor = useColorModeValue('white', '#252526')

  async function deleteTransaction() {
    await fetch(
      `${process.env.REACT_APP_API_URL}/transactions/${transaction.transactionId}`,
      {
        method: 'DELETE',
        credentials: 'include',
      }
    )
      .then((res) => {
        if (res.status === 401) {
          navigate('/login')
        }
        if (res.status === 200) {
          onClose()
          onSuccess()
        }
      })
      .catch((e) => {
        logger('error deleting transaction', e)
      })
  }

  return (
    <>
      <IconButton
        type="button"
        onClick={onOpen}
        bg={bgColor}
        icon={<FaTrashAlt />}
      />
      <AlertDialog
        isOpen={isOpen}
        leastDestructiveRef={cancelRef}
        onClose={onClose}
      >
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Delete Transaction
            </AlertDialogHeader>

            <AlertDialogBody>
              Are you sure you want to remove {transaction.name}?
            </AlertDialogBody>

            <AlertDialogFooter>
              <Button ref={cancelRef} onClick={onClose}>
                Cancel
              </Button>
              <Button
                color={'white'}
                bg={'#E63E3F'}
                _hover={{
                  bg: '#AA2E2F',
                }}
                onClick={async () => await deleteTransaction()}
                ml={3}
              >
                Delete
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </>
  )
}
