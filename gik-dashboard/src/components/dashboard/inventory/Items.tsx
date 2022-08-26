import {
    Group,
    InputWrapper,
    TextInput,
    Button,
    Center,
    Pagination,
    Space,
    Loader,
    Box,
    Table,
    Modal,
    ActionIcon,
    MultiSelect,
    useMantineTheme,
    Text,
} from "@mantine/core";
import { Dropzone, MIME_TYPES } from '@mantine/dropzone';
import { showNotification } from "@mantine/notifications";
import { useState, useEffect, Dispatch, SetStateAction } from "react";
import { CirclePlus, Tags, Trash, TableExport, TableImport, Upload, Photo, X} from "tabler-icons-react";
import { containerStyles } from "../../../styles/container";
import { Item } from "../../../types/item";

export const ItemRow = ({ item }: { item: Item }) => {
    return (
        <>
            <tr>
                <td>{item.name}</td>
                <td>{item.sku || "None"}</td>
                <td>{item.category || "None"}</td>
                <td>{item.price || "undefined"}</td>
                <td>{item.quantity}</td>
                <td>{item.size}</td>
            </tr>
        </>
    );
};

export const TagRow = ({ tag }: { tag: string }) => {
    return (
        <>
            <tr>
                <td>{tag}</td>
                <ActionIcon variant="default" /*onClick={}*/>
                    <Trash />
                </ActionIcon>
            </tr>
        </>
    );
};



const  UploadCSVModal = (
{
    opened,
    setOpened,
    refresh,

}: {
    opened: boolean;
    setOpened: Dispatch<SetStateAction<boolean>>;
    refresh: () => Promise<void>;
}) => {
    const importCSV = async (file: File) => {
        let data = new FormData()
        await data.append("file", file)
        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/items/import`,
            {
                credentials: "include",
                method: "POST",
                body: data
            }

        );
        await refresh();
        setOpened(false)
    };

    return (
        <>
            <Modal
                opened={opened}
                onClose={() => {
                    refresh();
                    setOpened(false);
                }}
            >
                <Dropzone
                    multiple={false}
                    onDrop={(file) => {importCSV(file[0])}}
                    maxSize={3 * 1024 ** 2}
                    accept={[MIME_TYPES.csv]}
                    children={() => {

                        return (
                            <Group position="center" spacing="xl" style={{ minHeight: 220, pointerEvents: 'none' }}>

                                <div>
                                    <Text size="xl" inline>
                                        Drag CSV here or click to select files
                                    </Text>
                                    <Text size="sm" color="dimmed" inline mt={7}>
                                        Select CSV file containing you upload data
                                    </Text>
                                </div>
                            </Group>
                        )
                    }}
                />
            </Modal>
        </>
    );
}

const CreateItemModal = ({
    opened,
    setOpened,
    refresh,
}: {
    opened: boolean;
    setOpened: Dispatch<SetStateAction<boolean>>;
    refresh: () => Promise<void>;
}) => {
    const [name, setName] = useState("");
    const [sku, setSku] = useState("");
    const [category, setCategory] = useState("");
    const [price, setPrice] = useState(0);
    const [quantity, setQuantity] = useState(0);
    const [size, setSize] = useState("");


    const doCreate = async () => {
        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/items/add`,
            {
                credentials: "include",
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    name,
                    sku,
                    category,
                    price,
                    quantity,
                    size,
                }),
            }
        );

        if (response.ok) {
            showNotification({
                color: "green",
                title: "Item created",
                message: "Item created successfully",
            });

            await refresh();
            setOpened(false);
            return;
        }

        const data = await response.json();

        showNotification({
            color: "red",
            title: "Error",
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
                title="Create Item"
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
                            placeholder="Women's All Season Short Pants"
                            onChange={(e) => setName(e.target.value)}
                        />
                    </InputWrapper>
                    <Space h="md" />
                    <InputWrapper label="SKU" required>
                        <TextInput
                            required
                            placeholder="XXXXXX"
                            onChange={(e) => setSku(e.target.value)}
                        />
                    </InputWrapper>
                    <Space h="md" />
                    <InputWrapper label="Category" required>
                        <TextInput
                            required
                            placeholder="Men, XL, Summer"
                            onChange={(e) => setCategory(e.target.value)}
                        />
                    </InputWrapper>
                    <Space h="md" />
                    <InputWrapper label="Price" required>
                        <TextInput
                            required
                            placeholder="10"
                            type="number"
                            onChange={(e) =>
                                setPrice(Number(e.target.value))
                            }
                        />
                    </InputWrapper>
                    <Space h="md" />
                    <InputWrapper label="Quantity" required>
                        <TextInput
                            required
                            placeholder="10"
                            type="number"
                            onChange={(e) =>
                                setQuantity(Number(e.target.value))
                            }
                        />
                    </InputWrapper>
                    <Space h="md" />
                    <InputWrapper label="Size" required>
                        <TextInput
                            required
                            placeholder="10/XL"
                            onChange={(e) => setSize(e.target.value)}
                        />
                    </InputWrapper>
                    <Space h="md" />
                    <Group position="right">
                        <Button color="green" type="submit">
                            Submit
                        </Button>
                    </Group>
                </form>
            </Modal>
        </>
    );
};


const EditTagsModal = ({
    opened,
    setOpened,
    refresh,
    tags,
      }: {
    opened: boolean;
    setOpened: Dispatch<SetStateAction<boolean>>;
    refresh: (search: string) => Promise<void>;
    tags: string[];
}) => {
    const [name, setName] = useState("");

    //const [tags, setTags] = useState<string[]>([]);
/*
    const fetchTags = async () => {

        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/tags/list?name=${name}`,
            {
                credentials: "include",
            }
        );

        const data: {
            success: boolean;
            data: string[];
        } = await response.json();

        if (data.success) {
            setTags(data.data);
        }
    };

    useEffect(() => {
        fetchTags();
    })*/

    return (
        <>
            <Modal
                opened={opened}
                onClose={() => {
                    refresh("");
                    setOpened(false);
                }}
                title="Edit Tags"
            >
                <Box sx={containerStyles}>
                    <Group>
                        <InputWrapper>
                            <TextInput
                                placeholder="Search Tags"
                                onChange={async (e: any) => {
                                    //await search()
                                    await refresh(e.target.value)
                                }
                                }
                            />
                        </InputWrapper>{/*
                        <Button
                            color="green"
                            onClick={() => {
                                setSearchQuery(searchQueryTyping);
                            }}
                            disabled={loading}
                        >
                            Search
                        </Button>*/}
                    </Group>
                    <Space h="md" />
                    <Table>
                        <thead>
                        <tr>
                            <th>Name</th>
                            <th>Action</th>
                        </tr>
                        </thead>
                        <tbody>
                        {tags.map((tag) => (
                            <TagRow tag={tag} />
                        ))}
                        </tbody>
                    </Table>

                </Box>
            </Modal>
        </>
    );
};

export const ItemsManager = () => {
    const [items, setItems] = useState<Item[]>([]);
    const [tags, setTags] = useState<string[]>([]);

    const [loading, setLoading] = useState<boolean>(false);

    const [currentPage, setCurrentPage] = useState(1);
    const [totalPage, setTotalPage] = useState(1);

    const [nameQuery, setNameQuery] = useState("");
    const [nameQueryTyping, setNameQueryTyping] = useState("");

    const [skuQuery, setSkuQuery] = useState("");
    const [skuQueryTyping, setSkuQueryTyping] = useState("");

    const [tagsQuery, setTagsQuery] = useState("");
    const [tagsQueryTyping, setTagsQueryTyping] = useState("");

    const [showCreationModal, setShowCreationModal] = useState(false);
    const [showTagsModal, setShowTagsModal] = useState(false);

    const [showImportModal, setShowImportModal] = useState(false);





    const exportCSV = async () => {
        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/items/export?&name=${nameQuery}&sku=${skuQuery}&tags=${tagsQuery}`,
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

    const fetchTags = async (search: string) => {

        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/tags/list?name=${search}`,
            {
                credentials: "include",
            }
        );

        const data: {
            success: boolean;
            data: string[];
        } = await response.json();

        if (data.success) {
            setTags(data.data);
        }
    };


    const fetchItems = async () => {
        setLoading(true);

        const response = await fetch(
            `${process.env.REACT_APP_API_URL}/items/list?page=${currentPage}&name=${nameQuery}&sku=${skuQuery}&tags=${tagsQuery}`,
            {
                credentials: "include",
            }
        );

        setLoading(false);

        const data: {
            success: boolean;
            message: string;
            data: {
                data: Item[];
                total: number;
                totalPages: number;
            };
        } = await response.json();

        if (data.success) {
            setItems(data.data.data);
            setTotalPage(data.data.totalPages);
        }
    };

    useEffect(() => {
        fetchItems();
        fetchTags("");
    }, [currentPage]);

    useEffect(() => {
        setCurrentPage(1);
        fetchItems();
        fetchTags("");
    }, [nameQuery, tagsQuery, skuQuery]);

    return (
        <>
            <CreateItemModal
                opened={showCreationModal}
                setOpened={setShowCreationModal}
                refresh={fetchItems}
            />
            <EditTagsModal
                opened={showTagsModal}
                setOpened={setShowTagsModal}
                refresh={fetchTags}
                tags={tags}
            />
            <UploadCSVModal
                opened={showImportModal}
                setOpened={setShowImportModal}
                refresh={fetchItems}
            />
            {/* @ts-ignore */}
            <Box sx={containerStyles}>
                <Group position="apart">
                    <h3>Items</h3>
                    <Group spacing={0}>
                        <ActionIcon
                            sx={{
                                height: "4rem",
                                width: "4rem",
                            }}
                            onClick={() => setShowCreationModal(true)}
                        >
                            <CirclePlus />
                        </ActionIcon>
                        <ActionIcon
                            sx={{
                                height: "4rem",
                                width: "4rem",
                            }}
                            onClick={() => setShowTagsModal(true)}
                        >
                            <Tags />
                        </ActionIcon>
                    </Group>
                </Group>
                <Space h="md" />


                <Group>
                    <InputWrapper>
                        <TextInput
                            placeholder="Search Items"
                            onChange={(e: any) =>
                                setNameQueryTyping(e.target.value)
                            }
                        />
                    </InputWrapper>
                    <InputWrapper>
                        <TextInput
                            placeholder="Search SKU"
                            onChange={(e: any) =>
                                setSkuQueryTyping(e.target.value)
                            }
                        />
                    </InputWrapper>
                    <MultiSelect
                        data={tags}
                        placeholder="Search Tags"
                        clearButtonLabel="Clear selection"
                        clearable
                        searchable
                        onChange={(e: any) => {
                                setTagsQueryTyping(e)
                            }
                        }
                    />
                    <Button
                        color="green"
                        onClick={() => {
                            setNameQuery(nameQueryTyping);
                            setSkuQuery(skuQueryTyping);
                            setTagsQuery(tagsQueryTyping);
                        }}
                        disabled={loading}
                    >
                        Search
                    </Button>
                </Group>


                <Space h="md" />
                <Table>
                    <thead>
                        <tr>
                            <th>Name</th>
                            <th>SKU</th>
                            <th>Category</th>
                            <th>Price</th>
                            <th>Quantity</th>
                            <th>Size</th>
                        </tr>
                    </thead>
                    <tbody>
                        {items.map((item) => (
                            <ItemRow key={item.id} item={item} />
                        ))}
                    </tbody>
                </Table>
                <Space h="md" />
                <Center>
                    {loading ? (
                        <Loader variant="dots" />
                    ) : (
                        <Pagination
                            boundaries={3}
                            withEdges
                            page={currentPage}
                            total={totalPage}
                            onChange={setCurrentPage}
                        />
                    )}
                </Center>
                <Space h="md" />
                <Group position="apart">
                    <Group spacing={0}>
                        <ActionIcon
                            sx={{
                                height: "2.5rem",
                                width: "2.5rem",
                            }}
                            onClick={exportCSV}
                        >
                            <TableExport size={"1.5rem"}/>
                        </ActionIcon>
                        <ActionIcon
                            sx={{
                                height: "2.5rem",
                                width: "2.5rem",
                            }}
                            onClick={() => {setShowImportModal(true)}}
                        >
                            <TableImport size={"1.5rem"}/>
                        </ActionIcon>
                    </Group>
                    <Button
                        onClick={fetchItems}
                        color="green"
                        disabled={loading}
                    >
                        Refresh
                    </Button>
                </Group>
            </Box>
        </>
    );
};
