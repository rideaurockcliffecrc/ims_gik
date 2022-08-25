import {
    Box,
    Table,
    Space,
    Group,
    Button,
    TextInput,
    Modal,
    InputWrapper,
    ActionIcon,
} from "@mantine/core";
import { useForm } from "@mantine/hooks";
import { showNotification } from "@mantine/notifications";
import { Dispatch, SetStateAction, useEffect, useState } from "react";
import {CirclePlus, TableExport, Trash} from "tabler-icons-react";
import { containerStyles } from "../../styles/container";
import { Client } from "../../types/client";

const ClientComponent = ({
    client,
    refresh,
}: {
    client: Client;
    refresh: () => Promise<void>;
}) => {
    const doDelete = async () => {
        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/client/delete?id=${client.ID}`,
            {
                method: "DELETE",
                credentials: "include",
            }
        );

        if (response.ok) {
            showNotification({
                message: "Client deleted",
                color: "green",
                title: "Success",
            });
            await refresh();
            return;
        }

        showNotification({
            message: "Failed to delete client",
            color: "red",
            title: "Error",
        });
    };



    return (
        <>
            <tr>
                <td>{client.ID}</td>
                <td>{client.name}</td>
                <td>{client.contact}</td>
                <td>{client.phone}</td>
                <td>{client.email}</td>
                <td>{client.address}</td>
                <td>{client.balance}</td>
                <td>
                    {" "}
                    <ActionIcon variant="default" onClick={doDelete}>
                        <Trash />
                    </ActionIcon>
                </td>
            </tr>
        </>
    );
};

const CreateClientModal = ({
    opened,
    setOpened,
    refresh,
}: {
    opened: boolean;
    setOpened: Dispatch<SetStateAction<boolean>>;
    refresh: () => Promise<void>;
}) => {
    const form = useForm({
        initialValues: {
            name: "",
            contact: "",
            phone: "",
            email: "",
            address: "",
            balance: 0,
        },
    });

    const doCreate = async () => {
        form.values.balance = Number(form.values.balance);

        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/client/add`,
            {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
                credentials: "include",
                body: JSON.stringify(form.values),
            }
        );

        const data: {
            success: boolean;
            message: string;
        } = await response.json();

        if (data.success) {
            showNotification({
                color: "green",
                title: "Client added",
                message: data.message,
            });

            await refresh();
            setOpened(false);
            return;
        }

        showNotification({
            color: "red",
            title: "Unable to add client",
            message: data.message,
        });
    };

    return (
        <>
            <Modal
                opened={opened}
                onClose={() => {
                    refresh();
                    setOpened(false);
                }}
                title="Create Client"
            >
                <form
                    onSubmit={(e) => {
                        e.preventDefault();
                        doCreate();
                    }}
                >
                    <InputWrapper label="Name" required>
                        <TextInput
                            required
                            placeholder="Gifts In Kind"
                            {...form.getInputProps("name")}
                        />
                    </InputWrapper>
                    <Space h="md" />
                    <InputWrapper label="Contact Person" required>
                        <TextInput
                            required
                            placeholder="Mr. John Doe"
                            {...form.getInputProps("contact")}
                        />
                    </InputWrapper>
                    <Space h="md" />
                    <Group grow>
                        <InputWrapper label="Phone" required>
                            <TextInput
                                required
                                type="tel"
                                placeholder="+1 (555) 555-5555"
                                {...form.getInputProps("phone")}
                            />
                        </InputWrapper>
                        <InputWrapper label="Email" required>
                            <TextInput
                                required
                                placeholder="john@giftsinkindottawa.ca"
                                type="email"
                                {...form.getInputProps("email")}
                            />
                        </InputWrapper>
                    </Group>
                    <Space h="md" />
                    <InputWrapper label="Address" required>
                        <TextInput
                            required
                            placeholder="123 Main St"
                            {...form.getInputProps("address")}
                        />
                    </InputWrapper>
                    <Space h="md" />
                    <InputWrapper label="Balance" required>
                        <TextInput
                            required
                            placeholder="50"
                            type="number"
                            {...form.getInputProps("balance")}
                        />
                    </InputWrapper>
                    <Space h="md" />
                    <Group position="right">
                        <Button color="green" type="submit">
                            Create
                        </Button>
                    </Group>
                </form>
            </Modal>
        </>
    );
};

const Clients = () => {
    const [clients, setClients] = useState<Client[]>([]);

    const [loading, setLoading] = useState<boolean>(true);

    const [nameQuery, setNameQuery] = useState<string>("");

    const [contactQuery, setContactQuery] = useState<string>("");

    const [phoneQuery, setPhoneQuery] = useState<string>("");

    const [emailQuery, setEmailQuery] = useState<string>("");

    const [addressQuery, setAddressQuery] = useState<string>("");

    const [showClientCreationModal, setShowClientCreationModal] =
        useState<boolean>(false);

    const exportCSV = async () => {
        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/client/export?name=${nameQuery}&contact=${contactQuery}&phone=${phoneQuery}&email=${emailQuery}&address=${addressQuery}`,
            {
                credentials: "include",
            }
        );

        // data is arraybuffer
        const data = await response.arrayBuffer();

        // convert to blob
        const blob = new Blob([data], { type: "text/csv" });

        // create a url from the blob
        const url = URL.createObjectURL(blob);

        // open the url with the pdf viewer
        window.open(url, "_blank");
    };

    const fetchClients = async () => {
        setLoading(true);

        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/client/list?name=${nameQuery}&contact=${contactQuery}&phone=${phoneQuery}&email=${emailQuery}&address=${addressQuery}`,
            {
                credentials: "include",
            }
        );

        setLoading(false);

        const data: {
            success: boolean;
            data: Client[];
        } = await response.json();

        if (data.success) {
            setClients(data.data);
        }
    };

    useEffect(() => {
        fetchClients();
    }, []);

    return (
        <>
            <Box sx={containerStyles}>
                <Group position="apart">
                    <h3>Clients</h3>
                    <ActionIcon
                        sx={{
                            height: "4rem",
                            width: "4rem",
                        }}
                        onClick={() => setShowClientCreationModal(true)}
                    >
                        <CirclePlus />
                    </ActionIcon>
                </Group>
                <Space h="md" />
                <Group>
                    <TextInput
                        placeholder="Name"
                        onChange={(e) => setNameQuery(e.target.value)}
                    />
                    <TextInput
                        placeholder="Contact"
                        onChange={(e) => setContactQuery(e.target.value)}
                    />
                    <TextInput
                        placeholder="Phone"
                        onChange={(e) => setPhoneQuery(e.target.value)}
                    />
                    <TextInput
                        placeholder="Email"
                        onChange={(e) => setEmailQuery(e.target.value)}
                    />
                    <TextInput
                        placeholder="Address"
                        onChange={(e) => setAddressQuery(e.target.value)}
                    />
                    <Button
                        color="green"
                        disabled={loading}
                        onClick={() => fetchClients()}
                    >
                        Filter
                    </Button>
                </Group>
                <Space h="md" />
                <Table>
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>Name</th>
                            <th>Contact</th>
                            <th>Phone</th>
                            <th>Email</th>
                            <th>Address</th>
                            <th>Balance</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        {clients.map((client) => (
                            <ClientComponent
                                refresh={fetchClients}
                                key={client.ID}
                                client={client}
                            />
                        ))}
                    </tbody>
                </Table>
                <Space h="md" />
                <Group position="apart">
                    <ActionIcon
                        sx={{
                            height: "4rem",
                            width: "4rem",
                        }}
                        onClick={exportCSV}
                    >
                        <TableExport />
                    </ActionIcon>
                    <Button
                        onClick={fetchClients}
                        color="green"
                        disabled={loading}
                    >
                        Refresh
                    </Button>
                </Group>
            </Box>
            <CreateClientModal
                opened={showClientCreationModal}
                setOpened={setShowClientCreationModal}
                refresh={fetchClients}
            />
        </>
    );
};

const ClientsDonors = () => {
    return (
        <>
            <Clients />
        </>
    );
};

export default ClientsDonors;
