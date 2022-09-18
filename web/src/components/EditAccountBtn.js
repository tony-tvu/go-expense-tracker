import {
  Popover,
  PopoverTrigger,
  PopoverContent,
  PopoverBody,
  PopoverArrow,
  IconButton,
  Button,
  Stack,
  useDisclosure,
  AlertDialog,
  AlertDialogOverlay,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogBody,
  AlertDialogFooter,
} from '@chakra-ui/react'

import { BsPencil, BsTrash } from 'react-icons/bs'
import React from 'react'
import logger from '../logger'

export default function EditAccountBtn(props) {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const cancelRef = React.useRef()

  async function deleteAccount() {
    await fetch(`${process.env.REACT_APP_API_URL}/items`, {
      method: 'DELETE',
      credentials: 'include',
      body: JSON.stringify({ id: props.item.id }),
    })
      .then((res) => {
        if (res.status === 200) {
          props.onSuccess()
          onClose()
        }
      })
      .catch((e) => {
        logger('error setting access token', e)
      })
  }

  return (
    <>
      <Popover placement="bottom" isLazy>
        <PopoverTrigger>
          <IconButton
            aria-label="More server options"
            icon={<BsPencil />}
            variant="solid"
            w="fit-content"
          />
        </PopoverTrigger>
        <PopoverContent w="fit-content" _focus={{ boxShadow: 'none' }}>
          <PopoverArrow />
          <PopoverBody>
            <Stack>
              <Button
                w="150px"
                variant="ghost"
                rightIcon={<BsTrash />}
                justifyContent="space-between"
                fontWeight="normal"
                colorScheme="red"
                fontSize="md"
                as="b"
                onClick={onOpen}
              >
                Delete
              </Button>
            </Stack>
          </PopoverBody>
        </PopoverContent>
      </Popover>

      <AlertDialog
        isOpen={isOpen}
        leastDestructiveRef={cancelRef}
        onClose={onClose}
      >
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Delete Account
            </AlertDialogHeader>

            <AlertDialogBody>
              Are you sure you want to remove {props.item.institution}?
            </AlertDialogBody>

            <AlertDialogFooter>
              <Button ref={cancelRef} onClick={onClose}>
                Cancel
              </Button>
              <Button colorScheme="red" onClick={deleteAccount} ml={3}>
                Delete
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </>
  )
}
