import {ActionIcon, Text} from "@mantine/core";
import { openConfirmModal } from '@mantine/modals';
import {Trash} from "tabler-icons-react";

export const ConfirmationModal = (command: () => Promise<void>, message: string) => {
    console.log("Hi")
    openConfirmModal({

    title: 'Please confirm your action',
    centered: true,
    children: (
        <Text size="sm">
            hi

        </Text>
    ),

    confirmProps: { color: 'red' },
    labels: { confirm: 'Confirm', cancel: 'Cancel' },
    onConfirm: () => {
        command()
    },

    onCancel: () => {}

    });}
