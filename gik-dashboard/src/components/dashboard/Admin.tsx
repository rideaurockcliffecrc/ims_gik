import {
    ActionIcon,
    Box,
    Button,
    Center,
    Group,
    Input,
    InputWrapper,
    Loader,
    Pagination,
    Space,
    Switch,
    Table,
} from "@mantine/core";

import { TransferList, TransferListData } from "@mantine/core";
import { hideNotification, showNotification } from "@mantine/notifications";
import { useEffect, useState } from "react";
import { SignupCode } from "../../types/signupCode";
import { User } from "../../types/user";

import { containerStyles } from "../../styles/container";
import {Trash} from "tabler-icons-react";

const UserManagement = () => {
    const [users, setUsers] = useState<User[]>([]);

    const [loading, setLoading] = useState<boolean>(true);

    const [disableToggle, setDisableToggle] = useState<boolean>(false);

    const [disableReset, setDisableReset] = useState<boolean>(false);

    const getUsers = async () => {
        setLoading(true);

        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/admin/users`,
            {
                credentials: "include",
            }
        );

        setLoading(false);

        const users = await response.json();

        if (!users.success) {
            showNotification({
                message: users.message,
                color: "red",
                title: "Unable to fetch users",
            });
            return;
        }

        setUsers(users.data);
    };

    useEffect(() => {
        getUsers();
    }, []);

    const toggleUser = async (id: number) => {
        setDisableToggle(true);

        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/admin/user/toggle?user_id=${id}`,
            {
                method: "PATCH",
                credentials: "include",
            }
        );

        setDisableToggle(false);

        if (!response.ok) {
            showNotification({
                message: "We were unable to process your request.",
                color: "red",
                title: "Toggle user error",
            });
            return;
        }

        await getUsers();
    };

    const doDelete = async (id: number, username: string) => {
        setDisableReset(true);
        showNotification({
            loading: true,
            message: "Please wait while we process your request.",
            title: "Deleting User...",
            id: "deleteUser",
        });

        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/admin/user/delete?user_id=${id}`,
            {
                method: "DELETE",
                credentials: "include",
            }
        );

        setDisableReset(false);

        hideNotification("deleteUser");

        if (!response.ok) {
            showNotification({
                message: "We were unable to process your request.",
                color: "red",
                title: "Delete User Error",
            });
            return;
        }

        showNotification({
            message: username + " has been deleted",
            color: "green",
            title: "User Deleted",
        });
        getUsers()
    };

    return (
        <>
            <Box style={containerStyles}>
                <h3>User Management</h3>
                <Space h="xl" />
                {loading ? (
                    <Center>
                        <Loader />
                    </Center>
                ) : (
                    <Table>
                        <thead>
                            <tr>
                                <th>Status</th>
                                <th>ID</th>
                                <th>Username</th>
                                <th>Registered At</th>
                                <th>Delete</th>
                            </tr>
                        </thead>
                        <tbody>
                            {users.map((user) => (
                                <tr key={user.ID}>
                                    <td>
                                        <Switch
                                            color={"green"}
                                            checked={!user.disabled}
                                            onChange={() => {
                                                toggleUser(user.ID);
                                            }}
                                            disabled={disableToggle}
                                        />
                                    </td>
                                    <td>{user.ID}</td>
                                    <td>{user.username}</td>
                                    <td>
                                        {new Date(
                                            user.registeredAt * 1000
                                        ).toLocaleString()}
                                    </td>
                                    <td>
                                        <ActionIcon variant="default" onClick={() => {doDelete(user.ID, user.username)}}>
                                            <Trash />
                                        </ActionIcon>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </Table>
                )}
                <Group position="right">
                    <Button color="green" disabled={loading} onClick={getUsers}>
                        Refresh
                    </Button>
                </Group>
            </Box>
        </>
    );
};

const SignupCodesManager = () => {
    const [codes, setCodes] = useState<SignupCode[]>([]);

    const [loading, setLoading] = useState<boolean>(false);

    const [totalPages, setTotalPages] = useState<number>(1);

    const [currentPage, setCurrentPage] = useState<number>(1);

    const [disableToggle, setDisableToggle] = useState<boolean>(false);



    const doDelete = async (id: number, username: string) => {
        showNotification({
            loading: true,
            message: "Please wait while we process your request.",
            title: "Deleting User...",
            id: "deleteUser",
        });

        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/admin/scode/delete?user_id=${id}`,
            {
                method: "DELETE",
                credentials: "include",
            }
        );


        hideNotification("deleteUser");

        if (!response.ok) {
            showNotification({
                message: "We were unable to process your request.",
                color: "red",
                title: "Delete User Error",
            });
            return;
        }

        showNotification({
            message: username + " has been deleted",
            color: "green",
            title: "User Deleted",
        });
        getCodes()
    };


    const getCodes = async () => {
        setLoading(true);

        showNotification({
            loading: true,
            message: "Please wait while we process your request.",
            title: "Loading codes...",
            id: "loadingNotification",
        });

        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/admin/scodes?page=${currentPage}`,
            {
                credentials: "include",
            }
        );

        setLoading(false);

        hideNotification("loadingNotification");

        const data: {
            success: boolean;
            message: string;
            data: {
                data: SignupCode[];
                totalPages: number;
                currentPage: number;
                total: number;
            };
        } = await response.json();

        setTotalPages(data.data.totalPages);

        if (data.success) {
            setCodes(data.data.data);
        }
    };

    useEffect(() => {
        getCodes();
    }, [currentPage]);

    const toggleCode = async (code: string) => {
        setDisableToggle(true);

        showNotification({
            loading: true,
            message: "Please wait while we process your request.",
            title: "Toggling code...",
            id: "loadingNotification",
        });

        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/admin/scode/toggle?code=${code}`,
            {
                method: "PATCH",
                credentials: "include",
            }
        );

        setDisableToggle(false);
        hideNotification("loadingNotification");

        if (response.status !== 200) {
            showNotification({
                message:
                    "We were unable to toggle the signup code. Maybe it's time expired?",
                title: "Signup Code Toggle Error",
                color: "red",
            });
            return;
        }

        await getCodes();
    };

    return (
        <>
            <Box sx={containerStyles}>
                <h3>Signup Codes</h3>
                <Space h="xl" />
                <Table>
                    <thead>
                        <tr>
                            <th>Status</th>
                            <th>Code</th>
                            <th>Username</th>
                            <th>Timestamp</th>
                            <th>Expiration</th>
                            <th>Delete</th>
                        </tr>
                    </thead>
                    <tbody>
                        {codes.map((code) => (
                            <tr key={code.ID}>
                                <td>
                                    <Switch
                                        color="green"
                                        checked={!code.expired}
                                        onChange={() => {
                                            toggleCode(code.code);
                                        }}
                                        disabled={disableToggle}
                                    />
                                </td>
                                <td>{code.code}</td>
                                <td>{code.designatedUsername}</td>
                                <td>
                                    {new Date(
                                        code.createdAt * 1000
                                    ).toLocaleString()}
                                </td>
                                <td>
                                    {new Date(
                                        code.expiration * 1000
                                    ).toLocaleString()}
                                </td>
                                <td>
                                    <ActionIcon variant="default" onClick={() => {doDelete(code.ID, code.designatedUsername)}}>
                                        <Trash />
                                    </ActionIcon>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </Table>
                <Space h="xl" />
                <Group position="right">
                    <Pagination
                        page={currentPage}
                        onChange={setCurrentPage}
                        total={totalPages}
                    />
                    <Button color="green" disabled={loading} onClick={getCodes}>
                        Refresh
                    </Button>
                </Group>
            </Box>
        </>
    );
};

const AdminManager = () => {
    const [values, setValues] = useState<TransferListData>(
        {} as TransferListData
    );

    const getData = async () => {
        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/admin/lists`,
            {
                credentials: "include",
            }
        );

        const data: {
            success: boolean;
            data: TransferListData;
        } = await response.json();

        if (data.success) {
            setValues(data.data);
        }
    };

    const doSave = async () => {
        showNotification({
            loading: true,
            id: "saveNotification",
            title: "Editing admins...",
            message: "Please wait while we process your changes.",
        });

        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/admin/admins`,
            {
                method: "PATCH",
                credentials: "include",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(values),
            }
        );

        hideNotification("saveNotification");

        const data: {
            success: boolean;
            message: string;
        } = await response.json();

        if (data.success) {
            showNotification({
                color: "green",
                message: data.message,
            });
            await getData();
        }
    };

    useEffect(() => {
        getData();
    }, []);

    return (
        <>
            <Box sx={containerStyles}>
                <h3>Admin Management</h3>
                <Space h="xl" />
                {values.length > 0 ? (
                    <TransferList
                        value={values}
                        onChange={setValues}
                        searchPlaceholder="Search Users"
                        nothingFound="No users found"
                        titles={["Users", "Admins"]}
                        breakpoint="sm"
                    />
                ) : (
                    <Center>
                        <Loader />
                    </Center>
                )}
                <Group position="right">
                    <Button
                        onClick={doSave}
                        color="green"
                        sx={{
                            marginTop: "1rem",
                        }}
                    >
                        Save
                    </Button>
                </Group>
            </Box>
        </>
    );
};

const SignupCodeGenerator = () => {
    const [newCode, setNewCode] = useState<string>("");

    const [newUsername, setNewUsername] = useState<string>("");

    const [loading, setLoading] = useState<boolean>(false);

    const doGeneration = async () => {
        setLoading(true);

        showNotification({
            loading: true,
            id: "generateNotification",
            title: "Generating code...",
            message: "Please wait while we process your request.",
        });

        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/admin/scode?username=${newUsername}`,
            {
                credentials: "include",
            }
        );

        hideNotification("generateNotification");

        setLoading(false);

        const data: {
            success: boolean;
            message: string;
            data: string;
        } = await response.json();

        if (data.success) {
            showNotification({
                color: "green",
                message: data.message,
            });

            setNewCode(data.data);

            return;
        }

        showNotification({
            color: "red",
            message: data.message,
        });
    };

    return (
        <>
            <Box sx={containerStyles}>
                <h3>Signup Code Generator</h3>
                <Space h="md" />
                <InputWrapper label="New Account Username">
                    <Input
                        onChange={(e: any) => setNewUsername(e.target.value.toLowerCase())}
                    />
                </InputWrapper>
                <Space h="md" />
                {newCode ? (
                    <p>
                        <b>New Signup Code</b>: {newCode}
                    </p>
                ) : (
                    <p>
                        Click "Generate" and you'll see a new code appear here.
                    </p>
                )}
                <Group position="right">
                    <Button
                        color="green"
                        onClick={doGeneration}
                        disabled={loading}
                    >
                        Generate
                    </Button>
                </Group>
            </Box>
        </>
    );
};

const Admin = () => {
    return (
        <>
            <Box sx={containerStyles}>
                <h1>Admin Panel</h1>
                <h2>Manage users, registration and more.</h2>
            </Box>
            <SignupCodeGenerator />
            <SignupCodesManager />
            <AdminManager />
            <UserManagement />
        </>
    );
};

export default Admin;
