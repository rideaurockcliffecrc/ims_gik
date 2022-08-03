import {
    Box,
    Modal,
    Container,
    Button,
    Group,
    InputWrapper,
    TextInput,
    Space,
} from "@mantine/core";
import { showNotification } from "@mantine/notifications";
import { Dispatch, SetStateAction, useEffect, useState } from "react";

import { QrReader } from "react-qr-reader";
import { useNavigate } from "react-router-dom";
import { Item } from "../../types/item";

const PermissionsModal = ({ opened }: { opened: boolean }) => {
    const navigate = useNavigate();

    return (
        <>
            <Modal
                onClose={() => {}}
                opened={opened}
                closeOnClickOutside={false}
                closeOnEscape={false}
                withCloseButton={false}
                title="Permissions Needed"
            >
                <h1>This application requires the camera.</h1>
                <p>Please activate the camera permission to enable scanning.</p>
                <Button
                    color="red"
                    onClick={() =>
                        navigate("/dashboard/analytics", { replace: true })
                    }
                >
                    Leave
                </Button>
            </Modal>
        </>
    );
};

const ResultsModal = ({
    selectedResult,
    setSelectedResult,
    modalOpened,
    setModalOpened,
}: {
    selectedResult: string;
    setSelectedResult: Dispatch<SetStateAction<string>>;
    modalOpened: boolean;
    setModalOpened: Dispatch<SetStateAction<boolean>>;
}) => {
    const [product, setProduct] = useState<Item>({} as Item);

    const [changingStock, setChangingStock] = useState<boolean>(false);

    const [jump, setJump] = useState<number>(1);

    const [netJump, setNetJump] = useState<number>(0);

    const submitStock = async () => {
        setChangingStock(true);
        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/items/stock/jump?product_id=${product.id}&diff=${netJump}`,
            {
                method: "PATCH",
                credentials: "include",
            }
        );

        const data: {
            success: boolean;
            message: string;
            data: number;
        } = await response.json();

        if (data.success) {
            setChangingStock(false);
            showNotification({
                message: `Stock updated to ${data.data}.`,
                color: "green",
                title: "Stock updated",
            });
            setNetJump(0);
            setJump(1);

            await lookupLocation();
            return;
        }

        showNotification({
            message: data.message,
            color: "red",
            title: "Unable to update stock",
        });
    };

    const lookupLocation = async () => {
        const splitResult = selectedResult.split("");

        const letter = splitResult[splitResult.length - 1];

        // get name which is everything except last character
        const name = selectedResult.slice(0, -1);

        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/location/scan?name=${name}&letter=${letter}`,
            {
                credentials: "include",
            }
        );

        const data: {
            success: boolean;
            data: Item;
        } = await response.json();

        setProduct(data.data);
    };

    useEffect(() => {
        lookupLocation();
    }, [selectedResult]);

    const decStock = async () => {
        setNetJump(netJump - jump);
    };

    const incStock = async () => {
        setNetJump(netJump + jump);
    };

    return (
        <>
            <Modal
                opened={modalOpened}
                onClose={() => {
                    setSelectedResult("");
                    setModalOpened(false);
                }}
                title="Location Details"
            >
                <p>
                    <b>Location ID</b>: {selectedResult}
                </p>
                <p>
                    <b>Item Name</b>: {product.name || "Loading..."}
                </p>
                <p>
                    <b>Current Stock</b>: {product.quantity || 0}
                </p>
                <b>Stock Control</b>
                <p>
                    <b>Net Jump</b>:{" "}
                    <span
                        style={{
                            color: netJump < 0 ? "red" : "green",
                        }}
                    >
                        {netJump}
                    </span>
                </p>
                <Group
                    grow
                    sx={{
                        alignItems: "flex-end",
                    }}
                >
                    <InputWrapper label="Jump">
                        <TextInput
                            value={jump.toString()}
                            type="number"
                            onChange={(e) => setJump(Number(e.target.value))}
                        />
                    </InputWrapper>
                    <Button
                        color="red"
                        onClick={decStock}
                        disabled={changingStock}
                    >
                        -{jump}
                    </Button>
                    <Button
                        color="green"
                        onClick={incStock}
                        disabled={changingStock}
                    >
                        +{jump}
                    </Button>
                </Group>
                <Space h="md" />
                <Button color="green" onClick={submitStock}>
                    Submit Stock
                </Button>
            </Modal>
        </>
    );
};

const Scanner = () => {
    const [newResult, setNewResult] = useState<string>("");

    const [results, setResults] = useState<string[]>([]);

    const [modalOpened, setModalOpened] = useState<boolean>(false);
    const [selectedResult, setSelectedResult] = useState<string>("");

    const [permissionsAllowed, setPermissionsAllowed] = useState<boolean>(true);

    useEffect(() => {
        if (!newResult) return;
        setResults([newResult, ...results]);
    }, [newResult]);

    useEffect(() => {
        if (!selectedResult) return;

        setModalOpened(true);
    }, [selectedResult]);

    useEffect(() => {
        checkPermission();
    }, []);

    const checkPermission = async () => {
        // let allowed = false;

        // try {
        //     allowed = !!(await navigator.mediaDevices.getUserMedia({
        //         video: true,
        //     }));
        // } catch (err) {
        //     /* handle the error */
        //     allowed = false;
        // }

        setPermissionsAllowed(true);
    };

    return (
        <>
            <Container>
                <h1>Item Scanner</h1>

                <h2>
                    Simply point your device at a QR code and wait for the
                    results to appear.
                </h2>
                {/* {permissionsAllowed && ( */}
                <QrReader
                    onResult={(result, err) => {
                        // @ts-ignore
                        if (!result?.text || results.includes(result?.text))
                            return;
                        // @ts-ignore
                        setNewResult(result?.text);
                    }}
                    constraints={{ facingMode: "environment" }}
                    scanDelay={500}
                />
                {/* )} */}
                <br />
                {results.map((result) => (
                    <Box
                        onClick={() => {
                            setSelectedResult(result);
                        }}
                        sx={{
                            borderRadius: "5px",
                            backgroundColor: "#fff",
                            padding: "1rem",
                            boxSizing: "border-box",
                            marginBottom: "1rem",
                        }}
                    >
                        <p style={{ margin: 0 }}>{result}</p>
                    </Box>
                ))}
            </Container>
            {results.length ? (
                <Button
                    color="red"
                    onClick={() => {
                        setResults([]);
                    }}
                >
                    Clear Scans
                </Button>
            ) : null}
            <ResultsModal
                selectedResult={selectedResult}
                setSelectedResult={setSelectedResult}
                modalOpened={modalOpened}
                setModalOpened={setModalOpened}
            />
            <PermissionsModal opened={!permissionsAllowed} />
        </>
    );
};

export default Scanner;
