import { ItemsManager } from "./inventory/Items";
import { LocationsManager } from "./inventory/Locations";
import { TransactionManager } from "./Transactions";

const Inventory = () => {
    return (
        <>
            <ItemsManager />
            <LocationsManager /> {/*
            <TransactionManager />*/}
        </>
    );
};

export default Inventory;
