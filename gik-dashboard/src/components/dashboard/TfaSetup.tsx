import {
    Button,
    Center,
    Group,
    Input,
    InputWrapper,
    LoadingOverlay,
    Modal,
    Space,
} from "@mantine/core";
import { showNotification } from "@mantine/notifications";
import {
    ChangeEventHandler,
    Dispatch,
    SetStateAction,
    useEffect,
    useState,
} from "react";

import QRCode from "react-qr-code";

const TfaSetup = ({
    opened,
    setOpened,
}: {
    opened: boolean;
    setOpened: Dispatch<SetStateAction<boolean>>;
}) => {
    const [secret, setSecret] = useState<string>("");

    const [secretSecret, setSecretSecret] = useState<string>("");

    const [verifyCode, setVerifyCode] = useState<string>("");

    const [loading, setLoading] = useState<boolean>(true);

    const finishSetup = async () => {
        setLoading(true);

        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/settings/tfa/setup?secret=${secretSecret}&code=${verifyCode}`,
            {
                method: "PATCH",
                credentials: "include",
            }
        );

        setLoading(false);

        const data = await response.json();

        if (data.success) {
            showNotification({
                color: "green",
                message: data.message,
            });

            setOpened(false);

            return;
        }

        showNotification({
            color: "red",
            message: data.message,
        });
    };

    const generateSecret = async () => {
        setLoading(true);

        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/settings/tfa/generate`,
            {
                credentials: "include",
            }
        );

        setLoading(false);

        const data: {
            success: boolean;
            message: string;
            data: {
                secret: string;
                url: string;
            }
        } = await response.json();

        if (data.success) {
            setSecret(data.data.url);
            setSecretSecret(data.data.secret);
        } else {
            showNotification({
                color: "red",
                message: data.message,
            });
        }
    };

    useEffect(() => {
        generateSecret();
    }, []);

    return (
        <>
            <Modal
                opened={opened}
                onClose={() => {}}
                withCloseButton={false}
                closeOnEscape={false}
                closeOnClickOutside={false}
            >
                <h1>Secure Your Account</h1>
                <p>
                    This application requires you to setup{" "}
                    <b>two-factor authentication</b> to prevent bad actors from
                    accessing your account. You'll need to follow the steps
                    below to continue. Once completed, we won't show this
                    message again.
                </p>
                <h3>Instructions</h3>
                <ol>
                    <li>
                        Download one of the following MFA applications to your
                        mobile device.{" "}
                        <a
                            href="https://www.tofuauth.com/"
                            target="_blank"
                            rel="noopener noreferrer"
                        >
                            Tofu (iPhone)
                        </a>{" "}
                        or{" "}
                        <a
                            href="https://getaegis.app/"
                            target="_blank"
                            rel="noopener noreferrer"
                        >
                            Aegis (Android)
                        </a>
                        .
                    </li>
                    <li>
                        Scan the QR code below with the newly installed
                        authenticator apps.
                    </li>
                    <li>
                        Enter the code displayed on your mobile device into the
                        box below the code.
                    </li>
                    <li>Click "Finish" to complete setup and continue.</li>
                </ol>
                <Center>{secret && <QRCode value={secret} />}</Center>
                <Space h="md" />
                <InputWrapper label="Step 3: 6 Digit Code">
                    <Input
                        placeholder="123456"
                        onChange={(e: any) => setVerifyCode(e.target.value)}
                        maxLength={6}
                    />
                </InputWrapper>
                <Space h="xl" />
                <Group position="right">
                    <Button color="green" onClick={finishSetup}>
                        Finish
                    </Button>
                </Group>
                <LoadingOverlay visible={loading} />
            </Modal>
        </>
    );
};

export default TfaSetup;
