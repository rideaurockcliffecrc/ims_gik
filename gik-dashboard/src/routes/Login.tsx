import {
    Container,
    Input,
    InputWrapper,
    Modal,
    PasswordInput,
    Space,
    Checkbox,
    Text,
} from "@mantine/core";

import styles from "../styles/Auth.module.scss";

import { TextInput, Button, Group } from "@mantine/core";
import { useForm } from "@mantine/form";

import { showNotification, hideNotification } from "@mantine/notifications";
import { useEffect, useState } from "react";
import { useNavigate, Link } from "react-router-dom";

const Login = () => {
    const navigate = useNavigate();

    const [loginEnabled, setLoginEnabled] = useState<boolean>(true);

    const [askTfa, setAskTfa] = useState<boolean>(false);

    const [topVerification, setTopVerification] = useState<string>("");

    const form = useForm({
        initialValues: {
            username: "",
            password: "",
            rememberMe: "",
        },
    });

    const checkAuthStatus = async () => {
        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/auth/status`,
            {
                credentials: "include",
            }
        );

        if (response.status === 200) navigate("/dashboard", { replace: true });
    };
/*
    const do2FACheck = async (verification: string) => {
        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/auth/tfa?verification=${verification}&username=${form.values.username}`
        );

        const data: {
            success: boolean;
            message: string;
            data: boolean;
        } = await response.json();

        if (!data.success) {
            showNotification({
                color: "red",
                title: "2FA Check Failed",
                message: data.message,
            });
            return;
        }

        if (data.data) {
            setTopVerification(verification);
            setAskTfa(true);
        } else {
            await doLogin(verification);
        }
    };*/

    const doPrelogin = async () => {
        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/auth/prelogin`,
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    username: form.values.username,
                    password: form.values.password,
                }),
            }
        );

        const data: {
            success: boolean;
            message: string;
            data: string;
        } = await response.json();

        if (!data.success) {
            showNotification({
                color: "red",
                title: "Login Failed",
                message: data.message,
            });
            return;
        }
        await doLogin(data.data);
        //await do2FACheck(data.data);
    };

    const doLogin = async (verification: string) => {
        setLoginEnabled(false);

        showNotification({
            id: "loginLoading",
            loading: true,
            title: "Logging you in...",
            message: "Please wait while we authenticate you.",
        });

        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/auth/login`,
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    username: form.values.username,
                    password: form.values.password,
                    //totp: form.values.totp || "",
                    verificationJWT: verification,
                    rememberMe: form.values.rememberMe,
                }),
                credentials: "include",
            }
        );

        hideNotification("loginLoading");

        if (response.status === 200) {
            showNotification({
                color: "green",
                title: "Logged in",
                message: "Registration successful. Redirecting to dashboard...",
            });

            setTimeout(() => {
                navigate("/dashboard", { replace: true });
            }, 1000);
        } else {
            const data = await response.json();

            showNotification({
                color: "red",
                title: "Login Failed",
                message: data.message,
            });
        }

        setLoginEnabled(true);
    };

    useEffect(() => {
        checkAuthStatus();
    }, []);

    return (
        <>
            <div className={`${styles.wrapper} ${styles.login}`}>
                <Container
                    sx={{
                        backgroundColor: "var(--inverted-text)",
                        padding: "2rem",
                        borderRadius: "5px",
                        width: "40rem",
                    }}
                >
                    <h1
                        style={{
                            marginTop: 0,
                        }}
                    >
                        Login
                    </h1>
                    <form
                        onSubmit={(e) => {
                            e.preventDefault();
                            doPrelogin();
                        }}
                    >
                        <TextInput
                            required
                            label="Username"
                            placeholder="amy"
                            {...form.getInputProps("username")}
                        />
                        <PasswordInput
                            required
                            placeholder="Password"
                            label="Password"
                            {...form.getInputProps("password")}
                        />
                        <Group position="apart">
                            <Checkbox
                                label="Remember me"
                                color="green"
                                {...form.getInputProps("rememberMe")}
                            />

                            <Group position="right" mt="md">
                                {window.location.protocol === "https:" ? (
                                    <p style={{ color: "green" }}>
                                        Your connection is secure.
                                    </p>
                                ) : (
                                    <p style={{ color: "red" }}>
                                        Your connection is insecure. Don't login
                                        unless you know what you're doing.
                                    </p>
                                )}
                                <Button
                                    type="submit"
                                    color="green"
                                    disabled={!loginEnabled}
                                >
                                    Login
                                </Button>
                            </Group>
                        </Group>
                    </form>
                    <Button component={Link} to="/register" compact variant="white">
                        <Text color="blue" size={"xs"}>Not Registered?</Text>
                    </Button>
                </Container>

            </div>
        </>
    );
};

export default Login;
